from pathlib import Path
import json
root=Path(__file__).resolve().parent.parent
rows=json.loads((root/'sources/discord-handlers.json').read_text())
events={'messageCreate':'message','messageUpdate':'message_changed','messageDelete':'message_deleted','reactionAdd':'reaction_added','reactionRemove':'reaction_removed','guildMemberAdd':'team_join','guildMemberUpdate':'user_change','guildMemberRemove':'user_change'}
renames={'ping':{'ping':'golem-ping','echo':'golem-echo','feedback':'golem-feedback','search':'golem-search','announce':'golem-announce'},'show-space':{'debug-roles':'debug-groups','debug-my-roles':'debug-my-permissions'},'moderation':{'mod-fetch-guild':'mod-fetch-workspace','mod-list-roles':'mod-list-groups','mod-fetch-role':'mod-fetch-group','mod-add-role':'mod-add-group-member','mod-kick':'mod-remove-workspace-user'}}
ui={'showcase:msg:primary':'demo.primary','showcase:msg:success':'demo.neutral','showcase:msg:danger':'demo.danger','showcase:msg:select':'demo.message.select','showcase:confirm:yes':'demo.confirm.yes','showcase:confirm:no':'demo.confirm.no','showcase:buy:confirm':'demo.card.buy','showcase:select:string':'demo.select.string','showcase:select:user':'demo.select.user','showcase:select:role':'demo.select.group','showcase:select:channel':'demo.select.channel','showcase:select:mentionable':'demo.select.users'}
expect=[];lines=[]
for row in rows:
 bot=row['file'].split('/')[2].removesuffix('.js');kind=row['kind'];name=row['name'];target='';note='Ported; registration covered, workflow tests listed below.';native=kind
 if kind=='event':
  if name in ('ready','guildCreate'): native='host';target='runtime load / Socket Mode authentication logs';note='Native lifecycle equivalent; no Discord gateway event emulation.'
  else:target=events[name]
 elif kind=='__verb__':
  if bot=='unified-demo' and name=='status':native='verb';target='status'
  else:native='host';target='bots run '+bot;note='Host-managed runtime configuration uses declared fields and the shared CLI.'
 elif bot=='moderation' and name in ('mod-set-slowmode','mod-timeout','mod-ban','mod-unban'):
  native='platform';target='No direct equivalent';note='No same-semantics public Slack bot API; do not substitute channel removal or organization deactivation.'
 elif bot=='support' and name in ('support-join-thread','support-leave-thread'):
  native='platform';target='No thread membership API';note='Native support-join-channel / support-leave-channel are explicitly different operations.'
 elif kind in ('command','subcommand'):
  native='command';target='/'+renames.get(bot,{}).get(name,name)
  if name=='mod-kick':note='Conditional Enterprise admin.users.remove; separate user token, actor allowlist and opt-in.'
  if 'role' in name:note='Native user-group membership workflow; Slack groups do not grant Discord role permissions.'
 elif kind=='userCommand':native='command';target='/show-avatar';note='Slack has no user context-menu registration; a native user selector preserves the avatar workflow.'
 elif kind=='messageCommand':native='shortcut';target={'Archive Thread':'archive-thread','Quote Message':'quote-message'}[name]
 elif kind=='autocomplete':
  native='options';target={'ping':'ping.search.query','custom-kb':'kb.suggest','knowledge-base':'knowledge.suggest','ui-showcase':'demo.search.suggest'}[bot];note='Native external select suggestions replace unsupported slash-option autocomplete.'
 elif kind=='modal':
  native='view';target={'ping':'ping.feedback','knowledge-base':'knowledge.submit','ui-showcase':'demo.form'}.get(bot,name.replace(':','.'))
 elif kind=='component':
  native='action'
  if bot=='poker':target=name.replace('poker:help:','poker.')
  elif bot=='show-space':target=name.replace('show-space:debug:','show.debug.').replace('.guild','.workspace')
  elif bot=='ui-showcase':
   if name=='showcase:buy:cancel':native='platform';target='Native confirmation cancellation';note='Cancelling Slack confirmation dismisses it without dispatching the destructive action.'
   else:target=ui[name]
  else:target=name.replace(':','.')
 else:raise ValueError(row)
 if native not in ('host','platform'):expect.append(dict(bot=bot,kind=native,name=target,source=f"{row['file']}:{row['line']}"))
 lines.append(f"| `{row['file']}:{row['line']}` | {kind} `{name}` | {native}: `{target}` | {note} |")
# Computed registrations that a literal-string regex does not capture.
extra={
 'custom-kb':{'action':['kb.select','kb.refresh'],'view':['kb.add','kb.search'],'options':['kb.suggest']},
 'knowledge-base':{'action':['knowledge.select','knowledge.verify','knowledge.stale','knowledge.reject','knowledge.source','knowledge.edit','knowledge.export','knowledge.previous','knowledge.next'],'view':['knowledge.submit','knowledge.open'],'options':['knowledge.suggest']},
 'ui-showcase':{'action':['demo.article.select','demo.article.verify','demo.article.stale','demo.article.reject','demo.article.edit','demo.article.source','demo.article.export','demo.article.details','demo.search.previous','demo.search.next','demo.pager.previous','demo.pager.next','demo.card.select','demo.card.buy','demo.card.info','demo.card.share'],'view':['demo.article.save','demo.search.open'],'options':['demo.search.suggest']},
}
for bot,groups in extra.items():
 for kind,names in groups.items():
  for name in names:expect.append(dict(bot=bot,kind=kind,name=name,source='Computed source registration; manual supplement'))
(root/'sources/port-registration-expectations.json').write_text(json.dumps(expect,indent=2)+'\n')
front='''---
Title: Source handler acceptance matrix
Ticket: SLACK-PORT-001
Status: active
Topics: [slack]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Explicit native mappings and local acceptance evidence; live qualification remains separate.
LastUpdated: 2026-09-15T23:00:00-04:00
WhatFor: Review parity without conflating registration, workflow tests and live Slack behavior.
WhenToUse: Reviewing any example port.
---

# Source handler acceptance matrix

This matrix maps literal source registrations and supplements computed registrations.
`TestSourceRegistrationParity` checks all mapped routes against inspected Slack descriptors.
Registration coverage is not a claim that every possible payload or Slack plan has
been tested. The workflow tests below verify substantive state and service behavior;
new ports have not been installed in the live workspace.

| Source | Registration | Native destination | Disposition |
| --- | --- | --- | --- |
'''
body=front+'\n'.join(lines)+'\n\n## Computed registrations\n\n'
for bot,groups in extra.items():
 body+='- **'+bot+'**: '+ '; '.join(kind+' '+', '.join('`'+n+'`' for n in names) for kind,names in groups.items())+'.\n'
body+='''

## Executable workflow evidence

- `pkg/slackcli/ports_test.go`: baseline command behavior, Poker deal/draw-once/rank/reset/user isolation, Support private follow-up and thread parent reference.
- `pkg/slackcli/port_workflows_test.go`: Custom KB upsert/reopen/workspace isolation, reviewer denial, Show Space manager denial and exact pin reference, paginated archive with attachments, moderation denial, UI review state and pager updates, modal validation, Ping modal and dynamic suggestions.
- `internal/jsslack/interactions_test.go`: message shortcuts and source content, command triggers, suggestion ACK payloads, modal-result ACK updates and original-message replacement.
- `internal/jsslack/database_test.go`: inspection does not open SQLite; runtime data survives independent host lifetimes.
- `internal/slacktransport/operations_test.go`: pagination cursors, stable permission failures, upload sequence and no bot token on content transfer, separate administrative token.
- `internal/slacktransport/port_decode_test.go`: user selections, shortcuts, option payloads, edits/deletes and user lifecycle decoding.
- `cmd/slack-bot/main_test.go`: CLI local verb invocation and error handling.

## Platform distinctions

Slack has no native Discord role permission hierarchy, independently joined threads,
colored embeds or green success-button style. User groups, parent-message threads,
Block Kit sections and neutral/primary/danger controls preserve supported workflows
without claiming those semantics are identical. Slash command names for Ping use
`/golem-*`; other bots share some names and should normally use separate installed
apps or a deliberately selected single app during development.

Workspace removal is implemented only as an explicit Enterprise-capable operation.
Discord bans/unbans, timeouts and slowmode are not silently approximated. User-group
membership updates preserve existing members before replacing Slack's full list;
a concurrent external edit can still race this local read-modify-write operation.
No distributed coordination mechanism is introduced.
'''
(root/'reference/02-source-handler-acceptance-matrix.md').write_text(body)
Path('pkg/slackcli/testdata/port-registration-expectations.json').write_text(json.dumps(expect,indent=2)+'\n')
print(len(rows),'literal source registrations;',len(expect),'mapped native registration expectations')
