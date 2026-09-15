---
Title: Investigation diary
Ticket: SLACK-UI-001
Status: active
Topics:
    - slack
    - architecture
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-14T22:50:30.724326162-04:00
WhatFor: ""
WhenToUse: ""
---

# Investigation diary

## Step 1: Establish the implementation baseline and source archive

Created a new research ticket and traced the current Slack dispatcher and Discord builders. The current Slack API is text-only; the Discord DSL provides useful construction patterns, but its wire payloads and interaction lifecycle cannot be copied directly.

### Prompt Context
**User prompt (verbatim):** "Let's do some deep research on the slack surfaces and how to port / adapt the discord bot UI dsl. Create a new docmgr ticket and research hard and Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Research and design with a source archive and reMarkable delivery.
**Inferred user intent:** Make a pragmatic Slack-native UI implementation understandable and actionable for an intern.

### What I did
- Created SLACK-UI-001 with docmgr; read both host paths and the pinned slack-go SDK.
- Read the vault runtime-ownership article and archived relevant source snapshots.
- Created scripts/01-capture-sources.py to capture official documentation using Defuddle.

### Why
- Current source determines the actual migration boundaries; older design documents can overstate implemented functionality.

### What worked
- Located the UI module exports, Proxy builders, normalization boundary, Slack invocation gate, and SDK acknowledgment methods.

### What didn't work
- Initial Defuddle batch failed with `Error loading content: getaddrinfo ENOTFOUND docs.slack.dev`. Retried the same capture with network escalation.
- A zsh glob search failed with `zsh:1: no matches found: internal/jsdiscord/*flow*`; searching ui_helpers.go located newFlowHelper.

### What I learned
- ui.flow is a namespace helper, not a persistent workflow engine.
- ui.message.file currently discards its content and stores a nil Reader placeholder.
- Interactive payloads are currently acknowledged and dropped by the Slack decoder.

### What was tricky to build
- The existing whole-invocation gate can delay UI work even while the VM owner is free during network I/O; the design must address that separately from VM ownership.

### What warrants a second pair of eyes
- Modal acknowledgments, expiring trigger IDs, and Go-to-JS value ownership.

### What should be done in the future
- Implement only after reviewing the proposed phases and acceptance criteria.

### Code review instructions
- Begin with internal/slacktransport/run.go, internal/jsslack/dispatch.go, and internal/jsdiscord/ui_module.go.

### Technical details
- Baseline revision is recorded in sources/code/REVISION.txt; SDK is v0.17.3.
- No runtime feature code or live bot process was changed.
