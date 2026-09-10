---
Title: Add Slack support to discord-bot
Ticket: DISCORD-SLACK-001
Status: active
Topics:
    - discord-bot
    - architecture
    - api-design
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-10T18:20:32.045367429-04:00
WhatFor: ""
WhenToUse: ""
---

# Add Slack support to discord-bot

Design an independent Slack implementation in the existing Go module, inspired by the Discord bot host. The documentation is complete; implementation remains planned.

- [Intern architecture and implementation guide](design-doc/01-slack-support-architecture-and-intern-implementation-guide.md)
- [Research sources and snapshots](sources/README.md)
- [Investigation diary](reference/01-investigation-diary.md)
- [Tasks](tasks.md)
- [Changelog](changelog.md)

Recommended first milestone: a separate `cmd/slack-bot`, Socket Mode, one workspace, slash ping and threaded mention replies. See guide section 12 for token and test access requirements.

## Delivery

[Reviewed 21-page PDF](artifacts/slack-intern-guide.pdf). Uploaded successfully to reMarkable cloud at `/ai/2026/09/10/DISCORD-SLACK-001` as `DISCORD-SLACK-001 Intern Implementation Guide`.
