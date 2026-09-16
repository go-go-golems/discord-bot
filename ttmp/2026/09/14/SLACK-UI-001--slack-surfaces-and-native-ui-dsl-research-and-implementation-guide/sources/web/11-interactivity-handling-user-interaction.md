---
Title: "11 interactivity handling user interaction"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/interactivity/handling-user-interaction/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

User interactions can blossom forth from the [seeds planted](https://docs.slack.dev/interactivity) in Slack apps. Cultivate the healthy growth of these interactions by preparing your app to understand and respond to them.

This guide will explain the nuances of user-triggered interactions and the steps necessary to handle the contextual interaction information your app will receive.

---

## The flow of user-triggered interactions

We explained in [invocation & interactivity](https://docs.slack.dev/interactivity) that interactions are essentially a trigger and a response. Apps can implement a number of interaction [entry points](https://docs.slack.dev/interactivity) that allow users to intentionally invoke a response from the app. Read our overview of those entry points to dig deeper into the field of options.

When one of those entry points is triggered, a new aspect is introduced to the interaction transaction — the interaction payload. This payload is a bundle of information that explains the context of the user action, giving the app enough information to construct a coherent response. We end up with an interaction flow that looks like this:

- A user **triggers** an interaction by engaging with one of your app's entry points.
- Your app **receives and processes** the interaction payload.
- Using this context, your app generates a **response** to the interaction.

Your app needs to be ready for the last two steps.

---

## Preparing your app for user interactions

In order for your app to receive interaction payloads, Slack needs to know where to send them. Each app can be configured with **Request URLs** that indicate a web endpoint belonging to the app.

To configure a **Request URL** for your app:

- Open the [App Management](https://api.slack.com/apps) dashboard.
- Select **Interactivity & Shortcuts**.
- Toggle **Interactivity** on.

You'll see a few new options appear. The ones relevant to us are:

- **Request URL**: the URL we'll send the request payload to when [interactive components](https://docs.slack.dev/block-kit#making-things-interactive) or [shortcuts](https://docs.slack.dev/interactivity/implementing-shortcuts) are used. You'll need to set up a URL to handle these payloads, as [we'll describe below](#payloads). Save your changes after you've added one.
	- If you are [distributing your app](https://docs.slack.dev/app-management/distribution), this request URL needs to be an HTTPS URL (self-signed certificates are not allowed). If you're just building a [single-workspace app](https://docs.slack.dev/app-management/distribution), it can be plain HTTP.
	- This **Request URL** is also used by [modals](https://docs.slack.dev/surfaces/modals) for [`view_submission` event payloads](https://docs.slack.dev/surfaces/modals#handling_submissions). Your app can distinguish between the different types of payload using the `type` field as [explained below](#payloads).
- **Options Load URL**: this is a setting used by some [Block Kit](https://docs.slack.dev/block-kit) interactive components, i.e. [select menus](https://docs.slack.dev/reference/block-kit/block-elements/select-menu-element) and [multi-select menus](https://docs.slack.dev/reference/block-kit/block-elements/multi-select-menu-element). These components can load menu options from an external source, and **Options Load URL** determines which URL is queried to return those menu options.
	- You can leave **Options Load URL** blank for now. We recommend reading the [reference guide for this component](https://docs.slack.dev/reference/block-kit/block-elements/select-menu-element).
	- If you create [**Slash Commands**](https://docs.slack.dev/interactivity/implementing-slash-commands) for your app, you'll find that each command has its own **Request URL**. Read our [guide to creating Slash Commands](https://docs.slack.dev/interactivity/implementing-slash-commands#creating_commands) to learn more.

You're all set! Your app can start receiving payloads. Now let's see how to process them.

---

## Handling interaction payloads

We mentioned that there were a few different types of interaction payloads your app might receive. They'll be sent to your specified **Request URL** in an HTTP POST request in the form `application/x-www-form-urlencoded`. For more information, refer to [Using the Slack Web API: Basics](https://docs.slack.dev/apis/web-api/#slack-web-api__basics). The body of the request will contain a `payload` parameter; your app should parse this `payload` parameter as JSON.

The resulting object can have different structures depending on the source. All those structures will have a `type` field that indicates the source of the interaction. Our reference docs have a more detailed look at the payload structures for the different `type` sources:

- [`block_actions` payloads](https://docs.slack.dev/reference/interaction-payloads/block_actions-payload) are received when a user clicks a [Block Kit interactive component](https://docs.slack.dev/block-kit#making-things-interactive).
- [`shortcut` and `message_actions` payloads](https://docs.slack.dev/reference/interaction-payloads/shortcuts-interaction-payload) are received when [global and message shortcuts](https://docs.slack.dev/interactivity/implementing-shortcuts) are used.
- [`view_submission` payloads](https://docs.slack.dev/reference/interaction-payloads/view-interactions-payload) are received when a [modal is submitted](https://docs.slack.dev/surfaces/modals#handling_submissions).
- [`view_closed` payloads](https://docs.slack.dev/reference/interaction-payloads/view-interactions-payload) are received when a [modal is canceled](https://docs.slack.dev/surfaces/modals#modal_cancellations).

In each case, the fields within the payload will provide a full context of the interaction. This will include the interacting user, the pre-defined state fields of any interactive component used, where the interaction happens, and more. Use this structure and these fields to interpret the request. You can use as much or as little of the info as your app needs. The payload types that your app can receive will depend on the other features your app implements. For example, if your app never publishes any interactive components, it will never receive a `block_actions` payload.

Read our [interaction payload reference docs](https://docs.slack.dev/reference/interaction-payloads) to examine the detailed structure for the features your app uses, and design the app to be able to process those fields.

Now that your app can receive and process interaction payloads, it needs to know what to do after one is sent.

---

## Responding to users

There are a bouquet of potential responses to the receipt of an interaction payload:

- [**Acknowledgment response**](#acknowledgment_response): a required, immediate response that confirms your app received the payload.
- [**Message responses**](#message_responses): optional responses that use a special URL from the payload to publish messages.
- [**Modal responses**](#modal_responses): optional responses that use a short-lived code from the payload to invoke a modal.
- [**Asynchronous responses**](#async_responses): optional responses that are inspired by the info in the payload, but not directly related.

### Acknowledgment response

All apps must, as a minimum, acknowledge the receipt of a valid interaction payload.

To do that, your app must reply to the HTTP POST request with an [HTTP 200 OK](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.2.1) response. This must be sent within 3 seconds of receiving the payload. If your app doesn't do that, the Slack user who interacted with the app will see an error message, so ensure your app responds quickly. Otherwise, the user won't see anything when your app only sends an acknowledgment response. If you want to do more, keep reading.

### Message responses

Depending on the source, the interaction payload your app receives may contain a `response_url`. This `response_url` is unique to each payload, and can be used to publish messages back to where the interaction happened.

Within these responses, you can include a [message payload](https://docs.slack.dev/messaging#payloads). This message payload can be composed according to the same [message composition guides](https://docs.slack.dev/messaging) as any other. You can use a `response_url` by making an HTTP POST directly to the URL and including a [message payload](https://docs.slack.dev/messaging#payloads) in the HTTP body.

If your app needs more than 30 minutes to respond with a message, you'll need to [publish it in the standard way](https://docs.slack.dev/messaging/sending-and-scheduling-messages).

✨ Go [here](https://docs.slack.dev/tools/bolt-js/concepts/actions) for more details about how to handle a `response_URL` in JavaScript, or [here](https://docs.slack.dev/tools/bolt-python/concepts/shortcuts) for more details about how to do this in Python.

Below are examples for different types of message responses.

#### Publishing ephemeral responses

By default, a message published via `response_url` will be sent as an [ephemeral message](https://docs.slack.dev/messaging#ephemeral):

```markdown
POST SLACK_REDACTED_EXAMPLE
Content-type: application/json
{
    "text": "Oh hey, this is a nifty ephemeral message response!"
}
```

#### Publishing responses in channel

If you want to publish a message to the same conversation as the interaction source, include an attribute `response_type` and set its value to `in_channel`.

```markdown
POST SLACK_REDACTED_EXAMPLE
Content-type: application/json
{
    "text": "Oh hey, this is a fun message in a channel!",
    "response_type": "in_channel"
}
```

#### Publishing responses in thread

If you want to publish a message to a specific thread, you'll need to include an attribute `response_type` and set its value to `in_channel`. Then, to specify the thread, include a `thread_ts`.

Also, be sure to set `replace_original` to `false` or you'll overwrite the message you're wanting to respond to!

```markdown
POST SLACK_REDACTED_EXAMPLE
Content-type: application/json
{
    "text": "Oh hey, this is a marvelous message in a thread!",
    "response_type": "in_channel",
    "replace_original": "false",
    "thread_ts": "1234567890"
}
```

#### Updating a source message in response

If your app received an interaction payload after an [interactive component](https://docs.slack.dev/block-kit#making-things-interactive) was used inside of a message, you can use `response_url` to update that source message.

Include an attribute `replace_original` and set it to `true`:

```markdown
POST SLACK_REDACTED_EXAMPLE
Content-type: application/json
{
    "replace_original": "true",
    "text": "Thanks for your request, we'll process it and get back to you."
}
```

Feel free to include [blocks](https://docs.slack.dev/block-kit) in your `response_url` update.

Non- [ephemeral](https://docs.slack.dev/messaging#ephemeral) messages can also be updated [using `chat.update`](https://docs.slack.dev/messaging/modifying-messages#deleting). A message's type cannot be changed from `ephemeral` to `in_channel`. Once a message is issued, it retains its visibility quality for life.

You cannot use `replace_original` to modify the user-posted message that invokes a [slash command](https://docs.slack.dev/interactivity/implementing-slash-commands).

#### Deleting a source message in response

You can also delete a source message of an interaction entirely using `response_url`.

Include `delete_original` as the *sole attribute* in your `response_url` JSON, with the value set to `true`:

```markdown
POST SLACK_REDACTED_EXAMPLE
Content-type: application/json
{
    "delete_original": "true"
}
```

Non- [ephemeral](https://docs.slack.dev/messaging#ephemeral) messages can also be deleted [using `chat.delete`](https://docs.slack.dev/messaging/modifying-messages#deleting).

If you include a new message payload *and* `delete_original`, the source message will be deleted, and your new message published.

The `delete_original` field is not supported for use with [slash commands](https://docs.slack.dev/interactivity/implementing-slash-commands), so sending it will not remove or retain a user-posted message that is invoked by the slash command.

### Modal responses

> The wonderful thing about triggers is that `trigger_id` s are wonderful things. They're attached to events and interactions. They expire in three seconds. They're fun fun fun fun *fun*! But the most wonderful thing about triggers is they can be used only once.

When certain events and interactions occur between users and your app, you'll receive a `trigger_id` as part of the [interaction payload](https://docs.slack.dev/reference/interaction-payloads). If you have a `trigger_id`, you can use the value of that field to [open a modal](https://docs.slack.dev/surfaces/modals#opening).

Triggers **expire in three seconds**. Use them before you lose them. You'll receive a `trigger_expired` error when using a method with an expired `trigger_id`.

Triggers **may only be used once**. You may perform just one operation with a `trigger_id`. Subsequent attempts are presented with a `trigger_exchanged` error.

#### Using a modal to generate a response\_url

When you're [composing your modal](https://docs.slack.dev/surfaces/modals#opening), you can use special parameters to generate a `response_url` when the modal is submitted. [Read our guide to using modals to learn more about this technique](https://docs.slack.dev/surfaces/modals#modal_response_url). You can then use this newly-generated `response_url` to publish a message [as described above](#message_responses).

### Asynchronous responses

The receipt of an interaction payload may furnish your app with a `response_url` or a `trigger_id`, but it also imbues a lot of contextual knowledge about the interaction.

That context can be used by your app to do, well, anything. For example, it can:

- respond to a user clicking a *Complete Task* button in an interactive message by a [adding a happy reaction emoji](https://docs.slack.dev/reference/methods/reactions.add) to the source message.
- use a modal submission as the trigger for an update of the app's [Home tab](https://docs.slack.dev/surfaces/app-home).
- send data to query and update an external service after the use of an [app action](https://docs.slack.dev/interactivity/implementing-shortcuts).
- order cheeseburgers (hold the onions) for the team after someone clicks the *Bring me Lunch* [button](https://docs.slack.dev/reference/block-kit/block-elements/button-element).

---

## Growing your garden

We're excited about the possibilities opened to your app, when your apps are open to possible interactions.

Get some inspiration by [reading our guide to planning your app](https://docs.slack.dev/concepts/app-design), and enable more opportunities for interaction by [using the available app surfaces](https://docs.slack.dev/surfaces).