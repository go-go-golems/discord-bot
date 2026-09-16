---
Title: "extra trigger object"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/composition-objects/trigger-object"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

**Defines an object containing trigger information.**

#### Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `url` | String | A [link trigger URL](https://docs.slack.dev/tools/deno-slack-sdk/guides/creating-link-triggers). Must be associated with a valid trigger. | Required |
| `customizable_input_parameters` | Object\[\] | An array of input parameter objects. Each specified name must match an input parameter defined on the workflow of the provided trigger (url), and the input parameter mapping on the trigger must be set as `customizable: true`. Each specified value must match the type defined by the workflow input parameter of the matching name. | Optional |

#### Example

A trigger object must be used inside of a [workflow](https://docs.slack.dev/reference/block-kit/composition-objects/workflow-object) object.

- JSON
- Python Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "section",
            "text": {
                "text": "A message *with some bold text* and _some italicized text_.",
                "type": "mrkdwn"
            },
            "accessory": {
                "type": "workflow_button",
                "text": {
                    "type": "plain_text",
                    "text": "Run Workflow"
                },
                "action_id": "workflowbutton123",
                "workflow": {
                    "trigger": {
                        "url": "https://slack.com/shortcuts/Ft0123ABC456/xyz...zyx",
                        "customizable_input_parameters": [
                            {
                                "name": "input_parameter_a",
                                "value": "Value for input param A"
                            },
                            {
                                "name": "input_parameter_b",
                                "value": "Value for input param B"
                            }
                        ]
                    }
                }
            }
        }
    ]
}
```

---