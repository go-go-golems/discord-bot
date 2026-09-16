---
Title: Slack UI research source catalog
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Source provenance and extraction limitations."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Locate archived evidence."
WhenToUse: "Review the research guide."
---

# Sources

Official Slack documentation captured 2026-09-14 local time / 2026-09-15 UTC. `catalog-complete.json` records URLs, capture status, and final Markdown hashes. Initial network failures were recovered with authorized network access. Two Defuddle URL timeouts were recovered by capturing browser-rendered HTML and running Defuddle on that HTML. Catalog pages require browser rendering to populate their component tables; use rendered captures rather than the short static versions.

Markdown wrappers add docmgr metadata and resolve root-relative Slack links. `.original.txt` files preserve Defuddle output. Rendered JSON retains HTML and plain text. Raw HTML catalog captures show why static extraction alone was inadequate. Official docs contain example token strings; no private workspace credentials were read or archived.

`code/REVISION.txt` pins the local code snapshots and SDK. The vault article is explanatory background, not authority to change dependencies. No upstream clone was necessary because the pinned module cache and existing repository contain the relevant source; snapshots make the evidence self-contained.

## Pages

| Page | Original | Archive | Status |
|---|---|---|---|
| 01-surfaces | [Slack](https://docs.slack.dev/surfaces/) | [capture](web/01-surfaces.md) | Captured |
| 02-surfaces-app-home | [Slack](https://docs.slack.dev/surfaces/app-home/) | [capture](web/02-surfaces-app-home.md) | Captured |
| 03-surfaces-modals | [Slack](https://docs.slack.dev/surfaces/modals/) | [capture](web/03-surfaces-modals.md) | Captured |
| 04-surfaces-canvases | [Slack](https://docs.slack.dev/surfaces/canvases/) | [capture](web/04-surfaces-canvases.md) | Captured |
| 05-surfaces-lists | [Slack](https://docs.slack.dev/surfaces/lists/) | [capture](web/05-surfaces-lists.md) | Captured |
| 06-surfaces-split-view | [Slack](https://docs.slack.dev/surfaces/split-view/) | [capture](web/06-surfaces-split-view.md) | Captured |
| 07-block-kit | [Slack](https://docs.slack.dev/block-kit/) | [capture](web/07-block-kit.md) | Captured |
| 08-reference-block-kit-blocks | [Slack](https://docs.slack.dev/reference/block-kit/blocks/) | [capture](web/08-reference-block-kit-blocks.md) | Captured |
| 09-reference-block-kit-block-elements | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/) | [capture](web/09-reference-block-kit-block-elements.md) | Captured |
| 10-reference-block-kit-composition-objects | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/) | [capture](web/10-reference-block-kit-composition-objects.md) | Captured |
| 11-interactivity-handling-user-interaction | [Slack](https://docs.slack.dev/interactivity/handling-user-interaction/) | [capture](web/11-interactivity-handling-user-interaction.md) | Captured |
| 12-apis-events-api-using-socket-mode | [Slack](https://docs.slack.dev/apis/events-api/using-socket-mode/) | [capture](web/12-apis-events-api-using-socket-mode.md) | Captured |
| 13-reference-interaction-payloads-block_actions-payload | [Slack](https://docs.slack.dev/reference/interaction-payloads/block_actions-payload/) | [capture](web/13-reference-interaction-payloads-block_actions-payload.md) | Captured |
| 14-reference-interaction-payloads-block_suggestion-payload | [Slack](https://docs.slack.dev/reference/interaction-payloads/block_suggestion-payload/) | [capture](web/14-reference-interaction-payloads-block_suggestion-payload.md) | Captured |
| 15-reference-interaction-payloads-view-interactions-payload | [Slack](https://docs.slack.dev/reference/interaction-payloads/view-interactions-payload/) | [capture](web/15-reference-interaction-payloads-view-interactions-payload.md) | Captured |
| 16-reference-views-modal-views | [Slack](https://docs.slack.dev/reference/views/modal-views/) | [capture](web/16-reference-views-modal-views.md) | Captured |
| 17-reference-methods-chat.postMessage | [Slack](https://docs.slack.dev/reference/methods/chat.postMessage/) | [capture](web/17-reference-methods-chat.postMessage.md) | Captured |
| 18-reference-methods-chat.update | [Slack](https://docs.slack.dev/reference/methods/chat.update/) | [capture](web/18-reference-methods-chat.update.md) | Captured |
| 19-reference-methods-chat.postEphemeral | [Slack](https://docs.slack.dev/reference/methods/chat.postEphemeral/) | [capture](web/19-reference-methods-chat.postEphemeral.md) | Captured |
| 20-reference-methods-views.open | [Slack](https://docs.slack.dev/reference/methods/views.open/) | [capture](web/20-reference-methods-views.open.md) | Captured |
| 21-reference-methods-views.update | [Slack](https://docs.slack.dev/reference/methods/views.update/) | [capture](web/21-reference-methods-views.update.md) | Captured |
| 22-reference-methods-views.publish | [Slack](https://docs.slack.dev/reference/methods/views.publish/) | [capture](web/22-reference-methods-views.publish.md) | Captured |
| 23-reference-block-kit-blocks-input-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/input-block/) | [capture](web/23-reference-block-kit-blocks-input-block.md) | Captured |
| 24-reference-block-kit-blocks-section-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/section-block/) | [capture](web/24-reference-block-kit-blocks-section-block.md) | Captured |
| 25-reference-block-kit-blocks-actions-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/actions-block/) | [capture](web/25-reference-block-kit-blocks-actions-block.md) | Captured |
| 26-reference-block-kit-block-elements-button-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/button-element/) | [capture](web/26-reference-block-kit-block-elements-button-element.md) | Captured |
| 27-reference-block-kit-block-elements-select-menu-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/select-menu-element/) | [capture](web/27-reference-block-kit-block-elements-select-menu-element.md) | Captured |
| 28-reference-block-kit-blocks-table-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/table-block/) | [capture](web/28-reference-block-kit-blocks-table-block.md) | Captured |
| 29-reference-block-kit-blocks-markdown-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/markdown-block/) | [capture](web/29-reference-block-kit-blocks-markdown-block.md) | Captured |
| 30-reference-block-kit-blocks-task-card-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/task-card-block/) | [capture](web/30-reference-block-kit-blocks-task-card-block.md) | Captured |
| 31-reference-block-kit-blocks-plan-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/plan-block/) | [capture](web/31-reference-block-kit-blocks-plan-block.md) | Captured |
| 32-reference-block-kit-blocks-rich-text-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block/) | [capture](web/32-reference-block-kit-blocks-rich-text-block.md) | Captured |
| 33-reference-block-kit-blocks-data-table-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/data-table-block/) | [capture](web/33-reference-block-kit-blocks-data-table-block.md) | Captured |
| 34-reference-block-kit-blocks-data-visualization-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/data-visualization-block/) | [capture](web/34-reference-block-kit-blocks-data-visualization-block.md) | Captured |
| 35-reference-block-kit-blocks-card-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/card-block/) | [capture](web/35-reference-block-kit-blocks-card-block.md) | Captured |
| 36-reference-block-kit-blocks-carousel-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/carousel-block/) | [capture](web/36-reference-block-kit-blocks-carousel-block.md) | Captured |
| 37-reference-block-kit-blocks-container-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/container-block/) | [capture](web/37-reference-block-kit-blocks-container-block.md) | Captured |
| 38-reference-block-kit-blocks-alert-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/alert-block/) | [capture](web/38-reference-block-kit-blocks-alert-block.md) | Captured |
| 39-reference-block-kit-blocks-context-actions-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/context-actions-block/) | [capture](web/39-reference-block-kit-blocks-context-actions-block.md) | Captured |
| 40-reference-methods-files.getUploadURLExternal | [Slack](https://docs.slack.dev/reference/methods/files.getUploadURLExternal/) | [capture](web/40-reference-methods-files.getUploadURLExternal.md) | Captured |
| 41-reference-methods-slackLists.create | [Slack](https://docs.slack.dev/reference/methods/slackLists.create) | [capture](web/rendered-lists-create.md) | Recovered via rendered HTML |
| 42-reference-methods-canvases.create | [Slack](https://docs.slack.dev/reference/methods/canvases.create/) | [capture](web/42-reference-methods-canvases.create.md) | Captured |
| 43-reference-events-app_home_opened | [Slack](https://docs.slack.dev/reference/events/app_home_opened/) | [capture](web/43-reference-events-app_home_opened.md) | Captured |
| 44-reference-app-manifest | [Slack](https://docs.slack.dev/reference/app-manifest/) | [capture](web/44-reference-app-manifest.md) | Captured |
| 45-messaging-work-objects | [Slack](https://docs.slack.dev/messaging/work-objects/) | [capture](web/45-messaging-work-objects.md) | Captured |
| 46-ai-developing-agents | [Slack](https://docs.slack.dev/ai/developing-agents/) | [capture](web/46-ai-developing-agents.md) | Captured |
| 47-ai-migrating-to-agent-messaging | [Slack](https://docs.slack.dev/ai/migrating-to-agent-messaging/) | [capture](web/47-ai-migrating-to-agent-messaging.md) | Captured |
| 48-messaging-work-objects-overview | [Slack](https://docs.slack.dev/messaging/work-objects-overview/) | [capture](web/48-messaging-work-objects-overview.md) | Captured |
| 49-messaging-work-objects-implementation | [Slack](https://docs.slack.dev/messaging/work-objects-implementation/) | [capture](web/rendered-work-objects-implementation.md) | Recovered via rendered HTML |
| 50-reference-methods-files.completeUploadExternal | [Slack](https://docs.slack.dev/reference/methods/files.completeUploadExternal/) | [capture](web/50-reference-methods-files.completeUploadExternal.md) | Captured |
| 51-reference-methods-chat.unfurl | [Slack](https://docs.slack.dev/reference/methods/chat.unfurl/) | [capture](web/51-reference-methods-chat.unfurl.md) | Captured |
| 52-reference-methods-views.push | [Slack](https://docs.slack.dev/reference/methods/views.push/) | [capture](web/52-reference-methods-views.push.md) | Captured |
| 53-reference-block-kit-blocks-video-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/video-block/) | [capture](web/53-reference-block-kit-blocks-video-block.md) | Captured |
| 54-reference-block-kit-block-elements-plain-text-input-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/plain-text-input-element/) | [capture](web/54-reference-block-kit-block-elements-plain-text-input-element.md) | Captured |
| 55-reference-block-kit-composition-objects-text-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/text-object/) | [capture](web/55-reference-block-kit-composition-objects-text-object.md) | Captured |
| 56-reference-block-kit-composition-objects-confirmation-dialog-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object/) | [capture](web/56-reference-block-kit-composition-objects-confirmation-dialog-object.md) | Captured |
| 57-interactivity-implementing-shortcuts | [Slack](https://docs.slack.dev/interactivity/implementing-shortcuts/) | [capture](web/57-interactivity-implementing-shortcuts.md) | Captured |
| 58-reference-methods-chat.startStream | [Slack](https://docs.slack.dev/reference/methods/chat.startStream/) | [capture](web/58-reference-methods-chat.startStream.md) | Captured |
| 59-reference-methods-assistant.threads.setStatus | [Slack](https://docs.slack.dev/reference/methods/assistant.threads.setStatus/) | [capture](web/59-reference-methods-assistant.threads.setStatus.md) | Captured |
| extra-agent-sessions | [Slack](https://docs.slack.dev/ai/agent-sessions/) | [capture](web/extra-agent-sessions.md) | Captured |
| extra-attachment-mention-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/attachment-mention-element) | [capture](web/extra-attachment-mention-element.md) | Captured |
| extra-broadcast-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/broadcast-element) | [capture](web/extra-broadcast-element.md) | Captured |
| extra-canvas-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/canvas-element) | [capture](web/extra-canvas-element.md) | Captured |
| extra-canvas-message-unfurl-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/canvas-message-unfurl-element) | [capture](web/extra-canvas-message-unfurl-element.md) | Captured |
| extra-canvas-user-mention-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/canvas-user-mention-element) | [capture](web/extra-canvas-user-mention-element.md) | Captured |
| extra-channel-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/channel-element) | [capture](web/extra-channel-element.md) | Captured |
| extra-checkboxes-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/checkboxes-element) | [capture](web/extra-checkboxes-element.md) | Captured |
| extra-citation-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/citation-element) | [capture](web/extra-citation-element.md) | Captured |
| extra-color-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/color-element) | [capture](web/extra-color-element.md) | Captured |
| extra-date-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/date-element) | [capture](web/extra-date-element.md) | Captured |
| extra-date-picker-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/date-picker-element) | [capture](web/extra-date-picker-element.md) | Captured |
| extra-datetime-picker-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/datetime-picker-element) | [capture](web/extra-datetime-picker-element.md) | Captured |
| extra-email-input-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/email-input-element) | [capture](web/extra-email-input-element.md) | Captured |
| extra-emoji-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/emoji-element) | [capture](web/extra-emoji-element.md) | Captured |
| extra-feedback-buttons-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/feedback-buttons-element) | [capture](web/extra-feedback-buttons-element.md) | Captured |
| extra-file-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/file-element) | [capture](web/extra-file-element.md) | Captured |
| extra-file-input-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/file-input-element) | [capture](web/extra-file-input-element.md) | Captured |
| extra-icon-button-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/icon-button-element) | [capture](web/extra-icon-button-element.md) | Captured |
| extra-image-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/image-element) | [capture](web/extra-image-element.md) | Captured |
| extra-link-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/link-element) | [capture](web/extra-link-element.md) | Captured |
| extra-list-record-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/list-record-element) | [capture](web/extra-list-record-element.md) | Captured |
| extra-message-mention-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/message-mention-element) | [capture](web/extra-message-mention-element.md) | Captured |
| extra-multi-select-menu-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/multi-select-menu-element) | [capture](web/extra-multi-select-menu-element.md) | Captured |
| extra-number-input-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/number-input-element) | [capture](web/extra-number-input-element.md) | Captured |
| extra-overflow-menu-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/overflow-menu-element) | [capture](web/extra-overflow-menu-element.md) | Captured |
| extra-radio-button-group-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/radio-button-group-element) | [capture](web/extra-radio-button-group-element.md) | Captured |
| extra-rich-text-input-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-input-element) | [capture](web/extra-rich-text-input-element.md) | Captured |
| extra-rich-text-list-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element) | [capture](web/extra-rich-text-list-element.md) | Captured |
| extra-rich-text-preformatted-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-preformatted-element) | [capture](web/extra-rich-text-preformatted-element.md) | Captured |
| extra-rich-text-quote-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-quote-element) | [capture](web/extra-rich-text-quote-element.md) | Captured |
| extra-rich-text-section-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element) | [capture](web/extra-rich-text-section-element.md) | Captured |
| extra-salesforce-data-field-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/salesforce-data-field-element) | [capture](web/extra-salesforce-data-field-element.md) | Captured |
| extra-tag-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/tag-element) | [capture](web/extra-tag-element.md) | Captured |
| extra-team-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/team-element) | [capture](web/extra-team-element.md) | Captured |
| extra-text-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/text-element) | [capture](web/extra-text-element.md) | Captured |
| extra-time-picker-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/time-picker-element) | [capture](web/extra-time-picker-element.md) | Captured |
| extra-url-input-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/url-input-element) | [capture](web/extra-url-input-element.md) | Captured |
| extra-url-source-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/url-source-element) | [capture](web/extra-url-source-element.md) | Captured |
| extra-user-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/user-element) | [capture](web/extra-user-element.md) | Captured |
| extra-usergroup-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/usergroup-element) | [capture](web/extra-usergroup-element.md) | Captured |
| extra-work-object-mention-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/work-object-mention-element) | [capture](web/extra-work-object-mention-element.md) | Captured |
| extra-workflow-button-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/workflow-button-element) | [capture](web/extra-workflow-button-element.md) | Captured |
| extra-workflow-mention-element | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/workflow-mention-element) | [capture](web/extra-workflow-mention-element.md) | Captured |
| extra-context-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/context-block) | [capture](web/extra-context-block.md) | Captured |
| extra-divider-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/divider-block) | [capture](web/extra-divider-block.md) | Captured |
| extra-file-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/file-block) | [capture](web/extra-file-block.md) | Captured |
| extra-header-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/header-block) | [capture](web/extra-header-block.md) | Captured |
| extra-image-block | [Slack](https://docs.slack.dev/reference/block-kit/blocks/image-block) | [capture](web/extra-image-block.md) | Captured |
| extra-conversation-filter-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/conversation-filter-object) | [capture](web/extra-conversation-filter-object.md) | Captured |
| extra-dispatch-action-configuration-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/dispatch-action-configuration-object) | [capture](web/extra-dispatch-action-configuration-object.md) | Captured |
| extra-option-group-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/option-group-object) | [capture](web/extra-option-group-object.md) | Captured |
| extra-option-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/option-object) | [capture](web/extra-option-object.md) | Captured |
| extra-slack-file-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/slack-file-object) | [capture](web/extra-slack-file-object.md) | Captured |
| extra-slack-icon-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/slack-icon-object) | [capture](web/extra-slack-icon-object.md) | Captured |
| extra-trigger-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/trigger-object) | [capture](web/extra-trigger-object.md) | Captured |
| extra-workflow-object | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/workflow-object) | [capture](web/extra-workflow-object.md) | Captured |
| extra-agents.sessions.setStatus | [Slack](https://docs.slack.dev/reference/methods/agents.sessions.setStatus/) | [capture](web/extra-agents.sessions.setStatus.md) | Captured |
| rendered-lists-create | [Slack](https://docs.slack.dev/reference/methods/slackLists.create/) | [capture](web/rendered-lists-create.md) | Captured |
| rendered-elements | [Slack](https://docs.slack.dev/reference/block-kit/block-elements/) | [capture](web/rendered-elements.md) | Captured |
| rendered-work-objects-implementation | [Slack](https://docs.slack.dev/messaging/work-objects-implementation/) | [capture](web/rendered-work-objects-implementation.md) | Captured |
| rendered-composition | [Slack](https://docs.slack.dev/reference/block-kit/composition-objects/) | [capture](web/rendered-composition.md) | Captured |
| rendered-blocks | [Slack](https://docs.slack.dev/reference/block-kit/blocks/) | [capture](web/rendered-blocks.md) | Captured |

## Local documentation

These captures supplement the code snapshots with the existing UI tutorial, offline testing plan, and vault ownership article. Local Markdown was converted to HTML with Pandoc and processed by Defuddle; exact originals are retained. `local-docs/catalog.json` records origin paths and hashes.

- [01-using-the-go-side-ui-dsl-for-discord-bots.md](local-docs/01-using-the-go-side-ui-dsl-for-discord-bots.md) — `/home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/pkg/doc/tutorials/using-the-go-side-ui-dsl-for-discord-bots.md`
- [02-02-full-local-testing-plan-and-slack-mock-evaluation.md](local-docs/02-02-full-local-testing-plan-and-slack-mock-evaluation.md) — `/home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/ttmp/2026/09/10/DISCORD-SLACK-001--add-slack-support-to-discord-bot/design-doc/02-full-local-testing-plan-and-slack-mock-evaluation.md`
- [03-goja-runtime-ownership-and-context-propagation.md](local-docs/03-goja-runtime-ownership-and-context-propagation.md) — `/home/manuel/code/wesen/go-go-golems/go-go-parc/Research/KB/Tribal/goja-runtime-ownership-and-context-propagation.md`
