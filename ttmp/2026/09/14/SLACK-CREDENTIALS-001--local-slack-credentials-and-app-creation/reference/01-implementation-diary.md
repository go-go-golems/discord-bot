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
