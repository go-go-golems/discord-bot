from pathlib import Path
import subprocess,concurrent.futures,json,hashlib
root=Path('ttmp/2026/09/14/SLACK-UI-001--slack-surfaces-and-native-ui-dsl-research-and-implementation-guide')
paths=['surfaces/','surfaces/app-home/','surfaces/modals/','surfaces/canvases/','surfaces/lists/','surfaces/split-view/','block-kit/','reference/block-kit/blocks/','reference/block-kit/block-elements/','reference/block-kit/composition-objects/','interactivity/handling-user-interaction/','apis/events-api/using-socket-mode/','reference/interaction-payloads/block_actions-payload/','reference/interaction-payloads/block_suggestion-payload/','reference/interaction-payloads/view-interactions-payload/','reference/views/modal-views/','reference/methods/chat.postMessage/','reference/methods/chat.update/','reference/methods/chat.postEphemeral/','reference/methods/views.open/','reference/methods/views.update/','reference/methods/views.publish/','reference/block-kit/blocks/input-block/','reference/block-kit/blocks/section-block/','reference/block-kit/blocks/actions-block/','reference/block-kit/block-elements/button-element/','reference/block-kit/block-elements/select-menu-element/','reference/block-kit/blocks/table-block/','reference/block-kit/blocks/markdown-block/','reference/block-kit/blocks/task-card-block/','reference/block-kit/blocks/plan-block/','reference/block-kit/blocks/rich-text-block/','reference/block-kit/blocks/data-table-block/','reference/block-kit/blocks/data-visualization-block/','reference/block-kit/blocks/card-block/','reference/block-kit/blocks/carousel-block/','reference/block-kit/blocks/container-block/','reference/block-kit/blocks/alert-block/','reference/block-kit/blocks/context-actions-block/','reference/methods/files.getUploadURLExternal/','reference/methods/slackLists.create','reference/methods/canvases.create/','reference/events/app_home_opened/','reference/app-manifest/','messaging/work-objects/']
paths += ['ai/developing-agents/','ai/migrating-to-agent-messaging/','messaging/work-objects-overview/','messaging/work-objects-implementation/','reference/methods/files.completeUploadExternal/','reference/methods/chat.unfurl/','reference/methods/views.push/','reference/block-kit/blocks/video-block/','reference/block-kit/block-elements/plain-text-input-element/','reference/block-kit/composition-objects/text-object/','reference/block-kit/composition-objects/confirmation-dialog-object/','interactivity/implementing-shortcuts/','reference/methods/chat.startStream/','reference/methods/assistant.threads.setStatus/']
(root/'sources/web').mkdir(parents=True,exist_ok=True)
def capture(item):
 i,path=item; url='https://docs.slack.dev/'+path; name=f'{i:02d}-'+path.strip('/').replace('/','-')+'.md'; out=root/'sources/web'/name
 if out.exists() and out.stat().st_size: p=subprocess.CompletedProcess([],0,'','')
 else:
  try: p=subprocess.run(['defuddle','parse',url,'--md','-o',str(out)],capture_output=True,text=True,timeout=35)
  except subprocess.TimeoutExpired: p=subprocess.CompletedProcess([],124,'','Defuddle timeout after 35 seconds')
 data=out.read_bytes() if out.exists() else b''
 return {'id':i,'url':url,'file':'web/'+name,'bytes':len(data),'exit':p.returncode,'sha256':hashlib.sha256(data).hexdigest(),'diagnostic':p.stderr[-400:] if p.returncode else ''}
with concurrent.futures.ThreadPoolExecutor(max_workers=6) as ex:
 rows=list(ex.map(capture,enumerate(paths,1)))
(root/'sources/catalog.json').write_text(json.dumps(rows,indent=2)+'\n')
for row in rows: print(row['id'],row['exit'],row['bytes'],row['file'])
