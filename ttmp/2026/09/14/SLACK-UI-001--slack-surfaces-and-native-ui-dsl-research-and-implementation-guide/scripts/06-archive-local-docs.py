"""Archive local Markdown originals and Defuddle conversions; run from repo root."""
from pathlib import Path
import subprocess,json,hashlib
r=Path(__file__).resolve().parents[1];out=r/'sources/local-docs';out.mkdir(exist_ok=True)
paths=[Path('pkg/doc/tutorials/using-the-go-side-ui-dsl-for-discord-bots.md'),Path('ttmp/2026/09/10/DISCORD-SLACK-001--add-slack-support-to-discord-bot/design-doc/02-full-local-testing-plan-and-slack-mock-evaluation.md'),Path('/home/manuel/code/wesen/go-go-golems/go-go-parc/Research/KB/Tribal/goja-runtime-ownership-and-context-propagation.md')]
rows=[]
for i,src in enumerate(paths,1):
 name=f'{i:02d}-'+src.stem;raw=out/(name+'.original.txt');raw.write_bytes(src.read_bytes());html=out/(name+'.html');md=out/(name+'.md')
 subprocess.run(['pandoc',str(src),'-f','markdown','-t','html','-s','-o',str(html)],check=True,capture_output=True)
 subprocess.run(['defuddle','parse',str(html),'--md','-o',str(md)],check=True,capture_output=True)
 rows.append({'source':str(src.resolve()),'original':raw.name,'markdown':md.name,'original_sha256':hashlib.sha256(raw.read_bytes()).hexdigest(),'note':'Local Markdown rendered by Pandoc, then extracted by Defuddle. Original retains exact content and metadata.'})
(out/'catalog.json').write_text(json.dumps(rows,indent=2)+'\n')
print('Archived',len(rows),'local documents with originals and Defuddle conversions')
