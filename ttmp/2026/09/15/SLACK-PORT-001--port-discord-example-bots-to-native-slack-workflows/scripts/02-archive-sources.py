from pathlib import Path
from concurrent.futures import ThreadPoolExecutor
import subprocess,json,hashlib
root=Path(__file__).resolve().parent.parent/'sources'
methods=['conversations.history','conversations.replies','chat.update','chat.delete','chat.postEphemeral','users.info','usergroups.users.update','conversations.kick','conversations.setTopic','pins.add','pins.remove','pins.list','files.getUploadURLExternal','files.completeUploadExternal','reactions.add','admin.users.remove']
urls=[f'https://docs.slack.dev/reference/methods/{m}/' for m in methods]+['https://docs.slack.dev/interactivity/implementing-shortcuts/','https://docs.slack.dev/reference/block-kit/block-elements/','https://docs.slack.dev/reference/events/message/','https://docs.slack.dev/reference/events/reaction_added/','https://docs.slack.dev/reference/interaction-payloads/block_suggestion-payload/','https://docs.slack.dev/apis/events-api/using-socket-mode/']
def save(pair):
 n,url=pair
 path=root/f'{n:02d}-{url.rstrip("/").split("/")[-1]}.md'
 r=subprocess.run(['defuddle','parse',url,'--md','-o',str(path)],capture_output=True,text=True)
 if r.returncode: return {'url':url,'error':r.stderr[-200:]}
 text=path.read_text()
 path.write_text(f'Source: {url}\nRetrieved: 2026-09-15\n\n'+text)
 return dict(url=url,file=path.name,sha256=hashlib.sha256(path.read_bytes()).hexdigest())
with ThreadPoolExecutor(max_workers=4) as pool: result=list(pool.map(save,enumerate(urls,1)))
(root/'index.json').write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps({'saved':sum('file' in r for r in result),'failures':[r for r in result if 'error' in r]},indent=2))
