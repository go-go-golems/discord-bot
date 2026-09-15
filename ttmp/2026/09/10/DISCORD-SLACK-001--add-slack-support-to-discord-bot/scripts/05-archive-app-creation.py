"""Archive app-setup references from this conversation without altering earlier snapshots."""
import concurrent.futures
import datetime
import hashlib
import json
from pathlib import Path
import subprocess

ROOT = Path(__file__).resolve().parents[1] / "sources" / "app-creation"
ITEMS = {
    "01-configuration-tokens": "https://docs.slack.dev/app-manifests/configuring-apps-with-app-manifests/",
    "02-create-api": "https://docs.slack.dev/reference/methods/apps.manifest.create/",
    "03-validate-api": "https://docs.slack.dev/reference/methods/apps.manifest.validate/",
    "04-rotate-token-api": "https://docs.slack.dev/reference/methods/tooling.tokens.rotate/",
    "05-app-setup": "https://docs.slack.dev/tools/bolt-python/creating-an-app/",
    "06-cli-quickstart": "https://docs.slack.dev/quickstart/",
    "07-cli-api": "https://docs.slack.dev/tools/slack-cli/reference/commands/slack_api/",
    "08-cli-install": "https://docs.slack.dev/tools/slack-cli/reference/commands/slack_app_install/",
    "09-cli-project-create": "https://docs.slack.dev/tools/slack-cli/reference/commands/slack_project_create/",
    "10-cli-manifest-validate": "https://docs.slack.dev/tools/slack-cli/reference/commands/slack_manifest_validate/",
    "11-cli-manifest-info": "https://docs.slack.dev/tools/slack-cli/reference/commands/slack_manifest_info/",
    "12-cli-commands": "https://docs.slack.dev/tools/slack-cli/guides/running-slack-cli-commands/",
    "13-cli-api-release": "https://docs.slack.dev/changelog/2026/05/19/slack-cli/",
    "14-app-settings-quickstart": "https://docs.slack.dev/app-management/quickstart-app-settings/",
    "15-cli-overview": "https://docs.slack.dev/tools/slack-cli/",
    "16-oauth-installation": "https://docs.slack.dev/authentication/installing-with-oauth/",
    "17-oauth-access-api": "https://docs.slack.dev/reference/methods/oauth.v2.access/",
    "18-oauth-token-rotation": "https://docs.slack.dev/authentication/using-token-rotation/",
}


def collect(item):
    name, url = item
    path = ROOT / (name + ".md")
    if path.exists():
        previous = ROOT / "manifest.json"
        if previous.exists():
            for entry in json.loads(previous.read_text()):
                if entry.get("name") == name and entry.get("sha256") == hashlib.sha256(path.read_bytes()).hexdigest():
                    return entry
    process = subprocess.run(["defuddle", "parse", url, "--md", "-o", str(path)],
                             capture_output=True, text=True, timeout=90)
    if process.returncode or not path.exists() or path.stat().st_size < 200:
        return {"name": name, "url": url, "error": process.stderr[-1000:]}
    data = path.read_bytes()
    return {"name": name, "url": url, "bytes": len(data),
            "sha256": hashlib.sha256(data).hexdigest(),
            "fetched_at": datetime.datetime.now(datetime.timezone.utc).isoformat()}


if __name__ == "__main__":
    ROOT.mkdir(parents=True, exist_ok=True)
    with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
        results = list(pool.map(collect, ITEMS.items()))
    (ROOT / "manifest.json").write_text(json.dumps(results, indent=2) + "\n")
    lines = ["# App creation and Slack CLI references", "",
             "Public documentation extracted with defuddle. Includes sources from earlier app-setup and CLI answers.",
             "Retrieval dates and hashes are in manifest.json. No credentials or authenticated pages are archived.", ""]
    for result in results:
        name = result["name"]
        lines.append(f'- [{name}]({name}.md): {result["url"]}')
        print(name, result.get("bytes", result.get("error")))
    lines += ["", "Previously archived Socket Mode, manifest schema, scope, and event references remain in the parent sources folder."]
    (ROOT / "README.md").write_text("\n".join(lines) + "\n")
    if any("error" in result for result in results):
        raise SystemExit(1)
