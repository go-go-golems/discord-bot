from pathlib import Path
import argparse
parser = argparse.ArgumentParser()
parser.add_argument('--input', default='design-doc/01-slack-support-architecture-and-intern-implementation-guide.md')
parser.add_argument('--output-name', default='01-Slack Support - Intern Implementation Guide.md')
args = parser.parse_args()
root=Path(__file__).resolve().parents[1]
s=(root/args.input).read_text()
# A print derivative: prevent code and ASCII diagrams splitting across pages,
# and convert tables to records so long code identifiers cannot overlap columns.
s=s.replace('../sources/', '../../sources/')
lines=s.splitlines(); out=[]; i=0; fence=False
while i<len(lines):
    line=lines[i]
    if line.startswith('```'):
        if not fence: out.extend(['','```{=latex}',r'\begin{samepage}','```','',line])
        else: out.extend([line,'','```{=latex}',r'\end{samepage}','```',''])
        fence=not fence; i+=1; continue
    if not fence and line.startswith('|') and i+1<len(lines) and lines[i+1].startswith('|---'):
        headers=[v.strip() for v in line.strip('|').split('|')]
        i+=2
        while i<len(lines) and lines[i].startswith('|'):
            values=[v.strip() for v in lines[i].strip('|').split('|')]
            out.extend(['',f'**{values[0]}**',''])
            for h,v in zip(headers[1:],values[1:]): out.append(f'- {h}: {v}')
            i+=1
        out.append(''); continue
    out.append(line); i+=1
outdir=root/'various'/'print'; outdir.mkdir(parents=True,exist_ok=True)
p=outdir/args.output_name
p.write_text('\n'.join(out)+'\n')
print(p)
