---
Title: Slack UI DSL
Slug: slack-ui-dsl
Short: Build Block Kit messages and connect buttons to modal submissions with the native Slack UI module.
Topics:
- slack
- javascript
- block-kit
- ui
Commands:
- bots
Flags:
- bot-repository
- event-file
- profile
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

The native `require("slack/ui")` module constructs Block Kit JSON in the Go-hosted
JavaScript runtime. The companion `require("slack")` module registers handlers
and supplies operations such as opening a modal and acknowledging its submission.
Building a value has no network effect. Returning a message from a handler,
calling `ctx.reply`, or calling a service sends it.

Start with the checked-in example at `examples/slack-bots/ui-showcase/index.js`.
It demonstrates one complete interaction: a slash command sends a message, its
button opens a title editor, and Save validates and accepts the submitted title.

## Run the showcase

The live runner needs an installed app and runtime credentials in a profile.
Follow `slack-bot help slack-bot-guide` to create and install one. From the
repository root, inspect and start the example:

```sh
go run ./cmd/slack-bot bots inspect ui-showcase
go run ./cmd/slack-bot bots run ui-showcase --profile go-go-golems --log-level debug
```

Replace the profile name with your own. Run the long-lived process in tmux.
Startup uploads the generated manifest to that profile's app by default, then
connects to Socket Mode. Switching from ping replaces its slash-command
registration with `/ui-showcase`; starting another JavaScript bot is not an
additional Slack app installation. `--skip-manifest-update` deliberately keeps
the existing remote configuration.

In Slack, invoke `/ui-showcase`, click **Acknowledge**, edit the title, and press
**Save**. A title shorter than three characters produces a field error.
A valid title closes the modal and writes a log entry. The current example
does not persist the title or update the original message. Acceptance is an
acknowledgment, not a storage operation.

## Message construction

A message combines fallback text with an array of blocks. The fallback text
remains mandatory when blocks are present. The host currently accepts 1–4000
characters and at most 50 message blocks.

```javascript
const ui = require("slack/ui");

const message = ui.message("UI showcase: choose an action")
  .block(ui.header("Slack UI showcase"))
  .block(ui.section(ui.mrkdwn("Edit the title with the button below.")))
  .block(ui.actions(
    "showcase-actions",
    ui.button("showcase.ack", "Acknowledge").value("showcase")
  ))
  .build();
```

Return this object from a command handler to send its implicit reply, or call
`await ctx.reply(message)`. Do one or the other: sending both fails with
`already_replied`. Slash-command replies use the response URL and are ephemeral,
so only the invoking user sees them.

For an explicit channel message, call
`await ctx.slack.messages.post({channelId: ctx.channelId, ...message})`.
Posting returns a channel ID and timestamp and does not consume the implicit
reply slot. The bot needs permission to post in that channel.

Builders produce detached JSON snapshots. Blocks can also be supplied as raw
objects, which allows unsupported builder shapes to be represented. Local
validation is partial: it checks basic structure and selected limits, rather
than the complete Slack schema. A raw block rendering successfully does not
add event routing or submission helpers for its interactive elements.

## Connect the button to a modal

The button's action ID selects the handler. Its block ID identifies the parent
layout block; its value is application data. Use stable IDs so emitted messages
continue to match registered handlers.

The following registration fragment belongs inside a `defineBot` callback
that destructures `action` and `view`:

```javascript
action("showcase.ack", async ctx => {
  await ctx.openModal(
    ui.modal("showcase.edit", "Edit showcase")
      .metadata("showcase")
      .input("title", "Title",
        ui.textInput("title_input").initial("Acknowledge"))
      .submit("Save")
      .build()
  );
});

view("showcase.edit", async ctx => {
  const title = ctx.values.text("title", "title_input");
  if (!title || title.trim().length < 3) {
    return ctx.ack.errors({title: "Use at least three characters."});
  }
  await ctx.ack.accept();
  ctx.log.info(`Saved ${title.trim()}`);
});
```

`ctx.openModal` uses the incoming action's trigger ID, kept in the Go-owned
invocation. Call it promptly from the action handler. The current method cannot
open a modal from a command or arbitrary background callback. It returns the
view's ID and hash.

The modal's callback ID selects `view("showcase.edit", ...)`. Submitted state is
indexed by input block ID and element action ID, hence
`ctx.values.text("title", "title_input")`. Missing or non-string values return
`undefined`. `ctx.values.all` exposes the detached state map for other decoding.

## Acknowledgment and state

An ACK tells Slack that an interaction was handled. Ordinary block actions get
an automatic empty ACK; their JavaScript handler can then open a modal or reply.
View submissions instead require exactly one explicit choice:

| Method | Effect |
| --- | --- |
| `await ctx.ack.accept()` | Accept the submission and close the modal. |
| `await ctx.ack.errors({blockId: "Explanation"})` | Keep the modal open and attach validation errors to input blocks. |

Error-map keys are block IDs, not element action IDs. Returning without choosing
an ACK produces `ack_required`. The Go receipt rejects duplicate or late choices.
Its deadline is three seconds from receipt creation; queueing and handler work
consume that time. Validate and choose the ACK before slow application work.
An accepted submission cannot subsequently be changed into field errors.

`ctx.view` contains `callbackId`, `privateMetadata`, `id`, `hash`, and
interaction `type`. Metadata can carry an application record key; it is not
persistent storage. `ctx.store` holds JSON values in memory, scoped to the host
and workspace, and loses them when the process restarts.

A modal submission has no automatic channel or response URL. Do not return a
message and expect it to appear in the original conversation. A follow-up
requires application context and an explicit supported service operation.
Message editing is not currently exposed.

## Builder reference

Chainable builders use `.build()` to produce JSON. APIs accepting builders
also convert them internally; explicit `.build()` is clearest at send boundaries.

| Entry point | Result and supported methods |
| --- | --- |
| `plain(text)` | Plain-text composition object. |
| `mrkdwn(text)` | Slack-formatted text object. |
| `header(text)` | Header block with plain text. |
| `divider()` | Divider block. |
| `section(textObject)` | Section containing a text object. |
| `button(actionId, label)` | Builder: `.value(text)`, `.style("primary" \| "danger")`, `.build()`. Default style is primary. |
| `actions(blockId, ...buttons)` | Actions block; accepts 1–25 button builders or built button objects. |
| `message(fallbackText)` | Builder: `.block(block)`, `.build()`; at most 50 blocks. |
| `textInput(actionId)` | Builder: `.initial(text)`, `.placeholder(text)`, `.required()`, `.optional()`, `.build()`. |
| `input(blockId, label, input)` | Input block containing the supplied text input. |
| `modal(callbackId, title)` | Builder: `.metadata(text)`, `.input(blockId, label, input)` or `.input(block)`, `.submit(text)`, `.close(text)`, `.build()`; at most 100 blocks. |

Known limitation: the current `textInput().optional()` builder emits
`optional` on the element rather than its containing input block. Avoid it
for live Slack views until corrected. For an optional field, create a normal
input block with `ui.input(...)`, set `block.optional = true`, and pass that
block to `modal.input(block)`.

Static-select actions are decoded and routed, but there is no static-select
builder yet, and `ui.actions` currently accepts buttons only. App Home,
shortcuts, external option loading, message updates, view updates/pushes, and
replacement-view ACK responses are not exposed by this API. The research
guide discusses these as designs or follow-ups, not current methods.

## Offline interaction tests

Simulation consumes normalized invocation JSON, not a raw Slack Socket Mode
envelope. It records operations without Slack credentials. The following
fixtures exercise the checked-in showcase in separate fresh runtimes:

```sh
go run ./cmd/slack-bot bots simulate ui-showcase \
  --event-file examples/slack-bots/fixtures/ui-command.json
go run ./cmd/slack-bot bots simulate ui-showcase \
  --event-file examples/slack-bots/fixtures/ui-action.json
go run ./cmd/slack-bot bots simulate ui-showcase \
  --event-file examples/slack-bots/fixtures/ui-view.json
```

The command fixture records an `ephemeral_reply` with blocks. The action
fixture supplies a synthetic trigger ID and records `open_view`. The view
fixture supplies nested form state and records `ack` with `kind: "accept"`.
Change its title to one character to exercise `kind: "errors"`.

Each simulate invocation starts with empty in-memory state. A synthetic
trigger or view ID proves local dispatch and serialization; it cannot be
used in the live Slack API.

## Troubleshooting

| Symptom | Check |
| --- | --- |
| Slack rejects `/ui-showcase` as unknown | Correct workspace, installed app, and successful startup manifest sync. |
| “No registered Slack handler” | The command/action/callback ID must match the selected script's registration. |
| Modal never opens | Invoke its button and open promptly; the operation needs the incoming action's trigger. |
| Save closes the modal without another message | Expected for the current showcase; acceptance only closes the dialog. |
| Field errors do not appear | Use input block IDs and choose the ACK before the deadline. |
| `ack_required`, `ack_already_sent`, or `ack_expired` | Choose exactly one ACK promptly in the view handler. |
| Slack rejects a raw block | Check Slack's schema; local validation is not a full schema validator. |

With `--log-level debug`, inspect incoming command/envelope IDs, dispatch
results, and outgoing reply fingerprints in tmux. Message text and response
URLs are not included in the host's reply diagnostics.

## See Also

- `slack-bot help slack-bot-guide` — credentials, app lifecycle, live execution and Go embedding.
- `examples/slack-bots/ui-showcase/index.js` — executable end-to-end example.
- `examples/slack-bots/slack.d.ts` — editor-facing method declarations.
- `internal/jsslack/ui_module.go` — Go implementation of the builders.
- `internal/jsslack/dispatch.go` — context methods and reply/ACK dispatch.
- Ticket `SLACK-UI-001`, `design-doc/01-slack-surfaces-and-ui-dsl-intern-guide.md` — deeper design, archived Slack references, diagrams and future surfaces.


## Native controls for example ports

`ui.input(blockId, label, element, {optional: true, hint: "Details"})` puts
optionality on the input block. Inputs are required by default. The former
`textInput().optional()` and `.required()` methods have been removed because
optionality is not an element property. The modal's three-argument `input`
method accepts the same fourth options argument. `textInput` supports
`.multiline()` and `.length(min, max)` for constraints up to 3,000 characters.

Select helpers return detached Slack elements and take an action ID plus an
options object using Slack wire keys:

```javascript
const choices = ui.staticSelect("article", {
  placeholder: "Choose an article",
  options: [ui.option("Getting started", "start"), ui.option("FAQ", "faq")]
});
const controls = ui.actions("browse", choices, ui.usersSelect("owner"));
```

Available helpers are `staticSelect`, `multiStaticSelect`, `externalSelect`,
`multiExternalSelect`, `usersSelect`, `multiUsersSelect`, `channelsSelect`,
`multiChannelsSelect`, `conversationsSelect`, `multiConversationsSelect`,
`datePicker`, `timePicker`, `checkboxes`, `radioButtons`, and `overflow`.
These helpers construct payloads; external option loading additionally requires
runtime support and a registered options handler. Construction alone does not
establish that routing capability.

`ui.section(text, accessory)` accepts an optional element. `ui.context(...elements)`
accepts up to ten text/image elements, and `ui.image(httpsURL, altText)` creates
an image. `ui.linkButton(id, label, httpsURL)` adds a URL button.
`ui.confirm(element, title, text, acceptLabel?, cancelLabel?)` adds a native
confirmation dialog. Buttons are neutral by default; use `.style("primary")`
or `.style("danger")` explicitly. `modal.block(block)` adds non-input blocks.

The helpers enforce a few common shape and count constraints, not the complete
Slack schema. Slack still validates surface compatibility and all field limits.
