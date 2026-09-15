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

## Step 2: Research surface contracts and constrain the implementation

Archived the Slack surface documentation and traced the interactions needed for a complete local note editor. The guide recommends three small implementation phases and keeps the rest of the surface catalog as optional reference material.

### Prompt Context
**User prompt (verbatim):** "save all the relevant pages and documentation to the sources folder of the ticket. Feel free to use browser tools or script them to process the documents. commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)"

**User prompt (verbatim):** "and despite me telling you to write a detailed / textbook style document doesn't mean this applies to the design itself, we don't want to astronaut engineer, as you are sometimes wont to do. we love solid engineering and take pride in being thorough, but we are also pragmatic"

**Assistant interpretation:** Preserve research thoroughly while keeping the recommended implementation small.
**Inferred user intent:** Make careful engineering decisions without creating unnecessary architecture.
**Commit (code):** d8d1027 — "docs(slack): establish UI research ticket and source baseline"

### What I did
- Captured official pages with Defuddle, browser-rendered catalog tables, and linked element/composition references.
- Ran scripts/02-run-sdk-probe.sh against slack-go v0.17.3.
- Wrote the intern guide with a full note-edit example, surface matrix, deadline pseudocode, file-level phases, API references, and test expectations.
- Kept messages, actions, and one modal as the core; Home and all other surfaces are independent follow-ups.

### Why
- The interaction acknowledgment protocol is a real correctness requirement even for one local bot.
- SDK breadth and documentation breadth must not force a matching implementation scope.

### What worked
- The SDK probe preserved section blocks but reduced table and task_card blocks to their type, demonstrating unknown-block data loss.
- Browser rendering populated compatibility tables absent from static extraction.
- Browser-captured HTML successfully passed through Defuddle for the two URL captures that timed out.

### What didn't work
- The first network-enabled batch stopped with `subprocess.TimeoutExpired` for `https://docs.slack.dev/reference/methods/slackLists.create/` after 60 seconds. Added bounded exception handling; lists and Work Objects implementation still timed out after 35 seconds, then were recovered through the browser.
- Browser code attempting filesystem import failed with `TypeError [ERR_VM_DYNAMIC_IMPORT_CALLBACK_MISSING]: A dynamic import callback was not specified.` Used browser_evaluate's supported filename output instead.
- Static HTML parsing showed empty catalog tables. Retained that evidence and captured rendered HTML and text.
- docmgr doctor warned `unknown_topics` for slack. Added the topic before final validation.

### What I learned
- Newer references include cards, containers, interactive tables, charts, and agent-session changes; the older Assistant examples are not the whole current API.
- Card and alert tables conflict with examples. Those features need a live serialization check before implementation.
- Ordinary Slack buttons do not expose Discord's disabled-button property.

### What was tricky to build
- The existing gate serializes complete invocations, independently of VM ownership. The design retains it with bounded admission and a retry limitation rather than adding a scheduler.
- Modal duplicate handling must preserve the original semantic ACK, because an empty duplicate ACK can accept an invalid form.

### What warrants a second pair of eyes
- The single-use ACK transition and cancellation before side effects.
- Surface restrictions and differences between wire documentation and the pinned SDK.

### What should be done in the future
- Implement core phases only after design review; verify optional newer blocks in the actual workspace when needed.

### Code review instructions
- Read the guide's note-editor example and phases, then compare sources/sdk-probe-output.txt with the probe code.
- Review sources/catalog-complete.json for source provenance and recovered captures.

### Technical details
- No production Go code, credentials, manifests, or running Slack process was changed.
- Plain message fallback is an explicit host policy; unsupported block fields are not silently accepted.

Source-integrity check during this step found two missing links to the original timed-out capture filenames: `web/41-reference-methods-slackLists.create.md` and `web/49-messaging-work-objects-implementation.md`. Both had recovered browser captures, but the catalog continuation skipped their recovery metadata. Corrected the catalog and README links to the rendered captures, then reran the check.
