---
Title: Pragmatic implementation plan
Ticket: SLACK-CREDENTIALS-001
Status: active
Topics:
    - discord-bot
    - api-design
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Small single-user file store with named profiles, explicit refresh, and reuse of the existing manifest/create-app verbs."
LastUpdated: 2026-09-14T21:31:19.157446063-04:00
WhatFor: ""
WhenToUse: ""
---

# Pragmatic implementation plan

## Outcome and scope

After this ticket, a developer can import Slack configuration tokens once, select a named profile, and create an app from the existing bot manifest. When a token expires, they run a refresh command. If refresh fails, they generate another pair and import it. Installation stays in the browser.

Use two files under `~/.config/go-go-slack/`: `config.yaml` for names and selections, and `credentials.json` for secrets. Support several management identities, apps, and workspace installations using ordinary maps. There is no daemon, database, keychain integration, operation journal, background renewal, or recovery subsystem.

Assume one person runs credential-writing commands sequentially on one computer. Concurrent writers are unsupported. A clear error and manual repair are acceptable outcomes for interrupted operations. The software should still avoid obvious mistakes: wrong profile selection, secret output, partial JSON writes, and automatic duplicate app creation.

## Understand the three credential groups

A local bot such as `ping` is JavaScript code. A Slack app is its remote registration. Installing that app in a workspace grants a bot token for that workspace. These objects have different identifiers and need different credentials.

- **Management identity:** configuration access/refresh tokens for a user and workspace. Used to create apps. One identity can manage several apps.
- **App:** app ID, OAuth client credentials, signing secret, and optional Socket Mode app token. These belong to the registration.
- **Installation:** bot token for one app in one workspace. The same app can have different bot tokens in different workspaces.

A profile selects these records by friendly names. It contains references rather than copies of shared secrets. Keep Slack IDs next to their records and compare known app/workspace IDs when linking an installation. Do not invent a workspace from a token prefix.

```text
profile: ping-dev
   management -> manuel-dev -> configuration token pair
   app -------> ping-dev ---> app ID + app credentials
   installation -> dev -----> workspace ID + bot token

profile: another-bot
   management -> manuel-dev   (same saved token pair)
   app -------> another-bot   (different registration)
```

## Current code to reuse

The following files exist in the repository or current working tree. The profile and storage commands described later are proposed, not available yet.

| File | Relevant behavior |
|---|---|
| `pkg/slackcli/commands.go` | Glazed fields, command dispatch, shared `Manifest` generator |
| `pkg/slackcli/create_app.go` | Direct creation API, explicit access-token file, private optional app-credential output |
| `pkg/slackcli/create_app_test.go` | Synthetic HTTP request tests and output checks |
| `cmd/slack-bot/main.go` | Root command registration |
| `pkg/slackdoc/slack-offline.md` | User-facing embedded help |
| `internal/slacktransport/client.go` | Existing loopback-only runtime; unchanged by this ticket |

`bots manifest ping` inspects the bot and prints Slack configuration without API traffic. `bots create-app ping --config-token-file PATH` uses the same generator and sends a direct Slack request. The official Slack CLI is not required. Preserve this shared generator when adding profile authentication.

## File format

Create the configuration directory with mode `0700` and both files with mode `0600`. Expose `--config-dir` for tests and alternate directories. Use the user's home-directory configuration location as the default; document how it is resolved. Do not discover tokens through environment variables or the Slack CLI's private login store.

This example is proposed `config.yaml`. IDs are illustrative, not real Slack values:

```yaml
default_profile: ping-dev
profiles:
  ping-dev:
    management: manuel-dev
    app: ping-dev
    installation: dev
management:
  manuel-dev:
    team_id: T_DEV
    user_id: U_MANUEL
apps:
  ping-dev:
    app_id: A_PING
installations:
  dev:
    app: ping-dev
    team_id: T_DEV
```

Before app creation, a profile may have only its management reference. Add the app reference after successful creation. Installation is optional until a bot token has been imported. Require `--profile` or an explicitly configured default; never choose a profile by map iteration order.

The corresponding `credentials.json` contains all secrets in one document:

```json
{
  "management": {
    "manuel-dev": {
      "access_token": "REDACTED",
      "refresh_token": "REDACTED",
      "expires_at": "2026-09-15T00:00:00Z"
    }
  },
  "apps": {
    "ping-dev": {
      "client_id": "EXAMPLE",
      "client_secret": "REDACTED",
      "signing_secret": "REDACTED",
      "app_token": "REDACTED"
    }
  },
  "installations": {
    "dev": {"bot_token": "REDACTED"}
  }
}
```

Expiry is optional when importing a token pair without metadata. Never guess it from a file timestamp. Status should say “expiry unknown” in that case. Refresh stores the actual returned expiry. This ticket supports manually imported non-rotating bot tokens; OAuth bot-token rotation is separate future work.

Use concrete Go structs and maps. A small `internal/slackconfig` package with `Load`, `Save`, and profile resolution is enough. It does not need a pluggable backend interface, record generations, or opaque local IDs. Friendly names are map keys inside fixed files, not filesystem paths.

## Commands and implementation steps

These command forms define the proposed interface. Implement them with Glazed fields and the existing root command patterns.

### 1. Import a management identity and create its profile

```sh
slack-bot credentials import-management \
  --profile ping-dev --management manuel-dev \
  --access-token-file /tmp/access-token.txt \
  --refresh-token-file /tmp/refresh-token.txt
```

Read and trim both files, reject empty input, and save the pair together. The explicit import command may replace the named management pair; it should say which name was updated, without printing values. Create the profile if absent, and set its management reference. Preserve unrelated records. Do not delete input files.

Team/user IDs can remain unknown on initial import. Store them when a rotation response supplies them. User-entered labels are names, not verified identities. Do not add a separate identity-verification protocol to bootstrap storage.

### 2. List profiles and show safe status

```sh
slack-bot profiles list
slack-bot credentials status --profile ping-dev
```

Print profile names, selected management/app/installation names, known Slack IDs, expiry, and whether required credentials are present. Status has no network side effects. Do not serialize the secret structs into CLI output or debug logs.

### 3. Refresh explicitly

```sh
slack-bot credentials refresh --profile ping-dev
```

Call `tooling.tokens.rotate` with the saved refresh token. Check HTTP and Slack's `ok` field, then require both replacement tokens. Save both in one update along with expiry and returned identity fields. If known workspace/user IDs conflict with the response, report the mismatch instead of silently relabeling the profile.

```text
load profile and management pair
result = rotateOnce(refreshToken, boundedTimeout)
if failure:
    return "Refresh failed; generate and import a new pair"
require replacement access and refresh tokens
update pair and returned metadata
save credentials, then config
print management name and expiry
```

No automatic refresh loop or retry machinery is needed. If the response is lost or saving fails, tell the user to import a new pair. This is acceptable recovery for this local tool.

### 4. Create an app using a profile

```sh
slack-bot bots manifest ping
slack-bot bots create-app ping --profile ping-dev
```

Reject simultaneous profile and explicit token-file inputs. Reject profile-based creation when that profile already references an app, explaining how to use a new profile for another registration. Keep the existing file-based path available as its explicit alternative.

Load the saved management access token, generate the existing manifest, and submit it once. No separate remote validation pass is necessary; creation already reports invalid manifests. If the token has expired, tell the user to run refresh. Save returned app credentials under the profile name, then add the app ID and profile reference to config. Refuse an existing destination app entry before making the request.

```text
load selected profile; require no existing app reference
manifest = Manifest(inspectBot(name))
body = {manifest: JSON.stringify(manifest)}
result = createOnce(savedAccessToken, body)
if accepted:
    save returned app credentials
    save app ID and profile reference
    print app ID, settings URL, installation URL
```

The API argument wraps the generated manifest as a serialized string. Preserve this request shape from `create_app.go`. On timeout or an unreadable response, tell the user to inspect Slack's app dashboard before retrying. If creation succeeded but local saving failed, include the app ID in the error; the user can recover through app settings. No journal, automatic reconciliation, or replay is required.

### 5. Save manually obtained runtime credentials

```sh
slack-bot credentials import-runtime \
  --profile ping-dev --installation dev \
  --team-id T_DEV \
  --bot-token-file /tmp/bot-token.txt \
  --app-token-file /tmp/app-token.txt
```

Require an existing app reference. Store the app token with that app and the bot token with the named installation. Check an existing installation's app/workspace before replacing it. Input IDs are explicit configuration; do not claim they were verified with Slack unless a check actually ran.

Installation remains manual: follow the OAuth URL, approve in the intended workspace, and obtain the bot token and Socket Mode token through Slack settings. This ticket only stores them. It does not enable external connections in `run-local`, implement a callback server, or send test messages.

## Saving files and handling errors

Read the current complete document, modify its maps, write a new temporary file in the same directory with mode `0600`, close it, and rename it over the destination. Use this simple helper for both files. Saving the entire credentials document keeps access and refresh tokens together and preserves other profiles.

Do not add locks or `fsync` machinery. Document that credential-writing commands should not run concurrently. The two files are not a transaction: save secrets before references, and surface any failure. Manual repair or re-import is sufficient if an interruption leaves them inconsistent. Reject malformed existing files rather than overwriting them with empty defaults.

The important failures are wrong selections, unreadable/empty token files, invalid JSON, failed writes, expired credentials, missing API permission, invalid manifests, and ambiguous creation responses. Other rare failures can be ordinary contextual errors. Secret values and raw credential-bearing API responses must stay out of those errors.

## Tests and definition of done

Use temporary directories and the existing injected HTTP transport. No Slack credentials are needed for automated tests.

1. Import two management identities and profiles; reload and select each correctly without losing the other.
2. Check new directory/file permissions and reject malformed existing JSON without overwriting it.
3. Refresh with a fake response; both tokens change together and status contains no secrets. Verify a rejected refresh leaves the previous file intact.
4. Create through a profile; assert the manifest wrapper and selected bearer token, save returned credentials, and print only safe fields.
5. Cover unknown profile, conflicting auth flags, existing app, expired token, invalid manifest, and a lost creation response. Assert no automatic second create request.
6. Import runtime tokens for different installations; verify separation and reject a known app/workspace mismatch.

Done means the documented command sequence works with fake responses, focused CLI/storage tests pass, help explains manual installation and recovery, and changes are committed with a diary checkpoint. A real app creation can be a separately recorded smoke test; it is not required to validate file storage. Do not expand the test suite into crash simulation or distributed concurrency testing.

## Source references

See [sources/README.md](../sources/README.md) for the relevant official API contracts already extracted with defuddle in the parent ticket. Reuse that archive. No new research or dependency is needed to begin this implementation.
