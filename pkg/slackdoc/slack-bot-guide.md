---
Title: Slack bot development guide
Slug: slack-bot-guide
Short: Create, install, run, and test JavaScript Slack bots with the Go host.
Topics:
- slack
- javascript
- development
Commands:
- bots
- credentials
Flags:
- profile
- skip-manifest-update
IsTopLevel: true
ShowPerDefault: true
SectionType: GeneralTopic
---

The `slack-bot` binary hosts JavaScript bots through go-go-goja and
`require("slack")`. It supports app creation, developer installation, stored
credentials, live Socket Mode, and offline inspection and simulation. Use
`bots run` for Slack and `bots run-local` for a prepared loopback mock.
Neither inspection nor simulation requires Slack credentials.

For Block Kit builders, buttons and modal submissions, read
[Slack UI DSL](slack-ui-dsl.md), also available as `slack-bot help slack-ui-dsl`.

Run these commands from the repository root:

```sh
go run ./cmd/slack-bot bots list
go run ./cmd/slack-bot bots inspect ping
go run ./cmd/slack-bot bots manifest ping
go run ./cmd/slack-bot bots simulate ping \
  --event-file examples/slack-bots/fixtures/command.json
go run ./cmd/slack-bot bots simulate ping \
  --event-file examples/slack-bots/fixtures/mention.json
```

Use `--bot-repository PATH` to select another repository, `--timeout-ms 5000` to bound inspection and each invocation, and `--log-level debug` for lifecycle diagnostics. Inspection and simulation commands emit one JSON document. `--bot-config-file PATH` on simulate reads a JSON object, for example `{"greeting":"hello"}`. Only declared fields reach `ctx.config`; no environment configuration is loaded. Fixture/config files must contain one JSON value and fit within 1 MiB.

Repositories contain root `.js` entries or immediate child `index.js` entries. Helpers below a bot directory are not discovered. A candidate entry executes during inspection; scripts are trusted local code, not an untrusted-code sandbox. The runtime exposes the Slack registration module and local CommonJS imports, without outbound host capabilities during inspection. Required runtime config does not prevent inspection.

## Create a Slack app

`bots manifest NAME` prints the generated manifest without contacting Slack.
`bots create-app NAME` sends that same manifest directly to `apps.manifest.create`;
the official Slack CLI is not required:

```sh
go run ./cmd/slack-bot bots create-app ping \
  --config-token-file /tmp/access-token.txt \
  --credentials-file /tmp/ping-app-credentials.json \
  --timeout-ms 30000
```

Get the configuration access token at <https://api.slack.com/apps>, below the app
list: **Your App Configuration Tokens → Generate Token**. Select the workspace.
Use the access token, not its refresh token. It expires after 12 hours. The
optional local profile store and explicit rotation command are available. These credentials belong to a user
and workspace; adding a bot OAuth scope does not issue a configuration token.

The command prints the app ID, settings URL and OAuth authorization URL. Follow
the latter to install the app. Creation does not install it or generate a Socket
Mode token. The optional credentials file captures the returned app credentials
with mode 0600 and refuses to overwrite an existing file. Without that flag,
credentials are omitted from stdout and can be obtained from app settings.
The configuration token never enters JavaScript or the output.

Each successful invocation creates a new app. Requests are not automatically
retried. If the outcome is unknown, check the dashboard before trying again.
A reserved credentials file can remain empty after failure; choose a new path
after resolving the failure. Explicit token files do not depend on environment
variables or the Slack CLI login store.

The manifest is generated from the selected JavaScript descriptor.
`bots manifest NAME` only prints it. `bots run NAME` updates the existing app
before connecting, unless `--skip-manifest-update` is supplied. See
“Automatic manifest sync on startup” below for token and reinstall requirements.

## Local credential profiles

The optional local store keeps profile names in `~/.config/go-go-slack/config.yaml`
and secrets in `credentials.json` (both private files). It is intended for one
developer using one workstation. Import a configuration access/refresh pair,
then use the profile for app creation:

```sh
slack-bot credentials import-management \
  --profile ping-dev --management owner \
  --access-token-file /tmp/access-token.txt \
  --refresh-token-file /tmp/refresh-token.txt
slack-bot profiles list
slack-bot credentials status --profile ping-dev
slack-bot credentials refresh --profile ping-dev
slack-bot bots create-app ping --profile ping-dev
```

For a local developer installation, the `bots install` verb uses the same
undocumented method observed in the open-source Slack CLI. It needs the
management access token in the selected profile and the target workspace ID:

```sh
slack-bot bots install ping --profile ping-dev --team-id T_DEV
```

On success it stores the workspace bot token under an installation record and
the Socket Mode app token on the app record. It prints only profile, app ID,
installation name, and team ID. This is a local convenience, not Slack's
documented public OAuth installation API; it does not start an HTTPS callback
server. If Slack rejects the method, use `slack app install` or the app
dashboard/OAuth flow, then import the resulting tokens with
`credentials import-runtime`.

Refresh is explicit and replaces both tokens together. If it fails, import a
new pair. Runtime tokens are saved by `bots install` or imported explicitly; this store
does not run a daemon or contact Slack during status/list commands.

After manual installation, runtime tokens can be stored with:

```sh
slack-bot credentials import-runtime --profile ping-dev \
  --installation dev --team-id T_DEV \
  --bot-token-file /tmp/bot-token.txt --app-token-file /tmp/app-token.txt
```

## JavaScript authoring

```javascript
const {defineBot} = require("slack");
module.exports = defineBot(({configure, command, event}) => {
  configure({name: "ping", run: {fields: {
    greeting: {type: "string", default: "pong"}
  }}});
  command("/golem-ping", {description: "Check the bot"}, async ctx => {
    return {text: ctx.config.greeting};
  });
  event("app_mention", async ctx => {
    await ctx.reply({text: "Hello"});
  });
});
```

Registration must be synchronous. `defineBot` and `configure` can each be called once during loading; duplicate handlers and unsupported events fail inspection. Config fields have lowerCamelCase names, types `string`, `bool`, or `number`, and optional `default`, `required`, and `help` fields.

`ctx` contains detached `id`, `teamId`, `channelId`, `userId`, `command`, `text`, `event`, and `config` values. A mention's event snapshot includes `ts` and `threadTs` strings. Timestamp IDs never pass through floating-point numbers.

- `await ctx.reply({text})`: one implicit reply. Slash replies use the injected ephemeral responder and return `{delivered:true,via:"response_url"}`. Mention replies post in the root/existing thread and return `{channelId,ts}`.
- Returning `{text}` sends the implicit reply automatically. Returning a response after calling `reply` fails with `already_replied`. Return `undefined` after explicitly replying.
- `await ctx.slack.messages.post({channelId,text,threadTs?})`: explicit message operation, returning `{channelId,ts}`. It does not consume the implicit reply slot.
- `ctx.store.get(key)`, `set(key,value)`, `delete(key)`, `keys()`: JSON values, copied on read/write, scoped by host and workspace. Missing keys return `undefined`; keys are sorted. State is in memory and disappears on restart.
- `ctx.log.debug/info/warn/error(message)`: structured logs with bot and invocation IDs. Transport credentials are never part of the JS context.

The Slack surface is intentionally small. The runtime supports plain text and
validated low-level Block Kit message payloads, plus the native
`require("slack/ui")` helpers for fallback text, sections, actions, buttons,
headers, and dividers. The builders return detached JSON objects, so a bot can
inspect or extend them before sending:

```javascript
const ui = require("slack/ui");
return ui.message("A notification fallback")
  .block(ui.section(ui.mrkdwn("*A section*")))
  .block(ui.actions("actions", ui.button("note.edit", "Edit").value("n1")))
  .build();
```

Block action callbacks are normalized from Socket Mode and routed by
`action_id`. Register one with `action("note.edit", handler)`. Ordinary button
and static-select actions are acknowledged by the transport before dispatch;
the handler receives `ctx.action` with the action ID, block ID, value, and
selected-option data. An action with a trigger ID can call `ctx.openModal` with
`ui.modal(...)`, `ui.input(...)`, and `ui.textInput(...)`. Register the modal
submission with `view(callbackId, handler)`; the handler reads
`ctx.values.text(blockId, actionId)` and must choose exactly one of
`ctx.ack.accept()` or `ctx.ack.errors({...})` before its invocation ends. The
transport rejects late and duplicate ACKs. Home tabs, shortcuts, message
actions, attachments, files, rich-text blocks, and canvases remain outside the
current contract. The Discord UI DSL is inspiration for construction style
only because Slack's payloads and acknowledgment rules are different.

Text must contain 1–4000 characters. Message options reject unknown fields. Context operations fail after an invocation closes. Native network operations are promises, backed by a bounded errgroup (16 concurrent operations per invocation) and settled on the VM owner. Await the operations whose result matters; only the handler's return/rejection determines the JS invocation result. Unawaited side effects should not be used for critical work.

## Fixture format and recording service

The input is a normalized invocation, **not a raw Slack envelope**:

```json
{
  "id":"mention-1", "teamId":"T-DEVELOPMENT",
  "channelId":"C-TEST", "userId":"U-TEST",
  "event":"app_mention", "text":"hello",
  "ts":"1741234567.000123", "threadTs":"1741234567.000001"
}
```

Use exactly one of `command`, `event`, `action`, or `interaction`. A mention requires `ts`. Simulation records `post`, `ephemeral_reply`, `open_view`, and `ack` operations. The UI DSL guide includes action and submission fixtures. Fake message references such as `offline.000001` are deliberately recognizable and cannot be used as Slack message IDs. A manifest is a review artifact; generation does not validate or install it with Slack.

## Go embedding

```go
recorder := &slackbot.Recorder{}
host, err := slackhost.Load(ctx, "examples/slack-bots/ping/index.js",
    slackhost.Options{Messages: recorder, Timeout: 5*time.Second})
if err != nil { return err }
defer host.Close(context.Background())
err = host.Dispatch(ctx, slackbot.Invocation{
    ID: "test", TeamID: "T", ChannelID: "C", UserID: "U",
    Command: "/golem-ping",
}, recorder)
```

Import `pkg/slackbot` and `pkg/slackhost` from this repository's module. Inject context-aware `MessageService` and `Responder` implementations; keep tokens and response URLs private to those implementations. The public host serializes whole invocations, and a deadline interrupts CPU-bound JavaScript. Go services must honor context cancellation; a service that blocks forever cannot be forcibly stopped safely by the host.

`slackbot.NewIngress` accepts a dispatcher, workspace/app policy, queue capacity, dedupe capacity and TTL. `Admit` takes a detached envelope and an `Acknowledger`. Receipt does not enter JS, duplicate events are keyed by workspace/event ID, and full queues/dedupe caches return `busy`. The worker uses host lifetime, not the ACK context. Shutdown cancels pending work. This is best-effort memory admission, not durable delivery. The local and remote transports connect this ingress to the pinned Slack SDK. Rate-limit retry policy remains intentionally out of scope.

## Run against the installed Slack app

The `run` verb resolves the profile's installation and connects to Slack's
public Web API and Socket Mode endpoints:

```sh
go run ./cmd/slack-bot bots run ping --profile go-go-golems --log-level debug
```

The profile must have an app ID, workspace installation, bot token, and
Socket Mode app token. The command verifies the workspace with `auth.test`,
automatically acknowledges commands, mentions and block actions before handler completion. Modal submissions require the handler to choose `ctx.ack.accept()` or `ctx.ack.errors(...)`. Press Ctrl-C to cancel the connection. Use `--bot-config-file` for
declared bot configuration fields; credentials never enter JavaScript.

## Validation and remaining work

Run `go test ./...`, `go build ./...`, `go vet ./...`, and race tests for the Slack packages. On the development workstation the ticket's `scripts/04-go-offline.sh` selects the cached matching Go 1.26.4 toolchain, disables downloads and bypasses the mismatched parent workspace.

A pinned SDK/mock probe, a complete local CLI scenario, and baseline HTTP/WebSocket fixtures exercise real network encoding. Messages, buttons, modal submissions and ACK handling are implemented; broader reconnect/deployment coverage and Slack-specific xgoja providers remain follow-up work. See the ticket's intern guide and local-testing plan for those phases. Offline tests require no Slack tokens or test-message authorization.

## Run against a prepared local mock

```sh
go run ./cmd/slack-bot bots run-local ping \
  --local-connection-file /tmp/slack-host-probe/config.json \
  --log-level debug
```

The connection file contains `apiURL`, `botToken`, `appToken`, `teamID`, `appID`, optional `userID` (probe metadata), and optional `allowedChannels`. Use synthetic mock values, not real Slack credentials. `apiURL` must use HTTP, a literal loopback IP, an explicit port and a trailing slash, for example `http://127.0.0.1:12345/api/`. The mock launcher in `testdata/slack/mock/probe.ts` writes this file with mode 0600. No host connection field enters JavaScript. `--bot-config-file` supplies declared bot configuration separately.

All HTTP and WebSocket dials are restricted to that exact host and port. HTTP proxies are disabled and redirects are rejected. A response URL from any other origin is rejected. The runner verifies the workspace using `auth.test`, decodes mention, command, block-action and view-submission envelopes, acknowledges receipt independently of handler completion, and dispatches through bounded ingress to the actual JS host. SIGINT and SIGTERM cancel the process lifetime.

Posting performs one SDK request. A 429 returns `rate_limited`; unknown or ambiguous failures return `delivery_unknown` without automatic retransmission. This conservative baseline avoids duplicating messages after a lost response; it does not yet implement method-scoped pacing or retry waits. Invalid input is rejected before HTTP.

Run the separate SDK gate and full process scenario described in `testdata/slack/mock/README.md`. Ordinary Go tests skip the external-server probe unless explicitly configured, while HTTP/WebSocket fixture tests need permission to bind loopback sockets.

## Automatic manifest sync on startup

`slack-bot bots run ui-showcase --profile go-go-golems` updates the selected
Slack app with the bot's generated manifest before connecting to Socket Mode.
This replaces the app configuration, including its display name, slash commands,
events and scopes. The profile's stored management access token is required.
The app ID and existing runtime tokens are reused.

Use `--skip-manifest-update` to connect using the current Slack configuration
without a management token. If a management token has expired, run
`slack-bot credentials refresh --profile go-go-golems` and retry.

If Slack returns `permissions_updated: true`, startup stops with an install
command. Run it to grant the changed scopes, then start the bot again. Command-only
changes normally do not require reinstalling. There is no automatic rate-limit retry.


## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| Slash command is not recognized | Slack has an older or different manifest | Start the intended bot without `--skip-manifest-update`; use the installed workspace. |
| Startup says `token_expired` | Management access token expired | Run `credentials refresh` with the same profile/config directory and retry. |
| Startup requests reinstall | Manifest changed permissions | Run the printed install command, then restart. |
| Command reaches the process but has no handler | Selected script and Slack command differ | Check bot, script and command fields in debug logs. |
| Modal Save only closes the dialog | Handler accepts without another operation | Add application state or an explicit follow-up; ACK is not a save operation. |
| API posting returns `not_in_channel` | Bot lacks channel membership | Invite the bot to the target channel. |

## See Also

- `slack-bot help slack-ui-dsl` — interactive UI tutorial and API reference.
- `slack-bot bots run --help` — live startup flags.
- `examples/slack-bots/slack.d.ts` — JavaScript-facing declarations.


## Bot services and local SQLite

`ctx.slack` exposes named message, conversation, user, user-group, pin, reaction,
workspace and file operations. The authoritative method list is
`pkg/slackbot/operations.go`. Except for the existing `messages.post` API, these
methods accept Slack wire keys and return detached Slack JSON results. Example:

```javascript
const page = await ctx.slack.conversations.history({channel: ctx.channelId, limit: 100});
const cursor = page.response_metadata.next_cursor;
await ctx.slack.messages.update({channel: ctx.channelId, ts, text: "Updated", blocks: []});
```

Pagination is explicit: callers must continue with `cursor` until it is empty.
Rate limits reject with `rate_limited`; the local framework does not retry or
silently return partial history. Additional API scopes belong in configure's
`scopes` array and require reinstallation. Administrative capabilities still
depend on Slack's token and plan requirements; having a named method does not
grant permission. `usergroups.setMembers` replaces the whole membership list.

`files.upload({channel_id, filename, content, thread_ts?})` uploads a generated
UTF-8 file through getUploadURLExternal, a token-free content transfer, and
completeUploadExternal. Files are limited to 8 MiB. This method requires
files:write and channel access. The older files.upload Web API is not used.

`require("database")` exposes the existing go-go-goja SQLite module. Configure
it lazily inside the first handler with a declared dbPath setting:

```javascript
const db = require("database");
let ready = false;
function ensure(ctx) {
  if (ready) return;
  db.configure("sqlite3", ctx.config.dbPath);
  const result = db.exec("CREATE TABLE IF NOT EXISTS notes (id TEXT PRIMARY KEY, body TEXT)");
  if (!result.success) throw new Error(result.error);
  ready = true;
}
```

The host owns this module and closes its connection after the runtime shuts down.
Inspection rejects configure calls during registration, so list/manifest commands
cannot create database files. Use a separate database path for each bot; scripts
must scope their records by workspace where one file serves multiple workspaces.
The in-memory ctx.store remains suitable for transient per-workspace UI state.



## Local verbs and conditional administrative credentials

Unified Demo exposes synchronous local verbs as well as Slack handlers:

```sh
go run ./cmd/slack-bot bots invoke unified-demo status
go run ./cmd/slack-bot bots invoke unified-demo run
```

Register these with `verb(name, {description}, handler)`. The handler receives
only declared configuration and returns JSON. Local verbs receive no transport
services and must be synchronous. The `run` metadata verb returns
`host-managed`; connect the bot with the regular `bots run` command.

Enterprise workspace removal is available through `ctx.slack.admin.removeUser`.
It uses a separately authorized user token with `admin.users:write`, never the
bot token. The Moderation example additionally requires an authorized actor and
`enableWorkspaceRemoval: true`. This removes a workspace membership; it is not
an implementation of Discord ban/unban or timeout.

`credentials import-runtime` accepts an optional `--user-token-file`. The token
is stored privately with the selected installation, retained during subsequent
developer installations, and represented only by `has_user_token` in status.
The developerInstall flow does not issue this Enterprise administrative token;
obtain it through the separately authorized Slack Admin API installation.
Normal bot startup does not require it or request admin scopes in its manifest.

Message handlers subscribe to public/private channels, direct messages and group
DMs where the installed app has access. The manifest includes their corresponding
history scopes. `message_changed`, `message_deleted` and `user_change` handlers
are also supported. A deletion may have no actor ID; use `ctx.event.data` for
its deleted timestamp. Registering these handlers does not confer access to
conversations outside the app's membership and authorization.
