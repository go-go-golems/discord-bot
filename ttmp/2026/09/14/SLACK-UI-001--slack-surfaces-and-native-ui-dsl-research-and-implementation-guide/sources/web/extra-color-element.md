---
Title: "extra color element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/color-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

This is a rich text element, compatible only with the [`rich_text`](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block) block. It must be used within the [`rich_text_list`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element), [`rich_text_quote`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-quote-element), or [`rich_text_section`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element) block element within the `rich_text` block's `elements` array.

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of object; in this case, "color". | Required |
| `value` | String | The hex value for the color. | Required |
| `style` | Object | An object of optional boolean properties that dictate style: `bold`, `italic`, `strike`, `highlight`, `client_highlight`, and `underline`. | Optional |

## Example

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "rich_text",
            "elements": [
                {
                    "type": "rich_text_section",
                    "elements": [
                        {
                            "type": "color",
                            "value": "#F405B3"
                        }
                    ]
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/T024BE7LD#%7B%22blocks%22:%5B%7B%22type%22:%22rich_text%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22color%22,%22value%22:%22#F405B3%22%7D%5D%7D%5D%7D%5D%7D)