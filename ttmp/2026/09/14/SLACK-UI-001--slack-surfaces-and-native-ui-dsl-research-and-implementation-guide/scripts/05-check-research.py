"""Check authored guide references and source hashes, without network or credentials."""
from pathlib import Path
import hashlib,json,re
r=Path(__file__).resolve().parents[1];guide=r/'design-doc/01-slack-surfaces-and-ui-dsl-intern-guide.md';text=guide.read_text();errors=[]
for p in [guide,r/'index.md',r/'sources/README.md']:
 for target in re.findall(r'\]\(([^)]+)\)',p.read_text()):
  if target.startswith(('http:','https:','#','mailto:')):continue
  path=target.split('#')[0]
  if path and not (p.parent/path).exists():errors.append(f'{p.name}: missing {path}')
used=set(re.findall(r'\[\^([^\]]+)\](?!:)',text));defined=re.findall(r'^\[\^([^\]]+)\]:',text,re.M)
if used-set(defined):errors.append('Undefined footnotes: '+str(used-set(defined)))
if len(defined)!=len(set(defined)):errors.append('Duplicate footnote definitions')
rows=json.loads((r/'sources/catalog-complete.json').read_text())
for row in rows:
 p=r/'sources'/row.get('recovered_by',row['file'])
 if not p.exists():errors.append('Missing archived source: '+str(p))
 elif row.get('archived_sha256') and hashlib.sha256((r/'sources'/row['file']).read_bytes()).hexdigest()!=row['archived_sha256']:errors.append('Hash mismatch: '+row['file'])
if errors:raise SystemExit('\n'.join(errors))
print(f'PASS: authored links, {len(defined)} footnotes, {len(set(x["url"].rstrip("/") for x in rows))} unique source URLs, and recorded source hashes')
print('This checks file integrity, not API correctness or visual rendering.')
