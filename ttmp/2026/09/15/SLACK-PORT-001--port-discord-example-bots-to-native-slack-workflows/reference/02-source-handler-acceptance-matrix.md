---
Title: Source handler acceptance matrix
Ticket: SLACK-PORT-001
Status: active
Topics:
- slack
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Every indexed Discord registration with a pending parity acceptance status.
LastUpdated: 2026-09-15T20:20:00-04:00
WhatFor: Prevent silent feature omissions while porting the examples.
WhenToUse: At each bot acceptance checkpoint.
---

# Source handler acceptance matrix

This checklist indexes source registrations. Every row starts pending even when a basic Slack counterpart exists; completion requires feature-level review and an executable acceptance test. Computed handlers and store/export helpers must also be checked.

| Source | Kind | Identifier | Status |
| --- | --- | --- | --- |
| `examples/discord-bots/announcements.js:10` | command | `announce-preview` | Pending |
| `examples/discord-bots/announcements.js:26` | event | `ready` | Pending |
| `examples/discord-bots/archive-helper/index.js:21` | event | `ready` | Pending |
| `examples/discord-bots/archive-helper/index.js:27` | command | `archive-channel` | Pending |
| `examples/discord-bots/archive-helper/index.js:101` | messageCommand | `Archive Thread` | Pending |
| `examples/discord-bots/custom-kb/index.js:23` | event | `ready` | Pending |
| `examples/discord-bots/custom-kb/index.js:28` | command | `kb-add` | Pending |
| `examples/discord-bots/custom-kb/index.js:48` | command | `kb-link` | Pending |
| `examples/discord-bots/custom-kb/index.js:62` | command | `kb-search` | Pending |
| `examples/discord-bots/custom-kb/index.js:67` | command | `kb-list` | Pending |
| `examples/discord-bots/hater/index.js:38` | command | `hate` | Pending |
| `examples/discord-bots/hater/index.js:60` | command | `roast` | Pending |
| `examples/discord-bots/hater/index.js:76` | command | `compliment` | Pending |
| `examples/discord-bots/hater/index.js:90` | component | `hater:roast` | Pending |
| `examples/discord-bots/hater/index.js:97` | component | `hater:mercy` | Pending |
| `examples/discord-bots/hater/index.js:105` | component | `hater:apology` | Pending |
| `examples/discord-bots/hater/index.js:114` | modal | `hater:apology:submit` | Pending |
| `examples/discord-bots/hater/index.js:131` | event | `ready` | Pending |
| `examples/discord-bots/hater/index.js:138` | event | `messageCreate` | Pending |
| `examples/discord-bots/interaction-types/index.js:11` | command | `hello` | Pending |
| `examples/discord-bots/interaction-types/index.js:18` | command | `echo` | Pending |
| `examples/discord-bots/interaction-types/index.js:28` | command | `fun` | Pending |
| `examples/discord-bots/interaction-types/index.js:48` | subcommand | `fun` | Pending |
| `examples/discord-bots/interaction-types/index.js:59` | subcommand | `fun` | Pending |
| `examples/discord-bots/interaction-types/index.js:70` | userCommand | `Show Avatar` | Pending |
| `examples/discord-bots/interaction-types/index.js:85` | messageCommand | `Quote Message` | Pending |
| `examples/discord-bots/interaction-types/index.js:98` | event | `ready` | Pending |
| `examples/discord-bots/knowledge-base/index.js:18` | event | `ready` | Pending |
| `examples/discord-bots/knowledge-base/index.js:28` | event | `messageCreate` | Pending |
| `examples/discord-bots/knowledge-base/index.js:40` | command | `remember` | Pending |
| `examples/discord-bots/knowledge-base/index.js:46` | command | `teach` | Pending |
| `examples/discord-bots/knowledge-base/index.js:52` | modal | `knowledge:submit` | Pending |
| `examples/discord-bots/knowledge-base/index.js:64` | command | `ask` | Pending |
| `examples/discord-bots/knowledge-base/index.js:83` | command | `kb-search` | Pending |
| `examples/discord-bots/knowledge-base/index.js:102` | command | `article` | Pending |
| `examples/discord-bots/knowledge-base/index.js:124` | command | `kb-article` | Pending |
| `examples/discord-bots/knowledge-base/index.js:146` | command | `review` | Pending |
| `examples/discord-bots/knowledge-base/index.js:169` | command | `kb-review` | Pending |
| `examples/discord-bots/knowledge-base/index.js:192` | command | `recent` | Pending |
| `examples/discord-bots/knowledge-base/index.js:206` | command | `kb-recent` | Pending |
| `examples/discord-bots/knowledge-base/index.js:220` | command | `kb-verify` | Pending |
| `examples/discord-bots/knowledge-base/index.js:243` | command | `kb-stale` | Pending |
| `examples/discord-bots/knowledge-base/index.js:266` | command | `kb-reject` | Pending |
| `examples/discord-bots/knowledge-base/lib/reactions.js:2` | event | `reactionAdd` | Pending |
| `examples/discord-bots/moderation/lib/register-channel-moderation-commands.js:2` | command | `mod-fetch-channel` | Pending |
| `examples/discord-bots/moderation/lib/register-channel-moderation-commands.js:25` | command | `mod-set-topic` | Pending |
| `examples/discord-bots/moderation/lib/register-channel-moderation-commands.js:39` | command | `mod-set-slowmode` | Pending |
| `examples/discord-bots/moderation/lib/register-events.js:2` | event | `messageCreate` | Pending |
| `examples/discord-bots/moderation/lib/register-events.js:9` | event | `messageUpdate` | Pending |
| `examples/discord-bots/moderation/lib/register-events.js:24` | event | `messageDelete` | Pending |
| `examples/discord-bots/moderation/lib/register-events.js:33` | event | `reactionAdd` | Pending |
| `examples/discord-bots/moderation/lib/register-events.js:42` | event | `reactionRemove` | Pending |
| `examples/discord-bots/moderation/lib/register-events.js:51` | event | `guildMemberAdd` | Pending |
| `examples/discord-bots/moderation/lib/register-events.js:60` | event | `guildMemberUpdate` | Pending |
| `examples/discord-bots/moderation/lib/register-events.js:70` | event | `guildMemberRemove` | Pending |
| `examples/discord-bots/moderation/lib/register-guild-role-lookup-commands.js:2` | command | `mod-fetch-guild` | Pending |
| `examples/discord-bots/moderation/lib/register-guild-role-lookup-commands.js:26` | command | `mod-list-roles` | Pending |
| `examples/discord-bots/moderation/lib/register-guild-role-lookup-commands.js:51` | command | `mod-fetch-role` | Pending |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:2` | command | `mod-fetch-member` | Pending |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:29` | command | `mod-list-members` | Pending |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:59` | command | `mod-add-role` | Pending |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:74` | command | `mod-timeout` | Pending |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:89` | command | `mod-kick` | Pending |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:104` | command | `mod-ban` | Pending |
| `examples/discord-bots/moderation/lib/register-member-moderation-commands.js:123` | command | `mod-unban` | Pending |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:2` | command | `mod-list-messages` | Pending |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:31` | command | `mod-fetch-message` | Pending |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:53` | command | `mod-pin` | Pending |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:67` | command | `mod-unpin` | Pending |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:81` | command | `mod-list-pins` | Pending |
| `examples/discord-bots/moderation/lib/register-message-moderation-commands.js:99` | command | `mod-bulk-delete` | Pending |
| `examples/discord-bots/moderation/lib/register-overview-commands.js:2` | command | `mod-summary` | Pending |
| `examples/discord-bots/moderation/lib/register-overview-commands.js:25` | command | `mod-guidelines` | Pending |
| `examples/discord-bots/ping/index.js:7` | command | `ping` | Pending |
| `examples/discord-bots/ping/index.js:56` | command | `echo` | Pending |
| `examples/discord-bots/ping/index.js:83` | command | `feedback` | Pending |
| `examples/discord-bots/ping/index.js:120` | command | `search` | Pending |
| `examples/discord-bots/ping/index.js:161` | command | `announce` | Pending |
| `examples/discord-bots/ping/index.js:183` | component | `ping:panel` | Pending |
| `examples/discord-bots/ping/index.js:190` | component | `ping:topic` | Pending |
| `examples/discord-bots/ping/index.js:198` | modal | `feedback:submit` | Pending |
| `examples/discord-bots/ping/index.js:217` | event | `ready` | Pending |
| `examples/discord-bots/ping/index.js:224` | event | `guildCreate` | Pending |
| `examples/discord-bots/ping/index.js:231` | event | `messageCreate` | Pending |
| `examples/discord-bots/poker/index.js:266` | command | `poker-help` | Pending |
| `examples/discord-bots/poker/index.js:270` | command | `poker-deal` | Pending |
| `examples/discord-bots/poker/index.js:274` | command | `poker-draw` | Pending |
| `examples/discord-bots/poker/index.js:301` | command | `poker-score` | Pending |
| `examples/discord-bots/poker/index.js:305` | command | `poker-reset` | Pending |
| `examples/discord-bots/poker/index.js:309` | command | `poker-rank` | Pending |
| `examples/discord-bots/poker/index.js:326` | command | `poker-action` | Pending |
| `examples/discord-bots/poker/index.js:360` | component | `poker:help:deal` | Pending |
| `examples/discord-bots/poker/index.js:362` | component | `poker:help:score` | Pending |
| `examples/discord-bots/poker/index.js:364` | component | `poker:help:rank` | Pending |
| `examples/discord-bots/poker/index.js:368` | component | `poker:help:action` | Pending |
| `examples/discord-bots/poker/index.js:372` | component | `poker:help:reset` | Pending |
| `examples/discord-bots/poker/index.js:374` | modal | `poker:rank:submit` | Pending |
| `examples/discord-bots/poker/index.js:382` | modal | `poker:action:submit` | Pending |
| `examples/discord-bots/poker/index.js:390` | event | `ready` | Pending |
| `examples/discord-bots/poker/index.js:397` | event | `messageCreate` | Pending |
| `examples/discord-bots/show-space/index.js:463` | event | `ready` | Pending |
| `examples/discord-bots/show-space/index.js:473` | command | `upcoming` | Pending |
| `examples/discord-bots/show-space/index.js:483` | command | `debug` | Pending |
| `examples/discord-bots/show-space/index.js:498` | command | `debug-roles` | Pending |
| `examples/discord-bots/show-space/index.js:513` | command | `debug-my-roles` | Pending |
| `examples/discord-bots/show-space/index.js:528` | component | `show-space:debug:summary` | Pending |
| `examples/discord-bots/show-space/index.js:541` | component | `show-space:debug:member` | Pending |
| `examples/discord-bots/show-space/index.js:554` | component | `show-space:debug:guild` | Pending |
| `examples/discord-bots/show-space/index.js:567` | component | `show-space:debug:config` | Pending |
| `examples/discord-bots/show-space/index.js:580` | component | `show-space:debug:checks` | Pending |
| `examples/discord-bots/show-space/index.js:593` | command | `announce` | Pending |
| `examples/discord-bots/show-space/index.js:624` | command | `add-show` | Pending |
| `examples/discord-bots/show-space/index.js:662` | command | `show` | Pending |
| `examples/discord-bots/show-space/index.js:678` | command | `cancel-show` | Pending |
| `examples/discord-bots/show-space/index.js:719` | command | `archive-show` | Pending |
| `examples/discord-bots/show-space/index.js:752` | command | `past-shows` | Pending |
| `examples/discord-bots/show-space/index.js:762` | command | `unpin-old` | Pending |
| `examples/discord-bots/show-space/index.js:797` | command | `archive-expired` | Pending |
| `examples/discord-bots/support/index.js:10` | command | `support-ticket` | Pending |
| `examples/discord-bots/support/index.js:41` | command | `support-status` | Pending |
| `examples/discord-bots/support/index.js:50` | command | `support-fetch-thread` | Pending |
| `examples/discord-bots/support/index.js:77` | command | `support-join-thread` | Pending |
| `examples/discord-bots/support/index.js:91` | command | `support-leave-thread` | Pending |
| `examples/discord-bots/support/index.js:105` | command | `support-start-thread` | Pending |
| `examples/discord-bots/support/index.js:155` | event | `guildCreate` | Pending |
| `examples/discord-bots/support/index.js:159` | event | `messageCreate` | Pending |
| `examples/discord-bots/ui-showcase/index.js:49` | event | `ready` | Pending |
| `examples/discord-bots/ui-showcase/index.js:55` | event | `messageCreate` | Pending |
| `examples/discord-bots/ui-showcase/index.js:66` | command | `demo-message` | Pending |
| `examples/discord-bots/ui-showcase/index.js:96` | component | `showcase:msg:primary` | Pending |
| `examples/discord-bots/ui-showcase/index.js:97` | component | `showcase:msg:success` | Pending |
| `examples/discord-bots/ui-showcase/index.js:98` | component | `showcase:msg:danger` | Pending |
| `examples/discord-bots/ui-showcase/index.js:100` | component | `showcase:msg:select` | Pending |
| `examples/discord-bots/ui-showcase/index.js:115` | command | `demo-form` | Pending |
| `examples/discord-bots/ui-showcase/index.js:128` | modal | `showcase:form:submit` | Pending |
| `examples/discord-bots/ui-showcase/index.js:152` | command | `demo-search` | Pending |
| `examples/discord-bots/ui-showcase/index.js:167` | command | `find` | Pending |
| `examples/discord-bots/ui-showcase/index.js:342` | command | `demo-review` | Pending |
| `examples/discord-bots/ui-showcase/index.js:510` | command | `demo-confirm` | Pending |
| `examples/discord-bots/ui-showcase/index.js:530` | component | `showcase:confirm:yes` | Pending |
| `examples/discord-bots/ui-showcase/index.js:542` | component | `showcase:confirm:no` | Pending |
| `examples/discord-bots/ui-showcase/index.js:550` | command | `demo-pager` | Pending |
| `examples/discord-bots/ui-showcase/index.js:602` | command | `demo-cards` | Pending |
| `examples/discord-bots/ui-showcase/index.js:611` | command | `browse` | Pending |
| `examples/discord-bots/ui-showcase/index.js:675` | component | `showcase:buy:confirm` | Pending |
| `examples/discord-bots/ui-showcase/index.js:685` | component | `showcase:buy:cancel` | Pending |
| `examples/discord-bots/ui-showcase/index.js:728` | command | `demo-selects` | Pending |
| `examples/discord-bots/ui-showcase/index.js:758` | component | `showcase:select:string` | Pending |
| `examples/discord-bots/ui-showcase/index.js:763` | component | `showcase:select:user` | Pending |
| `examples/discord-bots/ui-showcase/index.js:768` | component | `showcase:select:role` | Pending |
| `examples/discord-bots/ui-showcase/index.js:773` | component | `showcase:select:channel` | Pending |
| `examples/discord-bots/ui-showcase/index.js:778` | component | `showcase:select:mentionable` | Pending |
| `examples/discord-bots/ui-showcase/index.js:787` | command | `demo-alias` | Pending |
| `examples/discord-bots/ui-showcase/index.js:803` | command | `demo-alias-alt` | Pending |
| `examples/discord-bots/unified-demo/index.js:10` | command | `unified-ping` | Pending |
| `examples/discord-bots/unified-demo/index.js:21` | event | `ready` | Pending |
