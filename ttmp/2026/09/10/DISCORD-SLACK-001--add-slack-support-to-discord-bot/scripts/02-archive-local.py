from pathlib import Path
import subprocess, hashlib, json, shutil
root=Path(__file__).resolve().parents[1]
vault=Path('/home/manuel/code/wesen/go-go-golems/go-go-parc')
notes=[
('20-dsl-design','Projects/2026/06/22/ARTICLE - Designing DSLs with go-go-goja - Go-Backed JavaScript APIs.md'),
('21-discord-framework','Projects/2026/04/22/PROJ - JS Discord Bot Framework.md'),
('22-discord-ui-dsl','Projects/2026/04/22/ARTICLE - Go-Side JavaScript DSLs for Discord Bots - Types, Errors, and In-Place Updates.md'),
('23-xgoja-discord','Projects/2026/05/25/ARTICLE - xgoja Modules in Existing Runners - Discord Bot Case Study.md'),
('24-context-management','Projects/2026/05/15/ARTICLE - go-go-goja Context Management - Runtime Request and Async Call Context.md'),
]
records=[]
for name,rel in notes:
    src=vault/rel
    original=root/'sources'/f'{name}.original.md'
    shutil.copyfile(src,original)
    html=root/'sources'/f'{name}.html'
    subprocess.run(['pandoc',str(src),'-f','markdown','-t','html','-s','-o',str(html)],check=True,capture_output=True)
    subprocess.run(['defuddle','parse',str(html),'--md','-o',str(root/'sources'/f'{name}.md')],check=True,capture_output=True)
    html.unlink()
    records.append(dict(path=str(src),snapshot=original.name,extracted=f'{name}.md',sha256=hashlib.sha256(original.read_bytes()).hexdigest()))
repo=root.parents[4]
files=['README.md','go.mod','internal/bot/bot.go','internal/config/config.go','internal/jsdiscord/host.go','internal/jsdiscord/runtime.go','internal/jsdiscord/bot_dispatch.go','internal/jsdiscord/bot_context.go','internal/jsdiscord/descriptor.go','internal/jsdiscord/host_responses.go','internal/jsdiscord/ui_message.go','internal/jsdiscord/store.go','pkg/botcli/discover.go','pkg/botcli/command_run.go','pkg/botcli/runtime_helpers.go','pkg/botcli/run_description.go','pkg/botcli/command_root.go','pkg/framework/framework.go','pkg/xgoja/provider/provider.go']
for rel in files:
    src=repo/rel
    dest=root/'sources'/'code'/(rel+'.txt')
    dest.parent.mkdir(parents=True,exist_ok=True)
    shutil.copyfile(src,dest)
    records.append(dict(path=str(src),snapshot=str(dest.relative_to(root/'sources')),sha256=hashlib.sha256(dest.read_bytes()).hexdigest()))
(root/'sources'/'local-manifest.json').write_text(json.dumps(records,indent=2)+'\n')
with (root/'sources'/'README.md').open('a') as f:
    f.write('\n## Local evidence\n\nVault Markdown is preserved byte-for-byte as `.original.md`; a Pandoc HTML conversion was passed through defuddle for the companion extracted `.md`. Originals preserve code fences and Obsidian metadata. `local-manifest.json` records paths and hashes. Code snapshots are verbatim, not transformed.\n\n')
    for name,rel in notes: f.write(f'- [{name}]({name}.md) — [original]({name}.original.md)\n')
    f.write('\nDiscord code snapshots: `code/`, revision `2ea219e6d0a37bb27916dc29a2ba88c1f03b6182`.\n')
print('Archived',len(records),'local sources')
