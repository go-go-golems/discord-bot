---
Title: "extra rich text section element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of sub-element; in this case, `rich_text_section`. | Required |
| `elements` | Object \[\] | An array of rich text elements. Rich text elements include [`attachment_mention`](https://docs.slack.dev/reference/block-kit/block-elements/attachment-mention-element), [`broadcast`](https://docs.slack.dev/reference/block-kit/block-elements/broadcast-element), [`canvas`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-element), [`canvas_user_mention`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-user-mention-element), [`canvas_message_unfurl`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-message-unfurl-element), [`channel`](https://docs.slack.dev/reference/block-kit/block-elements/channel-element), [`citation`](https://docs.slack.dev/reference/block-kit/block-elements/citation-element), [`color`](https://docs.slack.dev/reference/block-kit/block-elements/color-element), [`date`](https://docs.slack.dev/reference/block-kit/block-elements/date-element), [`emoji`](https://docs.slack.dev/reference/block-kit/block-elements/emoji-element), [`file`](https://docs.slack.dev/reference/block-kit/block-elements/file-element), [`link`](https://docs.slack.dev/reference/block-kit/block-elements/link-element), [`list_record`](https://docs.slack.dev/reference/block-kit/block-elements/list-record-element), [`message_mention`](https://docs.slack.dev/reference/block-kit/block-elements/message-mention-element), [`salesforce_data_field`](https://docs.slack.dev/reference/block-kit/block-elements/salesforce-data-field-element), [`tag`](https://docs.slack.dev/reference/block-kit/block-elements/tag-element), [`team`](https://docs.slack.dev/reference/block-kit/block-elements/team-element), [`text`](https://docs.slack.dev/reference/block-kit/block-elements/text-element), [`user`](https://docs.slack.dev/reference/block-kit/block-elements/user-element), [`usergroup`](https://docs.slack.dev/reference/block-kit/block-elements/usergroup-element), [`work_object_mention`](https://docs.slack.dev/reference/block-kit/block-elements/work-object-mention-element), [`workflow_mention`](https://docs.slack.dev/reference/block-kit/block-elements/workflow-mention-element). | Required |

## Example

![An example of a rich\_text\_section block](https://docs.slack.dev/assets/images/bk_rich_text_section_example-2ab6aae732c461ab4f39e2da4802f7a0.png)

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
                            "type": "text",
                            "text": "Hello there, I am a basic rich text block!"
                        }
                    ]
                }
            ]
        },
        {
            "type": "rich_text",
            "elements": [
                {
                    "type": "rich_text_section",
                    "elements": [
                        {
                            "type": "text",
                            "text": "Hello there, "
                        },
                        {
                            "type": "text",
                            "text": "I am a bold rich text block!",
                            "style": {
                                "bold": true
                            }
                        }
                    ]
                }
            ]
        },
        {
            "type": "rich_text",
            "elements": [
                {
                    "type": "rich_text_section",
                    "elements": [
                        {
                            "type": "text",
                            "text": "Hello there, "
                        },
                        {
                            "type": "text",
                            "text": "I am an italic rich text block!",
                            "style": {
                                "italic": true
                            }
                        }
                    ]
                }
            ]
        },
        {
            "type": "rich_text",
            "elements": [
                {
                    "type": "rich_text_section",
                    "elements": [
                        {
                            "type": "text",
                            "text": "Hello there, "
                        },
                        {
                            "type": "text",
                            "text": "I am a strikethrough rich text block!",
                            "style": {
                                "strike": true
                            }
                        }
                    ]
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/#%7B%22blocks%22:%5B%7B%22type%22:%22rich_text%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Hello%20there,%20I%20am%20a%20basic%20rich%20text%20block!%22%7D%5D%7D%5D%7D,%7B%22type%22:%22rich_text%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Hello%20there,%20%22%7D,%7B%22type%22:%22text%22,%22text%22:%22I%20am%20a%20bold%20rich%20text%20block!%22,%22style%22:%7B%22bold%22:true%7D%7D%5D%7D%5D%7D,%7B%22type%22:%22rich_text%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Hello%20there,%20%22%7D,%7B%22type%22:%22text%22,%22text%22:%22I%20am%20an%20italic%20rich%20text%20block!%22,%22style%22:%7B%22italic%22:true%7D%7D%5D%7D%5D%7D,%7B%22type%22:%22rich_text%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Hello%20there,%20%22%7D,%7B%22type%22:%22text%22,%22text%22:%22I%20am%20a%20strikethrough%20rich%20text%20block!%22,%22style%22:%7B%22strike%22:true%7D%7D%5D%7D%5D%7D%5D%7D)