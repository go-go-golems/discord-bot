---
Title: "extra emoji element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/emoji-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

This is a rich text element, compatible only with the [`rich_text`](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block) block. It must be used within the [`rich_text_list`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element), [`rich_text_quote`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-quote-element), or [`rich_text_section`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element) block element within the `rich_text` block's `elements` array.

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of object; in this case, "emoji". | Required |
| `name` | String | The name of the emoji; i.e. "wave" or "wave::skin-tone-2". | Required |
| `unicode` | String | Represents the unicode code point of the emoji, where applicable. | Optional |

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
                            "type": "emoji",
                            "name": "basketball"
                        },
                        {
                            "type": "text",
                            "text": " "
                        },
                        {
                            "type": "emoji",
                            "name": "snowboarder"
                        },
                        {
                            "type": "text",
                            "text": " "
                        },
                        {
                            "type": "emoji",
                            "name": "checkered_flag"
                        }
                    ]
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/#%7B%2522blocks%2522%3A%255B%257B%2522type%2522%3A%2522rich_text%2522%2C%2522elements%2522%3A%255B%257B%2522type%2522%3A%2522rich_text_section%2522%2C%2522elements%2522%3A%255B%257B%2522type%2522%3A%2522emoji%2522%2C%2522name%2522%3A%2522basketball%2522%257D%2C%257B%2522type%2522%3A%2522text%2522%2C%2522text%2522%3A%2522%2520%2522%257D%2C%257B%2522type%2522%3A%2522emoji%2522%2C%2522name%2522%3A%2522snowboarder%2522%257D%2C%257B%2522type%2522%3A%2522text%2522%2C%2522text%2522%3A%2522%2520%2522%257D%2C%257B%2522type%2522%3A%2522emoji%2522%2C%2522name%2522%3A%2522checkered_flag%2522%257D%255D%257D%255D%257D%255D%257D)