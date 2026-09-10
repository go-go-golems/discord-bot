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

The separate `slack-bot` binary uses the existing Go module and dependencies. It hosts JavaScript through go-go-goja and `require("slack")`. This release implements offline behavior and a bounded Go ingress seam. It does not open Socket Mode, call Slack, install an app, or load credentials.

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

Use `--bot-repository PATH` to select another repository, `--timeout-ms 5000` to bound inspection and each invocation, and `--log-level debug` for lifecycle diagnostics. Commands emit one JSON document. `--bot-config-file PATH` on simulate reads a JSON object, for example `{"greeting":"hello"}`. Only declared fields reach `ctx.config`; no environment configuration is loaded. Fixture/config files must contain one JSON value and fit within 1 MiB.

Repositories contain root `.js` entries or immediate child `index.js` entries. Helpers below a bot directory are not discovered. A candidate entry executes during inspection; scripts are trusted local code, not an untrusted-code sandbox. The runtime exposes the Slack registration module and local CommonJS imports, without outbound host capabilities during inspection. Required runtime config does not prevent inspection.

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

`slackbot.NewIngress` accepts a dispatcher, workspace/app policy, queue capacity, dedupe capacity and TTL. `Admit` takes a detached envelope and an `Acknowledger`. Receipt does not enter JS, duplicate events are keyed by workspace/event ID, and full queues/dedupe caches return `busy`. The worker uses host lifetime, not the ACK context. Shutdown cancels pending work. This is best-effort memory admission, not durable delivery. It still needs a real transport decoder, SDK connection, rate-limit handling and reconnect tests.

## Validation and remaining work

Run `go test ./...`, `go build ./...`, `go vet ./...`, and race tests for the Slack packages. On the development workstation the ticket's `scripts/04-go-offline.sh` selects the cached matching Go 1.26.4 toolchain, disables downloads and bypasses the mismatched parent workspace.

Live Socket Mode, HTTP/WebSocket protocol fixtures, SDK integration, rate-limit/reconnect behavior, buttons, modals and xgoja providers are not implemented. See the ticket's intern guide and local-testing plan for those phases. Offline tests require no Slack tokens or test-message authorization.
