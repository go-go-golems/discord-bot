---
Title: "extra url source element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/url-source-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of element. In this case `type` is always `url`. | Required |
| `url` | String | The URL type source. | Required |
| `text` | String | Display text for the URL. | Required |

## Usage info

The URL source element is used to display clickable URL references within a [task card block](https://docs.slack.dev/reference/block-kit/blocks/task-card-block). It cannot be used within other blocks. Note that whether the URL actually resolves via DNS is not validated.

## Examples

A URL source element:

- JSON
- Python Slack SDK
- Node Slack SDK

```json
{
  "type": "url",
  "url": "https://docs.slack.dev/",
  "text": "Slack API docs"
}
```

Example within a task card block:

- JSON
- Python Slack SDK
- Node Slack SDK

```json
{
  "type": "task_card",
  "task_id": "task_1",
  "title": "Scientific findings",
  "status": "complete",
  "sources": [
    {
      "type": "url",
      "url": "https://docs.example.com/",
      "text": "Tracy's delightful docs"
    },
    {
      "type": "url",
      "url": "https://research.example.com/",
      "text": "Haley's resourceful research"
    }
  ]
}
```