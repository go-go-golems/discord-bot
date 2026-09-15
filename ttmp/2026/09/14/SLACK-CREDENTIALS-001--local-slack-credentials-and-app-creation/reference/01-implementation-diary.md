---
Title: Implementation diary
Ticket: SLACK-CREDENTIALS-001
Status: active
Topics:
    - discord-bot
    - api-design
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-14T21:31:19.295803347-04:00
WhatFor: ""
WhenToUse: ""
---

# Implementation diary

## Step 1: Reset scope for a local tool

Created a separate implementation ticket after the user rejected the earlier reliability-heavy design. The replacement uses two files, named maps, explicit refresh, and manual installation/recovery. The credential manager is not implemented in this planning step.

### Prompt Context

**User prompt (verbatim):** "ok, start again, make a pragmatic implementation ticket."

**Follow-up (verbatim):** "mark the previous one deprecated"

**Assistant interpretation:** Replace the implementation scope with a small local workflow and retire the old credentials design.

**Inferred user intent:** Make progress on a useful local bot without production reliability infrastructure.

### What I did

- Created a plan, ordered tasks, command forms, and bounded acceptance tests.
- Linked the existing defuddle source archive.
- Deprecated the old credentials design while keeping the broader Slack-support ticket active.

### Why

Detailed explanation does not imply automatic recovery from every rare failure. Manual re-import is acceptable here.

### What worked

Existing manifest/create-app source provides a concrete starting point. A single secrets document keeps replacement tokens together without generations or locks.

### What didn't work

The preceding design exceeded the user's needs. Also, the first combined documentation patch failed with `apply_patch verification failed: Failed to find expected lines` because a diary template comment differed; reread the actual template and applied the corrected patch. No partial changes from that patch were present.

### What I learned

Documentation depth should explain a small implementation, not expand its reliability requirements.

### What was tricky to build

App secrets and workspace installation tokens still need distinct map entries. Shared management identities require references, but not a generalized storage framework.

### What warrants a second pair of eyes

Ensure implementation does not add automatic refresh/retries implicitly. The current create-app source is an uncommitted foundation; profile commands remain proposed.

### What should be done in the future

Implement the ordered tasks. Keep production transport separate.

### Code review instructions

Read the plan/tasks, then `pkg/slackcli/create_app.go` and `commands.go`. Validate ticket metadata and relative links before committing.

### Technical details

Storage is `config.yaml` plus `credentials.json`, with private permissions, temporary-file rename, and one writer at a time. This planning milestone makes no API requests and reads no credentials.

## Step 2: Implement the pragmatic profile store and explicit refresh

The first implementation milestone adds a small `internal/slackconfig` package and exposes profile, import, status, and refresh commands. Profile-based `bots create-app` now uses the stored management access token, saves returned app credentials, and links the app to the profile. The deprecated design's journals, generations, locks, and automatic recovery are not present.

### Prompt Context

See Step 1 for the original implementation and deprecation prompts. The current request authorizes implementing that pragmatic ticket with commits and a detailed diary.

**Assistant interpretation:** Implement the seven ticket tasks in reviewable milestones while preserving explicit token-file app creation.

**Inferred user intent:** Have a usable local workflow for several profiles/apps without cloud-scale credential infrastructure.

### What I did

- Added `internal/slackconfig/store.go` with YAML metadata, JSON secrets, bounded reads, private directory/file modes, atomic temporary-file replacement, and profile resolution.
- Added `pkg/slackcli/credentials.go` with `credentials import-management`, `credentials status`, `credentials refresh`, and root `profiles` listing.
- Added `--profile` and `--config-dir` to `bots create-app`; managed creation saves app ID and returned app credentials.
- Added storage and refresh/profile creation tests, and expanded embedded help with the command sequence.

### Why

Two ordinary files are sufficient for one developer on one workstation. Saving access and refresh tokens together avoids the most likely local inconsistency without introducing a database or background service.

### What worked

- Focused `04-go-offline.sh test ./pkg/slackcli ./internal/slackconfig ./cmd/slack-bot` passed after the profile round-trip fix.
- Tests verify private modes, refresh request form data, token replacement, safe status output, manifest bearer selection, and app linking.
- Existing explicit token-file create-app tests continue to pass.

### What didn't work

- The first profile creation test saved the app under an empty map key because `profileName, profile, err :=` shadowed the outer `profileName`. The test showed `profiles: "": {}` and `apps: "":`; changing this to assignment fixed the issue.
- The first config-directory flag implementation used a pointer returned from a helper instead of binding the command flag to the settings variable. The test attempted `/home/manuel/.config/go-go-slack` and failed with `mkdir /home/manuel/.config/go-go-slack: read-only file system`; binding with `StringVar` fixed it.

### What I learned

Glazed's default middleware can overwrite a configured default value in this command shape; leaving `config-dir` without a parser default and applying the user-config-directory fallback in execution avoids that ambiguity.

### What was tricky to build

The profile path must select a management credential before the shared create request, then persist app metadata after Slack returns. The existing explicit-file path still reserves its optional output file separately, so the two modes have distinct persistence behavior.

### What warrants a second pair of eyes

Review the simple two-file save order and the decision to report manual recovery after a failed refresh. Confirm the `tooling.tokens.rotate` form request matches the current Slack contract before using it against a real token.

### What should be done in the future

Add runtime-token import and profile-based production transport only after this milestone is reviewed. Do not add automatic refresh loops, retries, journals, or multi-process coordination without a new requirement.

### Code review instructions

Read `internal/slackconfig/store.go`, `pkg/slackcli/credentials.go`, and the new tests. Run the focused command above, then inspect `slack-bot credentials --help` and `slack-bot bots create-app --help` using the offline wrapper.

### Technical details

The store caps JSON/YAML reads at 1 MiB, writes modes `0700`/`0600`, and uses same-directory `CreateTemp` plus rename. Refresh posts `refresh_token` as form data to `https://slack.com/api/tooling.tokens.rotate`, requires replacement `token` and `refresh_token`, and records `exp` as RFC3339. Errors omit token contents and advise importing a new pair.

## Step 3: Add manual runtime-token import and finish the local command surface

The final implementation milestone adds `credentials import-runtime`, which associates a bot token with an installation and a Socket Mode token with its app. It deliberately performs no remote verification or connection; the user supplies the workspace ID after completing Slack's browser installation/settings flow.

### Prompt Context

See Steps 1–2 for the verbatim user prompts and scope. The current request is to implement the pragmatic ticket.

**Assistant interpretation:** Complete the remaining local storage task, update help, validate, and commit.

**Inferred user intent:** Store enough credentials to run a selected local bot while keeping setup understandable and manual.

### What I did

- Added `credentials import-runtime --profile --installation --team-id --bot-token-file --app-token-file`.
- Checked an existing installation's app/workspace association before replacement and saved profile/credential maps together.
- Updated the implementation plan and embedded help with runtime import instructions.

### Why

Runtime bot and Socket Mode tokens are separate from management tokens and need explicit association, but a local tool does not need an OAuth callback server or remote identity probe to store them.

### What worked

Focused tests pass for management import, refresh, profile app creation, runtime import, private modes, and secret-free output.

### What didn't work

No new failures in this milestone. Earlier profile shadowing and config-directory binding failures are recorded in Step 2.

### What I learned

The local credential lifecycle is complete enough for manual Slack setup: management pair, app creation, browser installation, and runtime token import are separate explicit actions.

### What was tricky to build

The Socket Mode app token belongs to the app record, while the bot token belongs to an `(app, workspace)` installation. The import command updates both references without duplicating the app identity.

### What warrants a second pair of eyes

Review the one-writer assumption and ensure production transport reads only the selected installation. Confirm future OAuth work does not bypass the explicit profile association.

### What should be done in the future

Run the full repository validation gate, then perform an optional real-workspace smoke test with the user's supplied token files. OAuth callback automation and bot-token rotation remain out of scope.

### Code review instructions

Review `internal/slackconfig/store.go`, `pkg/slackcli/credentials.go`, and `pkg/slackcli/create_app.go`; run focused tests and inspect CLI help. Check task/changelog updates and the deprecated predecessor link.

### Technical details

Runtime import rejects missing profile/app, missing files, and known app/workspace mismatch. It writes `config.yaml` and `credentials.json` via private atomic replacement and emits only profile, installation, and team ID.

Full validation for this milestone completed with the repository test suite in tmux: all packages passed with loopback socket access. The sandbox-only run had failed at the existing `httptest` IPv6 listener; the escalated tmux run passed `internal/slacktransport` as well. Full `build -buildvcs=false ./...`, `vet ./...`, focused tests, CLI help checks, `git diff --check`, and `docmgr doctor --ticket SLACK-CREDENTIALS-001` all pass.

## Step 4: Commit the completed pragmatic implementation

The implementation and ticket bookkeeping are committed. All seven ticket tasks are checked. Generated build and PDF artifacts from earlier work remain untracked and were intentionally excluded from this source commit.

### Prompt Context

See Steps 1–3 for the verbatim request and implementation decisions.

**Assistant interpretation:** Finish the requested implementation with an auditable commit and diary record.

**Inferred user intent:** Leave a usable, reviewable local Slack credential workflow in the repository.

**Commit (code):** `23e819b` — `feat(slack): add pragmatic local credentials and profiles`

### What I did

- Committed the store, commands, profile-based app creation, runtime import, tests, help, and ticket updates.
- Preserved the deprecated predecessor design as historical context.

### Why

The commit is the review boundary for the complete local scope.

### What worked

Full tests passed in tmux with loopback access; build and vet passed; CLI help and docmgr doctor passed.

### What didn't work

The unprivileged full-test attempt was blocked by the sandbox's IPv6 listener restriction, then passed with the required socket-enabled tmux run.

## Step 5: Record the local developer-install design

The next requested capability is an install verb that follows the local Slack CLI flow. Source inspection of the open-source Slack CLI showed that its local install path calls the Web API method `apps.developerInstall`, receives both `api_access_tokens.bot` and `api_access_tokens.app_level`, and exports them as `SLACK_BOT_TOKEN` and `SLACK_APP_TOKEN`. Slack's public API reference does not document this method, so the design keeps it explicitly local and provides a manual fallback.

### Prompt Context

**User prompt (verbatim):** "ok, implement it for our purposes, add a design doc documenting all this information you found, add the tasks to the ticket, then implment, commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)"

**Assistant interpretation:** Add a pragmatic profile-based install command, document the evidence and limitations, track the work in the active ticket, and implement it with focused tests and commits.

**Inferred user intent:** Make local Slack app setup require as little manual token copying as practical without turning this project into a hosted OAuth service.

### What I did

- Read the Slack CLI source under `/home/manuel/code/others/slack-cli`, including `internal/api/app.go` and `internal/pkg/apps/install.go`.
- Added design doc `design-doc/02-local-app-installation-with-slack-cli-developerinstall.md` with the observed request/response contract, credential mapping, pseudocode, test plan, and browser-OAuth comparison.
- Added four ticket tasks for documenting, implementing, persisting tokens, and testing the feature.

### Why

The endpoint is useful for the user's single-workstation workflow, but its undocumented status must be visible to anyone maintaining the command. Recording the Slack CLI evidence also makes it possible to update the isolated call if Slack changes it.

### What worked

Documentation and task additions were committed as `d1831cf` (`docs(slack): design local developer app installation`).

### What didn't work

No design-time command or Slack request was executed. The design intentionally avoids treating the private endpoint as a supported public OAuth contract.

### What I learned

The Slack CLI does not need an HTTPS callback for local developer installation. It sends the app ID, bot scopes, outgoing domains, and optional team ID directly with a developer token, then receives runtime tokens in the response.

### What was tricky to build

The app-level token belongs to the app record, while the bot token belongs to a specific app/workspace installation. The design preserves that distinction in the existing two-file store.

### What warrants a second pair of eyes

Review whether relying on `apps.developerInstall` is acceptable for this local tool and whether the fallback text is clear enough if Slack removes or restricts the method.

### What should be done in the future

Implement the isolated request, persist both token types, update help, and validate without contacting Slack in automated tests.

### Code review instructions

Read the new design doc and compare its wire example with `/home/manuel/code/others/slack-cli/internal/api/app.go`. Confirm that no user or app token is included in the ticket or command output.

### Technical details

The planned request is `POST https://slack.com/api/apps.developerInstall` with a bearer management access token and JSON fields `app_id`, `bot_scopes`, `outgoing_domains`, and `team_id`. The implementation will derive `bot_scopes` from `Manifest(d)` and use a deterministic `<profile>-<team-id>` installation key.

## Step 6: Implement and test profile-based app installation

The implementation adds `slack-bot bots install NAME`, which resolves an existing profile and app, calls the observed developer-install endpoint once, and saves the returned bot and app-level tokens in the existing store. The output contains only identifiers and names; the tokens remain in the private credentials file.

### Prompt Context

See Step 5 for the verbatim implementation request and the design decision it authorized.

**Assistant interpretation:** Complete the install command and its documentation/tests, then record a reviewable commit and validation results.

**Inferred user intent:** After creating an app, install it into a selected workspace with one concise local command and have the credentials ready for future runtime wiring.

### What I did

- Added `team-id` settings and an `install` operation to `pkg/slackcli/commands.go`.
- Added `pkg/slackcli/install_app.go`, including manifest-derived scopes, bounded response handling, app-ID/token checks, deterministic installation records, and secret-free output.
- Added `pkg/slackcli/install_app_test.go` with wire assertions, persistence checks, setup failures, Slack errors, mismatches, missing tokens, and no-overwrite behavior.
- Updated `pkg/slackdoc/slack-offline.md` with the command and its undocumented-endpoint fallback.

### Why

The command reuses the existing profile store and manifest generator, avoiding duplicate scope configuration and keeping management, app, and installation credentials separate.

### What worked

The focused offline command passed:

```text
ttmp/2026/09/10/DISCORD-SLACK-001--add-slack-support-to-discord-bot/scripts/04-go-offline.sh test ./pkg/slackcli ./internal/slackconfig ./cmd/slack-bot
ok github.com/go-go-golems/discord-bot/pkg/slackcli
ok github.com/go-go-golems/discord-bot/internal/slackconfig
ok github.com/go-go-golems/discord-bot/cmd/slack-bot
```

### What didn't work

The first test run failed because invalid-setup cases omitted `--bot-repository`; bot discovery therefore failed before the intended profile/team validation with `open .../pkg/slackcli/examples/slack-bots: no such file or directory`. Adding the repository fixture path made those tests exercise the intended errors. No production code change was required for that failure.

### What I learned

Glazed resolves the bot descriptor before operation-specific validation, so tests for command argument errors still need a valid descriptor path.

### What was tricky to build

`Manifest` stores scopes as a concrete `[]string` inside `map[string]any`; the install helper must preserve that representation while extracting `oauth_config.scopes.bot`.

### What warrants a second pair of eyes

Check that saving the app-level token in `creds.Apps[profile.App]` and also mirroring it in the installation credential is useful for the current runtime and does not create ambiguity for future transport wiring.

### What should be done in the future

Run the full offline build, vet, focused help checks, and docmgr validation. A real Slack install should be treated as an optional smoke test because the endpoint is undocumented and mutates a workspace.

### Code review instructions

Review `pkg/slackcli/install_app.go`, `pkg/slackcli/install_app_test.go`, `pkg/slackcli/commands.go`, and the updated help. Verify the test fake asserts the bearer token and that output assertions reject both synthetic runtime token values.

### Technical details

The successful response requires `ok:true`, a matching or omitted `app_id`, and non-empty `api_access_tokens.bot` and `api_access_tokens.app_level`. HTTP failures, malformed responses, and timeouts return “outcome unknown” guidance without retrying. The installation key is `<profile>-<team-id>` and the profile's installation reference is updated only after a successful response and local save.

### What I learned

No cloud-scale reliability feature was needed to meet the local workflow: manual re-import is an adequate recovery path.

### What was tricky to build

Keep generated artifacts out of the source commit while retaining the ticket's implementation evidence and diary.

### What warrants a second pair of eyes

Review API request details for a real Slack workspace before using refresh or profile-based creation with production credentials.

### What should be done in the future

Optional real-workspace smoke test and production transport integration remain separate work.

### Code review instructions

Start at `23e819b`, inspect the seven checked tasks, and run `04-go-offline.sh test ./...` with loopback permission if needed.

### Technical details

Commit contains 14 files and 1,169 insertions. No supplied token file was read during automated validation, and no real Slack API mutation was performed by the tests.

## Step 7: Commit the developer-install implementation and close validation tasks

The install implementation is now committed, and the four follow-up ticket tasks are checked. The remaining validation is repository-level hygiene; no live Slack request is required because the endpoint mutates a workspace and the automated tests use a fake transport.

### Prompt Context

See Step 5 for the verbatim user request. This step records the commit and validation boundary requested there.

**Assistant interpretation:** Commit the implementation at a meaningful boundary, run the required offline checks, and leave the ticket with an auditable task list and diary.

**Inferred user intent:** Be able to review and use the local install flow without exposing credentials or depending on a live Slack test during development.

### What I did

- Committed the implementation as `457b257` (`feat(slack): install local apps via developerInstall`).
- Checked tasks `e7em`, `46bm`, `7wca`, and `aytz` with `docmgr task check`.
- Added a changelog entry linking the implementation and design doc.
- Ran focused tests, offline build, vet, CLI help, and `git diff --check`.

### Why

The code commit and documentation/task updates form a reviewable boundary. The ticket remains active for future OAuth or production-install work; all current pragmatic tasks are complete.

### What worked

The focused package gate passed after the fixture correction:

```text
04-go-offline.sh test ./pkg/slackcli ./internal/slackconfig ./cmd/slack-bot
ok github.com/go-go-golems/discord-bot/pkg/slackcli
ok github.com/go-go-golems/discord-bot/internal/slackconfig
ok github.com/go-go-golems/discord-bot/cmd/slack-bot
```

`04-go-offline.sh build -buildvcs=false ./...`, `04-go-offline.sh vet ./...`, the `bots install --help` smoke check, and `git diff --check` passed. The task command reported all four new tasks complete.

### What didn't work

No live Slack installation was attempted in this implementation pass. That is deliberate: `apps.developerInstall` is undocumented and workspace-mutating, while fake transport tests cover the request and persistence behavior deterministically.

### What I learned

The local command can provide the same low-copy workflow as Slack CLI while remaining explicit about its private API dependency and preserving a manual OAuth/dashboard fallback.

### What was tricky to build

The command must validate the local app/profile association before making the request, but cannot verify workspace identity remotely without adding another API dependency. The explicit `--team-id` keeps that boundary visible.

### What warrants a second pair of eyes

Review whether the mirrored `AppToken` in the installation credential should remain for future runtime selection or be reduced to the app record only. Also review Slack API changes before relying on this private method for routine setup.

### What should be done in the future

Wire the selected installation into a real Socket Mode transport, or create a separate ticket for documented browser OAuth installation if the workflow needs to support other users/workspaces.

### Code review instructions

Start with commits `d1831cf` and `457b257`, then inspect the design doc, `pkg/slackcli/install_app.go`, and its fake-transport tests. Run the offline wrapper commands above and `docmgr doctor --ticket SLACK-CREDENTIALS-001`.

### Technical details

The persisted installation key is `<profile>-<team-id>`. The app record receives `app_token`; the installation record receives `bot_token` and a copy of `app_token` for current runtime convenience. Output is JSON metadata only. HTTP, JSON, Slack `ok:false`, app-ID mismatch, missing token, and save failures are surfaced without token values or retries.

## Step 8: Align install profile selection with the documented default

Review found that the first implementation required `--profile` even though the design and existing store support a configured `default_profile`. The command now passes an empty profile name to `ResolveProfile`, allowing the default while still returning `no profile selected` when neither a flag nor a default exists.

### Prompt Context

See Step 5 for the verbatim implementation request. This is a follow-up correction within the same requested scope.

**Assistant interpretation:** Remove an unnecessary typing requirement and keep command behavior consistent with the profile store.

**Inferred user intent:** Minimize local setup commands while retaining explicit selection when multiple profiles exist.

### What I did

- Removed the unconditional `--profile` check in `pkg/slackcli/install_app.go`.
- Updated the missing-profile test to use an isolated config directory and assert the resolver error.

### Why

`default_profile` is already written by management import and is safe for a single-user workflow. Requiring the flag would contradict the documented profile resolution behavior.

### What worked

The focused offline gate passed again after the correction:

```text
04-go-offline.sh test ./pkg/slackcli ./internal/slackconfig ./cmd/slack-bot
ok github.com/go-go-golems/discord-bot/pkg/slackcli
ok github.com/go-go-golems/discord-bot/internal/slackconfig
ok github.com/go-go-golems/discord-bot/cmd/slack-bot
```

The correction was committed as `a3bd15f` (`fix(slack): honor default profile during install`).

### What didn't work

The initial test of missing-profile behavior omitted `--config-dir`, so it loaded the developer's real default store and reached the fake HTTP response instead of the intended resolver error. Isolating the test directory fixed the test without changing runtime behavior.

### What I learned

Tests for profile absence must control both the profile flag and config-directory source; otherwise a user's default profile can make the test non-deterministic.

### What was tricky to build

The command must keep `--team-id` mandatory even when the profile is defaulted, because the workspace target cannot be inferred safely from local token metadata.

### What warrants a second pair of eyes

Review the default-profile behavior with a multi-workspace configuration before using the command without `--profile`.

### What should be done in the future

If no default is desired, expose a profile listing/status workflow rather than guessing by map iteration order.

### Code review instructions

Compare the profile resolver call in `install_app.go` with the design doc's CLI contract and run the focused test command above.

### Technical details

The correction changes only pre-request validation. No Slack request is made when profile resolution fails, and all token redaction assertions remain unchanged.

## Step 9: Correct the standalone developer-install request shape

The first live install attempt returned Slack error `invalid_argument`. No tokens were returned and no local credentials were changed. Reviewing the Slack CLI implementation showed that it deliberately clears `team_id` for non-enterprise apps before serializing the `apps.developerInstall` request. The local command now follows that wire shape while retaining `--team-id` for local installation identity and storage.

### Prompt Context

**User prompt (verbatim):** "i ran go work use"

**Assistant interpretation:** Retry the requested live install after correcting the Go workspace setup and the request mismatch discovered from the previous response.

**Inferred user intent:** Complete the `ping` installation in the intended workspace, with minimal manual intervention.

### What I did

- Removed `team_id` from the standalone developer-install JSON request.
- Updated the fake transport assertion and design documentation.
- Preserved explicit team ID validation and `<profile>-<team-id>` local storage mapping.

### Why

The Slack CLI source sends `team_id` only as an organization grant identifier. For a standalone app, including the workspace ID produced the observed `invalid_argument` response.

### What worked

The focused offline tests passed after the request-shape correction. The earlier live attempt completed with a safe `invalid_argument` error and did not write runtime tokens.

### What didn't work

The first live request included `team_id` and Slack rejected it. The response contained only the error code; no secret-bearing response was logged.

### What I learned

The workspace ID is needed by this local credential store even when it is not an API argument for standalone developer installation.

### What was tricky to build

The same field name has two meanings in Slack CLI internals: an organization grant target in the API request and a local workspace selector in this tool. Keeping those roles separate avoids sending an invalid standalone argument.

### What warrants a second pair of eyes

Review the request body against the current Slack CLI source before future live installs, especially if enterprise-grid support is added.

### What should be done in the future

Add an explicit enterprise-install mode only if a real organization workflow requires it; do not infer one from a workspace ID.

### Code review instructions

Inspect `install_app.go` request serialization and compare it with `DeveloperAppInstall` in `/home/manuel/code/others/slack-cli/internal/api/app.go`.

### Technical details

The corrected JSON contains `app_id`, `bot_scopes`, and `outgoing_domains`; it omits `team_id`. The command still saves credentials under `profile-teamID` after a successful response.

## Step 10: Install the ping app in the requested workspace

After the request-shape correction and the user's `go work use` update, the authorized live installation completed successfully. Slack returned the expected app ID and the command persisted the bot and Socket Mode credentials without exposing their values.

### Prompt Context

**User prompt (verbatim):** "ok, install ping"

**Assistant interpretation:** Execute the new profile-based install command for the configured `go-go-golems` app and its known workspace, then verify the local status.

**Inferred user intent:** Finish setup of the real ping app so its runtime credentials are available to the local bot workflow.

### What I did

- Ran `04-go-offline.sh run ./cmd/slack-bot bots install ping --profile go-go-golems --team-id T0C1UJMCPGA --timeout-ms 30000`.
- Verified the safe status command and inspected metadata-only `config.yaml`.

### Why

The user explicitly requested the workspace mutation. The command uses the stored management credential, sends the corrected standalone request, and writes runtime credentials to the private credential file.

### What worked

Slack returned:

```json
{"app_id":"A0C1YJCCCP6","installation":"go-go-golems-T0C1UJMCPGA","profile":"go-go-golems","team_id":"T0C1UJMCPGA"}
```

Safe status confirms `has_access_token:true`, `has_refresh_token:true`, app `go-go-golems`, and installation `go-go-golems-T0C1UJMCPGA`. The metadata file now links that installation to app `A0C1YJCCCP6` and team `T0C1UJMCPGA`.

### What didn't work

The first live attempt, recorded in Step 9, returned `invalid_argument` because the request included `team_id`. The corrected retry succeeded.

### What I learned

The local profile store can now complete the management-token, app-creation, and developer-install portions of the setup without copying runtime tokens through the shell.

### What was tricky to build

The live command had to use the repository wrapper after plain `go run` encountered the stale Go workspace version declarations. The wrapper selected the cached toolchain while preserving the explicit network request.

### What warrants a second pair of eyes

Confirm the newly stored runtime credentials are used by the eventual real Socket Mode runner before relying on this installation for production messages.

### What should be done in the future

Add a separate runtime command that resolves this installation and connects to Socket Mode, with a local dry-run path retained for tests.

### Code review instructions

Review the successful command output above, the metadata diff, and Step 9's request-shape correction. Do not inspect or print the private credentials file contents.

### Technical details

The successful API response was accepted only after matching app ID and non-empty bot/app-level token fields. The command wrote `config.yaml` metadata and `credentials.json` secrets with the existing private atomic writer.

## Step 11: Add the real Slack runtime command

The installed credentials were previously stored but unused by the CLI runtime. This milestone adds a remote Slack client constructor and a `bots run` command that loads the selected profile installation, hosts the JavaScript bot, and reuses the existing Socket Mode ingress loop. The deterministic `run-local` mock path remains unchanged.

### Prompt Context

**User prompt (verbatim):** "ok, let's do it."

**Assistant interpretation:** Implement the next step identified after installation: connect the stored credentials to a real Slack runtime and make the bot runnable.

**Inferred user intent:** Run `ping` against Slack with one profile-based command instead of manually constructing a mock connection file.

### What I did

- Added `slacktransport.NewRemote`, using the pinned Slack SDK with explicit bot and Socket Mode tokens.
- Added `bots run NAME`, which resolves app/installation records, loads optional declared bot configuration, and dispatches through `slackhost.Host`.
- Allowed Slack's HTTPS response URLs for slash-command ephemeral replies while retaining strict host checks; local response URLs remain loopback-only.
- Added remote credential-resolution and response-capability tests.
- Updated the embedded help with the real runtime command and profile requirements.
- Added three runtime tasks to the ticket.

### Why

The existing `Client.Run` already verifies `auth.test`, acknowledges envelopes, filters ingress, and dispatches JavaScript. A remote constructor lets that reviewed behavior operate against Slack without duplicating transport logic.

### What worked

Focused CLI and transport tests passed. The remote runtime command rejects missing bot/app tokens before opening a network connection. Slack response-capability tests accept `https://hooks.slack.com/...` and reject non-Slack, non-HTTPS, and empty-path URLs.

### What didn't work

No live runtime session was started. Starting it would wait for events and could send real messages; the implementation is validated through local tests and is ready for an explicit smoke run.

### What I learned

The stored app-level token can be selected from the app record, with the installation copy as a fallback for credentials imported through the manual runtime path.

### What was tricky to build

The local client previously assumed every response URL matched a loopback origin. Remote slash commands require a controlled Slack webhook host instead, while arbitrary redirect or exfiltration URLs must remain rejected.

### What warrants a second pair of eyes

Review the allowed Slack webhook host list for GovSlack or future Slack response URL domains, and verify the SDK's reconnect behavior during a supervised live run.

### What should be done in the future

Run the process with `bots run ping --profile go-go-golems`, then add a small live smoke procedure for an app mention and slash command. Keep the mock transport tests as the default validation path.

### Code review instructions

Inspect `internal/slacktransport/client.go`, `pkg/slackcli/run_remote.go`, and the new tests. Compare `NewRemote` with `NewLocal` and confirm credentials are passed only into SDK clients and never into JavaScript configuration or logs.

### Technical details

`bots run` uses the profile's installation team ID and app ID, calls `auth.test` through the Slack SDK, and then invokes `Client.Run`. It accepts `--bot-config-file` for declared fields and `--timeout-ms` for host/invocation deadlines. Shutdown follows the existing context cancellation path.

## Step 12: Validate the real runtime command

The remote runtime implementation has passed the repository validation gate. The runtime-specific tasks are checked, and the ticket now contains the complete path from credential import and app installation through real Socket Mode execution.

### Prompt Context

See Step 11 for the verbatim request that authorized the runtime implementation.

**Assistant interpretation:** Finish validation and task bookkeeping, leaving a command that can be run explicitly against the installed app.

**Inferred user intent:** Have a dependable local development command for receiving Slack events and invoking the JavaScript bot.

### What I did

- Ran the full offline repository tests with loopback access.
- Ran offline build and vet.
- Ran `bots run --help` to verify the runtime command surface.
- Checked runtime tasks `y3sa`, `8qti`, and `axvn`.
- Added the runtime changelog entry.

### Why

The remote path reuses existing transport behavior, so full tests are the appropriate regression gate after focused CLI and transport tests.

### What worked

All packages passed `04-go-offline.sh test ./...`; `build -buildvcs=false ./...`, `vet ./...`, CLI help, and `docmgr doctor --ticket SLACK-CREDENTIALS-001` also passed.

### What didn't work

No validation failures occurred in this gate. A live runtime session was intentionally deferred because it waits for events and can send real Slack messages.

### What I learned

The installed profile is now sufficient to start the runtime with one command; no token file arguments or Slack CLI installation are required.

### What was tricky to build

The remote constructor needed to preserve strict response URL validation while allowing Slack's HTTPS webhook host for slash-command replies.

### What warrants a second pair of eyes

Run a supervised mention and slash-command smoke test and inspect the bot's response in the intended channel.

### What should be done in the future

Add an operator runbook for stopping/restarting the process and a small live smoke script if repeated local testing becomes routine.

### Code review instructions

Review commit `f4f36ee`, run `bots run --help`, and start the bot with the installed profile only when ready to receive real events.

### Technical details

The command is:

```sh
go run ./cmd/slack-bot bots run ping --profile go-go-golems --log-level debug
```

It resolves `go-go-golems-T0C1UJMCPGA`, uses the stored app ID `A0C1YJCCCP6`, verifies `auth.test`, and then enters the Socket Mode event loop. Credentials remain private to the Go transport.
