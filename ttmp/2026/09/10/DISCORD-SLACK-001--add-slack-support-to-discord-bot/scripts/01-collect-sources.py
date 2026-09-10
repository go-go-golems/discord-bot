from pathlib import Path
import concurrent.futures, subprocess, hashlib, json, datetime
root = Path(__file__).resolve().parents[1]
items = [
('01-socket-mode','https://docs.slack.dev/apis/events-api/using-socket-mode/'),
('02-connections-scope','https://docs.slack.dev/reference/scopes/connections.write/'),
('03-slash-commands','https://docs.slack.dev/interactivity/implementing-slash-commands/'),
('04-post-message','https://docs.slack.dev/reference/methods/chat.postMessage/'),
('05-interactions','https://docs.slack.dev/interactivity/handling-user-interaction/'),
('06-rate-limits','https://docs.slack.dev/apis/web-api/rate-limits/'),
('07-app-manifest','https://docs.slack.dev/reference/app-manifest/'),
('08-app-mention','https://docs.slack.dev/reference/events/app_mention/'),
('09-request-signing','https://docs.slack.dev/authentication/verifying-requests-from-slack/'),
('10-slack-go','https://github.com/slack-go/slack'),
('11-auth-test','https://docs.slack.dev/reference/methods/auth.test/'),
('12-chat-update','https://docs.slack.dev/reference/methods/chat.update/'),
('13-views-open','https://docs.slack.dev/reference/methods/views.open/'),
('14-message-im','https://docs.slack.dev/reference/events/message.im/'),
]
def collect(item):
    name,url=item
    path=root/'sources'/f'{name}.md'
    p=subprocess.run(['defuddle','parse',url,'--md','-o',str(path)],capture_output=True,text=True,timeout=90)
    if p.returncode or not path.exists() or path.stat().st_size < 200:
        return dict(name=name,url=url,error=p.stderr[-1500:])
    data=path.read_bytes()
    return dict(name=name,url=url,bytes=len(data),sha256=hashlib.sha256(data).hexdigest(),fetched_at=datetime.datetime.now(datetime.timezone.utc).isoformat())
with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
    results=list(pool.map(collect,items))
(root/'sources'/'manifest.json').write_text(json.dumps(results,indent=2)+'\n')
lines=['# Archived research sources','','HTML pages extracted with `defuddle parse URL --md -o FILE`. Retrieval metadata and SHA-256 hashes are in `manifest.json`. These are research snapshots; linked upstream pages remain authoritative.','','| Snapshot | Original URL | Bytes |','|---|---|---|']
for r in results:
    print(r['name'],r.get('bytes',r.get('error')))
    lines.append(f"| [{r['name']}]({r['name']}.md) | {r['url']} | {r.get('bytes','FAILED')} |")
(root/'sources'/'README.md').write_text('\n'.join(lines)+'\n')
if any('error' in r for r in results): raise SystemExit(1)
