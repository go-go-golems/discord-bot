---
Title: "31 reference block kit blocks plan block"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/plan-block/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of block. For this block, type will always be `plan`. | Required |
| `title` | String | Title of the plan in plain text. | Required |
| `tasks` | Array | A sequence of [task card blocks](https://docs.slack.dev/reference/block-kit/blocks/task-card-block), maximum of 50 tasks. Each one contains task-like objects without a type, and each `task_id` in a plan must be unique. | Required |
| `block_id` | String | A unique identifier for a block. If not specified, one will be generated. Maximum length for this field is 255 characters. `block_id` should be unique for each message and each iteration of a message. If a message is updated, use a new `block_id`. | Optional |

## Usage info

## Examples

- JSON
- Python Slack SDK
- Node Slack SDK

```json
{
  "blocks": [
    {
      "type": "plan",
      "title": "Thinking completed",
      "tasks": [
        {
          "task_id": "call_001",
          "title": "Fetched user profile information",
          "status": "in_progress",
          "details": {
            "type": "rich_text",
            "block_id": "viMWO",
            "elements": [
              {
                "type": "rich_text_section",
                "elements": [
                  {
                    "type": "text",
                    "text": "Searched database..."
                  }
                ]
              }
            ]
          },
          "output": {
            "type": "rich_text",
            "block_id": "viMWO",
            "elements": [
              {
                "type": "rich_text_section",
                "elements": [
                  {
                    "type": "text",
                    "text": "Profile data loaded"
                  }
                ]
              }
            ]
          }
        },
        {
          "task_id": "call_002",
          "title": "Checked user permissions",
          "status": "pending"
        },
        {
          "task_id": "call_003",
          "title": "Generated comprehensive user report",
          "status": "complete",
          "output": {
            "type": "rich_text",
            "block_id": "crsk",
            "elements": [
              {
                "type": "rich_text_section",
                "elements": [
                  {
                    "type": "text",
                    "text": "15 data points compiled"
                  }
                ]
              }
            ]
          }
        }
      ]
    }
  ]
}
```