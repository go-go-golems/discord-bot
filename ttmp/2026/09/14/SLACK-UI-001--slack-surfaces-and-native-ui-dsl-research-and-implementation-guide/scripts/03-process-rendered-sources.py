from pathlib import Path
import json,subprocess,concurrent.futures,hashlib
from bs4 import BeautifulSoup
r=Path(__file__).resolve().parents[1];web=r/'sources/web'; catalog=json.loads((r/'sources/catalog.json').read_text()); urls={x['url'].rstrip('/') for x in catalog}; extra=set()
for p in web.glob('rendered-*.json'):
 d=json.loads(p.read_text()); html=p.with_suffix('.html');html.write_text(d['html']);p.with_suffix('.txt').write_text(d['text'])
 subprocess.run(['defuddle','parse',str(html),'--md','-o',str(p.with_suffix('.md'))],check=True,capture_output=True)
 for a in BeautifulSoup(d['html'],'html.parser').find_all('a',href=True):
  url=a['href'];url='https://docs.slack.dev'+url if url.startswith('/') else url
  if '/reference/block-kit/' in url and url.rstrip('/') not in urls and '#' not in url: extra.add(url.rstrip("/"))
for url in ['https://docs.slack.dev/reference/methods/agents.sessions.setStatus/','https://docs.slack.dev/ai/agent-sessions/']:
 extra.add(url.rstrip("/"))
def capture(url):
 name='extra-'+url.rstrip('/').split('/')[-1]+'.md';p=web/name
 if not p.exists():
  try: result=subprocess.run(['defuddle','parse',url,'--md','-o',str(p)],capture_output=True,text=True,timeout=30);status=result.returncode
  except subprocess.TimeoutExpired:status=124
 else:status=0
 return {'url':url,'file':'web/'+name,'exit':status,'bytes':p.stat().st_size if p.exists() else 0}
with concurrent.futures.ThreadPoolExecutor(max_workers=6) as pool: results=list(pool.map(capture,sorted(extra)))
(r/'sources/extra-catalog.json').write_text(json.dumps(results,indent=2)+'\n')
for x in results:print(x['exit'],x['bytes'],x['url'])
