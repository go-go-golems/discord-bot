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

## Step 5: Re-render and deliver the clarified scope

Regenerated the reading copy after the scope edit and uploaded it under a new name so no prior reMarkable document or annotations were overwritten. The v2 rendering contains the explicit layers 1–7 scope and the deferred layer 8 boundary.

### Prompt Context
**User prompt (verbatim):** "Mark 8 explicitly deferred and remove the worker reservation and bounded interactive execution. We basically want to do 1-7, because the ACK handling is at the core of it and warrants the compelxity."
**Assistant interpretation:** Re-render and redeliver the edited research guide after committing the scope correction.
**Inferred user intent:** Ensure the document on reMarkable matches the agreed implementation scope.
**Commit (code):** bb959d1 — "docs(slack): defer UI scheduling and retain ACK scope"

### What I did
- Refreshed the print copy from the revised primary guide while retaining the Graphviz diagrams and corrected archive links.
- Rendered and visually inspected the v2 diagram, scope/table page, modal sequence, pseudocode, and JavaScript layout.
- Uploaded v2 as `SLACK-UI-001 Slack UI Research Guide v2` without `--force`.

### Why
- The primary design changed materially; the previous uploaded PDF was stale.
- A new remote name preserved earlier documents and their annotations.

### What worked
- Final v2 PDF: 21 pages, SHA-256 `6bb3cc6156adc95ff9830428d1d2a53413b6a046915ffdfe5dd0e082697e3853`.
- Upload receipt: `OK: uploaded SLACK-UI-001 Slack UI Research Guide v2.pdf -> /ai/2026/09/14/SLACK-UI-001`.
- Source integrity and docmgr checks remained valid after the primary-document edit.

### What didn't work
- The earlier uploaded documents cannot be replaced without `--force`; v2 was intentionally uploaded under a new name.

### What I learned
- Scope decisions affecting pseudocode require a fresh reading-copy render, even when no production code changes.

### What was tricky to build
- Keeping the source archive links and Graphviz assets valid after copying the revised Markdown into the print directory.

### What warrants a second pair of eyes
- Review the ACK state machine during implementation; the scheduler boundary is now explicit and deferred.

### What should be done in the future
- Implement layers 1–7 in phases 1–3. Do not add layer 8 unless local measurements demonstrate a need.

### Code review instructions
- Compare the primary guide, print copy, and v2 receipt. Run the research integrity script and docmgr doctor.

### Technical details
- v2 was uploaded with `remarquee upload bundle`, `--toc-depth 2`, 0.72-inch margins, DejaVu fonts, and `--non-interactive`.
- Cloud receipt proves cloud delivery, not physical device synchronization.

## Step 6: Implement the first rich-message and builder slice

The first production code checkpoint extends the existing text-only message contract without introducing a second runtime or scheduling subsystem. The implementation keeps Block Kit as detached JSON at the domain boundary, converts it losslessly at the Slack transport edge, and exposes a small native `slack/ui` module that produces ordinary JSON-compatible objects.

### Prompt Context
**User prompt (verbatim):** "continue."
**Assistant interpretation:** Proceed from the completed research handoff into the agreed implementation phases, starting with Phase 1.
**Inferred user intent:** Begin implementing layers 1–7 while preserving the explicit deferral of layer 8 scheduling machinery.

### What I did
- Added `slackbot.Block`, `MessagePayload`, and rich `PostMessage` validation with a mandatory fallback text and a 50-block cap.
- Added an optional `RichResponder` capability so response URLs can carry blocks while existing text-only responders remain valid.
- Extended the Socket/Web API transport to preserve arbitrary validated block fields through a small lossless SDK wrapper instead of `slack.UnknownBlock`.
- Registered `require("slack/ui")` with plain-text, mrkdwn, section, action, divider, header, button, and fluent message builders.
- Added a runtime test that constructs a section and a danger button, posts the resulting message, and inspects the detached payload.

### Why
- The Slack SDK's `UnknownBlock` round-trip drops fields for newer blocks. A transport-owned JSON wrapper is sufficient for the low-level message boundary and avoids an SDK upgrade for the core slice.
- Rich message support must work for both returned handler values and explicit `ctx.slack.messages.post` calls, while the existing text-only response contract should not be widened unnecessarily.
- The builder API follows the documented Slack-native construction style and keeps builders separate from the transport client.

### What worked
- `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./internal/jsslack` passed, including the new rich-message builder test.
- `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./... -run '^$'` compiled every package successfully.
- `git diff --check` passed.

### What didn't work
- The first test invocation used the repository's default Go cache and failed because `/home/manuel/.cache/go-build` is read-only in this environment.
- The transport test package cannot start its existing `httptest` IPv6 listener under the sandbox (`listen tcp6 [::1]:0: operation not permitted`); compilation succeeds and the failure is environmental rather than a test assertion.

### What I learned
- Rich payloads can be added without changing ordinary message service ownership: the existing `MessageService.Post` remains the network boundary, and only the payload shape grows.
- Goja's object spread works with the detached objects returned by the builders, which allows `ctx.slack.messages.post({channelId, ...ui.message(...).build()})`.

### What was tricky to build
- The response URL path has a narrower existing `Responder` interface. The additive `RichResponder` check lets block responses work where supported and returns a clear capability error otherwise.
- Builder objects need a `build` method but must become plain JSON before entering the strict Go decoder; a small `valueMap` helper handles both builders and raw objects.

### What warrants a second pair of eyes
- Review block-shape validation before adding action routing. The current low-level layer verifies type and count but intentionally does not claim to validate every Block Kit schema.
- Review the response URL JSON fixture once the local HTTP test harness can bind a loopback listener.

### What should be done in the future
- Add TypeScript declarations and an example bot for the new `slack/ui` module, then mark Phase 1 complete.
- Implement interactive envelope decoding and action registration next; keep ACK handling in the transport boundary and do not add a scheduler.

### Code review instructions
- Start with `pkg/slackbot/model.go`, `internal/slacktransport/client.go`, `internal/jsslack/ui_module.go`, and `internal/jsslack/dispatch.go`.
- Run `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./internal/jsslack` and the compile-only repository check. The existing transport test requires a sandbox that permits an IPv6 loopback listener.

### Technical details
- Checkpoint tests: `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./internal/jsslack`; `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./... -run '^$'`; `git diff --check`.
- No Slack credentials, network calls, or live bot processes were used.

## Step 7: Route button and static-select actions

The second implementation checkpoint adds the first interactive path. Socket Mode `interactive` envelopes containing `block_actions` are decoded into a small Slack-independent action value, acknowledged through the existing receipt object, and dispatched to a JavaScript handler registered by action ID. This keeps ordinary action acknowledgment in transport while leaving modal response-bearing ACKs for the next phase.

### Prompt Context
**Prompt source:** Continued implementation after the Phase 1 checkpoint; the active user request remains the earlier instruction to implement the agreed layers 1–7.
**Assistant interpretation:** Implement Phase 2 from the ticket guide: button and static-select routing using the existing ingress and host lifecycle.

### What I did
- Added a normalized `slackbot.Action` attached to an invocation, including action type, action ID, block ID, value, selected option(s), and source message timestamps.
- Added `action(actionId, handler)` registration and `ctx.action` exposure in the JavaScript runtime.
- Extended Socket Mode decoding for `interactive` `block_actions` requests and preserved a validated response URL when Slack supplies one.
- Included interactive requests in the receipt's response-payload capability and routed handler results through the existing reply/message path.
- Added unit coverage for action dispatch and interactive payload decoding, plus TypeScript and offline-guide updates.

### Why
- Button and static-select callbacks use the same Slack `block_actions` envelope, so one normalized action contract covers both without a general interaction framework.
- The action ID is the stable application routing key; values and selected options remain data supplied by Slack and are not interpreted by the transport.
- Ordinary actions can use the existing empty Socket Mode ACK. Modal submissions still need explicit response-bearing ACK choices and are intentionally deferred to Phase 3.

### What worked
- `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./internal/jsslack ./internal/slacktransport -run 'TestAction|TestDecodeBlockAction|TestSlackUIBuilders'` passed.
- `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./... -run '^$'` compiled every package successfully.
- `git diff --check` passed.

### What didn't work
- No new code failure remained after correcting the action type mapping; the first test exposed that the transport had used the outer `block_actions` type instead of the element's `static_select` type.

### What I learned
- Slack's outer interactive type and inner action element type are separate discriminators. The runtime must preserve both concepts and route on `action_id`.
- Response URLs can be carried privately by the transport while the JavaScript context receives only detached action data.

### What was tricky to build
- Channel IDs may be present under either `channel.id` or `container.channel_id`; message timestamps likewise have message and container forms. The decoder chooses the explicit message value and falls back to the container snapshot.
- Existing invocation validation assumed only commands and events. It now counts command, event, and action kinds and allows channel-less action contexts for later view work.

### What warrants a second pair of eyes
- Verify the ACK-before-handler timing on a real Socket Mode fixture; the existing ingress still owns the admission queue and receipt boundary.
- Review the action response behavior when Slack omits `response_url`; the current fallback posts through the source channel when a handler returns a message.

### What should be done in the future
- Add modal builders, view registration, submitted-state decoding, and explicit single-use deadline-aware ACK responses in Phase 3.
- Add a normalized local interactive fixture to `bots simulate` so action payloads and acknowledgments can be replayed without a socket server.

### Code review instructions
- Review `pkg/slackbot/model.go`, `internal/slacktransport/run.go`, `internal/jsslack/module.go`, and `internal/jsslack/dispatch.go` together because they define one cross-package action contract.
- Run the focused action tests and compile-only repository check above. The existing full Socket Mode integration test still needs an environment that permits its IPv6 `httptest` listener.

### Technical details
- Action routing key: `action:<action_id>`.
- Automatic receipt capability: `accepts_response_payload` is honored for both slash commands and interactive envelopes.
- No scheduler, worker reservation, durable state, or new execution budget was added.

## Step 8: Complete modal flow and offline interactive simulation

The third implementation checkpoint completes the agreed first UI path. A button action can open a modal through the existing Slack client, a view submission is routed by callback ID, and JavaScript chooses a single acceptance or field-error acknowledgment. The offline recorder now makes these choices and modal opens visible in `bots simulate`, which provides a repeatable test without Slack credentials.

### Prompt Context
**Prompt source:** Continued implementation after the Phase 2 checkpoint; the active user request remains the earlier instruction to implement the agreed layers 1–7.
**Assistant interpretation:** Finish the basic modal editor and deadline-aware ACK handling, then update CLI and documentation surfaces.

### What I did
- Added `ModalView`, `ViewService`, `Interaction`, and `InteractionAcknowledger` contracts to the SDK-independent Slack domain package.
- Added modal and plain-text-input builders to `require("slack/ui")`, plus `ctx.openModal`, `ctx.view`, `ctx.values.text`, `ctx.ack.accept`, and `ctx.ack.errors`.
- Added Socket Mode `view_submission` decoding and carried its receipt through ingress without an automatic ACK.
- Implemented a mutex-protected transport receipt that rejects duplicate and late ACK choices and emits Slack's `response_action: errors` payload when requested.
- Added Slack `views.open` conversion using the same lossless block wrapper as messages.
- Extended the offline recorder and `bots simulate` to capture `ack` and `open_view` operations, and updated the showcase bot, declarations, help, and design guide.

### Why
- Modal submission ACKs are protocol responses, not ordinary messages. They must be selected exactly once and before Slack's deadline, while validation errors must be keyed by input block ID.
- A local developer needs to test the complete shape of a modal without repeatedly installing or interacting with a live app. Recorded ACK and view operations make that path inspectable and deterministic.
- The existing host owner and invocation timeout are sufficient for this local framework. Adding scheduling would increase scope without solving a demonstrated problem.

### What worked
- `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./... -run '^$'` compiled all packages.
- Full focused gates passed: `go test ./pkg/slackbot ./internal/jsslack ./pkg/slackcli ./pkg/slackhost` and `go test ./internal/slacktransport` with loopback networking enabled.
- `go run ./cmd/slack-bot bots inspect ui-showcase` reports the command, action, and view registrations.
- `go run ./cmd/slack-bot bots simulate ui-showcase` with action and view fixtures records `open_view` and `ack` operations.
- The research integrity script, `docmgr doctor`, and `git diff --check` passed.

### What didn't work
- The default sandbox cannot run the existing transport integration test because it disallows an IPv6 `httptest` listener; the same test passed with loopback networking enabled.

### What I learned
- The simplest reliable ACK boundary is a Go-owned receipt attached to the normalized invocation. JavaScript receives only narrow methods and never sees the Socket Mode envelope or response URL.
- Offline simulation needs service capabilities, not special fake JavaScript branches: the ordinary `ViewService` and `InteractionAcknowledger` interfaces can record exact operations.

### What was tricky to build
- A view submission has no channel and no response URL, so invocation validation and dispatch cannot reuse command or action assumptions. The callback ID is its routing key, and its state is indexed by block ID and action ID.
- The receipt needs to reserve its single-use transition before sending the network ACK. Otherwise a handler and a timeout could both emit responses.

### What warrants a second pair of eyes
- Validate the real Socket Mode wire shape for `view_submission` and `response_action: errors` with the local mock before a live manual demonstration.
- Review whether any future handler should be allowed to return a message after a view ACK; the current contract requires handlers to choose an ACK and leaves follow-up messaging explicit.

### What should be done in the future
- Exercise the modal editor against the development Slack app and add a concrete message-update service if the note-edit example needs to mutate its source message.
- Keep App Home, shortcuts, external options, newer Block Kit surfaces, files, canvases, Lists, Work Objects, and agent sessions as separate use-case-driven tickets.

### Code review instructions
- Start with `pkg/slackbot/model.go`, `pkg/slackbot/recording.go`, `internal/slacktransport/run.go`, `internal/slacktransport/client.go`, `internal/jsslack/dispatch.go`, and `internal/jsslack/ui_module.go`.
- Replay `/tmp/slack-ui-action.json` and `/tmp/slack-ui-view.json` or create equivalent fixtures under `examples/slack-bots/fixtures/` to inspect the offline operation JSON.
- Run the full package gates listed under “What worked”; use a loopback-enabled environment for the Socket Mode integration test.

### Technical details
- ACK kinds: `accept` maps to an empty Socket Mode payload; `errors` maps to `{"response_action":"errors","errors":{...}}`.
- ACK deadline: three seconds from receipt creation, with the caller context additionally respected.
- Explicit view submissions bypass the ingress automatic ACK; ordinary block actions continue to use the existing automatic empty ACK.
- Successful interactive ACK payloads are retained in a five-minute process-local replay map keyed by Socket Mode envelope ID; a duplicate envelope replays the same payload without re-running JavaScript.

## Step 9: Refresh and upload the implemented guide

The implementation changed the authoritative design from a research-only handoff into a guide that documents the shipped seven-layer UI slice. I refreshed the printable reading copy, generated the v3 PDF, and uploaded that copy to reMarkable.

### Prompt Context

**User prompt (verbatim):** "just upload, it's fine"

**Assistant interpretation:** Upload the current rendered guide immediately and skip another round of visual inspection or other delivery checks.

### What I did

- Ran `scripts/07-refresh-reading-copy.py` to synchronize the reading copy with the implemented guide and replace Mermaid blocks with the archived diagrams.
- Generated `artifacts/final-v3/SLACK-UI-001-Slack-UI-Implementation-Guide-v3.pdf`.
- Uploaded the reading-copy bundle with `remarquee upload bundle` to `/ai/2026/09/15/SLACK-UI-001`.
- Recorded the upload result and artifact checksum in `artifacts/delivery-receipt.txt`.

### What worked

- ReMarkable reported: `OK: uploaded SLACK-UI-001 Slack UI Implementation Guide v3.pdf -> /ai/2026/09/15/SLACK-UI-001`.
- The local v3 PDF is 19 pages and has SHA-256 `70cc99cdf523c40839edc8377ad7e6cd3bc86250d476f88d372904620e29e0b6`.

### What didn't work

- Nothing failed during this delivery step. Additional checks were intentionally skipped at the user's direction.

### What I learned

- The printable copy must be refreshed after implementation because the guide's API references and file references change with the runtime surface.

### What was tricky to build

- The bundle needs the refreshed relative links and the rendered diagram assets to remain readable when Pandoc is run from the ticket's print directory.

### What warrants a second pair of eyes

- None for this upload step; the user accepted the current artifact without another inspection pass.

### What should be done in the future

- Keep the v3 receipt alongside the ticket and regenerate the bundle if the implementation guide changes.

### Code review instructions

- Review the final guide, reading copy, and delivery receipt together with commit `docs(slack): refresh implemented UI guide delivery`.

### Technical details

- Upload destination: `/ai/2026/09/15/SLACK-UI-001`.
- Local PDF: `artifacts/final-v3/SLACK-UI-001-Slack-UI-Implementation-Guide-v3.pdf`.

## Step 10: Log live identity, connection, and command dispatch

The live test exposed insufficient diagnostics: the user saw a pong while the showcase process only logged an unidentified receipt and generic failure. Added metadata logs so the next invocation can be traced to its local bot, installed app, and handler route.

### Prompt Context
**User prompt (verbatim):** "can you log more?"
**Assistant interpretation:** Add practical diagnostics and restart the live bot.
**Inferred user intent:** Determine which application and handler are answering Slack commands.

### What I did
- Logged profile, script, bot, app ID and workspace at startup.
- Logged successful authentication and Socket Mode lifecycle event types.
- Added envelope IDs and command/event names to admission logs.
- Added dispatch start/completion, duration, classified failures, and missing handler route logs.
- Restarted the showcase in tmux session `slack-ui-showcase`.

### Why
- Startup alone did not prove a socket connection. Generic failure logs did not identify the missing route.
- Raw payloads, tokens, submitted values and arbitrary JavaScript error messages are unnecessary for this diagnosis.

### What worked
- `GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./internal/jsslack ./pkg/slackcli` passed.
- Focused transport ACK tests passed.
- Live startup reported app `A0C1YJCCCP6`, workspace `T0C1UJMCPGA`, bot user `U0C1NREFWBV`, then `connecting`, `connected`, and `hello`.

### What didn't work
- The earlier attribution of pong to another process was unsupported. The source of that reply remains unconfirmed until a new interaction is traced.

### What I learned
- The selected local script and the registered Slack app manifest are separate state; a successful socket connection does not register new commands.

### What was tricky to build
- Dispatch errors can contain arbitrary script data. Logging classifies errors and records the route separately instead of dumping raw errors.

### What warrants a second pair of eyes
- Confirm the next user command's receipt and dispatch logs agree with the Slack response.

### What should be done in the future
- Register the showcase slash command in the installed app before testing that route.

### Code review instructions
- Review `pkg/slackcli/run_remote.go`, `internal/slacktransport/run.go`, and `internal/jsslack/dispatch.go`.
- Capture the tmux pane after a user-triggered command.

### Technical details
- Dispatch duration is logged as `duration_ms`; unknown command routes produce `error_code=not_found`.
- Logs omit tokens, response URLs, message text, and form values.

## Step 11: Sync the selected bot manifest before connecting

The runner previously selected JavaScript without changing Slack's app configuration. This left the app advertising /golem-ping while the process registered /ui-showcase. Remote startup now updates the app from the selected descriptor before connecting, with an explicit opt-out.

### Prompt Context
**User prompt (verbatim):** "it should force update when starting by default i think?"
**Assistant interpretation:** Make manifest synchronization the default startup behavior.
**Inferred user intent:** Switch local bot implementations without manually editing Slack app settings.

### What I did
- Added a bounded apps.manifest.update request using the selected profile's stored management access token.
- Added --skip-manifest-update to the run verb. Default startup sends the complete generated manifest and reuses the app ID and runtime tokens.
- Stop before connecting if Slack reports permissions_updated=true, with an install command to grant changed scopes.
- Added API request, error, permission-change and default CLI wiring tests.
- Fixed Client.Close to close the HTTP client's idle connections: remote clients do not set the local transport field.
- Updated the CLI discovery test for the second example, README, and embedded help.
- Archived manifest update and app lifecycle documentation with Defuddle in sources/118 and sources/119.
- Refreshed the expired management token using credentials refresh, then restarted the tmux showcase.

### Why
- Slack command registration is remote configuration; loading the JavaScript script does not register it.
- Sending the generated manifest makes the selected bot authoritative, including display name, events and scopes.
- Permission changes need an explicit install; errors stop startup instead of silently running against stale configuration.

### What worked
- Full repository tests passed with GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./... -count=1.
- go build -buildvcs=false ./... and go vet -buildvcs=false ./... passed with the same cache/workspace settings.
- Help displays --skip-manifest-update.
- Live sync returned app_id=A0C1YJCCCP6 and permissions_updated=false; authentication and Socket Mode connecting/connected/hello followed.
- Stored credentials refresh succeeded without exposing either token.

### What didn't work
- Initial run-default test panicked at net/http.(*Transport).CloseIdleConnections(0x0). Fixed the remote cleanup path through http.Client.CloseIdleConnections.
- TestOfflineCLI expected one example but discovered two. Updated the expected inventory.
- Build and vet initially reported "error obtaining VCS status: exit status 128", even outside the sandbox. Disabling VCS stamping resolved these checks.
- make glazed-lint reports existing raw Cobra flags in pkg/slackcli/credentials.go.
- make govulncheck reports eight reachable vulnerabilities in existing Go 1.26.4, golang.org/x/text v0.37.0, and excelize/v2 v2.10.0 dependencies. Dependency upgrades are outside this manifest change.
- First live sync failed with token_expired. Ran credentials refresh --profile go-go-golems and restarted successfully.

### What I learned
- Slack confirmed this command/configuration switch needed no reinstall.
- Default manifest sync requires a current management token even when runtime tokens remain valid.

### What was tricky to build
- Runtime credentials and script configuration must be checked before changing remote configuration; the sync runs after host loading and before client.Run.
- A changed-scope response is already an applied manifest update. The error explicitly says to reinstall instead of implying that the update failed.

### What warrants a second pair of eyes
- The update replaces the full app configuration with the generated manifest; manual settings not represented there are replaced.
- Arbitrary Slack response bodies are not logged. Safe error codes are returned, with refresh guidance for expired authentication.

### What should be done in the future
- Trace the user's next /ui-showcase interaction in tmux.
- Address existing lint and dependency findings independently.

### Code review instructions
- Start with pkg/slackcli/run_remote.go, update_manifest.go, and update_manifest_test.go.
- Review the new run flag in commands.go and the cleanup fix in internal/slacktransport/client.go.
- Use the focused CLI tests and inspect the live startup logs.

### Technical details
- POST https://slack.com/api/apps.manifest.update with bearer management token, app_id, and manifest as a JSON-encoded string.
- Sources: https://docs.slack.dev/reference/methods/apps.manifest.update/ and https://docs.slack.dev/app-management/distribution/.
- Session: slack-ui-showcase. Bot: ui-showcase. Profile: go-go-golems.

## Step 12: Investigate the unexpected pong reply

The user observed that pong stopped when the showcase process stopped. That is evidence of a connection to our runner, but does not identify the generating code path. Earlier assertions attributing the reply to another instance were unsupported. This investigation has not established the root cause.

### Prompt Context
**User prompt (verbatim):** "ok, figure it out"
**Assistant interpretation:** Trace the unexplained pong response and reproduce it.
**Inferred user intent:** Account for the response with evidence rather than another speculative explanation.

### What I did
- Inspected bot discovery, runtime loading, command registration, dispatch, response URL delivery, ACK construction, and pinned slack-go v0.17.3 Socket Mode parsing.
- Checked process names while the runner was stopped; no other Go, Slack, Node or Deno runner appeared in that process listing.
- Added a local wire-level regression reproducing CLI discovery (including ping inspection), loading the showcase, and delivering /golem-ping followed by /ui-showcase over a real local WebSocket.
- Captured both Socket Mode ACK payloads and response-URL HTTP messages in the test.
- Added reply fingerprints, byte counts, block counts, and delivery success logs correlated with bot and invocation. Startup now includes PID.
- Exported the live manifest read-only: app A0C1YJCCCP6 is named ui-showcase and registers only /ui-showcase.
- Restarted the instrumented showcase in tmux with --skip-manifest-update for a fresh user-triggered reproduction.

### Why
- Local source searches cannot explain an observed live reply by themselves. The reproduction exercises the SDK and actual response encoding.
- Inspecting both scripts before loading the showcase tests the hypothesis that discovery leaks ping handlers into the runtime.
- Text hashes identify a literal pong without logging user message bodies or credential-bearing response URLs.

### What worked
- TestShowcaseWireDoesNotReplyPong passed: both ACKs had no payload, and the sole HTTP reply contained the showcase fallback text and blocks.
- Host and CLI tests passed after adding outbound diagnostics.
- Live manifest export confirmed that /golem-ping is absent from the current app configuration.
- Instrumented PID 132981 authenticated and received Socket Mode connected/hello.

### What didn't work
- First wire test imported slackcli from an internal slacktransport test and triggered "import cycle not allowed in test"; changed to the external slacktransport_test package.
- Next fixture timed out waiting for an ACK because slack-go SlashCommand.UnmarshalJSON requires is_enterprise_install even when false. Adding that real payload field fixed the fixture.
- No local test has reproduced the unexplained pong. This is not a root-cause fix.

### What I learned
- The current showcase rejects /golem-ping and emits only its own UI message in the tested discovery/runtime/transport sequence.
- The SDK's WebSocket PONG control frames do not implement slash-command text replies.
- The existing live log contained a successful showcase command and modal submission, but no attributed pong emission.

### What was tricky to build
- A missing SDK-required slash-command field prevents the event from reaching our decoder; the mock must match the real wire payload, not just the framework's normalized invocation.

### What warrants a second pair of eyes
- The fresh live reproduction remains necessary; do not interpret the passing mock as disproving the user's observation.

### What should be done in the future
- Capture the next /ui-showcase and any still-accepted /golem-ping attempt, then compare incoming command, envelope, and outgoing fingerprint.
- Keep the cause explicitly unresolved until evidence identifies it.

### Code review instructions
- Review internal/slacktransport/showcase_wire_test.go and the sendReply logs in internal/jsslack/dispatch.go.
- Run GOCACHE=/tmp/go-build-cache-slack-ui GOWORK=off go test ./internal/slacktransport -run TestShowcaseWireDoesNotReplyPong -v with loopback access.

### Technical details
- SHA-256 of literal pong: 9795c5ff8937f23526ccb207a5684c1fc94a7854e19c021b39d944e51f5baef2.
- SHA-256 of showcase fallback: ca43944b37f349fbf2841f100250fee5d35e7a800c2c6c29cf548c1a513a38ff.
- Logs report text_sha256, text_bytes, blocks and delivered; no message bodies or response URLs.

## Step 13: Split general Slack documentation from the UI DSL guide

The offline help title no longer represented the installed-app lifecycle or live interactive runtime. Renamed the general guide and added an embedded tutorial/reference specifically for the implemented UI API.

The user also reported that /golem-ping now fails and /ui-showcase works, and suggested the earlier pong observation may have been a misreading. There is no reproduced code defect behind that report; the investigation's test and diagnostics remain useful, but no root-cause fix is claimed.

### Prompt Context
**User prompt (verbatim):** "ok, do it"
**Assistant interpretation:** Rename slack-offline to slack-bot-guide and write a separate slack-ui-dsl page, as proposed in the preceding response.
**Inferred user intent:** Make current runtime documentation discoverable under accurate names and give bot authors a practical UI reference.

### What I did
- Renamed pkg/slackdoc/slack-offline.md and its help slug to slack-bot-guide.
- Added slack-ui-dsl.md with message construction, button/modal routing, form state, ACK semantics, builder tables, supported limits, troubleshooting and source references.
- Embedded both pages and updated README links and help commands.
- Removed stale claims that live connections, manifest updates, buttons and modals were unimplemented.
- Corrected fixture documentation to cover action and interaction invocations.
- Added checked-in UI command/action/view fixtures and filled missing TypeScript declarations for view registration, ui.input and modal input overloads.

### Why
- General operations and UI authoring have different entry points; the old offline label obscured both.
- API examples must describe the actual running implementation, especially that accepting a modal does not persist data or update a message.

### What worked
- Both go run ./cmd/slack-bot help slack-bot-guide and help slack-ui-dsl rendered the intended embedded pages.
- All three documented bots simulate ui-showcase commands ran and recorded ephemeral_reply, open_view and ack respectively.
- go test ./cmd/slack-bot ./pkg/slackdoc passed with the established cache and GOWORK=off.
- No slack-offline references remain outside historical ticket files.

### What didn't work
- No failures in this documentation validation.
- Reading the builder exposed an existing limitation: textInput().optional() puts optional on the element rather than the containing input block. The new guide documents this and the raw-input-block workaround.

### What I learned
- The declarations had omitted APIs already implemented in Go; view, ui.input and the modal input-block overload are now represented.
- Research proposals and runtime API documentation must identify their different implementation status explicitly.

### What was tricky to build
- Avoiding claims of complete Block Kit validation or unsupported operations while documenting the raw JSON escape hatch.

### What warrants a second pair of eyes
- Compare the builder reference with internal/jsslack/ui_module.go when adding new methods.
- Fix the optional-input builder separately with a wire-shape regression test.

### What should be done in the future
- Complete a stateful editor with a visible save result; the tutorial accurately documents the current log-only example.

### Code review instructions
- Read pkg/slackdoc/slack-bot-guide.md and slack-ui-dsl.md, then inspect doc.go's embedding and the README links.
- Replay examples/slack-bots/fixtures/ui-*.json with bots simulate.

### Technical details
- New slugs: slack-bot-guide and slack-ui-dsl.
- The old slug was removed without a compatibility alias.
- This is a documentation/declaration/fixture change; the live runtime was not restarted.
