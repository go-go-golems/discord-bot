---
Title: "extra image element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/image-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of element. In this case `type` is always `image`. | Required |
| `alt_text` | String | A plain-text summary of the image. This should not contain any markup. | Required |
| `image_url` | String | The URL for a publicly hosted image. You must provide either an `image_url` or `slack_file`. Maximum length for this field is 3000 characters. | Optional |
| `slack_file` | Object | A [Slack image file object](https://docs.slack.dev/reference/block-kit/composition-objects/slack-file-object) that defines the source of the image. | Optional |

## Usage info

Use the [`image`](https://docs.slack.dev/reference/block-kit/blocks/image-block) block if you want a block with *only* an image in it.

## Examples

The image element must be used inside of the [section](https://docs.slack.dev/reference/block-kit/blocks/section-block) block or the [context](https://docs.slack.dev/reference/block-kit/blocks/context-block) block.

![An example of an image element](https://docs.slack.dev/assets/images/bk_image_ele_example-cfe2dcdf8ce5cc28217eb77cec112c27.png)

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "section",
            "block_id": "section567",
            "text": {
                "type": "mrkdwn",
                "text": "This is a section block with an accessory image."
            },
            "accessory": {
                "type": "image",
                "image_url": "https://pbs.twimg.com/profile_images/625633822235693056/lNGUneLX_400x400.jpg",
                "alt_text": "cute cat"
            }
        }
    ]
}
```

[View in Block Kit Builder](https://api.slack.com/tools/block-kit-builder?blocks=%5B%0A%09%7B%0A%09%09%22type%22%3A%20%22section%22%2C%0A%09%09%22block_id%22%3A%20%22section567%22%2C%0A%09%09%22text%22%3A%20%7B%0A%09%09%09%22type%22%3A%20%22mrkdwn%22%2C%0A%09%09%09%22text%22%3A%20%22This%20is%20a%20section%20block%20with%20an%20accessory%20image.%22%0A%09%09%7D%2C%0A%09%09%22accessory%22%3A%20%7B%0A%09%09%09%22type%22%3A%20%22image%22%2C%0A%09%09%09%22image_url%22%3A%20%22https%3A%2F%2Fpbs.twimg.com%2Fprofile_images%2F625633822235693056%2FlNGUneLX_400x400.jpg%22%2C%0A%09%09%09%22alt_text%22%3A%20%22cute%20cat%22%0A%09%09%7D%0A%09%7D%0A%5D)

An image block using `slack_file` with a `url`:

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "section",
            "block_id": "section567",
            "text": {
                "type": "mrkdwn",
                "text": "This is a section block with an accessory image."
            },
            "accessory": {
                "type": "image",
                "slack_file": {
                    "url": "https://files.slack.com/files-pri/T0123456-F0123456/xyz.png"
                },
                "alt_text": "Slack file object."
            }
        }
    ]
}
```

An image block using `slack_file` with a `id`:

- JSON
- Python Slack SDK
- Node Slack SDK
- Java Slack SDK

```json
{
    "blocks": [
        {
            "type": "section",
            "block_id": "section567",
            "text": {
                "type": "mrkdwn",
                "text": "This is a section block with an accessory image."
            },
            "accessory": {
                "type": "image",
                "slack_file": {
                    "id": "F01234567"
                },
                "alt_text": "Slack file object."
            }
        }
    ]
}
```