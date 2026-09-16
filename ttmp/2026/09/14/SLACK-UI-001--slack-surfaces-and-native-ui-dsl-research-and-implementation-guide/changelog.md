# Changelog

## 2026-09-14

- Initial workspace created

## 2026-09-14

Researched Slack surfaces and Discord UI construction, archived official and local sources, measured pinned SDK unknown-block data loss, and authored a pragmatic intern guide. Core implementation remains open.

## 2026-09-14

Validated references and docmgr metadata; archived three local background documents; rendered and delivered the corrected 20-page Slack UI Research Guide to reMarkable. Receipt and PDF stored in artifacts; implementation tasks remain open.

## 2026-09-14

Clarified implementation scope: layers 1–7 are in scope, including raw Block Kit values and core ACK handling; layer 8 scheduling is explicitly deferred. Removed worker reservation and bounded interactive execution from the design, refreshed the print copy, and uploaded v2 to reMarkable.

## 2026-09-15

Implemented Phase 1 rich messages and native Slack UI builders; added detached Block Kit payloads, lossless transport encoding, TypeScript declarations, and a showcase bot.

### Related Files

- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/examples/slack-bots/ui-showcase/index.js — Offline showcase acceptance example
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/jsslack/ui_module.go — Native slack/ui builder module
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/slacktransport/client.go — Lossless Block Kit transport
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/pkg/slackbot/model.go — Rich message and block contracts

## 2026-09-15

Implemented Phase 2 block-action routing for button and static-select interactions; Socket Mode now normalizes action payloads, acknowledges interactive envelopes, and exposes ctx.action/action registrations.

### Related Files

- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/pkg/slackbot/model.go — Normalized action contract
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/jsslack/module.go — JavaScript action registration
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/slacktransport/run.go — Interactive block_actions decoding and acknowledgment capability

## 2026-09-15

Implemented Phase 3 modal flow and offline simulation support: modal/text-input builders, views.open, view callback routing, values extraction, single-use deadline-aware ACKs, and recorded open_view/ack operations.

### Related Files

- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/jsslack/dispatch.go — Modal open, view state extraction, and ACK methods
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/slacktransport/run.go — Interactive receipt state and view submission decoding
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/pkg/slackbot/recording.go — Offline ACK and modal operation recording
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/pkg/slackcli/commands.go — Interactive fixture simulation wiring

## 2026-09-15

Added process-local replay of successful interactive ACK payloads for duplicate Socket Mode envelope IDs, preventing duplicate view handlers while preserving the five-minute in-memory scope.

### Related Files

- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/slacktransport/client.go — ACK replay map
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/slacktransport/run.go — Duplicate envelope replay
