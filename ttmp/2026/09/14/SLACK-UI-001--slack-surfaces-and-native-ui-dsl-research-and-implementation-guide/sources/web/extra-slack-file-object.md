---
Title: "extra slack file object"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/composition-objects/slack-file-object"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

**Defines an object containing Slack file information to be used in an image block or image element.**

This [file](https://docs.slack.dev/reference/objects/file-object) must be an image and you must provide either the URL or ID. In addition, the user posting these blocks must have access to this file. If both are provided then the payload will be rejected. Currently only `png`, `jpg`, `jpeg`, and `gif` Slack image files are supported.

#### Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `url` | string | This URL can be the `url_private` or the `permalink` of the Slack file. | Optional |
| `id` | string | Slack ID of the file. | Optional |

#### Example

The Slack file object must be used within the [image](https://docs.slack.dev/reference/block-kit/blocks/image-block) block.

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "image",
            "title": {
                "type": "plain_text",
                "text": "Please enjoy this photo of a kitten"
            },
            "block_id": "image4",
            "slack_file": {
                "url": "https://files.slack.com/files-pri/T0123456-F0123456/xyz.png"
            },
            "alt_text": "An incredibly cute kitten."
        }
    ]
}
```

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "image",
            "title": {
                "type": "plain_text",
                "text": "Please enjoy this photo of a kitten"
            },
            "block_id": "image4",
            "slack_file": {
                "id": "F0123456"
            },
            "alt_text": "An incredibly cute kitten."
        }
    ]
}
```