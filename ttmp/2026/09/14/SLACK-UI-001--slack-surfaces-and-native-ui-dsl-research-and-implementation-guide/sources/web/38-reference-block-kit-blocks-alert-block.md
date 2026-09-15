---
Title: "38 reference block kit blocks alert block"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/alert-block/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of block. For an alert block, `type` is always `alert`. | Required |
| `text` | String | The alert message, using `plain_text` or `mrkdwn` formatting. Maximum 200 characters. | Required |
| `level` | Array | One of `default`, `info`, `warning`, `error`, or `success`. Will be `default` if omitted. | Optional |
| `block_id` | String | A unique identifier for a block. If not specified, a `block_id` will be generated. | Optional |

Alert blocks are currently only supported in modals.

## Example

A sample alert block:

- JSON
- Python Slack SDK
- Node Slack SDK

```json
{
    "blocks": [
        {
            "type": "alert",
            "text": {
                "type": "mrkdwn",
                "text": "The work is mysterious and important.",
                "verbatim": false
            },
            "level": "info"
        }
    ]
}
```

[View in Block Kit Builder](https://app.slack.com/block-kit-builder/T29KZ003T#%7B%22type%22:%22modal%22,%22title%22:%7B%22type%22:%22plain_text%22,%22text%22:%22My%20App%22,%22emoji%22:true%7D,%22submit%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Submit%22,%22emoji%22:true%7D,%22close%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Cancel%22,%22emoji%22:true%7D,%22blocks%22:%5B%7B%22type%22:%22alert%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22The%20work%20is%20mysterious%20and%20important.%22,%22verbatim%22:false%7D,%22level%22:%22info%22%7D%5D%7D)