---
Title: "extra rich text list element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of sub-element; in this case, `rich_text_list`. | Required |
| `style` | String | Either `bullet` or `ordered`, the latter meaning a numbered list. | Required |
| `elements` | Object \[\] | An array of [`rich_text_section`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element) objects containing two properties: `type`, which is "rich\_text\_section", and `elements`, which is an array of rich text elements. Rich text elements include [`attachment_mention`](https://docs.slack.dev/reference/block-kit/block-elements/attachment-mention-element), [`broadcast`](https://docs.slack.dev/reference/block-kit/block-elements/broadcast-element), [`canvas`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-element), [`canvas_user_mention`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-user-mention-element), [`canvas_message_unfurl`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-message-unfurl-element), [`channel`](https://docs.slack.dev/reference/block-kit/block-elements/channel-element), [`citation`](https://docs.slack.dev/reference/block-kit/block-elements/citation-element), [`color`](https://docs.slack.dev/reference/block-kit/block-elements/color-element), [`date`](https://docs.slack.dev/reference/block-kit/block-elements/date-element), [`emoji`](https://docs.slack.dev/reference/block-kit/block-elements/emoji-element), [`file`](https://docs.slack.dev/reference/block-kit/block-elements/file-element), [`link`](https://docs.slack.dev/reference/block-kit/block-elements/link-element), [`list_record`](https://docs.slack.dev/reference/block-kit/block-elements/list-record-element), [`message_mention`](https://docs.slack.dev/reference/block-kit/block-elements/message-mention-element), [`salesforce_data_field`](https://docs.slack.dev/reference/block-kit/block-elements/salesforce-data-field-element), [`tag`](https://docs.slack.dev/reference/block-kit/block-elements/tag-element), [`team`](https://docs.slack.dev/reference/block-kit/block-elements/team-element), [`text`](https://docs.slack.dev/reference/block-kit/block-elements/text-element), [`user`](https://docs.slack.dev/reference/block-kit/block-elements/user-element), [`usergroup`](https://docs.slack.dev/reference/block-kit/block-elements/usergroup-element), [`work_object_mention`](https://docs.slack.dev/reference/block-kit/block-elements/work-object-mention-element), [`workflow_mention`](https://docs.slack.dev/reference/block-kit/block-elements/workflow-mention-element). | Required |
| `indent` | Number | Sub-list indent level. | Optional |
| `offset` | Number | Number to offset the first number in the list. For example, if the `offset = 4`, the first number in the ordered list would be 5. | Optional |
| `border` | Number | Turn the border on or off. | Optional |

## Example

![An example of a rich\_text\_list block](https://docs.slack.dev/assets/images/bk_rich_text_list_example-b045fe20977becf5a4cb94ffe90bb561.png)

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "rich_text",
            "block_id": "block1",
            "elements": [
                {
                    "type": "rich_text_section",
                    "elements": [
                        {
                            "type": "text",
                            "text": "My favorite Slack features (in no particular order):"
                        }
                    ]
                },
                {
                    "type": "rich_text_list",
                    "elements": [
                        {
                            "type": "rich_text_section",
                            "elements": [
                                {
                                    "type": "text",
                                    "text": "Huddles"
                                }
                            ]
                        },
                        {
                            "type": "rich_text_section",
                            "elements": [
                                {
                                    "type": "text",
                                    "text": "Canvas"
                                }
                            ]
                        },
                        {
                            "type": "rich_text_section",
                            "elements": [
                                {
                                    "type": "text",
                                    "text": "Developing with Block Kit"
                                }
                            ]
                        }
                    ],
                    "style": "bullet",
                    "indent": 0,
                    "border": 1
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/T024BE7LD#%7B%22blocks%22:%5B%7B%22type%22:%22rich_text%22,%22block_id%22:%22block1%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22My%20favorite%20Slack%20features%20\(in%20no%20particular%20order\):%22%7D%5D%7D,%7B%22type%22:%22rich_text_list%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Huddles%22%7D%5D%7D,%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Canvas%22%7D%5D%7D,%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Developing%20with%20Block%20Kit%22%7D%5D%7D%5D,%22style%22:%22bullet%22,%22indent%22:0,%22border%22:1%7D%5D%7D%5D%7D)

Let's say we want to create a nested list, for example something that looks like this:

Breakfast foods I enjoy:

- Hashbrowns
- Eggs
	- Scrambled
	- Over easy
- Pancakes, extra syrup

To create that in rich text, create three instances of `rich_text_list`, the middle one using the `indent` property to indent the types of eggs into that sub-list.

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "rich_text",
            "block_id": "block1",
            "elements": [
                {
                    "type": "rich_text_section",
                    "elements": [
                        {
                            "type": "text",
                            "text": "Breakfast foods I enjoy:"
                        }
                    ]
                },
                {
                    "type": "rich_text_list",
                    "style": "bullet",
                    "elements": [
                        {
                            "type": "rich_text_section",
                            "elements": [
                                {
                                    "type": "text",
                                    "text": "Hashbrowns"
                                }
                            ]
                        },
                        {
                            "type": "rich_text_section",
                            "elements": [
                                {
                                    "type": "text",
                                    "text": "Eggs"
                                }
                            ]
                        }
                    ]
                },
                {
                    "type": "rich_text_list",
                    "style": "bullet",
                    "indent": 1,
                    "elements": [
                        {
                            "type": "rich_text_section",
                            "elements": [
                                {
                                    "type": "text",
                                    "text": "Scrambled"
                                }
                            ]
                        },
                        {
                            "type": "rich_text_section",
                            "elements": [
                                {
                                    "type": "text",
                                    "text": "Over easy"
                                }
                            ]
                        }
                    ]
                },
                {
                    "type": "rich_text_list",
                    "style": "bullet",
                    "elements": [
                        {
                            "type": "rich_text_section",
                            "elements": [
                                {
                                    "type": "text",
                                    "text": "Pancakes, extra syrup"
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/T024BE7LD#%7B%22blocks%22:%5B%7B%22type%22:%22rich_text%22,%22block_id%22:%22block1%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Breakfast%20foods%20I%20enjoy:%22%7D%5D%7D,%7B%22type%22:%22rich_text_list%22,%22style%22:%22bullet%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Hashbrowns%22%7D%5D%7D,%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Eggs%22%7D%5D%7D%5D%7D,%7B%22type%22:%22rich_text_list%22,%22style%22:%22bullet%22,%22indent%22:1,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Scrambled%22%7D%5D%7D,%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Over%20easy%22%7D%5D%7D%5D%7D,%7B%22type%22:%22rich_text_list%22,%22style%22:%22bullet%22,%22elements%22:%5B%7B%22type%22:%22rich_text_section%22,%22elements%22:%5B%7B%22type%22:%22text%22,%22text%22:%22Pancakes,%20extra%20syrup%22%7D%5D%7D%5D%7D%5D%7D%5D%7D)