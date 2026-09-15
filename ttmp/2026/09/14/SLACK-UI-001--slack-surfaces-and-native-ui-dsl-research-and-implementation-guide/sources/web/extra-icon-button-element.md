---
Title: "extra icon button element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/icon-button-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of element. In this case `type` is always `icon_button`. | Required |
| `icon` | String | The icon to show. The `trash` icon is the only icon available at this time. | Required |
| `text` | Object | A [text object](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) that defines the button's text. Can only be of `type: plain_text`. | Required |
| `action_id` | String | An identifier for this action. You can use this when you receive an interaction payload to [identify the source of the action](https://docs.slack.dev/interactivity/handling-user-interaction#payloads). Should be unique among all other `action_id` s in the containing block. Maximum length is 255 characters. | Optional |
| `value` | String | The value to send along with the [interaction payload](https://docs.slack.dev/interactivity/handling-user-interaction#payloads). Maximum length is 2000 characters. | Optional |
| `confirm` | Object | A [confirm object](https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object) that defines an optional confirmation dialog after the button is clicked. | Optional |
| `accessibility_label` | String | A label for longer descriptive text about a button element. This label will be read out by screen readers instead of the button `text` object. Maximum length is 75 characters. | Optional |
| `visible_to_user_ids` | Array | An array of user IDs for which the icon button appears. If not provided, the button is visible to all users. | Optional |

## Examples

The icon button must be used inside of the [context actions](https://docs.slack.dev/reference/block-kit/blocks/context-actions-block) block, like this:

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

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
          "action_id": "delete_button",
          "value": "delete_item"
        }
      ]
    }
  ]
}
```

[Preview in Block Kit Builder](https://app.slack.com/block-kit-builder/T024BE7LD#%7B%22blocks%22:%5B%7B%22type%22:%22context_actions%22,%22elements%22:%5B%7B%22type%22:%22icon_button%22,%22icon%22:%22trash%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Delete%22%7D,%22action_id%22:%22delete_button%22,%22value%22:%22delete_item%22%7D%5D%7D%5D%7D)