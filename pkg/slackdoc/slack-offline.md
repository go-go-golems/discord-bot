---
Title: Offline Slack bot development
Slug: slack-offline
Short: Inspect Slack bots, generate manifests, and replay invocation fixtures without credentials.
Topics:
- slack
- javascript
- offline
IsTopLevel: true
ShowPerDefault: true
SectionType: GeneralTopic
---

# Offline Slack bot development

The separate `slack-bot` binary uses the existing Go module and dependencies. It hosts JavaScript through go-go-goja and `require("slack")`. This release implements offline behavior, bounded ingress, and an explicit loopback-only Socket Mode runner. Inspection and simulation require no credentials or network. The local runner reads synthetic connection settings from a file; it cannot connect to real Slack.

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

Refresh is explicit and replaces both tokens together. If it fails, import a
new pair. Installation and runtime bot tokens remain manual; this store does
not run a daemon or contact Slack during status/list commands.

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

Use exactly one of `command` or `event`. A mention requires `ts`. Simulation returns recorded `post` and `ephemeral_reply` operations. Fake message references such as `offline.000001` are deliberately recognizable and cannot be used as Slack message IDs. A manifest is a review artifact; generation does not validate or install it with Slack.

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

Import `pkg/slackbot` and `pkg/slackhost` from this repository's module. Inject context-aware `MessageService` and `Responder` implementations; keep any future tokens and response URLs private to those implementations. The public host serializes whole invocations, and a deadline interrupts CPU-bound JavaScript. Go services must honor context cancellation; a service that blocks forever cannot be forcibly stopped safely by the host.

`slackbot.NewIngress` accepts a dispatcher, workspace/app policy, queue capacity, dedupe capacity and TTL. `Admit` takes a detached envelope and an `Acknowledger`. Receipt does not enter JS, duplicate events are keyed by workspace/event ID, and full queues/dedupe caches return `busy`. The worker uses host lifetime, not the ACK context. Shutdown cancels pending work. This is best-effort memory admission, not durable delivery. The local transport connects this ingress to the pinned Slack SDK. Rate-limit retry policy and reconnect tests remain pending.

## Validation and remaining work

Run `go test ./...`, `go build ./...`, `go vet ./...`, and race tests for the Slack packages. On the development workstation the ticket's `scripts/04-go-offline.sh` selects the cached matching Go 1.26.4 toolchain, disables downloads and bypasses the mismatched parent workspace.

A pinned SDK/mock probe, a complete local CLI scenario, and baseline HTTP/WebSocket fixtures now exercise real network encoding. External Slack connections, rate-limit retries, reconnect acceptance tests, buttons, modals and xgoja providers remain pending. See the ticket's intern guide and local-testing plan for those phases. Offline tests require no Slack tokens or test-message authorization.

## Run against a prepared local mock

```sh
go run ./cmd/slack-bot bots run-local ping \
  --local-connection-file /tmp/slack-host-probe/config.json \
  --log-level debug
```

The connection file contains `apiURL`, `botToken`, `appToken`, `teamID`, `appID`, optional `userID` (probe metadata), and optional `allowedChannels`. Use synthetic mock values, not real Slack credentials. `apiURL` must use HTTP, a literal loopback IP, an explicit port and a trailing slash, for example `http://127.0.0.1:12345/api/`. The mock launcher in `testdata/slack/mock/probe.ts` writes this file with mode 0600. No host connection field enters JavaScript. `--bot-config-file` supplies declared bot configuration separately.

All HTTP and WebSocket dials are restricted to that exact host and port. HTTP proxies are disabled and redirects are rejected. A response URL from any other origin is rejected. The runner verifies the workspace using `auth.test`, decodes mention/command envelopes, acknowledges receipt independently of handler completion, and dispatches through bounded ingress to the actual JS host. SIGINT and SIGTERM cancel the process lifetime.

Posting performs one SDK request. A 429 returns `rate_limited`; unknown or ambiguous failures return `delivery_unknown` without automatic retransmission. This conservative baseline avoids duplicating messages after a lost response; it does not yet implement method-scoped pacing or retry waits. Invalid input is rejected before HTTP.

Run the separate SDK gate and full process scenario described in `testdata/slack/mock/README.md`. Ordinary Go tests skip the external-server probe unless explicitly configured, while HTTP/WebSocket fixture tests need permission to bind loopback sockets.
