---
Title: "extra url input element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/url-input-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Fields | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of element. In this case `type` is always `url_text_input`. | Required |
| `action_id` | String | An identifier for the input value when the parent modal is submitted. You can use this when you receive a `view_submission` payload [to identify the value of the input element](https://docs.slack.dev/surfaces/modals#interactions). Should be unique among all other `action_id` s in the containing block. Maximum length is 255 characters. | Optional |
| `initial_value` | String | The initial value in the URL input when it is loaded. | Optional |
| `dispatch_action_config` | Object | A [dispatch configuration object](https://docs.slack.dev/reference/block-kit/composition-objects/dispatch-action-configuration-object) that determines when during text input the element returns a [`block_actions` payload](https://docs.slack.dev/reference/interaction-payloads/block_actions-payload). | Optional |
| `focus_on_load` | Boolean | Indicates whether the element will be set to auto focus within the [`view object`](https://docs.slack.dev/reference/views). Only one element can be set to `true`. Defaults to `false`. | Optional |
| `placeholder` | Object | A [`plain_text`](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) only text object that defines the placeholder text shown in the URL input. Maximum length for the `text` in this field is 150 characters. | Optional |

## Example

The URL input element must be used inside of the [input](https://docs.slack.dev/reference/block-kit/blocks/input-block) block, like this:

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "input",
            "element": {
                "type": "url_text_input",
                "action_id": "url_text_input-action"
            },
            "label": {
                "type": "plain_text",
                "text": "Label",
                "emoji": true
            }
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/T024BE7LD#%7B%22blocks%22:%5B%7B%22type%22:%22input%22,%22element%22:%7B%22type%22:%22url_text_input%22,%22action_id%22:%22url_text_input-action%22%7D,%22label%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Label%22,%22emoji%22:true%7D%7D%5D%7D)