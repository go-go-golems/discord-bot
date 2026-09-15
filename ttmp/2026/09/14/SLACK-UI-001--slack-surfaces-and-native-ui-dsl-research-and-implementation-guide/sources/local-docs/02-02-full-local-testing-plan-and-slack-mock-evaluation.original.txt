---
Title: Full local testing plan and Slack mock evaluation
Ticket: DISCORD-SLACK-001
Status: active
Topics:
    - discord-bot
    - architecture
    - api-design
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/code/others/slack/desplega-slack-mock/src/server.ts
      Note: Stateful mock API and observed fidelity gaps
    - Path: abs:///home/manuel/code/others/slack/desplega-slack-mock/src/socket-mode.ts
      Note: Envelope acknowledgment and retry model
    - Path: abs:///home/manuel/code/others/slack/slack-go/slacktest/server.go
      Note: Go HTTP test-server seams and constructor overrides
    - Path: abs:///home/manuel/code/others/slack/slack-go/socketmode/socket_mode_managed_conn_test.go
      Note: Custom Go WebSocket fixtures
    - Path: repo://ttmp/2026/09/10/DISCORD-SLACK-001--add-slack-support-to-discord-bot/design-doc/01-slack-support-architecture-and-intern-implementation-guide.md
      Note: Architecture and MVP requirements
ExternalSources: []
Summary: Source-backed mock evaluation and complete offline test architecture, scenarios, harness, CI, and acceptance gates.
LastUpdated: 2026-09-10T18:46:18.691557971-04:00
WhatFor: Implement credential-free local testing of the Slack host.
WhenToUse: Before building Slack test infrastructure or validating transport integration.
---

# Full local testing plan and Slack mock evaluation

## 1. Decision and scope

Use three complementary testing mechanisms: Go fakes for bot behavior and runtime ownership, small Go HTTP/WebSocket fixtures for exact protocol failures, and the existing `desplega-ai/slack-mock` server for stateful end-to-end scenarios. The third component is a real local server: the actual Go Slack client can connect through HTTP and Socket Mode while the test driver acts as a human user. Its compatibility with our selected Go SDK version must be demonstrated by an initial probe; the upstream integration tests primarily exercise Bolt for JavaScript.

This document supplements the [intern implementation guide](01-slack-support-architecture-and-intern-implementation-guide.md). It specifies a test implementation, not a completed harness. The work in this side conversation was source inspection and documentation only. No mock servers, dependency installation scripts, upstream test suites, production bots, or Slack messages were run. The main-thread implementation continues independently.

The MVP under test is a single-workspace Slack host, one JavaScript bot per process, `/golem-ping`, mention events, threaded replies, offline discovery, explicit configuration, and clean lifecycle handling. Later sections define tests for buttons, modals, direct messages, and persistence so that these features have acceptance gates before they are added. They are not requirements to implement those features now.

The practical goal is to let an intern answer four questions without Slack credentials:

1. Does the JavaScript API produce the correct typed operations?
2. Does the transport encode, route, acknowledge, retry, and cancel those operations correctly?
3. Does the complete process produce the expected visible conversation state?
4. When a test fails, can its artifacts identify the failing boundary?

## 2. Research baseline and checkout locations

The following repositories were cloned, without installing dependencies or modifying their source, under `/home/manuel/code/others/slack/` on 2026-09-10. Commit IDs, rather than branch names, identify the evidence behind this recommendation.

| Checkout directory | Project | Inspected commit |
|---|---|---|
| `slack-go` | `slack-go/slack` | `4b9cec9bdb1eaca3b66968b02c9f08fecd860285` |
| `desplega-slack-mock` | `desplega-ai/slack-mock` | `6397b31a9c5e04a3ba52dd6e82f16ab8b5b10eac` |
| `slack-testing-library` | `chrishutchinson/slack-testing-library` | `3f450bb5abcf00991a5ebe8cbeef98a60e1d4530` |
| `slack-mock` | `Skellington-Closet/slack-mock` | `7bf3eb5073a8a2c071ece360a7b8729390a37024` |

The application HEAD observed during this analysis was `1ee4d4ac535e1b69971f5bf0022a3b20041299f1`. Uncommitted files under `pkg/slackbot/` and `internal/jsslack/` belonged to the active main thread. Their names and contracts informed the test plan, but this document does not assert that their current contents are final, committed, or tested. Reconcile test signatures with the implementation at the start of each milestone.

## 3. Evaluation of existing mocks

### 3.1 Recommended end-to-end server: desplega-ai/slack-mock

The checked-out package declares version 0.4.0 and Bun >=1.2. It has no declared production dependencies; its development dependencies include Bolt, TypeScript, and Biome. Its server uses `Bun.serve`, so the project should be treated as a Bun application rather than assumed runnable under ordinary Node.js. Pin the source revision and runtime version before making it a reproducible test dependency.

The implementation supports the exact broad path we need:

```text
Go Slack client -- HTTP apps.connections.open --> mock
Go Slack client <-------- WebSocket hello ------- mock
Go Slack client <------ Socket Mode envelope ---- mock
Go Slack client -------- envelope ACK ----------> mock
Go Slack client ------- chat.postMessage -------> mock
                                                   |
                                                   v
                                      stored messages and threads
                                                   |
                                                   v
                                    test assertions and browser UI
```

Concrete source evidence:

- `src/server.ts:243`, `listen`: loopback by default and dynamically allocated ports through `port: 0`.
- `src/web-api.ts:119`: `apps.connections.open` issues a mock WebSocket ticket.
- `src/socket-mode.ts:104`, `open`: sends `hello`, connection metadata, and periodic WebSocket ping frames.
- `src/socket-mode.ts:125`, `message`: matches envelope ACKs to pending deliveries and records ACK payloads.
- `src/socket-mode.ts:170`, `send`, and `:195`, `attempt`: emits envelopes and redelivers with an incremented retry count when receipt is missing.
- `src/server.ts:397`, `handleApi`: authentication checks, method dispatch, API-call records, and injected error responses.
- `src/server.ts:941`, `slashCommand`: constructs a command payload, supplies a response URL, and waits for the envelope ACK.
- `src/server.ts:1032`, `submitView`: applies modal ACK responses such as errors, update, push, and clear.
- `src/server.ts:1271`, `handleAdmin`: HTTP endpoints for driving a separate application process.
- `test/smoke.test.ts`: a real Bolt app receives a mention and creates a threaded reply.
- `test/socket-mode.test.ts`: source examples for redelivery, reconnect, and missing connections.
- `test/faults.test.ts`: error and 429 injection exercised with Slack's JavaScript WebClient.

These findings correct the earlier search-based uncertainty: the repository is available and Socket Mode is implemented. They do not establish that the Go client passes the same scenarios; that is the first proposed experiment.

### 3.2 Important limits of that server

The mock is useful precisely because it implements a bounded Slack model. A passing test proves conformance to that model, not every behavior of Slack itself.

- **OAuth scopes are not fully enforced.** `actorFor` accepts the configured bot or app token as a bot actor, and `handleApi` specifically restricts `apps.connections.open` to the app token. It does not generally reject the app token for ordinary bot Web API calls. Verify correct token selection separately with strict Go request fixtures.
- **Manifest import is JSON.** `loadManifest` uses `JSON.parse(readFileSync(...))`. The YAML manifest in the original design must be converted to JSON by the harness, or represented directly as an object. Passing YAML to this CLI is not supported by the inspected code.
- **The default event set is broad.** Supply the actual MVP subscriptions rather than accepting `DEFAULT_EVENTS`; otherwise the mock may emit traffic a real installation would not receive.
- **Retries preserve `envelope_id`.** The retry loop changes attempt metadata on the same envelope. Test identical Events API `event_id` under a different envelope using `mock.hub.send` or the Go wire fixture; the default timeout retry alone is not enough to verify business-event deduplication.
- **ACK recording lacks ACK timestamps.** `DeliveryRecord` stores `sentAt`, attempt count, ACK status, and ACK payload. It does not expose a monotonic `ackedAt`. Use a Go wire fixture to assert exact ACK ordering and measure receipt latency.
- **Fault injection is application-level.** The public `Fault` shape supplies method, error code, status, retry delay, count, and extras. It does not implement “commit a message and then lose the HTTP response,” arbitrary malformed frames, or a controlled blocked request. Those belong in Go fixtures.
- **Some controls are TypeScript-only.** `injectFault` and `deliveries` are public methods, but the inspected admin router has no `/mock/faults`, `/mock/deliveries`, or `/mock/reset` endpoints. Do not invent those URLs in scripts.
- **`flush()` is a convenience, not completion proof.** It waits for tracked deliveries and a period of HTTP quiet; delivery failures are settled internally. Its `maxMs` check occurs after awaiting inflight promises, so it is not a hard outer timeout. Always assert ACK state and business results explicitly, with a test-owned deadline.
- **Message waits can match old state.** `waitForMessage` and `waitForApiCall` search existing records first. Use unique scenario markers, references, and per-test instances or record offsets.
- **Timestamp comparisons inside waits use Number.** `waitForMessage(..., {after})` converts timestamps to numbers. Exact identifier-preservation tests should inspect captured strings rather than trusting numeric ordering in the mock.
- **UI rendering is illustrative.** Its HTML renderer is suitable for reviewing message shape and thread flow. It does not prove pixel parity, accessibility behavior, or layout limits in Slack's official clients.

Each of these limits has a test placement below. None requires changing the production bot to imitate an emulator quirk.

### 3.3 Useful Go building block: slack-go/slacktest

The current checkout is more capable than its RTM-oriented exported names imply. `slacktest/server.go:37` constructs an `httptest` server and registers `/auth.test`, `/chat.postMessage`, and `/apps.connections.open` among other methods. `SendToWebsocket` can send arbitrary JSON, and constructor binders can install a custom WebSocket handler before the defaults are registered.

However, the default `slacktest/handlers.go:371` WebSocket loop parses RTM event types and responds to RTM ping messages. It does not supply the stateful Socket Mode delivery/ACK/retry ledger that this plan needs. Default message handling also uses timestamp seconds and parses form bodies. It should not be treated as a comprehensive modern Slack emulator.

The useful reuse point is demonstrated by the SDK itself: `socketmode/socket_mode_managed_conn_test.go:154` constructs `slacktest.NewTestServer` with a custom `/ws` handler. Register overrides through that constructor: `Server.Handle` silently ignores a pattern that has already been registered. An override attempted after construction may therefore leave the default behavior in place.

Use this package if it makes the small wire fixture easier. Plain `httptest.NewServer` plus the already-selected WebSocket library is also sufficient. There is no need to create a new generic mock framework or adapt Discord interfaces.

### 3.4 Not selected: slack-testing-library

This TypeScript library runs a local HTTP interception server and has an attractive user-oriented vocabulary such as `openHome`, `mentionApp`, and `interactWith`. Source inspection shows a narrower fit than the README summary suggests:

- `src/util/server.ts:26` implements HTTP request handling, not Socket Mode.
- The server parses request bodies as query strings and defaults to `{ok: true}` for unhandled calls.
- `src/slack-testing-library.ts:377` explicitly rejects interacting with buttons when the active screen is a channel.
- Several observation paths use delayed retry polling.

Its checkout's last commit is dated 2021-09-25. That date alone is not a disqualifier; the missing Socket Mode path and channel-interaction restriction are the practical reasons not to choose it for this project. Retain its user-oriented scenario vocabulary as inspiration.

### 3.5 Not selected: Skellington-Closet/slack-mock

`src/mocker/web.js:13` uses `nock('https://slack.com')` to intercept requests inside a Node process. That interception does not capture HTTP requests made by a separate Go process. `src/index.js` exposes a singleton and RTM helpers; `package.json` declares older Node dependencies and a pretest script that runs a formatter with `--fix`.

The inspected last commit is dated 2018-03-16. This is useful historical prior art, but integrating it would require extra machinery while still leaving the Socket Mode requirement unmet. Do not run its default tests as a harmless read-only probe: its `pretest` mutates source formatting.

## 4. Testing architecture and ownership

The word “offline” means no real Slack account, credentials, or external application traffic during tests. Initial checkout and dependency preparation may use the network. Once dependencies are prepared, test execution should work with external networking unavailable.

Separate five layers. Most assertions live in the fastest layer that can establish the behavior, with a smaller cross-layer suite proving that the pieces are wired together.

| Layer | Executes | Replaces | Main evidence |
|---|---|---|---|
| L0: domain | Go validation and configuration | Network and JS | Exact values and typed errors |
| L1: JS host | Real Goja and real bot scripts | Message and response services | Captured operations, lifecycle barriers |
| L2: wire | Actual Go Slack SDK and our transport | Slack HTTP/WebSocket service | Wire frames, headers, request counts, ACK timing |
| L3: process | Actual CLI, transport, JS host | Slack using desplega mock | Conversation state, command receipts, process exit |
| L4: visual | L3 plus local browser | Slack client rendering | Screenshots paired with structural assertions |

```text
L0: input ----------------------> domain validation
L1: fixture script -> Goja ------> fake MessageService / Responder
L2: transport + real SDK --------> precise local HTTP/WS fixture
L3: actual slack-bot process ----> desplega SlackMock
                                      ^          |
                                      |          v
                               scenario driver  state / journal
L4: browser --------------------> local mock UI
```

Do not replace our transport in L3. A process test that injects `Invocation` directly into the JS host cannot establish that Socket Mode parsing or credentials were wired correctly. Conversely, L1 should not require Bun or a socket listener merely to test a rejected promise.

## 5. Required seams in the implementation

These are proposed requirements for the implementation owner, not changes made by this document.

### 5.1 Domain and JS seams

The observed main-thread files already introduce `slackbot.MessageService.Post`, `slackbot.Responder.Reply`, typed `Invocation`, and `jsslack.Options.Messages`. Tests should use those narrow contracts directly. A responder is per invocation because its response URL is a capability tied to that interaction; do not put one mutable response URL on a shared fake.

A useful fake records immutable copies of calls and offers explicit barriers:

```go
// Sketch: final names follow the implementation.
type RecordedPost struct {
    Input slackbot.PostMessage
    Sequence uint64
}
type BlockingMessages struct {
    Entered chan RecordedPost
    Release chan struct{}
    Result slackbot.MessageRef
    Err error
}
var _ slackbot.MessageService = (*BlockingMessages)(nil)

// Post sends an observation, then selects on Release or ctx.Done().
// It must remain context-aware even when the test fails early.
```

Use a buffered channel or bounded recorder so observation cannot accidentally block the behavior being tested. Return detached data; never pass Goja values through the fake service. Protect recorder state when several workers can write to it.

### 5.2 HTTP and WebSocket endpoint injection

The real Slack SDK needs a configurable API base URL and HTTP client. In the inspected SDK, `slack.OptionAPIURL` and `slack.OptionHTTPClient` supply these seams; `socketmode.OptionDialer` supplies the WebSocket dialer. Verify the same options in the version actually added to the application.

The HTTP client must cover Web API and response-URL requests. The WebSocket dialer must cover the URL returned by `apps.connections.open`. Replacing only the REST base URL is insufficient if another path can still dial an external address.

For test execution, supply a loopback-only transport/dialer with proxy lookup disabled and redirects checked at every hop. Validate the actual dial destination, not merely the configured base URL. A test should deliberately return an external WebSocket or response URL and assert rejection without a network attempt. Local TLS tests use a private test CA or the `httptest` client; do not disable certificate verification globally.

The standalone CLI needs explicit local endpoint configuration or a clearly defined test entry point using the real composition root. Prefer useful host flags such as `--api-base-url` and explicitly enabled local mode, while retaining normal production defaults. These flags are proposed; they are not claimed to exist yet. Credentials remain explicit token-file inputs. No new environment reads are necessary.

### 5.3 Time, metrics, and observations

Use a controllable clock for host-owned dedupe expiry, pacing, and retry scheduling. SDK-owned timers may need real-time tests; do not assume the SDK exposes a fake clock.

Record receive, admission, ACK-written, handler-started, handler-finished, API-attempted, and shutdown-finished events using monotonic time or a sequence counter in the test recorder. Production logs can expose the corresponding correlation IDs, but assertions should avoid parsing prose log messages where a structured event is available.

Expose queue capacity, worker limits, and timeouts as explicit settings. Tests should not use unsafe access to private fields to force these conditions. An observation hook must report behavior rather than change it.

## 6. Harness layout and reproducibility

Suggested new test paths, all inside the existing Go module:

```text
pkg/slackbot/model_test.go                 L0 rules
internal/jsslack/host_test.go              L1 script loading
internal/jsslack/dispatch_test.go          L1 behavior and ownership
internal/slacktransport/transport_test.go  L2 HTTP and lifecycle
internal/slacktransport/socket_test.go     L2 socket behavior
internal/slacktestfixture/                 narrow Go test helpers
testdata/slack/wire/                      exact protocol fixture JSON
testdata/slack/bots/                      scenario scripts
testdata/slack/mock/manifest.json         minimal mock app manifest
testdata/slack/mock/revisions.json        exact mock and runtime pins
testdata/slack/mock/scenarios.ts          L3 test driver
testdata/slack/mock/mock-runner.ts        local mock launcher
```

Do not create another `go.mod`. A small pinned Bun dependency manifest for the test driver is acceptable if that becomes the chosen dependency workflow; it must not introduce a second application runtime architecture. Alternatively, the first probe can import the already-cloned mock by a command-line path. CI must not depend on a developer's home directory.

`revisions.json` should record the application commit, selected Slack SDK version, mock commit, Bun version, and optional browser version. Synthetic fixture IDs remain stable (`TTEST`, `ATEST`, `CTEST`, `UTEST` where the parser permits them); IDs generated by the emulator are obtained from its state instead of guessed.

Use one fresh mock per independent scenario or scenario group. There is no documented reset route in this checkout. Reusing a server without clearing messages, pending deliveries, faults, and waits can make tests pass on yesterday's state. A new instance and separate journal file are simpler than inventing reset semantics.

## 7. First experiment: prove Go SDK interoperability

Before integrating the full CLI, write a bounded Go probe using the selected SDK. It should execute under a test or short-lived command, while the mock server is managed in tmux. Its purpose is to confirm the HTTP body, authentication, socket, and receipt contracts with the Go client.

1. Start the pinned mock on loopback with a minimal JSON manifest containing `/golem-ping` and only `app_mention`.
2. Read its base URL and synthetic tokens from the test-owned launcher result. Pass values directly or through temporary files; do not inspect the operator's environment.
3. Call `auth.test` with the bot token and assert the workspace/user identity matches the mock state.
4. Start Socket Mode with the app token; wait for a connection and `hello` observation with a bounded deadline.
5. Inject one mention from the mock's human interface.
6. Decode the received event, ACK the exact envelope, then call `chat.postMessage` with the received channel and thread root.
7. Assert a stored reply with the exact text and `thread_ts`, and an acknowledged delivery record.
8. Inject `/golem-ping`; ACK empty and post an ephemeral response to its response URL. Assert the recipient and absence of a public channel message.
9. Cancel the probe and assert its process/worker exit before stopping the mock.

**Pass gate:** actual Go SDK, no real credentials, one mention reply and one ephemeral command response, all expected ACKs, zero external requests, clean shutdown. A Bolt test passing upstream does not satisfy this gate. If the probe fails, record raw sanitized requests/frames and identify whether the problem is an emulator mismatch, SDK behavior, or our integration assumption before choosing a fix.

## 8. Domain and JavaScript test catalog

Each case should execute a real fixture script when the JavaScript boundary matters. Do not manufacture the handler return value directly in Go and call it a JS test.

### L0: pure rules

| ID | Input or operation | Required assertion |
|---|---|---|
| D01 | Empty/whitespace and maximum-length text | Explicit boundary policy; no off-by-one errors |
| D02 | Unicode, combining characters, emoji | Document whether limit counts runes; preserve content |
| D03 | Missing channel or ambiguous command/event | Stable `invalid_argument`, no service call |
| D04 | Fractional timestamp strings | Exact byte-for-byte identifier preservation |
| D05 | Required/defaulted bot config | Correct defaults, missing-required failure, type checks |
| D06 | Undeclared host/token fields mixed into input | Only declared bot fields appear in config |
| D07 | Invalid names, duplicates, unsupported events | Clear registration failure with name/path |
| D08 | Nested unknown message option | Reject before network; point to offending field |

The observed MVP text validator uses 1–4000 runes; treat that as an application policy to test, not a claim about every Slack payload limit.

### L1: runtime and behavior

| ID | Scenario | Required assertion |
|---|---|---|
| J01 | Load valid bot without services | Descriptor succeeds and no network call occurs |
| J02 | Throw during import or export wrong value | Helpful error and runtime cleanup |
| J03 | Infinite import loop | Bounded interruption; next test starts cleanly |
| J04 | Return `{text}` from command | One call to that invocation's responder |
| J05 | `await ctx.reply(...)` | One response, correct receipt shape |
| J06 | Explicit reply then return response | `already_replied`, exactly one side effect |
| J07 | Return undefined | No implicit response |
| J08 | Mention outside/inside existing thread | Use root `ts`/existing `threadTs` respectively |
| J09 | Explicit post | Correct `channelId`, text, reference, optional thread |
| J10 | Throw Error/reject Error/reject primitive | Useful message and available stack preserved |
| J11 | API pending while another owner callback settles | No owner deadlock; promise resolves |
| J12 | More async calls than worker limit | Deterministic overload; no unbounded goroutines |
| J13 | Cancel pending API, handler, and queued invocation | Context reaches service; bounded cleanup |
| J14 | CPU-bound handler times out, then next invocation | VM interrupt cleared; later handler succeeds |
| J15 | Retain prior `ctx` and use in later invocation | Closed context cannot post under later identity |
| J16 | Concurrent invocations from different users | Per-invocation responder/config/state isolation |
| J17 | Close twice and dispatch after close | Idempotence and explicit closed-host result |
| J18 | Mutate descriptor/config snapshot in JS/caller | No unintended mutation of future invocations |
| J19 | Service error contains synthetic token/URL markers | Public error and logs do not expose secrets |
| J20 | Async fire-and-forget operation fails | Chosen error policy is observable; no silently successful delivery claim |

J20 requires an explicit decision: the host may await registered work, reject unsupported fire-and-forget use, or expose a failure event. Do not define the expected result merely by copying whichever behavior the first implementation happens to have.

### Deterministic cancellation example

```text
start Dispatch in an errgroup
wait until fake Post reports Entered
cancel invocation context
wait until fake Post reports context cancellation
wait for Dispatch to return
assert zero successful sends and no pending workers
Close host with a bounded cleanup context
```

A timeout in the test is a failure bound, not a substitute for the `Entered` barrier. Do not use a fixed sleep to guess whether the worker started.

## 9. Wire test catalog

### L2 HTTP contracts

| ID | Fixture behavior | Required assertion |
|---|---|---|
| H01 | Accept Web API calls | Correct method, path, bot token, body encoding |
| H02 | Accept connections.open | Uses app token, never bot token |
| H03 | Accept response URL | Correct interaction route, ephemeral body, no unrelated bearer credential |
| H04 | HTTP 200 with `{ok:false,error}` | Treated as failure, mapped to useful domain error |
| H05 | HTTP 429 with Retry-After | One retry owner; not retried before delay; bounded attempts |
| H06 | Cancel during rate-limit backoff | Immediate cancellation rather than waiting full delay |
| H07 | Unknown/malformed Retry-After | Explicit fallback policy; no tight retry loop |
| H08 | Request blocked before any response | Context stops the request and cleanup completes |
| H09 | Server records accepted post, then drops connection | `delivery_unknown` or equivalent; no blind duplicate post |
| H10 | Truncated/malformed JSON or HTML response | Decode error, no fabricated message reference |
| H11 | Wrong token type accepted by broad emulator | Strict fixture still rejects it and catches wiring mistake |
| H12 | External redirect or external returned socket URL | Dial denied; no real traffic |
| H13 | Invalid payload supplied through JS | Zero HTTP requests |
| H14 | Response URL expires or belongs to prior invocation | Clear failure and no reroute to another user |
| H15 | Retry one method while another is available | Intended method/channel scoping of pacing verified |

For H09, record acceptance in the fixture and close the response connection before returning its body. That distinguishes uncertain delivery from a request that never reached the server. Do not model this with `{ok:false}`: that would test a different failure.

### L2 Socket Mode contracts

| ID | Injected traffic | Required assertion |
|---|---|---|
| W01 | Hello plus valid mention envelope | Correct identity and typed event snapshot |
| W02 | Command envelope with response capability | ACK ID and command routing correct |
| W03 | Handler/API blocked after admission | ACK observed while business work remains blocked |
| W04 | Same envelope sent twice | Both deliveries receipted, one invocation admitted |
| W05 | Same event_id under two envelope IDs | One business invocation, both ACKs |
| W06 | Same event ID in distinct workspace context | Correct tenancy validation/key policy |
| W07 | Queue full | Bounded admission and documented busy/drop behavior |
| W08 | Invalid JSON, missing envelope ID, unsupported type | Defined handling, no panic or accidental side effect |
| W09 | Self/bot event, wrong workspace/app | Filtered; valid receiptable rejected envelopes are ACKed according to policy |
| W10 | Disconnect refresh request | Reconnect and new delivery without rebuilding duplicate JS hosts |
| W11 | Abrupt close or failed WebSocket upgrade | Bounded recovery; observable error |
| W12 | Cancel while reconnecting | Receiver and reconnect timer stop |
| W13 | ACK send fails after admission | Failure visible; dedupe/retry semantics defined |
| W14 | Dedupe cache reaches capacity/TTL | Bounded size, deterministic eviction, documented replay window |
| W15 | Shutdown during decode/admission/handler/API | No channel-close panic; no orphan process |

For W03, use a fixture sequence counter to assert `ACK written < work released`. Retain measured latency as diagnostic evidence. A functional CI guard can require ACK before an adjustable test budget below Slack's deadline, but do not mistake a busy CI machine's submillisecond timing for a product guarantee.

W13 needs a written delivery policy. If work has already been admitted but ACK fails, the duplicate must not start another handler while the first remains active. If admission failed, its dedupe reservation must not suppress a later usable delivery. Test both paths with explicit barriers.

## 10. Stateful process scenarios with the deployed mock

L3 executes the actual CLI command and selected bot, with real Slack SDK transport pointed at the emulator. Scenario setup and inspection use the emulator's TypeScript API or documented admin endpoints. These are tests of the complete composition, not additional product commands.

### Scenario P01: offline discovery and configuration

Run list/help/inspect with no tokens and an empty test-owned credential directory. Assert deterministic bot names and fields, a useful malformed-script diagnostic, and zero connection attempts. Confirm that help and manifest generation do not open Socket Mode. Run from both repository root and a different working directory with an explicit repository path.

### Scenario P02: slash ping

Start the app, wait for a mock connection, inject `/golem-ping` as Alice in the test channel, and wait for an ephemeral message containing a unique scenario marker. Assert the response recipient, the envelope ACK, and absence of an unintended public message. The return from `slashCommand` may contain only the immediate ACK response; a later response-URL message must be observed separately.

### Scenario P03: threaded mentions

Create a top-level mention, wait for exactly one bot reply under its root timestamp, and verify the full thread structurally. Then have Bob mention the bot within the same thread. Assert that the next bot reply stays in the original thread rather than starting a nested or unrelated thread. Exact `channel` and `thread_ts` comparisons matter more than substring matches on text.

### Scenario P04: no self-reply loop

Use an explicit scenario subscription set that includes the message traffic necessary to exercise the filter, with bot echo enabled. Verify that bot-generated traffic is ignored by our application. The mock's Bolt examples rely on Bolt's ignoreSelf middleware; our Go host must implement its own filtering. Finish with an invocation-completion barrier and assert no further posts, rather than waiting an arbitrary short interval.

### Scenario P05: membership and permissions failures

Create a channel without the bot and attempt an explicit post through the app's normal operation. Expect `not_in_channel`, one failed API call, no stored post, and a usable public error. Use injected `missing_scope` to test mapping, while clearly labeling it as an error-handling test rather than real scope enforcement.

### Scenario P06: error and rate-limit recovery

Using the TypeScript driver, inject one 429 for `chat.postMessage` with a short Retry-After. Deliver a mention. Assert one rejected attempt followed by one successful message, with no duplicate JS invocation. Then exercise cancellation during a longer retry delay in L2; do not make the process suite wait minutes.

### Scenario P07: reconnect

Send a successful mention, request `disconnectSockets("refresh_requested")`, and observe the old connection disappear and a new connection arrive. Merely seeing `connections == 1` can match the old connection. After reconnection, deliver a new marker and verify one reply. Retain connection lifecycle observations; do not copy the Bolt test's assumed five-second backoff into the Go test.

### Scenario P08: process cancellation

Send SIGTERM to the actual bot process while a handler is pending. Assert the documented exit status, bounded shutdown, and no remaining socket connection. Test normal startup, failed startup, active operation, and disconnected states. Stop the mock after the bot so teardown does not manufacture unrelated transport errors.

### Scenario P09: exact credentials and endpoint routing

Run with synthetic bot/app tokens stored in test-owned files. Reject swapping them in a strict fixture. Inspect stdout/stderr and error output for raw token markers, response URLs, and socket tickets. Confirm that only bot-declared config enters JavaScript. The broad emulator's permissive app-token handling makes the L2 token test essential.

### Scenario P10: bad manifest and unknown capability

Use a JSON manifest declaring one command and one event, then attempt an undeclared command through the mock. Verify the setup failure is reported as a harness/configuration failure, not a bot runtime failure. Add a fixture that calls an unimplemented API method and require a visible failure; never use a mock that defaults all unknown methods to success as the sole oracle.

## 11. Driver design and concrete mock API usage

Prefer a small Bun scenario driver that imports `SlackMock` and controls a tmux-managed application process through a launcher helper. This gives access to `injectFault`, `deliveries`, and `hub.send` without modifying the upstream mock to add routes. The helper is test orchestration, not a compatibility adapter around the production Slack API.

```typescript
// Pseudocode: runBotInTmux/stopBot are proposed harness functions.
const mock = await SlackMock.start({
  host: "127.0.0.1", port: 0,
  manifest: minimalManifestObject,
  subscribedEvents: ["app_mention"],
  echoBotMessages: true,
});
let bot;
try {
  bot = await runBotInTmux({
    apiBaseUrl: mock.apiUrl,
    botToken: mock.bot.token,
    appToken: mock.bot.appToken,
    teamId: mock.team.id,
  });
  await mock.waitForConnection(10_000);
  const root = await mock.postMessage({
    channel: "general", user: "alice",
    text: `<@${mock.bot.userId}> scenario-P03-unique`,
  });
  const reply = await mock.waitForMessage({
    channel: "general", thread_ts: root.ts, from: "bot",
    text: "expected response",
  }, {timeoutMs: 10_000});
  assert(reply.thread_ts === root.ts);
  // Also assert exact result count after app completion.
  assert(mock.deliveries("app_mention").every(d => d.acked));
} finally {
  if (bot) await stopBot(bot);
  await mock.stop();
}
```

Use new instances or a scenario-specific marker to avoid an old result satisfying the wait. Use `mock.hub.send("events_api", payload)` twice with the same payload event ID for W05-style end-to-end coverage; each call generates a new envelope ID. The Go fixture remains preferable for races and malformed bytes.

These admin endpoints exist in the inspected source and can support manual exploration from any language:

| Method/path | Purpose |
|---|---|
| `GET /mock/state` | Workspace, app, bot IDs, connections, users, channels |
| `POST /mock/messages` | Human message; accepts channel, user, text, thread_ts |
| `GET /mock/messages?...` | Query stored messages; check exact query semantics |
| `POST /mock/commands` | Invoke a slash command |
| `POST /mock/actions` | Click a button on a stored message |
| `POST /mock/views/submit` | Submit a modal |
| `GET /mock/api-calls?method=chat.postMessage` | Inspect API records |
| `POST /mock/disconnect` | Request disconnect/reconnect |
| `GET /mock/channels/CHANNEL/threads/TS` | Retrieve a full stored thread |

Do not use `ephemeral=false` as though the query parser coerces it to boolean: the admin query is cast from URL strings. Prefer the typed API or direct structural inspection for assertions involving ephemeral messages.

### Manual runbook

The following mock commands are supported by the inspected checkout. Run them only when implementing or exercising the harness; they were not executed during this design task. Reserve a dedicated local port and session name first. Do not kill a main-thread server simply because it occupies the example port.

```bash
# From a dedicated terminal, after verifying port 14040 is free:
tmux new-session -d -s slack-mock-local-test \
  -c /home/manuel/code/others/slack/desplega-slack-mock \
  'bun src/cli.ts serve --host 127.0.0.1 --port 14040'
tmux capture-pane -p -t slack-mock-local-test -S -100
curl --fail http://127.0.0.1:14040/mock/state
```

The test driver should create synthetic credential files and pass their paths to the actual `go run ./cmd/slack-bot ...` invocation in another tmux session. The exact application flags are finalized by the implementation owner; this guide does not supply a command that pretends the transport already exists.

```bash
curl --fail -X POST http://127.0.0.1:14040/mock/commands \
  -H 'content-type: application/json' \
  -d '{"command":"/golem-ping","channel":"general","user":"alice"}'
```

For human review, open `http://127.0.0.1:14040/c/general`. Stop the test bot first. For the owned mock process, follow the workspace's cleanup convention `lsof-who -p 14040 -k`, then remove the dedicated tmux session if still present. Automated launchers must register equivalent cleanup on failure and signals. Prefer dynamically assigned ports in automated runs; expose readiness through a machine-readable file or line, not a sleep.

## 12. Later UI and feature coverage

Add these only when their corresponding product features enter scope.

- **Buttons:** create a persistent Block Kit message; inject an action using its exact `action_id`; assert ACK, authorized user context, and `chat.update` targeting the same reference. Unknown action IDs fail without an accidental post.
- **Ephemeral updates:** distinguish response-URL replacement from persistent `chat.update`. A test must not permit a persistent-message operation merely because the emulator's state model happens to accept it.
- **Modals:** open with a fresh trigger, submit invalid values, assert `response_action: errors` in the ACK and preservation of the view. Submit valid values and assert close/update plus exactly one business side effect. Expired triggers and validation deadline failures belong in dedicated cases.
- **Direct messages:** include `message.im` in the explicit manifest, create/open a DM, and verify correct channel identity and absence of channel-wide output.
- **History:** fixture pagination, missing access, deleted root, empty thread, and rate-limited continuation. Do not infer full history API fidelity from a single stored thread example.
- **Durable inbox:** restart after admission and before execution; restart after external acceptance and before local completion; specify replay and uncertain-delivery behavior. The mock's journal stores fake Slack state, not our application's durable inbox.

For visual review, pair each screenshot with the exact message JSON and scenario ID. Use the emulator's `screenshot` or `frames` command after the structural assertions pass. Pin viewport, browser, theme, and seed data. A screenshot demonstrates the emulator rendering; it is not evidence that Slack received an ACK or that the official Slack UI will render identically.

## 13. Test execution, CI gates, and failure artifacts

### Fast gate

Run L0 and L1 on every change. Run L2 when transport or lifecycle changes. Focused Go commands can be:

```bash
go test ./pkg/slackbot ./internal/jsslack
go test -race ./pkg/slackbot ./internal/jsslack
# Once the proposed package exists:
go test ./internal/slacktransport
```

These are implementation-time commands, not reported results of this document. Run `go test ./...`, `go build ./...`, and the repository's lint gate at integration milestones. No new runtime code was changed by this side task, so executing the main thread's unfinished tests here would not establish a useful baseline.

### Process gate

Provision a pinned Bun runtime, mock checkout/dependency cache, Go toolchain/dependency cache, and tmux in the test environment. Cache keys include lockfile, tool versions, architecture, and mock revision. Test execution disables external egress but permits loopback. It must not install dependencies midway through a scenario.

Keep the L3 suite separate from the default Go unit suite so missing Bun is not disguised as a passed skipped integration test. In the designated integration job, a missing dependency is a setup failure. A documented local command may skip optional L3 tests explicitly, but the CI gate must report that omission.

Use a modest suite-level timeout and per-step deadlines. Start with 10-second readiness/result bounds and a separately bounded reconnect test; tune after measurement rather than copying another SDK's timers. Run stateful process scenarios serially until isolation is demonstrated. Parallel Go unit tests must have separate host instances and ephemeral ports.

### Flake policy

A failure is retained on the first attempt. A diagnostic rerun can establish reproducibility, but does not erase the initial failure. Prefer barriers, completion receipts, request counts, and unique IDs over fixed sleeps. Do not repeatedly widen timeouts to mask a deadlocked owner or a missing ACK.

For negative assertions, wait for a known completion state or release-and-drain barrier before asserting zero additional calls. A short quiet period alone cannot prove that a delayed goroutine will never post.

### Artifacts

Write artifacts to a unique per-run directory, for example `ttmp/.../artifacts/local-testing/<run-id>/` for manual investigation or the CI artifact directory. Retain on failure:

- `versions.json`: application, SDK, mock, runtime, scenario, and seed revisions.
- `result.json`: pass/fail/setup-failure classification, expected and observed values.
- `app.jsonl` and `mock.log`: sanitized structured output.
- `wire.jsonl`: received frames, ACKs, and HTTP calls with correlation IDs and monotonic offsets where available.
- `state.json`: relevant channels, messages, ephemerals, and views.
- `journal.jsonl`: emulator state history if enabled.
- `processes.json`: owned session/process IDs, start/exit status, cleanup outcome.
- Optional screenshots and their corresponding message snapshots.

The observed mock API recorder removes body token arguments, but that is not a universal secrecy guarantee. Redact headers, response URLs, socket tickets, and credential-file contents in the harness. All scenario text and identities are synthetic. Full traces should never need real workspace data.

## 14. Implementation sequence and completion criteria

### T1: pin and validate the mock contract

Create the dependency inventory and minimal JSON manifest. Execute the interoperability probe in section 7. Record which SDK methods and mock endpoints actually passed. This milestone prevents a large test harness from being built on an assumed client compatibility.

### T2: implement domain and runtime assertions

Add the L0/L1 scenarios using the current public domain interfaces and real JS fixtures. Prioritize cancellation, worker saturation, stale context, rejected promises, and reply duplication. Run focused tests and the race detector. Any missing semantic decision, such as fire-and-forget error propagation, is resolved explicitly before writing its expected assertion.

### T3: implement the precise wire fixture

Support only `auth.test`, `apps.connections.open`, `chat.postMessage`, response URLs, and scripted socket frames initially. Add request recording, context-aware blocking, dropped responses, 429s, explicit reconnect control, and ACK barriers. Use `slacktest` constructor overrides where they help, or a direct httptest fixture. Do not implement a second stateful Slack clone.

### T4: connect the full process

Build the Bun driver and tmux lifecycle helpers, with dynamic ports and readiness receipts. Implement P01–P10 using actual application commands. Keep synthetic credentials isolated and ensure no endpoint escapes loopback. Save failure artifacts before cleanup.

### T5: install the CI gates and optional visual checks

Make setup failures explicit, document cache preparation, and add the integration job. Add screenshot tests only for UI features whose structure already has assertions. Record emulator limitations alongside the supported feature list.

The local-testing milestone is complete when a fresh prepared environment can execute the defined MVP layers without Slack credentials; wrong-token routing, missing ACKs, duplicate events, blocked work, reconnect, and shutdown each have a failing-then-passing test; the process suite uses the real composition root; and a failed run produces enough evidence to reproduce the boundary failure.

Real Slack validation remains a later, separately authorized step for installation scopes, platform behavior, and official-client rendering. Local test success must be described as local test success.

## 15. Source references and analysis record

### Immutable upstream references

- [slack-go test server](https://github.com/slack-go/slack/blob/4b9cec9bdb1eaca3b66968b02c9f08fecd860285/slacktest/server.go)
- [slack-go default handlers](https://github.com/slack-go/slack/blob/4b9cec9bdb1eaca3b66968b02c9f08fecd860285/slacktest/handlers.go)
- [slack-go custom socket tests](https://github.com/slack-go/slack/blob/4b9cec9bdb1eaca3b66968b02c9f08fecd860285/socketmode/socket_mode_managed_conn_test.go)
- [desplega server and admin router](https://github.com/desplega-ai/slack-mock/blob/6397b31a9c5e04a3ba52dd6e82f16ab8b5b10eac/src/server.ts)
- [desplega Socket Mode hub](https://github.com/desplega-ai/slack-mock/blob/6397b31a9c5e04a3ba52dd6e82f16ab8b5b10eac/src/socket-mode.ts)
- [desplega Web API model](https://github.com/desplega-ai/slack-mock/blob/6397b31a9c5e04a3ba52dd6e82f16ab8b5b10eac/src/web-api.ts)
- [desplega protocol tests](https://github.com/desplega-ai/slack-mock/blob/6397b31a9c5e04a3ba52dd6e82f16ab8b5b10eac/test/socket-mode.test.ts)
- [desplega fault tests](https://github.com/desplega-ai/slack-mock/blob/6397b31a9c5e04a3ba52dd6e82f16ab8b5b10eac/test/faults.test.ts)
- [HTTP testing library server](https://github.com/chrishutchinson/slack-testing-library/blob/3f450bb5abcf00991a5ebe8cbeef98a60e1d4530/src/util/server.ts)
- [HTTP testing library interaction limitations](https://github.com/chrishutchinson/slack-testing-library/blob/3f450bb5abcf00991a5ebe8cbeef98a60e1d4530/src/slack-testing-library.ts)
- [Skellington Nock interception](https://github.com/Skellington-Closet/slack-mock/blob/7bf3eb5073a8a2c071ece360a7b8729390a37024/src/mocker/web.js)

Official Slack protocol references and archived copies remain in the [ticket source inventory](../sources/README.md), particularly Socket Mode, interaction handling, and Web API rate limits. The mock repositories are examples and test infrastructure, not the authority for production protocol rules.

### Research receipt

The active user request was: “Create a separate design doc that writes out a full local testing plan. you can clone and analyze the projects in ~/code/others/slack/ if you want.” Four shallow clones succeeded at the paths and commits listed above. Source searches and targeted reads established the recommendations and limitations. A few combined inspection commands referenced a file in the wrong checkout and returned `No such file or directory`; the corresponding files were then read from the correct projects. No conclusions depend on those failed reads.

No upstream tests or runtime experiments were executed. No application files, shared implementation diary, existing guide, or Git index were edited by this side-conversation design task. This separate document is the handoff to the main implementation thread.

Document validation: docmgr frontmatter and ticket doctor passed; Pandoc parsed the Markdown successfully; relative links and fenced-block balance were checked. These are documentation checks, not mock interoperability results.
