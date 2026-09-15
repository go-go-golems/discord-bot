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
Summary: "Chronological research decisions, failures, evidence, and delivery."
LastUpdated: 2026-09-14T22:50:30.724326162-04:00
WhatFor: "Review how the Slack UI recommendation was derived."
WhenToUse: "Review or resume SLACK-UI-001."
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

## Step 3: Validate and deliver the research handoff

Completed the source archive and delivered the corrected 20-page reading copy to reMarkable. The research guide is approximately 8,000 words with 41 footnotes; it recommends a small core implementation and leaves all implementation tasks open.

### Prompt Context
**User prompt (verbatim):** "continue"
**Assistant interpretation:** Finish the interrupted research delivery and commits.
**Inferred user intent:** Receive the completed ticket and readable guide without additional approval loops.
**Commit (code):** d2ef195 — "docs(slack): research surfaces and design pragmatic UI DSL"

### What I did
- Added three relevant local documents: the existing Discord UI tutorial, the prior local-testing plan, and the vault ownership article, with exact originals and Pandoc-to-Defuddle conversions.
- Ran scripts/05-check-research.py: authored links, 41 footnotes, 117 distinct official source URLs, and recorded hashes passed.
- Ran docmgr doctor --ticket SLACK-UI-001 --stale-after 30: all checks passed.
- Delegated mechanical rendering and upload under the remarkable-upload skill; kept the primary guide frozen during finalization.
- Created a reading copy with Graphviz-rendered versions of the two Mermaid diagrams and visually checked representative pages.

### Why
- A self-contained archive makes future implementation independent of live documentation changes.
- The final delivered rendering must use the same settings as the validated local reading artifact.

### What worked
- SDK probe, reference integrity, and docmgr validation passed.
- The corrected final PDF is 20 pages with diagrams, tables, prose, pseudocode, and API references.
- Cloud upload reported: `OK: uploaded SLACK-UI-001 Slack UI Research Guide.pdf -> /ai/2026/09/14/SLACK-UI-001`.

### What didn't work
- The first uploaded reading copy retained archive links relative to the primary guide, rather than its print directory. The local correction therefore differed from the initial cloud artifact.
- Corrected those links and regenerated with the same settings used by the bundle uploader. Uploaded under a new name instead of overwriting the earlier document or risking annotation loss. The earlier upload remains; the final name above is authoritative.
- The initial custom render was 17 pages; the actual bundle settings produce 20. Validation was repeated on the corrected bundle output rather than relying on the earlier rendering.

### What I learned
- A successful upload alone does not establish that the locally inspected PDF used the uploader's render settings.
- The surface catalog needed browser rendering, while individual references generally worked with direct Defuddle extraction.

### What was tricky to build
- Preserving readable diagrams and relative source links when moving a Markdown document into a print directory.
- Keeping documentation breadth separate from implementation commitments: only research and delivery tasks are checked.

### What warrants a second pair of eyes
- Review the proposed ACK lifecycle before implementation and verify newer surface schemas only when adopting them.
- Visual review sampled contents, diagrams, a narrow table, pseudocode, and JavaScript pages; it was not an exhaustive page-by-page typography audit.

### What should be done in the future
- Implement the first rich-message phase when authorized, then action routing and the modal editor.
- Treat the newer surfaces as optional use-case-driven work.

### Code review instructions
- Start with design-doc/01-slack-surfaces-and-ui-dsl-intern-guide.md and tasks.md.
- Run `python <ticket>/scripts/05-check-research.py` and `docmgr doctor --ticket SLACK-UI-001 --stale-after 30` from the repository root.
- Consult artifacts/delivery-receipt.txt for exact render/upload settings and final PDF SHA-256.

### Technical details
- Final local artifact: `artifacts/final/SLACK-UI-001 Slack UI Research Guide.pdf`.
- SHA-256: `ee19bcb3898de628e2551974ad8480a5c0dff2fb82007f44e25e25486855c4ba`.
- The uploader regenerates from the same reading copy and settings; metadata timestamps may differ between PDF generations.
- No runtime feature code, private credential files, or live bot process was changed. No repo push was requested for this ticket.

Final link-check correction: adding the PDF link exposed that the small checker did not recognize CommonMark angle-bracket destinations containing spaces. It reported `index.md: missing <artifacts/final/SLACK-UI-001 Slack UI Research Guide.pdf>`. Updated the checker to remove the angle delimiters before resolving the path; the PDF link itself was valid.

## Step 4: Narrow the execution design while keeping ACK semantics

The implementation scope was clarified after reviewing the complexity scale. Layers 1–7 remain in scope because raw Block Kit values, helpers, builders, interaction routing, modal behavior, and ACK handling form one coherent local UI path. Scheduler behavior is a separate performance concern and is explicitly deferred.

### Prompt Context
**User prompt (verbatim):** "Mark 8 explicitly deferred and remove the worker reservation and bounded interactive execution. We basically want to do 1-7, because the ACK handling is at the core of it and warrants the compelxity."
**Assistant interpretation:** Update the design, tasks, and delivery copy to preserve the ACK state machine while removing scheduler-like execution machinery.
**Inferred user intent:** Keep protocol correctness where Slack requires it, and avoid inventing a second worker system for a local bot.

### What I did
- Made layers 1–7 explicit in the scope section and added a task-file scope decision.
- Removed worker reservation, interactive worker-slot pseudocode, bounded interactive execution, busy worker responses, and errgroup tracking specific to interactions.
- Kept the Go-owned single-use ACK, response-kind validation, deadline rejection, duplicate semantic ACK replay, and existing invocation lifecycle.
- Updated the design decision record to name “deadline-aware acknowledgments with existing execution.”

### Why
- Slack's modal and suggestion protocols require response-bearing ACK handling. That complexity belongs in the transport contract.
- Worker pools, admission budgets, and scheduling solve a different problem and are unnecessary for the current local bot.

### What worked
- The design now says explicitly that layers 1–7 are in scope and layer 8 is deferred.
- The pseudocode dispatches interactive work through the existing host path and lets the ACK object reject late responses.

### What didn't work
- No implementation or live Slack behavior was changed in this documentation-only adjustment.

### What I learned
- The right boundary is “complex ACK, simple execution.” A late ACK should fail normally; the framework does not need to pre-reserve a worker or manufacture a busy response.

### What was tricky to build
- Removing scheduler language without weakening the deadline and duplicate-response invariants required changing both the prose and the pseudocode, not only the decision record.

### What warrants a second pair of eyes
- Review the semantic duplicate ACK behavior once the transport is implemented. It is part of the ACK contract, while execution scheduling remains deferred.

### What should be done in the future
- Re-upload the revised guide and implement phases 1–3 under the clarified scope. Revisit scheduling only after measurements from real local use.

### Code review instructions
- Search the guide for “worker”, “bounded”, and “scheduler”; scheduler language should describe the deferred boundary, not an implementation requirement.
- Run the ticket integrity script, `docmgr doctor`, and the same reMarkable rendering workflow.

### Technical details
- The existing invocation gate and ordinary handler timeout remain unchanged.
- Interactive ACK deadline state belongs to the single-use ACK object; it is not a new handler execution budget.
