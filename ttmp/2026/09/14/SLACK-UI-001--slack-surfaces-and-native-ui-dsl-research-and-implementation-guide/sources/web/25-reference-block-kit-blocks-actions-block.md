---
Title: "25 reference block kit blocks actions block"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/actions-block/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of block. For an actions block, `type` is always `actions`. | Required |
| `elements` | Object\[\] | An array of interactive [element objects](https://docs.slack.dev/reference/block-kit/block-elements) - [buttons](https://docs.slack.dev/reference/block-kit/block-elements/button-element), [select menus](https://docs.slack.dev/reference/block-kit/block-elements/select-menu-element), [overflow menus](https://docs.slack.dev/reference/block-kit/block-elements/overflow-menu-element), or [date pickers](https://docs.slack.dev/reference/block-kit/block-elements/date-picker-element). There is a maximum of 25 elements in each action block. | Required |
| `block_id` | String | A unique identifier for a block. If not specified, a `block_id` will be generated. You can use this `block_id` when you receive an interaction payload to [identify the source of the action](https://docs.slack.dev/interactivity/handling-user-interaction#payloads). Maximum length for this field is 255 characters. `block_id` should be unique for each message and each iteration of a message. If a message is updated, use a new `block_id`. | Optional |

## Examples

**Example 1**: An actions block with a select menu and a button:

![An action block containing multiple interactive elements](data:image/png;base64,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)

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "actions",
            "block_id": "actions1",
            "elements": [
                {
                    "type": "static_select",
                    "placeholder": {
                        "type": "plain_text",
                        "text": "Which witch is the witchiest witch?"
                    },
                    "action_id": "select_2",
                    "options": [
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "Matilda"
                            },
                            "value": "matilda"
                        },
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "Glinda"
                            },
                            "value": "glinda"
                        },
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "Granny Weatherwax"
                            },
                            "value": "grannyWeatherwax"
                        },
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "Hermione"
                            },
                            "value": "hermione"
                        }
                    ]
                },
                {
                    "type": "button",
                    "text": {
                        "type": "plain_text",
                        "text": "Cancel"
                    },
                    "value": "cancel",
                    "action_id": "button_1"
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/#%7B%22blocks%22:%5B%7B%22type%22:%22actions%22,%22block_id%22:%22actions1%22,%22elements%22:%5B%7B%22type%22:%22static_select%22,%22placeholder%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Which%20witch%20is%20the%20witchiest%20witch?%22%7D,%22action_id%22:%22select_2%22,%22options%22:%5B%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Matilda%22%7D,%22value%22:%22matilda%22%7D,%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Glinda%22%7D,%22value%22:%22glinda%22%7D,%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Granny%20Weatherwax%22%7D,%22value%22:%22grannyWeatherwax%22%7D,%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Hermione%22%7D,%22value%22:%22hermione%22%7D%5D%7D,%7B%22type%22:%22button%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Cancel%22%7D,%22value%22:%22cancel%22,%22action_id%22:%22button_1%22%7D%5D%7D%5D%7D)

---

**Example 2**: An actions block with a datepicker, an overflow, and a button:

![An action block containing multiple interactive elements](data:image/png;base64,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)

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "actions",
            "block_id": "actionblock789",
            "elements": [
                {
                    "type": "datepicker",
                    "action_id": "datepicker123",
                    "initial_date": "1990-04-28",
                    "placeholder": {
                        "type": "plain_text",
                        "text": "Select a date"
                    }
                },
                {
                    "type": "overflow",
                    "options": [
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "*this is plain_text text*"
                            },
                            "value": "value-0"
                        },
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "*this is plain_text text*"
                            },
                            "value": "value-1"
                        },
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "*this is plain_text text*"
                            },
                            "value": "value-2"
                        },
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "*this is plain_text text*"
                            },
                            "value": "value-3"
                        },
                        {
                            "text": {
                                "type": "plain_text",
                                "text": "*this is plain_text text*"
                            },
                            "value": "value-4"
                        }
                    ],
                    "action_id": "overflow"
                },
                {
                    "type": "button",
                    "text": {
                        "type": "plain_text",
                        "text": "Click Me"
                    },
                    "value": "click_me_123",
                    "action_id": "button"
                }
            ]
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/#%7B%22blocks%22:%5B%7B%22type%22:%22actions%22,%22block_id%22:%22actionblock789%22,%22elements%22:%5B%7B%22type%22:%22datepicker%22,%22action_id%22:%22datepicker123%22,%22initial_date%22:%221990-04-28%22,%22placeholder%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Select%20a%20date%22%7D%7D,%7B%22type%22:%22overflow%22,%22options%22:%5B%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22*this%20is%20plain_text%20text*%22%7D,%22value%22:%22value-0%22%7D,%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22*this%20is%20plain_text%20text*%22%7D,%22value%22:%22value-1%22%7D,%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22*this%20is%20plain_text%20text*%22%7D,%22value%22:%22value-2%22%7D,%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22*this%20is%20plain_text%20text*%22%7D,%22value%22:%22value-3%22%7D,%7B%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22*this%20is%20plain_text%20text*%22%7D,%22value%22:%22value-4%22%7D%5D,%22action_id%22:%22overflow%22%7D,%7B%22type%22:%22button%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Click%20Me%22%7D,%22value%22:%22click_me_123%22,%22action_id%22:%22button%22%7D%5D%7D%5D%7D)