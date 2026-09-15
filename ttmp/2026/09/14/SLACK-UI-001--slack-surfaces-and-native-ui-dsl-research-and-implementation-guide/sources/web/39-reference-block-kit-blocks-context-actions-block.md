---
Title: "39 reference block kit blocks context actions block"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/context-actions-block/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of block. For a context actions block, `type` is always `context_actions`. | Required |
| `elements` | Object\[\] | An array of [feedback buttons elements](https://docs.slack.dev/reference/block-kit/block-elements/feedback-buttons-element) and [icon button elements](https://docs.slack.dev/reference/block-kit/block-elements/icon-button-element). Maximum number of items is 5. | Required |
| `block_id` | String | A unique identifier for a block. You can use this `block_id` when you receive an interaction payload to [identify the source of the action](https://docs.slack.dev/interactivity/handling-user-interaction#payloads). If not specified, `block_id` will be generated. Maximum length for this field is 255 characters. `block_id` should be unique for each message and each iteration of a message. If a message is updated, use a new `block_id`. | Optional |

## Examples

**Example 1**: Context actions block with [feedback buttons](https://docs.slack.dev/reference/block-kit/block-elements/feedback-buttons-element):

- JSON
- Python Slack SDK
- Node Slack SDK

```json
{
    "blocks": [
        {
            "type": "context_actions",
            "elements": [
                {
                    "type": "feedback_buttons",
                    "action_id": "feedback_buttons_1",
                    "positive_button": {
                        "text": {
                            "type": "plain_text",
                            "text": "👍"
                        },
                        "value": "positive_feedback"
                    },
                    "negative_button": {
                        "text": {
                            "type": "plain_text",
                            "text": "👎"
                        },
                        "value": "negative_feedback"
                    }
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/T024BE7LD#%7B%22blocks%22:%5B%7B%22type%22:%22context_actions%22,%22elements%22:%5B%7B%22type%22:%22feedback_buttons%22,%22action_id%22:%22feedback_buttons_1%22,%22positive_button%22:%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22%F0%9F%91%8D%22%7D,%22value%22:%22positive_feedback%22%7D,%22negative_button%22:%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22%F0%9F%91%8E%22%7D,%22value%22:%22negative_feedback%22%7D%7D%5D%7D%5D%7D)

---

**Example 2**: Context actions block with an [icon button](https://docs.slack.dev/reference/block-kit/block-elements/icon-button-element):

- JSON
- Python Slack SDK
- Node Slack SDK

```json
{
    "blocks": [
        {
            "type": "context_actions",
            "elements": [
                {
                    "type": "icon_button",
                    "icon": "trash",
                    "text": {
                        "type": "plain_text",
                        "text": "Delete"
                    },
                    "action_id": "delete_button_1",
                    "value": "delete_item"
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/T024BE7LD#%7B%22blocks%22:%5B%7B%22type%22:%22context_actions%22,%22elements%22:%5B%7B%22type%22:%22icon_button%22,%22icon%22:%22trash%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Delete%22%7D,%22action_id%22:%22delete_button_1%22,%22value%22:%22delete_item%22%7D%5D%7D%5D%7D)