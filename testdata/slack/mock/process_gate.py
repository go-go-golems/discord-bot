#!/usr/bin/env python3
"""Explicit local process gate. Run both server and application in owned tmux sessions."""
import argparse
import json
import os
from pathlib import Path
import shlex
import signal
import subprocess
import sys
import time
from urllib.parse import urlparse

ROOT = Path(__file__).resolve().parents[3]


def wait_file(path, seconds):
    deadline = time.monotonic() + seconds
    while time.monotonic() < deadline:
        if path.exists():
            return path.read_text()
        time.sleep(0.05)
    raise RuntimeError(f"timed out waiting for {path.name}")


def tmux(*args):
    return subprocess.check_output(["tmux", *args], text=True).strip()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--mock", required=True)
    parser.add_argument("--bun", required=True)
    parser.add_argument("--go", required=True, help="Go executable or the ticket offline Go wrapper")
    parser.add_argument("--run-directory", required=True, type=Path)
    parser.add_argument("--child", action="store_true", help=argparse.SUPPRESS)
    args = parser.parse_args()
    run = args.run_directory.resolve()
    if args.child:
        with (run / "application.log").open("w") as log:
            process = subprocess.Popen([
                args.go, "run", "-buildvcs=false", "./cmd/slack-bot", "bots", "run-local", "ping",
                "--local-connection-file", str(run / "config.json"), "--log-level", "debug",
            ], cwd=ROOT, stdout=log, stderr=subprocess.STDOUT)
            (run / "go.pid").write_text(str(process.pid))
            code = process.wait()
            (run / "exit-code.json.tmp").write_text(json.dumps({"goRunExitCode": code}) + "\n")
            (run / "exit-code.json.tmp").rename(run / "exit-code.json")
        return
    run.mkdir(mode=0o700, parents=True, exist_ok=False)
    prefix = f"slack-process-{os.getpid()}"
    sessions = []
    app_pid = None
    try:
        mock_command = [args.bun, str(ROOT / "testdata/slack/mock/probe.ts"), args.mock, str(run), "host"]
        tmux("new-session", "-d", "-s", prefix + "-mock", shlex.join(mock_command) + "; read")
        sessions.append(prefix + "-mock")
        wait_file(run / "config.json", 5)
        child = [sys.executable, str(Path(__file__).resolve()), "--child", "--mock", args.mock,
                 "--bun", args.bun, "--go", args.go, "--run-directory", str(run)]
        tmux("new-session", "-d", "-s", prefix + "-app", shlex.join(child) + "; read")
        sessions.append(prefix + "-app")
        result = json.loads(wait_file(run / "result.json", 25))
        go_pid = wait_file(run / "go.pid", 1).strip()
        children = subprocess.check_output(["ps", "--ppid", go_pid, "-o", "pid=,comm="], text=True)
        pids = [int(line.split()[0]) for line in children.splitlines() if line.split()[1] == "slack-bot"]
        if len(pids) != 1:
            raise RuntimeError("expected exactly one slack-bot child")
        app_pid = pids[0]
        os.kill(app_pid, signal.SIGTERM)
        receipt = json.loads(wait_file(run / "exit-code.json", 5))
        app_pid = None
        shutdown = json.loads(wait_file(run / "shutdown.json", 5))
        if receipt["goRunExitCode"] != 0 or shutdown["connections"] != 0:
            raise RuntimeError("application did not terminate cleanly")
        print(json.dumps({"result": result, "shutdown": shutdown, "process": receipt}, indent=2))
        for session in sessions:
            (run / (session.rsplit("-", 1)[-1] + "-pane.log")).write_text(tmux("capture-pane", "-p", "-t", session) + "\n")
    finally:
        if app_pid is not None:
            try:
                os.kill(app_pid, signal.SIGTERM)
            except ProcessLookupError:
                pass
        config = run / "config.json"
        if config.exists():
            port = urlparse(json.loads(config.read_text())["apiURL"]).port
            subprocess.run(["lsof-who", "-p", str(port), "-k"], check=False)
        for session in reversed(sessions):
            subprocess.run(["tmux", "kill-session", "-t", session], check=False)


if __name__ == "__main__":
    main()
