---
Title: "extra rich text preformatted element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/rich-text-preformatted-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of the sub-element; in this case, `rich_text_preformatted`. | Required |
| `elements` | Object \[\] | An array of [text](https://docs.slack.dev/reference/block-kit/block-elements/text-element) or [link](https://docs.slack.dev/reference/block-kit/block-elements/link-element) elements. | Required |
| `border` | Number | Turn the border on or off. | Optional |
| `language` | String | The language of the code block, used for syntax highlighting (e.g., `"python"`, `"javascript"`, `"json"`). | Optional |

## Example

![An example of a rich\_text\_preformatted block](https://docs.slack.dev/assets/images/bk_rich_text_preformatted_example-e2104fccd56714686d7d8d4d5e3a774b.png)

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
                    "type": "rich_text_preformatted",
                    "elements": [
                        {
                            "type": "text",
                            "text": "{\n  \"object\": {\n    \"description\": \"this is an example of a json object\"\n  }\n}"
                        }
                    ],
                    "border": 0
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/T024BE7LD#%7B%22blocks%22:%5B%7B%22type%22:%22rich_text%22,%22elements%22:%5B%7B%22type%22:%22rich_text_preformatted%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22%7B%5Cn%20%20%5C%22object%5C%22:%20%7B%5Cn%20%20%20%20%5C%22description%5C%22:%20%5C%22this%20is%20an%20example%20of%20a%20json%20object%5C%22%5Cn%20%20%7D%5Cn%7D%22%7D%5D,%22border%22:0%7D%5D%7D%5D%7D)