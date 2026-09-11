# Local SDK interoperability gate

This probe runs the actual pinned Go Slack SDK against the separately checked-out
Bun mock. It verifies identity, Socket Mode hello and two deliveries, exact ACKs
as observed by the mock, a threaded Web API post, an ephemeral response URL post,
and socket cancellation. It does not yet execute the application JavaScript host
or prove production permission enforcement, retry policy, or precise ACK timing.

Versions are recorded in `revisions.json`; prepare those dependencies before
running offline. The tested mock checkout is
`/home/manuel/code/others/slack/desplega-slack-mock`. No installation happens in
the probe. The mock launcher imports `src/server.ts` directly.

Run from the repository root, using a fresh run directory each time:

```sh
tmux new-session -d -s slack-sdk-probe \
  '/home/manuel/.bun/bin/bun testdata/slack/mock/probe.ts /home/manuel/code/others/slack/desplega-slack-mock /tmp/slack-sdk-probe-NEW; read'
# Wait for config.json, without printing its synthetic credentials.
go test -count=1 -v ./internal/slackprobe -run TestSDKInteroperability \
  -args -slack-probe-directory /tmp/slack-sdk-probe-NEW
tmux capture-pane -p -t slack-sdk-probe
```

On this workstation, use the ticket's `scripts/04-go-offline.sh` instead of `go`
to select the cached matching toolchain and disable downloads. The launcher has a
45-second outer deadline; the Go test has a 30-second operation context and waits
for explicit result/shutdown receipts. The server normally stops itself after
client disconnection. Clean the owned port with `lsof-who -p PORT -k`, then remove
the owned tmux session with `tmux kill-session -t slack-sdk-probe`.

The synthetic `config.json` stays outside Git with mode 0600. Do not archive it.
`result.json` and `shutdown.json` are sanitized evidence suitable for the ticket.
Writes use atomic rename so readiness cannot expose a partially written file.

Ordinary `go test ./...` explicitly skips the external-server test unless the
run-directory argument is supplied; it still runs the destination rejection
test. Therefore a successful ordinary suite is **not** a successful mock gate.
An explicitly selected nonexistent directory fails, rather than silently skipping.
Future CI must run this gate separately and fail when its prepared runtime is absent.

The Go HTTP and WebSocket dialers reject all non-literal/non-loopback addresses
before DNS or dialing. HTTP proxies are unset in the constructed transport,
redirects are rejected, and the same HTTP client handles response URLs. Only
synthetic credentials are loaded from the explicit file; no environment values
are read by the probe.

## Complete JavaScript application scenario

Pass `host` as the launcher's third argument to assert the actual example bot's
reply strings. In a second tmux session, run:

```sh
go run ./cmd/slack-bot bots run-local ping \
  --local-connection-file /tmp/slack-host-probe-NEW/config.json --log-level debug
```

Wait for `result.json`, then send SIGINT (tmux `send-keys ... C-c`) or SIGTERM to
the application. The mock writes `shutdown.json` only after observing zero socket
connections. Capture output before closing the sessions. The launcher starts with
a fresh store and emits one human mention and one slash command; `host` mode
asserts the example's threaded acknowledgment and private `pong` response. This is
a baseline process scenario, not the complete P01–P10 catalog.

Strict Go fixtures in `internal/slacktransport` additionally verify app/bot token
routing, request contents, private response bodies, invalid-input zero traffic,
429/error mapping, accepted-then-lost-response behavior without retry, and ACKs
for duplicate business events while a handler blocks. The WebSocket fixture
validates the SDK's fixed `Origin: https://api.slack.com` header explicitly.

For a repeatable full process gate with explicit exit-status capture, use:

```sh
python3 testdata/slack/mock/process_gate.py \
  --mock /home/manuel/code/others/slack/desplega-slack-mock \
  --bun /home/manuel/.bun/bin/bun \
  --go /absolute/path/to/the/ticket/scripts/04-go-offline.sh \
  --run-directory /tmp/slack-host-probe-NEW
```

The directory must be new. The driver starts both processes in tmux, reads only
synthetic settings, waits for business assertions, sends SIGTERM to the actual
application child, checks its Go runner's exit code and mock connection count,
and cleans owned sessions and the mock port. It does not depend on tmux's
`pane_dead_status`, which was empty on this workstation despite process exit.
Archive sanitized results and logs, never `config.json` or `go.pid`.
