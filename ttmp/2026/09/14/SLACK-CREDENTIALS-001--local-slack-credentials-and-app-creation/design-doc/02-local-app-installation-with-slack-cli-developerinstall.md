---
Title: Local app installation with Slack CLI developerInstall
Ticket: SLACK-CREDENTIALS-001
Status: active
Topics:
    - discord-bot
    - api-design
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://pkg/slackcli/commands.go
      Note: install operation and Glazed flags
    - Path: repo://pkg/slackcli/install_app.go
      Note: isolated developerInstall API call and token persistence
    - Path: repo://pkg/slackcli/install_app_test.go
      Note: fake transport request and persistence tests
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-14T22:10:02.469741992-04:00
WhatFor: ""
WhenToUse: ""
---

# Local app installation with `apps.developerInstall`

## Decision

Add a small `slack-bot bots install` verb for the local developer workflow. It uses the management access token already stored in a selected profile, calls the same Slack Web API method used by the open-source Slack CLI, and stores the returned bot token and Socket Mode app token in the local credential store. The command is intentionally a single request with bounded input and output, no retry loop, and no browser callback server.

This is a pragmatic convenience for one developer working with a known Slack workspace. Slack does not document `apps.developerInstall` as a general public API. The implementation must therefore describe the endpoint as an observed Slack CLI contract, keep it isolated, and provide a manual fallback when Slack changes or rejects it.

## Evidence from Slack CLI

The Slack CLI source at `/home/manuel/code/others/slack-cli` is Apache-2.0 open source. Its local installation path calls `DeveloperAppInstall` in `internal/api/app.go`. The relevant source names are:

- `appDeveloperInstallMethod = "apps.developerInstall"`.
- `DeveloperAppInstallResult.APIAccessTokens.Bot`, `.AppLevel`, and `.User`.
- `internal/pkg/apps/install.go`, which exports the returned app-level token as `SLACK_APP_TOKEN` and the bot token as `SLACK_BOT_TOKEN`.

The observed wire contract is:

```http
POST https://slack.com/api/apps.developerInstall
Authorization: Bearer <management access token>
Content-Type: application/json

{"app_id":"A123","bot_scopes":["chat:write","commands"],"outgoing_domains":[]}
```

The response used by the CLI has this shape:

```json
{
  "ok": true,
  "app_id": "A123",
  "api_access_tokens": {
    "bot": "xoxb-…",
    "app_level": "xapp-…",
    "user": "xoxp-…"
  }
}
```

The local command only needs `bot` and `app_level`. A user token, when present, is ignored and never printed. The request's `bot_scopes` are derived from the existing `Manifest` generator, so app creation and installation use the same declared permissions. This avoids a second scope configuration that can drift from the bot code. Slack CLI omits `team_id` for standalone workspace apps; this command still requires the ID to choose the local installation record, but does not send it in the request.

The public references explain adjacent, supported APIs and behavior: [Slack CLI app install](https://docs.slack.dev/tools/slack-cli/reference/commands/slack_app_install/), [Slack CLI environment variables](https://docs.slack.dev/tools/slack-cli/guides/using-environment-variables-with-the-slack-cli/), [apps.manifest.create](https://api.slack.com/methods/apps.manifest.create), [OAuth v2 access](https://api.slack.com/methods/oauth.v2.access), and [Socket Mode connections.open](https://api.slack.com/methods/apps.connections.open). None of the public API pages documents `apps.developerInstall`; that absence is a material limitation of this feature.

## Why this is not browser OAuth

Browser OAuth is the correct flow when another person installs an app, when a service must support arbitrary workspaces, or when an app is distributed publicly. It requires a client ID and secret, a registered HTTPS redirect URI, a state value, a browser authorization step, and a server-side code exchange at `oauth.v2.access`. The callback server must validate state and retain the installation's team and bot token.

The local developer-install path already has a management token and a concrete team ID. Slack CLI performs the installation directly and receives both tokens, so adding a callback server would introduce work without solving a local requirement. Socket Mode still needs an app-level token; the observed developer-install response supplies it, while the normal dashboard flow asks a developer to create it manually. If this private method stops working, use `slack app install`, the Slack app dashboard, or a documented OAuth installation flow and then run `credentials import-runtime`.

## Credential and profile mapping

The existing store has three layers. A profile names the management record and app record and points at one selected installation. An app record owns registration metadata and its app-level Socket Mode token. An installation record owns a bot token for one app/workspace pair.

```text
profile go-go-golems
  management  -> go-go-golems (management access token)
  app         -> ping          (app_id A123, app token xapp-…)
  installation -> go-go-golems-T123
                                  (team T123, bot token xoxb-…)
```

For `bots install`, the implementation performs these checks before making a request:

- A profile is selected explicitly or through the configured default.
- The profile references an app and the app record has a non-empty app ID.
- A team ID is supplied. It is not guessed from a token prefix.
- The management record has an access token.

After a successful response it creates or replaces the deterministic installation key `<profile>-<team-id>`, links it from the profile, stores the bot token there, and stores the app-level token on the app record. Existing records for unrelated profiles and workspaces remain untouched. A known installation with a different app or team is rejected before replacement.

## CLI contract

```sh
slack-bot bots install ping \
  --profile go-go-golems \
  --team-id T0C1UJMCPGA \
  --config-dir ~/.config/go-go-slack
```

The command emits only safe metadata:

```json
{
  "profile": "go-go-golems",
  "installation": "go-go-golems-T0C1UJMCPGA",
  "app_id": "A123",
  "team_id": "T0C1UJMCPGA"
}
```

It must never include token values in output, errors, logs, or test diagnostics. `--timeout-ms` bounds the request. An unsuccessful HTTP response, Slack `ok:false`, a malformed response, a missing returned token, or a mismatched app ID is an ordinary command error. There is no automatic retry because a timeout can leave Slack's installation state uncertain; the user can inspect the app or run the command again after deciding whether an installation occurred.

## Implementation outline

```text
resolve bot descriptor and manifest
load config.yaml and credentials.json
resolve selected profile and its management credential
require app_id and explicit team_id
scopes = Manifest(descriptor).oauth_config.scopes.bot
POST developerInstall with bearer management access token (omit team_id for standalone apps)
decode bounded response
require ok, matching app_id, bot token, app-level token
installation = profile + "-" + team_id
update config and credentials maps
save both files with existing private atomic writer
emit profile, installation, app_id, team_id
```

The code belongs in `pkg/slackcli/install_app.go`, with command registration and settings in `pkg/slackcli/commands.go`. It reuses `internal/slackconfig.Store`, `slackconfig.ResolveProfile`, `Manifest`, and `appHTTPClient` from `pkg/slackcli/create_app.go`. No new dependency or transport abstraction is needed.

## Tests and review boundary

Tests use the existing injected `http.RoundTripper`; they never contact Slack. The fake server asserts the URL, bearer header, JSON fields, manifest-derived scopes, and one-request behavior. Success tests reload the store and verify the app token and installation bot token are separated and persisted. Failure tests cover missing profile/app/team, HTTP and Slack errors, app-ID mismatch, missing returned tokens, and secret-free output. A test also proves an existing installation with a different workspace or app is not overwritten.

The feature is complete when the focused package tests, offline build/vet, CLI help, and `docmgr doctor` pass. The endpoint remains an implementation detail with an explicit fallback; public OAuth support is a separate future ticket.

# Local app installation with Slack CLI developerInstall

## Executive Summary

<!-- Provide a high-level overview of the design proposal -->

## Problem Statement

<!-- Describe the problem this design addresses -->

## Proposed Solution

<!-- Describe the proposed solution in detail -->

## Design Decisions

<!-- Document key design decisions and rationale -->

## Alternatives Considered

<!-- List alternative approaches that were considered and why they were rejected -->

## Implementation Plan

<!-- Outline the steps to implement this design -->

## Open Questions

<!-- List any unresolved questions or concerns -->

## References

<!-- Link to related documents, RFCs, or external resources -->
