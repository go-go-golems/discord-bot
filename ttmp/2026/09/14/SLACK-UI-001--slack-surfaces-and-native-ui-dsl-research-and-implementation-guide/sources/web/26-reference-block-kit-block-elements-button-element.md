---
Title: "26 reference block kit block elements button element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/button-element/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

Example:

![Three buttons showing default, primary, and danger color styles](data:image/png;base64,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)

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of element. In this case `type` is always `button`. | Required |
| `text` | Object | A [text object](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) that defines the button's text. Can only be of `type: plain_text`. `text` may truncate with ~30 characters. Maximum length for the `text` in this field is 75 characters. | Required |
| `action_id` | String | An identifier for this action. You can use this when you receive an interaction payload to [identify the source of the action](https://docs.slack.dev/interactivity/handling-user-interaction#payloads). Should be unique among all other `action_id` s in the containing block. Maximum length is 255 characters. | Optional |
| `url` | String | A URL to load in the user's browser when the button is clicked. Maximum length is 3000 characters. If you're using `url`, you'll still receive an [interaction payload](https://docs.slack.dev/interactivity/handling-user-interaction#payloads) and will need to [send an acknowledgement response](https://docs.slack.dev/interactivity/handling-user-interaction#acknowledgment_response). | Optional |
| `value` | String | The value to send along with the [interaction payload](https://docs.slack.dev/interactivity/handling-user-interaction#payloads). Maximum length is 2000 characters. | Optional |
| `style` | String | Decorates buttons with alternative visual color schemes. Use this option with restraint.`primary` gives buttons a green outline and text, ideal for affirmation or confirmation actions. `primary` should only be used for one button within a set.`danger` gives buttons a red outline and text, and should be used when the action is destructive. Use `danger` even more sparingly than `primary`.If you don't include this field, the default button style will be used. | Optional |
| `confirm` | Object | A [confirm object](https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object) that defines an optional confirmation dialog after the button is clicked. | Optional |
| `accessibility_label` | String | A label for longer descriptive text about a button element. This label will be read out by screen readers *instead of* the button [`text`](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) object. Maximum length is 75 characters. | Optional |
| `agent_prompt` | String | A prompt to hand off to [Slackbot](https://docs.slack.dev/ai), Slack's AI assistant, when the button is clicked. Any app can set this field. When the user who clicks the button has Slackbot AI available to them, clicking opens Slackbot seeded with this prompt instead of sending an [interaction payload](https://docs.slack.dev/interactivity/handling-user-interaction#payloads) to your app; otherwise the button sends a standard interaction payload. Maximum length is 4000 characters. | Optional |

## Examples

A regular interactive button:

A button with a `primary` `style` attribute:

A link button:

A button that hands off to [Slackbot](https://docs.slack.dev/ai) with a pre-filled prompt:

```json
{
  "type": "button",
  "text": {
    "type": "plain_text",
    "text": "Summarize with Slackbot"
  },
  "action_id": "summarize_incident",
  "agent_prompt": "Summarize this incident and suggest next steps."
}
```

The button element must be used inside either the [section](https://docs.slack.dev/reference/block-kit/blocks/section-block) or [actions](https://docs.slack.dev/reference/block-kit/blocks/actions-block) block, like this: