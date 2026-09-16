from pathlib import Path
import re,json
root=Path('examples/discord-bots')
out=[]
for p in sorted(root.rglob('*.js')):
 for match in re.finditer(r'\b(command|subcommand|userCommand|messageCommand|component|modal|event)\(\s*["\']([^"\']+)',p.read_text()):
  out.append(dict(file=str(p),line=p.read_text()[:match.start()].count('\n')+1,kind=match[1],name=match[2]))
target=Path(__file__).resolve().parent.parent/'sources'/'discord-handlers.json'
target.write_text(json.dumps(out,indent=2)+'\n')
print(len(out),'registrations archived')
