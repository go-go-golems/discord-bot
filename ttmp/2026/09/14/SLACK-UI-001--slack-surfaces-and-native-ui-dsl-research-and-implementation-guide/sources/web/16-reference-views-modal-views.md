---
Title: "16 reference views modal views"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/views/modal-views/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

Modal view objects are used within the following [Web API](https://docs.slack.dev/apis/web-api/) methods:

- [`views.open`](https://docs.slack.dev/reference/methods/views.open)
- [`views.update`](https://docs.slack.dev/reference/methods/views.update)
- [`views.push`](https://docs.slack.dev/reference/methods/views.push)

Non-standard characters (including characters with diacritics) within view objects are converted and sent in unicode format when you receive the view callback payloads.

## Modal view object fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of view. Set to `modal` for modals. | Required |
| `title` | Object | The title that appears in the top-left of the modal. Must be a [`plain_text` text element](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) with a max length of 24 characters. | Required |
| `blocks` | Array | An array of [blocks](https://docs.slack.dev/reference/block-kit/blocks) that defines the content of the view. Max of 100 blocks. | Required |
| `close` | Object | A [`plain_text` element](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) that defines the text displayed in the close button at the bottom-right of the view. Max length of 24 characters. | Optional |
| `submit` | Object | A [`plain_text` element](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) that defines the text displayed in the submit button at the bottom-right of the view. `submit` is required when an `input` block is within the `blocks` array. Max length of 24 characters. | Optional |
| `private_metadata` | String | A string that will be sent to your app in `view_submission` and `block_actions` events. Max length of 3000 characters. | Optional |
| `callback_id` | String | An identifier to recognize interactions and submissions of this particular view. Don't use this to store sensitive information (use `private_metadata` instead). Max length of 255 characters. | Optional |
| `clear_on_close` | Boolean | When set to `true`, clicking on the `close` button will clear all views in a modal and close it. Defaults to `false`. | Optional |
| `notify_on_close` | Boolean | Indicates whether Slack will send your request URL a `view_closed` event when a user clicks the `close` button. Defaults to `false`. | Optional |
| `external_id` | String | A custom identifier that must be unique for all views on a per-team basis. | Optional |
| `submit_disabled` | Boolean | When set to `true`, disables the `submit` button until the user has completed one or more inputs. *This property is for [configuration modals](https://docs.slack.dev/changelog/2023-08-workflow-steps-from-apps-step-back).* |  |

## Modal view example

```json
{
    "type": "modal",
    "title": {
        "type": "plain_text",
        "text": "Modal title"
    },
    "blocks": [
        {
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": "It's Block Kit...but _in a modal_"
            },
            "block_id": "section1",
            "accessory": {
                "type": "button",
                "text": {
                    "type": "plain_text",
                    "text": "Click me"
                },
                "action_id": "button_abc",
                "value": "Button value",
                "style": "danger"
            }
        },
        {
            "type": "input",
            "label": {
                "type": "plain_text",
                "text": "Input label"
            },
            "element": {
                "type": "plain_text_input",
                "action_id": "input1",
                "placeholder": {
                    "type": "plain_text",
                    "text": "Type in here"
                },
                "multiline": false
            },
            "optional": false
        }
    ],
    "close": {
        "type": "plain_text",
        "text": "Cancel"
    },
    "submit": {
        "type": "plain_text",
        "text": "Save"
    },
    "private_metadata": "Shhhhhhhh",
    "callback_id": "view_identifier_12"
}
```