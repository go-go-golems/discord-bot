from pathlib import Path
import urllib.request,re,subprocess,concurrent.futures
r=next(Path('ttmp/2026/09/14').glob('SLACK-UI-001*')); out=r/'sources/web'
urls=['https://docs.slack.dev/reference/block-kit/blocks/','https://docs.slack.dev/reference/block-kit/block-elements/','https://docs.slack.dev/reference/block-kit/composition-objects/']
for i,url in enumerate(urls,1):
 data=urllib.request.urlopen(url,timeout=25).read(); (out/f'catalog-{i}.html').write_bytes(data)
 links=sorted(set(re.findall(r'href="(/reference/block-kit/[^"#?]+)"',data.decode())))
 (out/f'catalog-{i}-links.txt').write_text('\n'.join('https://docs.slack.dev'+x for x in links)+'\n')
 print(i,len(data),len(links))
