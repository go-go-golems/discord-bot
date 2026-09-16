---
Title: "extra file block"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/file-block"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of block. For a file block, `type` is always `file`. | Required |
| `external_id` | String | The external unique ID for this file. | Required |
| `source` | String | At the moment, `source` will always be `remote` for a remote file. | Required |
| `block_id` | String | A unique identifier for a block. If not specified, one will be generated. Maximum length for this field is 255 characters. `block_id` should be unique for each message and each iteration of a message. If a message is updated, use a new `block_id`. | Optional |

## Usage info

You can't add this block to app surfaces directly, but it will show up when [retrieving messages](https://docs.slack.dev/messaging/retrieving-messages) that contain remote files.

If you want to add remote files to messages, [follow our guide](https://docs.slack.dev/messaging/working-with-files#remote).

![An example of a file block](https://docs.slack.dev/assets/images/file_upload_remote_file-cf719274fa2f07195a391c5e868e3820.png)