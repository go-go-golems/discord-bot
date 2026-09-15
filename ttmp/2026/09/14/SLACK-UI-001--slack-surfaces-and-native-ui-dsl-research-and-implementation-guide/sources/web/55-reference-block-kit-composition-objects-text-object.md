---
Title: "55 reference block kit composition objects text object"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/composition-objects/text-object/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

**Defines an object containing some text.**

Formatted either as `plain_text` or using [`mrkdwn`](https://docs.slack.dev/messaging/formatting-message-text), our proprietary contribution to the much beloved [Markdown standard](https://xkcd.com/927/).

#### Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The formatting to use for this text object. Can be one of `plain_text` or `mrkdwn`. | Required |
| `text` | String | The text for the block. This field accepts any of the standard [text formatting markup](https://docs.slack.dev/messaging/formatting-message-text) when `type` is `mrkdwn`. The minimum length is 1 and maximum length is 3000 characters. | Required |
| `emoji` | Boolean | Indicates whether emojis in a text field should be escaped into the colon emoji format. This field is only usable when `type` is `plain_text`. | Optional |
| `verbatim` | Boolean | When set to `false` (as is default) URLs will be auto-converted into links, conversation names will be link-ified, and certain mentions will be [automatically parsed](https://docs.slack.dev/messaging/formatting-message-text#automatic-parsing). When set to `true`, Slack will continue to process all markdown formatting and [manual parsing strings](https://docs.slack.dev/messaging/formatting-message-text#advanced), but it won’t modify any plain-text content. For example, channel names will not be hyperlinked. This field is only usable when `type` is `mrkdwn`. |  |

#### Example

The text object must be used within another block or element, such as the [header](https://docs.slack.dev/reference/block-kit/blocks/header-block) block, [section](https://docs.slack.dev/reference/block-kit/blocks/section-block) block, [button](https://docs.slack.dev/reference/block-kit/block-elements/button-element) element, [icon](https://docs.slack.dev/reference/block-kit/block-elements/icon-button-element) button element, or [workflow button](https://docs.slack.dev/reference/block-kit/block-elements/workflow-button-element) element. This example shows a section block containing a text object.

---