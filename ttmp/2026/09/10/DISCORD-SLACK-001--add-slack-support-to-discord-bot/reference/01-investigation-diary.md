---
Title: Investigation diary
Ticket: DISCORD-SLACK-001
Status: active
Topics:
    - discord-bot
    - architecture
    - api-design
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://internal/jsslack/dispatch.go
      Note: Owner-thread async dispatch and reply lifecycle
    - Path: repo://pkg/slackbot/ingress.go
      Note: Bounded offline admission and receipt contract
    - Path: repo://pkg/slackcli/commands.go
      Note: Offline inspect manifest and simulation
    - Path: repo://pkg/slackdoc/slack-offline.md
      Note: Implemented API and operator commands
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-10T18:20:32.32301754-04:00
WhatFor: ""
WhenToUse: ""
---

# Investigation diary

## Step 1: Repository investigation and Slack design

Mapped the existing Discord host, discovery, configuration, runtime owner, response handling, and xgoja integration before designing Slack support. The user clarified that a separate implementation following the same concepts is welcome, so the guide uses new Slack packages and a new binary in the existing module.

### Prompt Context

**User prompt (verbatim):**

> Create a new docmgr ticket to add slack support to discord-bot . Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.
>
> Then tell me what tokens and stuff you need to test / develop yourself

**Follow-up (verbatim):**

> save all useful resources to the ticket's sources folder, with defuddle for html/markdown

**Follow-up (verbatim):**

> yeah, use inspiration from discord-bot but feel free to make a separate implementation following the same concepts. You can also look at ~/code/wesen/go-go-golems/go-go-parc for more information about go-go-goja JS api design and also the work behind discord-bot

**Assistant interpretation:** Deliver design, sources, testing prerequisites, and reMarkable upload; no runtime implementation or live Slack messaging requested.

**Inferred user intent:** Give an intern enough architectural and operational context to build Slack support independently.

**Code revision:** `2ea219e6d0a37bb27916dc29a2ba88c1f03b6182`; initial git status clean.

### What I did

- Read README, configuration, lifecycle, JS dispatch, discovery, provider, and runtime ownership source.
- Created ticket and guide through docmgr.
- Researched official Slack APIs and the slack-go repository.
- Archived 14 remote references with defuddle and 5 vault notes with original and defuddle-extracted Markdown copies.
- Archived 19 current code/document inputs with original paths and hashes.
- Wrote the implementation guide with typed contracts, diagrams, pseudocode, phases, test cases, and credential setup.

### Why

Slack acknowledgment, command registration, identity, threading, and rendering differ from Discord. A separate implementation makes those contracts explicit without a compatibility layer. The vault reinforces Go ownership of lifecycle and domain state, and JS ownership of bot composition.

### What worked

`docmgr ticket create-ticket` and document creation succeeded. Public source extraction succeeded after enabling network access for the bounded source collection command. Local source archiving preserved original Markdown as well as normalized reading copies.

### What did not work

- `cat .docmgr.yaml`: `No such file or directory`. Repository configuration is `.ttmp.yaml`; used the installed docmgr configuration successfully.
- Initial `defuddle parse` calls failed with `Error loading content: getaddrinfo ENOTFOUND docs.slack.dev` and `getaddrinfo ENOTFOUND github.com`. Re-ran the collector with approved network access; all 14 sources succeeded.
- Large combined code/note reads exceeded output budgets; used targeted symbol searches and section reads for review-critical details.

### What I learned

Discovery and descriptors are coupled to Discord, not only transport. Current runtime config uses a host-field blacklist, so the Slack design uses a declared-field allowlist. Inspection executes scripts. Adjacent go-go-goja source and vault APIs may differ from the pinned dependency.

### What was tricky to design

Acknowledging before arbitrary JS work is necessary, but in-memory admission is not durable. Modal validation cannot use the same unconditional early ACK flow. Promise settlement must return to the VM owner while preserving invocation cancellation.

### What warrants a second pair of eyes

Review ingress dedupe/reservation ordering, owner-thread promise settlement, credential projection, and the selected SDK's retry/ACK behavior before implementation.

### What should be done in the future

Implement phases 0–4 before expanding to buttons, modals, xgoja integration, or persistence. Obtain the credentials and explicit test-message authorization described in guide section 12.

### Code review instructions

Start with the guide and source inventory. Follow the current file map and compare it to the archived revision. Use `docmgr doctor --ticket DISCORD-SLACK-001 --stale-after 30` for documentation hygiene. Runtime tests are future implementation gates, not claimed research results.

### Technical details

Source scripts are `scripts/01-collect-sources.py` and `scripts/02-archive-local.py`. Remote URL/time/hash records are in `sources/manifest.json`; vault and repository paths/hashes are in `sources/local-manifest.json`. No environment credentials were read and no Slack messages were sent.

## Step 2: Document validation and print layout

The guide contains approximately 6,500 words and the print edition has 21 pages. The first render exposed overlapping long code paths in table columns and page breaks inside diagrams. A generated print derivative converts tables to labeled records and keeps fenced examples together; the canonical guide retains its Markdown tables.

- `docmgr doctor --ticket DISCORD-SLACK-001 --stale-after 30`: all checks passed.
- Source SHA-256 verification, canonical guide relative links, and fenced-block balance passed.
- Initial print-wrapper attempt failed with `pandoc failed: Error producing PDF. ! Missing $ inserted.` at `event("app_`. Raw LaTeX environments consumed Markdown code; corrected by isolating wrappers in Pandoc raw-LaTeX fences.
- Final `remarquee upload bundle ... --pdf-only`: `OK: generated /tmp/slack-ticket-pdf-final/DISCORD-SLACK-001 Intern Implementation Guide.pdf`.
- Inspected rendered pages containing both architecture diagrams, file references, and test records. The overlaps and split diagrams were resolved.
- Saved the reviewed PDF to `artifacts/slack-intern-guide.pdf`. Print derivation script: `scripts/03-prepare-print.py`.
- Code snapshots use `.txt` extensions so the archive cannot create Go packages or nested modules. Added two adjacent go-go-goja runtime snapshots, each with revision provenance.
- No application code changed; Go build/test results are not claimed. Runtime validation belongs to the implementation tasks.

## Step 3: reMarkable delivery

Uploaded the completed guide using remarquee to `/ai/2026/09/10/DISCORD-SLACK-001`. The command returned `OK: uploaded DISCORD-SLACK-001 Intern Implementation Guide.pdf -> /ai/2026/09/10/DISCORD-SLACK-001`. A transient `remote tree has changed, refresh the file tree` warning preceded the successful upload. This is evidence of cloud delivery, not a physical-device synchronization check.

The final doctor pass initially reported `missing_numeric_prefix` for the generated print Markdown. Renamed that derivative with a `01-` prefix and updated the generator. The delivered content is unchanged. Implementation remains an unchecked future task.

## Step 4: Offline baseline and native Slack runtime

Started the implementation using the pinned, cached dependencies and explicit offline Go settings. Committed the completed research package first as `1ee4d4a` so implementation changes have a clean review boundary. The new code separates Slack domain contracts, VM glue, and public host composition.

### Prompt Context

**User prompt (verbatim):** "Ok, work offline for now. commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)"

**Assistant interpretation:** Implement the offline phases with fake services, incremental commits, and investigation-format diary entries. No Slack credentials, online SDK fetch, messages, or new reMarkable upload.

**Inferred user intent:** Make useful implementation progress before live workspace access is ready.

**Commit (prior milestone):** `1ee4d4a` — docs: design independent Slack bot host and archive research.

### What I did

- Read the diary, Glazed command-authoring, and Goja module-authoring skills and the pinned module source.
- Ran the full existing test suite with downloads disabled.
- Added Slack-specific domain values and configuration projection in `pkg/slackbot`.
- Added an explicitly restricted native module, descriptor loading, owner-thread calls, async fake-service operations, invocation deadlines, and reply-state enforcement in `internal/jsslack`.
- Added public composition in `pkg/slackhost` and the ping example.
- Added tests exercising real JS registration/dispatch with fake message and response services.

### Why

The user approved an independent implementation. The existing engine provides the owner and lifecycle primitives, while Slack response semantics and the domain services remain separate from Discord. The public composition package avoids the import cycle identified in the guide.

### What worked

The unchanged repository passed `go test ./...` using the cached Go 1.26.4 toolchain, an explicit matching GOROOT, GOWORK=off, and GOPROXY/GOSUMDB=off. The reproduction command is in `scripts/04-go-offline.sh`.

### What didn't work

- Initial offline `go test ./...` failed with `go: golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64: verifying module: checksum database disabled by GOSUMDB=off`.
- Invoking the cached Go executable alone exposed a conflicting toolchain root: `compile: version "go1.25.5" does not match go tool version "go1.26.4"`.
- Explicitly setting the matching cached GOROOT fixed the toolchain issue without a network request or dependency change.

### What I learned

The pinned engine supports runtime-aware module registration and owner scheduling, but cancellation of an owner call alone does not interrupt a CPU-bound handler. Each active JS entry needs a cancellation watcher that interrupts the VM and clears that interrupt before the next entry.

### What was tricky to build

The host must not hold the VM owner while waiting for an outbound request. Requests decode into Go values on the owner, run in a bounded errgroup, and return to the owner for promise settlement. Complete invocations are serialized while the owner remains available for promise completion.

### What warrants a second pair of eyes

Cancellation/interrupt ordering, retained contexts, unknown service errors, and whether a returned response can accidentally duplicate an explicit reply. Async operations require services that honor context cancellation.

### What should be done in the future

Finish offline discovery, CLI simulation, manifest output, ingress admission tests, and user-facing API documentation. Live Socket Mode remains pending network/SDK and workspace validation.

### Code review instructions

Read `pkg/slackbot/model.go`, then `internal/jsslack/host.go`, `module.go`, and `dispatch.go`. Run the focused package tests with `scripts/04-go-offline.sh`; review integration tests for actual side effects and deadline behavior.

### Technical details

No new module or dependency was introduced. The research source snapshot files retain `.txt` suffixes and do not participate in Go compilation. Script config is an allowlist projection, and no token fields are exposed to JavaScript.

Step 4 validation: the host/domain tests passed under the race detector, covering actual fake-service replies and posts, retained contexts, detached workspace-scoped storage, and CPU/network cancellation. A follow-up split between inspection and execution initially introduced `load redeclared in this block` because a test helper had that name. Renamed the helper to `loadTestHost`; inspection can now describe required config without requiring its value. Runtime loading still validates required fields.

## Step 5: Offline CLI and manifest generation

Committed the native runtime milestone as `e21a06d`, then implemented a separate `slack-bot` command with list, inspect, manifest, and simulate operations. Simulation executes the actual example JS bot against a recording service and emits the resulting operations as JSON.

### Prompt Context

See Step 4 for the implementation request. The current scope is offline work with incremental commits and detailed diary entries.

**Commit (previous milestone):** `e21a06d` — feat(slack): add offline JavaScript host and typed services.

### What I did

- Added `pkg/slackcli` discovery with explicit root-file / immediate-child-index conventions.
- Added JSON manifest generation and fixture-based simulation with `--event-file` and `--bot-config-file`.
- Added `cmd/slack-bot` with zerolog, `--log-level`, cancellation, and application-owned errors.
- Added command and mention fixtures and CLI tests that execute Cobra commands and inspect actual JSON output.

### Why

An offline end-to-end command makes the host useful immediately and gives later transport work a stable behavioral reference. The CLI does not claim that simulation proves Slack wire compatibility or that a live runner exists.

### What worked

`04-go-offline.sh test ./cmd/slack-bot ./pkg/slackcli ./pkg/slackbot` passed. List, inspect, manifest, command replay, and threaded mention replay produce the expected output; missing names and invalid deadlines return errors.

### What didn't work

- The first milestone staging attempt failed with `Unable to create .../.git/worktrees/discord-bot/index.lock: Read-only file system`. The worktree's Git metadata lives outside the workspace. The authorized commit succeeded with filesystem escalation; no network was needed.
- CLI construction initially failed with `Flag 'config-file' ... already exists`; the pinned framework reserves that name. Renamed the bot input to `bot-config-file`.
- The pinned high-level Glazed builder writes to `os.Stdout` and invokes `cobra.CheckErr` on domain failures. The missing-bot test therefore exited the test process with `Error: bot "missing" not found` rather than returning an error.

### What I learned

The installed command-authoring skill describes newer Glazed behavior than this repository's v1.3.6. The public parser API is sufficient to keep Glazed schemas and parsing while giving the new application its own `RunE` and output writer. No dependency upgrade or compatibility wrapper was needed.

### What was tricky to build

Inspection must not require runtime config values. It loads registration with a deadline, obtains metadata, and closes the runtime. Simulation then loads the selected bot with configuration validation and injected services. The Glazed source chain explicitly contains flags, arguments, and defaults; it does not load environment credentials or config files implicitly.

### What warrants a second pair of eyes

Review command-name conventions, discovery's exclusion of nested helpers, and the clear distinction between normalized invocation fixtures and Slack wire payloads. The JSON manifest still needs validation against Slack when live setup is authorized.

### What should be done in the future

Add embedded API help, TypeScript declarations, full offline validation, and the bounded ingress test seam. A separate side-conversation testing plan appeared as untracked `design-doc/02-full-local-testing-plan-and-slack-mock-evaluation.md`; it is left untouched and unstaged.

### Code review instructions

Read `pkg/slackcli/commands.go` and `discover.go`, then `cmd/slack-bot/main_test.go`. Run the example fixture commands after the tests and verify that stdout contains JSON operations only.

### Technical details

The parser is constructed with `cli.NewCobraParserFromSections` and mounted with `AddToCobraCommand`; `cli.NewCobraCommandFromCommandDescription` supplies the command shape. Command settings are disabled, and explicit Glazed source middleware handles flags/arguments/defaults.

## Step 6: Bounded ingress, integration coverage, and offline handoff

Extended the offline implementation with bounded admission, event deduplication, filtering, overload receipts, and host shutdown. The ingress integration test drives an actual JavaScript mention handler through the public host and verifies the recorded thread destination while canceling the independent ACK context.

Added an embedded help page, TypeScript declarations, README entry points, fixture validation, and generated package logging registration. The offline phase is usable and validated; the broader MVP task remains open because SDK transport, real wire-protocol tests, rate limiting, and reconnect support have not been implemented.

### Prompt Context

See Step 4. The user authorized offline implementation and commits, not live Slack communication or a new upload.

**Commit (previous milestone):** `6cb0b79` — feat(slack): add offline discovery manifest and fixture replay CLI.

### What I did

- Added `pkg/slackbot/ingress.go`: workspace/app/channel/self-bot filtering, queue bounds, event/envelope identity, dedupe TTL/capacity, receipt dispatch, metrics, and cancellation.
- Added tests with a blocked dispatcher proving ACK independence, duplicate suppression across delivery IDs, explicit overload, TTL expiry, filtering, and closure.
- Added `pkg/slackhost/host_test.go` to connect ingress to the real example JS handler and recording service.
- Added `pkg/slackdoc` with embedded `help slack-offline`, a TypeScript module declaration, and README commands.
- Tightened fixture size checks and fixed null/scalar values in the detached memory store.
- Generated logcopter registration for new library packages using the pinned repository tool.
- Built the Glazed analyzer from the selected v1.3.6 module into `/tmp/slack-glazed-lint` and ran it on the new CLI.

### Why

The in-memory ingress seam establishes lifecycle and admission behavior before a Socket Mode SDK is connected. An end-to-end host test distinguishes “our interfaces compile” from “the actual JS bot produces the intended effect.” Help explicitly separates normalized fixture replay from Slack wire compatibility.

### What worked

- Full `04-go-offline.sh test ./...`: passed, including existing Discord packages.
- `04-go-offline.sh build -buildvcs=false ./...`: passed.
- `04-go-offline.sh vet ./...`: passed.
- Race tests for new host, runtime, domain, and CLI packages: passed after the shutdown behavior fix described below.
- Focused `golangci-lint run` on all new Slack packages: `0 issues.`
- Pinned Glazed analyzer on `pkg/slackcli` and `cmd/slack-bot`: passed after expressing log-level as a Glazed field.
- Actual `go run` command replay emitted one `ephemeral_reply` with `pong`; mention replay emitted one post with `threadTs: 1741234567.000001` and an explicit fake reference.
- `go run ... help slack-offline` rendered the embedded guide.
- `docmgr doctor` and `git diff --check`: passed.

### What didn't work

- Build/VCS discovery initially failed with `error obtaining VCS status: exit status 128`. This sandbox cannot reliably inspect the external worktree metadata from the Go builder. `-buildvcs=false` resolves build-only stamping without changing source or Git configuration. Lint package loading needed the same `GOFLAGS=-buildvcs=false` setting.
- Initial focused lint found unchecked deferred `Close` results and an implicit fulfilled-promise default branch. Cleanup discards are now explicit and the promise switch names `PromiseStateFulfilled`.
- The Glazed analyzer rejected a raw root `StringVar` flag: `define CLI flags with cmds.WithFlags(fields.New(...)) instead of raw Cobra/pflag/flag APIs`. Moved `log-level` into each Glazed command's settings and removed the custom shared level writer.
- The new integration test initially observed `context canceled` on the ingress error channel during normal shutdown, after the message had already been recorded. The worker now suppresses dispatch failures once its own lifetime is canceled; unexpected active-lifetime failures remain observable.
- No `tsc` executable was installed at the checked workstation path. The declaration was reviewed alongside the JS API, but a TypeScript compiler run is not claimed.

### What I learned

A successful side effect does not mean the handler has fully settled yet. Shutdown may legitimately cancel the last promise-completion step. That cancellation should not be reported as an operational ingress failure. A bounded dedupe cache should reject admission when full rather than evict a live key and silently weaken duplicate suppression.

### What was tricky to build

There are three independently meaningful contexts: ACK receipt, admitted work, and runtime lifetime. The integration test cancels ACK context immediately and still obtains the JS response, then shuts down the ingress before the host. The queue's capacity and dedupe capacity are separate bounds; neither may grow to absorb overload indefinitely.

### What warrants a second pair of eyes

The ingress consumes normalized snapshots, not Slack wire payloads. Its ACK interface must be wired to the chosen SDK and tested against HTTP/WebSocket fixtures before claiming protocol support. Injected Go services must honor cancellation. Process-local storage/admission are not durable. Critical scripts should await their side effects.

### What should be done in the future

Follow the separate local-testing plan for SDK selection, exact HTTP/WebSocket fixtures, and a stateful Slack mock server. Implement token-file configuration and Socket Mode only with a fake-compatible transport and no implicit network access. Live app setup and actual workspace checks remain separately required. No new reMarkable upload was requested.

### Code review instructions

Start at `pkg/slackdoc/slack-offline.md` for the shipped contract, then `pkg/slackbot/ingress.go` and `pkg/slackhost/host_test.go`. Reproduce checks with `scripts/04-go-offline.sh`; add `-buildvcs=false` to build/run if the sandbox blocks VCS metadata. Review the independent testing-plan document without staging or overwriting it as part of this implementation.

### Technical details

Dependency downloads remained disabled with GOPROXY/GOSUMDB=off and GOTOOLCHAIN=local. Tests used cached Go 1.26.4, the pinned module graph, local fixtures and fake services. No token values, environment credentials, or real Slack messages were accessed. Vulnerability-database refresh and dependency-installing Make targets were not run because this phase is offline; existing cached lint tooling and a locally built version-matched analyzer were used instead.

## Step 7: Publish the implementation deep dive to the Obsidian vault

The user redirected the next phase from local transport experiments to a detailed project report. Wrote a 5,211-word technical article explaining the completed offline implementation at `fd84234`, its concurrency and delivery contracts, the evidence behind the tests, and the separate planned transport verification. No SDK/mock compatibility experiment was completed or claimed during this reporting step.

Published a new append-only vault note at `Projects/2026/09/10/PROJ - Slack Bot - Runtime Ownership Admission and Local Verification.md`. The go-go-parc commit is `80b1f8f` and the push to `origin/main` succeeded. The existing research guide's reMarkable delivery remains a separate historical deliverable.

### Prompt Context

**User prompt (verbatim):**

> Actually  write a detailed project report for the obsidian vault as a deep dive technical analysis blog post using a textbook writing style (no analogies, see skill).
> Commit and push the bsidian vault when done (go-go-parc vault).

**Assistant interpretation:** Replace the immediate transport implementation step with an evidence-backed vault article, then commit and push that article. Preserve the earlier detailed diary requirement.

**Inferred user intent:** Produce a durable technical explanation of the project and its current engineering state for future readers.

**Commit (implementation baseline):** `fd84234` — feat(slack): validate bounded ingress and document offline workflow.

**Commit (vault publication):** `80b1f8f` — docs: publish Slack bot runtime technical deep dive.

### What I did

- Applied the Obsidian vault-writing and textbook-authoring skills, including direct technical explanations without analogies.
- Read the diary, runtime, domain, ingress, CLI, example, tests, pinned dependencies, and local-testing plan.
- Wrote the report with three Mermaid diagrams, real code, algorithm pseudocode, package references, a captured simulation result, and related vault links.
- Distinguished implemented behavior from planned SDK, wire-fixture, stateful mock, and live workspace verification.
- Validated YAML frontmatter, fenced blocks, referenced local source paths, and all four vault wikilinks.
- Staged only the new article in a previously clean vault and pushed its focused commit.

### Why

The report should teach ownership, cancellation, admission, and response semantics through the actual implementation. A project narrative must also preserve the limits of its evidence so that an intern does not treat normalized replay as successful Socket Mode integration.

### What worked

- `04-go-offline.sh test ./pkg/slackbot ./pkg/slackhost ./internal/jsslack ./pkg/slackcli ./cmd/slack-bot` passed with cached test results.
- `04-go-offline.sh run -buildvcs=false ./cmd/slack-bot bots simulate ping --event-file examples/slack-bots/fixtures/mention.json` freshly produced one post to `C-TEST`, preserving `threadTs: 1741234567.000001` with fake reference `offline.000001`.
- The report's YAML, fence parity, source paths, and wikilinks validated. Vault staged diff whitespace checks passed.
- `git push origin main` returned `9e0ebe7..80b1f8f main -> main` for go-go-parc.

### What didn't work

No publication or validation failures occurred. Some combined source reads were truncated by output limits; targeted follow-up reads recovered the specific code needed for claims. No new full lint, race, or TypeScript compiler run is claimed for this documentation-only step.

### What I learned

Ingress ACK independence is weaker than strict ACK-before-handler-start ordering: queue insertion precedes the acknowledgment call, so the worker can begin first. The report explicitly describes the actual guarantee. Native workers reject promises but return nil to the errgroup, so waiting for workers alone does not surface all unawaited service failures.

### What was tricky to build

The report needed to distinguish VM serialization from whole-invocation serialization, receipt context from worker lifetime, and reply-slot consumption from confirmed remote delivery. It also distinguishes the separate untracked test plan from the committed application revision rather than attributing the plan to `fd84234`.

### What warrants a second pair of eyes

Review the ACK ordering explanation, unawaited-operation policy, typed error trust boundary, store capacity limitations, and conservative reply-slot consumption before implementing the transport. The article identifies these constraints without modifying application behavior.

### What should be done in the future

Resume the local testing plan with a pinned Go SDK/mock interoperability probe, then precise wire fixtures and the full process harness. The report lists this as future work. No real Slack credentials are needed for the implemented tests or planned local synthetic-credential harness.

### Code review instructions

Read the new vault article at commit `80b1f8f`, then compare its key explanations against `internal/jsslack/dispatch.go`, `internal/jsslack/host.go`, and `pkg/slackbot/ingress.go`. Use the diary's earlier validation entries for the complete offline milestone checks; this step records its narrower fresh verification separately.

### Technical details

Vault root: `/home/manuel/code/wesen/go-go-golems/go-go-parc`. The article is 419 lines and 5,211 whitespace-delimited words, with three native Mermaid diagrams and no external image dependencies. Only the article was included in the vault commit. Filesystem escalation enabled the authorized vault write/commit; network escalation enabled the explicitly requested push. No environment credentials were inspected and no Slack traffic was sent.

## Step 8: Prove actual Go SDK and local Socket Mode interoperability

Resumed offline implementation with the test plan's first experiment. Pinned cached slack-go v0.17.3 and its WebSocket dependency, wrote an explicit external-server test, and started the pinned Bun mock in tmux. The real SDK authenticated, received two Socket Mode deliveries, acknowledged them, posted a threaded message, and sent an ephemeral command response. The mock independently verified conversation state and socket closure.

The probe is intentionally smaller than a complete application process test. It establishes SDK/mock interoperability without claiming that the production transport or JavaScript composition has been connected. Sources and sanitized receipts are archived in the ticket; synthetic credentials remain in the explicit /tmp run directory.

### Prompt Context

**User prompt (verbatim):** "Ok, continue.commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)"

**Assistant interpretation:** Resume the local testing and implementation plan, preserving offline execution, focused commits, and detailed investigation entries.

**Inferred user intent:** Progress from the validated offline runtime toward a working Slack transport using local evidence before real workspace access.

**Commit (previous checkpoint):** `7cc01a0` — docs: record Slack deep dive vault publication.

### What I did

- Added cached `github.com/slack-go/slack v0.17.3`; upgraded its required `gorilla/websocket` to v1.5.3 without downloads.
- Added `internal/slackprobe/probe_test.go`, a real SDK test with explicit fixture-directory selection and literal-loopback-only HTTP/WebSocket dialers.
- Added `testdata/slack/mock/probe.ts`, pinned revisions, and a reproduction runbook.
- Ran Bun 1.2.13 against mock commit `6397b31a9c5e04a3ba52dd6e82f16ab8b5b10eac` in tmux.
- Archived SDK/mock source snapshots with hashes under `sources/sdk-probe` and sanitized results under `artifacts/sdk-probe-001`.

### Why

The mock's upstream tests mainly establish JavaScript/Bolt compatibility. A real Go SDK experiment is necessary before depending on this server for the process gate. Independent mock observations prevent successful local SDK calls from being mistaken for correct thread, recipient, or ACK behavior.

### What worked

`04-go-offline.sh test -count=1 -v ./internal/slackprobe -run TestSDKInteroperability -args -slack-probe-directory /tmp/slack-sdk-probe-001` passed in 0.178 seconds. The mock recorded exactly two deliveries, each acknowledged on its first attempt. Thread identity was preserved, the private response belonged to Alice, no public private-response message existed, the slash ACK was empty, and cancellation reduced connection count to zero.

The actual server used loopback port 45255. It stopped itself after cancellation; `lsof-who -p 45255 -k` reported `no process listening/bound on port 45255`, then the owned tmux session was removed. No outside endpoints were dialed.

### What didn't work

One exploratory source read used nonexistent `socketmode/event_type.go` and returned `No such file or directory`; the probe uses the verified request type from `request.go` and did not require that guessed filename. No compile, execution, or compatibility failures occurred in the first probe.

### What I learned

The cached SDK supports endpoint injection, an explicit HTTP client, a WebSocket dialer, and context-aware ACK scheduling. The mock supports these actual SDK calls without source changes. ACK scheduling success still needs mock/wire observation to establish receipt, which this probe checks through the mock delivery ledger.

### What was tricky to build

All three network paths must obey offline constraints: Web API, returned WebSocket URL, and response URL. The probe uses a no-proxy HTTP transport, rejects redirects, and rejects hostnames and non-loopback literal addresses before dialing. Readiness/results use atomic file rename to avoid partial JSON observations. Cancellation is verified by a separate server-side shutdown receipt.

### What warrants a second pair of eyes

The selected mock is permissive about some token usage, so this test alone cannot prove correct token routing. Strict HTTP fixtures remain required. Ordinary repository tests explicitly skip the external-server gate unless a directory is supplied; CI must run the gate separately and must not count the default skip as coverage.

### What should be done in the future

Implement the actual transport using these endpoint seams, strict request fixtures, normalized ingress, and real JS host composition. Add rate-limit, uncertain-delivery, reconnect, ACK failure, and process lifecycle cases per the plan.

### Code review instructions

Read `testdata/slack/mock/README.md`, then compare the SDK test with the independent TS assertions. Review the two sanitized receipts and pinned revisions. Reproduce with a fresh run directory and the prepared dependencies; do not archive config.json.

### Technical details

SDK v0.17.3, WebSocket v1.5.3, Go 1.26.4, Bun 1.2.13, mock package v0.4.0. The test has a 30-second operation deadline and the launcher a 45-second watchdog. No production secrets or environment credentials are read. Source archives are original code snapshots with .txt suffixes; no new HTML/Markdown web sources were fetched.

Step 8 final validation: full offline `go test ./...` passed after the SDK dependency change, including existing Discord packages. The previously separate local-testing plan is included unchanged in this milestone so the implementation gate and its source design are reviewable together.
