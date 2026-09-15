---
Title: "extra workflow object"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/composition-objects/workflow-object"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

**Defines an object containing workflow information.**

#### Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `trigger` | Object | A [`trigger`](https://docs.slack.dev/reference/block-kit/composition-objects/trigger-object) object that contains information about a workflow's trigger. | Required |

#### Example

A workflow object must be used inside of a [workflow button](https://docs.slack.dev/reference/block-kit/block-elements/workflow-button-element) element.

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