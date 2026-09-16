---
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
| `examples/discord-bots/announcements.js:10` | command `announce-preview` | command: `/announce-preview` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/announcements.js:26` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/archive-helper/index.js:21` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/archive-helper/index.js:27` | command `archive-channel` | command: `/archive-channel` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/archive-helper/index.js:101` | messageCommand `Archive Thread` | shortcut: `archive-thread` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/custom-kb/index.js:23` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/custom-kb/index.js:28` | command `kb-add` | command: `/kb-add` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/custom-kb/index.js:48` | command `kb-link` | command: `/kb-link` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/custom-kb/index.js:62` | command `kb-search` | command: `/kb-search` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/custom-kb/index.js:67` | command `kb-list` | command: `/kb-list` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/custom-kb/index.js:78` | autocomplete `kb-search` | options: `kb.suggest` | Native external select suggestions replace unsupported slash-option autocomplete. |
| `examples/discord-bots/hater/index.js:38` | command `hate` | command: `/hate` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/hater/index.js:60` | command `roast` | command: `/roast` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/hater/index.js:76` | command `compliment` | command: `/compliment` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/hater/index.js:90` | component `hater:roast` | action: `hater.roast` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/hater/index.js:97` | component `hater:mercy` | action: `hater.mercy` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/hater/index.js:105` | component `hater:apology` | action: `hater.apology` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/hater/index.js:114` | modal `hater:apology:submit` | view: `hater.apology.submit` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/hater/index.js:131` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/hater/index.js:138` | event `messageCreate` | event: `message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/interaction-types/index.js:11` | command `hello` | command: `/hello` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/interaction-types/index.js:18` | command `echo` | command: `/echo` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/interaction-types/index.js:28` | command `fun` | command: `/fun` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/interaction-types/index.js:48` | subcommand `fun` | command: `/fun` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/interaction-types/index.js:59` | subcommand `fun` | command: `/fun` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/interaction-types/index.js:70` | userCommand `Show Avatar` | command: `/show-avatar` | Slack has no user context-menu registration; a native user selector preserves the avatar workflow. |
| `examples/discord-bots/interaction-types/index.js:85` | messageCommand `Quote Message` | shortcut: `quote-message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/interaction-types/index.js:98` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/knowledge-base/index.js:18` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/knowledge-base/index.js:28` | event `messageCreate` | event: `message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:40` | command `remember` | command: `/remember` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:46` | command `teach` | command: `/teach` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:52` | modal `knowledge:submit` | view: `knowledge.submit` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:64` | command `ask` | command: `/ask` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:83` | command `kb-search` | command: `/kb-search` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:102` | command `article` | command: `/article` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:124` | command `kb-article` | command: `/kb-article` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:146` | command `review` | command: `/review` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:169` | command `kb-review` | command: `/kb-review` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:192` | command `recent` | command: `/recent` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:206` | command `kb-recent` | command: `/kb-recent` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:220` | command `kb-verify` | command: `/kb-verify` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:243` | command `kb-stale` | command: `/kb-stale` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:266` | command `kb-reject` | command: `/kb-reject` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/knowledge-base/index.js:427` | autocomplete `ask` | options: `knowledge.suggest` | Native external select suggestions replace unsupported slash-option autocomplete. |
| `examples/discord-bots/knowledge-base/index.js:432` | autocomplete `kb-search` | options: `knowledge.suggest` | Native external select suggestions replace unsupported slash-option autocomplete. |
| `examples/discord-bots/knowledge-base/index.js:437` | autocomplete `article` | options: `knowledge.suggest` | Native external select suggestions replace unsupported slash-option autocomplete. |
| `examples/discord-bots/knowledge-base/index.js:442` | autocomplete `kb-article` | options: `knowledge.suggest` | Native external select suggestions replace unsupported slash-option autocomplete. |
| `examples/discord-bots/knowledge-base/index.js:531` | __verb__ `run` | host: `bots run knowledge-base` | Host-managed runtime configuration uses declared fields and the shared CLI. |
| `examples/discord-bots/knowledge-base/lib/reactions.js:2` | event `reactionAdd` | event: `reaction_added` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-channel-moderation-commands.js:2` | command `mod-fetch-channel` | command: `/mod-fetch-channel` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-channel-moderation-commands.js:25` | command `mod-set-topic` | command: `/mod-set-topic` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-channel-moderation-commands.js:39` | command `mod-set-slowmode` | platform: `No direct equivalent` | No same-semantics public Slack bot API; do not substitute channel removal or organization deactivation. |
| `examples/discord-bots/moderation/lib/register-events.js:2` | event `messageCreate` | event: `message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-events.js:9` | event `messageUpdate` | event: `message_changed` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-events.js:24` | event `messageDelete` | event: `message_deleted` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-events.js:33` | event `reactionAdd` | event: `reaction_added` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-events.js:42` | event `reactionRemove` | event: `reaction_removed` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-events.js:51` | event `guildMemberAdd` | event: `team_join` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-events.js:60` | event `guildMemberUpdate` | event: `user_change` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-events.js:70` | event `guildMemberRemove` | event: `user_change` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-guild-role-lookup-commands.js:2` | command `mod-fetch-guild` | command: `/mod-fetch-workspace` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-guild-role-lookup-commands.js:26` | command `mod-list-roles` | command: `/mod-list-groups` | Native user-group membership workflow; Slack groups do not grant Discord role permissions. |
| `examples/discord-bots/moderation/lib/register-guild-role-lookup-commands.js:51` | command `mod-fetch-role` | command: `/mod-fetch-group` | Native user-group membership workflow; Slack groups do not grant Discord role permissions. |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:2` | command `mod-fetch-member` | command: `/mod-fetch-member` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:29` | command `mod-list-members` | command: `/mod-list-members` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:59` | command `mod-add-role` | command: `/mod-add-group-member` | Native user-group membership workflow; Slack groups do not grant Discord role permissions. |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:74` | command `mod-timeout` | platform: `No direct equivalent` | No same-semantics public Slack bot API; do not substitute channel removal or organization deactivation. |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:89` | command `mod-kick` | command: `/mod-remove-workspace-user` | Conditional Enterprise admin.users.remove; separate user token, actor allowlist and opt-in. |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:104` | command `mod-ban` | platform: `No direct equivalent` | No same-semantics public Slack bot API; do not substitute channel removal or organization deactivation. |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:123` | command `mod-unban` | platform: `No direct equivalent` | No same-semantics public Slack bot API; do not substitute channel removal or organization deactivation. |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:2` | command `mod-list-messages` | command: `/mod-list-messages` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:31` | command `mod-fetch-message` | command: `/mod-fetch-message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:53` | command `mod-pin` | command: `/mod-pin` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:67` | command `mod-unpin` | command: `/mod-unpin` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:81` | command `mod-list-pins` | command: `/mod-list-pins` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:99` | command `mod-bulk-delete` | command: `/mod-bulk-delete` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-overview-commands.js:2` | command `mod-summary` | command: `/mod-summary` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/moderation/lib/register-overview-commands.js:25` | command `mod-guidelines` | command: `/mod-guidelines` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:7` | command `ping` | command: `/golem-ping` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:56` | command `echo` | command: `/golem-echo` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:83` | command `feedback` | command: `/golem-feedback` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:120` | command `search` | command: `/golem-search` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:161` | command `announce` | command: `/golem-announce` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:183` | component `ping:panel` | action: `ping.panel` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:190` | component `ping:topic` | action: `ping.topic` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:198` | modal `feedback:submit` | view: `ping.feedback` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ping/index.js:207` | autocomplete `search` | options: `ping.search.query` | Native external select suggestions replace unsupported slash-option autocomplete. |
| `examples/discord-bots/ping/index.js:217` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/ping/index.js:224` | event `guildCreate` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/ping/index.js:231` | event `messageCreate` | event: `message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:266` | command `poker-help` | command: `/poker-help` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:270` | command `poker-deal` | command: `/poker-deal` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:274` | command `poker-draw` | command: `/poker-draw` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:301` | command `poker-score` | command: `/poker-score` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:305` | command `poker-reset` | command: `/poker-reset` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:309` | command `poker-rank` | command: `/poker-rank` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:326` | command `poker-action` | command: `/poker-action` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:360` | component `poker:help:deal` | action: `poker.deal` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:362` | component `poker:help:score` | action: `poker.score` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:364` | component `poker:help:rank` | action: `poker.rank` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:368` | component `poker:help:action` | action: `poker.action` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:372` | component `poker:help:reset` | action: `poker.reset` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:374` | modal `poker:rank:submit` | view: `poker.rank.submit` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:382` | modal `poker:action:submit` | view: `poker.action.submit` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/poker/index.js:390` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/poker/index.js:397` | event `messageCreate` | event: `message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:463` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/show-space/index.js:473` | command `upcoming` | command: `/upcoming` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:483` | command `debug` | command: `/debug` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:498` | command `debug-roles` | command: `/debug-groups` | Native user-group membership workflow; Slack groups do not grant Discord role permissions. |
| `examples/discord-bots/show-space/index.js:513` | command `debug-my-roles` | command: `/debug-my-permissions` | Native user-group membership workflow; Slack groups do not grant Discord role permissions. |
| `examples/discord-bots/show-space/index.js:528` | component `show-space:debug:summary` | action: `show.debug.summary` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:541` | component `show-space:debug:member` | action: `show.debug.member` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:554` | component `show-space:debug:guild` | action: `show.debug.workspace` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:567` | component `show-space:debug:config` | action: `show.debug.config` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:580` | component `show-space:debug:checks` | action: `show.debug.checks` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:593` | command `announce` | command: `/announce` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:624` | command `add-show` | command: `/add-show` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:662` | command `show` | command: `/show` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:678` | command `cancel-show` | command: `/cancel-show` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:719` | command `archive-show` | command: `/archive-show` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:752` | command `past-shows` | command: `/past-shows` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:762` | command `unpin-old` | command: `/unpin-old` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/show-space/index.js:797` | command `archive-expired` | command: `/archive-expired` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/support/index.js:10` | command `support-ticket` | command: `/support-ticket` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/support/index.js:41` | command `support-status` | command: `/support-status` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/support/index.js:50` | command `support-fetch-thread` | command: `/support-fetch-thread` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/support/index.js:77` | command `support-join-thread` | platform: `No thread membership API` | Native support-join-channel / support-leave-channel are explicitly different operations. |
| `examples/discord-bots/support/index.js:91` | command `support-leave-thread` | platform: `No thread membership API` | Native support-join-channel / support-leave-channel are explicitly different operations. |
| `examples/discord-bots/support/index.js:105` | command `support-start-thread` | command: `/support-start-thread` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/support/index.js:155` | event `guildCreate` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/support/index.js:159` | event `messageCreate` | event: `message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:49` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/ui-showcase/index.js:55` | event `messageCreate` | event: `message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:66` | command `demo-message` | command: `/demo-message` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:96` | component `showcase:msg:primary` | action: `demo.primary` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:97` | component `showcase:msg:success` | action: `demo.neutral` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:98` | component `showcase:msg:danger` | action: `demo.danger` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:100` | component `showcase:msg:select` | action: `demo.message.select` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:115` | command `demo-form` | command: `/demo-form` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:128` | modal `showcase:form:submit` | view: `demo.form` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:152` | command `demo-search` | command: `/demo-search` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:167` | command `find` | command: `/find` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:330` | autocomplete `demo-search` | options: `demo.search.suggest` | Native external select suggestions replace unsupported slash-option autocomplete. |
| `examples/discord-bots/ui-showcase/index.js:334` | autocomplete `find` | options: `demo.search.suggest` | Native external select suggestions replace unsupported slash-option autocomplete. |
| `examples/discord-bots/ui-showcase/index.js:342` | command `demo-review` | command: `/demo-review` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:510` | command `demo-confirm` | command: `/demo-confirm` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:530` | component `showcase:confirm:yes` | action: `demo.confirm.yes` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:542` | component `showcase:confirm:no` | action: `demo.confirm.no` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:550` | command `demo-pager` | command: `/demo-pager` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:602` | command `demo-cards` | command: `/demo-cards` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:611` | command `browse` | command: `/browse` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:675` | component `showcase:buy:confirm` | action: `demo.card.buy` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:685` | component `showcase:buy:cancel` | platform: `Native confirmation cancellation` | Cancelling Slack confirmation dismisses it without dispatching the destructive action. |
| `examples/discord-bots/ui-showcase/index.js:728` | command `demo-selects` | command: `/demo-selects` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:758` | component `showcase:select:string` | action: `demo.select.string` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:763` | component `showcase:select:user` | action: `demo.select.user` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:768` | component `showcase:select:role` | action: `demo.select.group` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:773` | component `showcase:select:channel` | action: `demo.select.channel` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:778` | component `showcase:select:mentionable` | action: `demo.select.users` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:787` | command `demo-alias` | command: `/demo-alias` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/ui-showcase/index.js:803` | command `demo-alias-alt` | command: `/demo-alias-alt` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/unified-demo/index.js:10` | command `unified-ping` | command: `/unified-ping` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/unified-demo/index.js:21` | event `ready` | host: `runtime load / Socket Mode authentication logs` | Native lifecycle equivalent; no Discord gateway event emulation. |
| `examples/discord-bots/unified-demo/index.js:37` | __verb__ `status` | verb: `status` | Ported; registration covered, workflow tests listed below. |
| `examples/discord-bots/unified-demo/index.js:46` | __verb__ `run` | host: `bots run unified-demo` | Host-managed runtime configuration uses declared fields and the shared CLI. |

## Computed registrations

- **custom-kb**: action `kb.select`, `kb.refresh`; view `kb.add`, `kb.search`; options `kb.suggest`.
- **knowledge-base**: action `knowledge.select`, `knowledge.verify`, `knowledge.stale`, `knowledge.reject`, `knowledge.source`, `knowledge.edit`, `knowledge.export`, `knowledge.previous`, `knowledge.next`; view `knowledge.submit`, `knowledge.open`; options `knowledge.suggest`.
- **ui-showcase**: action `demo.article.select`, `demo.article.verify`, `demo.article.stale`, `demo.article.reject`, `demo.article.edit`, `demo.article.source`, `demo.article.export`, `demo.article.details`, `demo.search.previous`, `demo.search.next`, `demo.pager.previous`, `demo.pager.next`, `demo.card.select`, `demo.card.buy`, `demo.card.info`, `demo.card.share`; view `demo.article.save`, `demo.search.open`; options `demo.search.suggest`.


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
