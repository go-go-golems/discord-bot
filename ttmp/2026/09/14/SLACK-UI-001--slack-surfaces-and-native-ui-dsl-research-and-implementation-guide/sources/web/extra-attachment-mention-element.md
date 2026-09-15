---
Title: "extra attachment mention element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/attachment-mention-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

This is a rich text element, compatible only with the [`rich_text`](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block) block. It must be used within the [`rich_text_list`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element), [`rich_text_quote`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-quote-element), or [`rich_text_section`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element) block element within the `rich_text` block's `elements` array.

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of object; in this case, "attachment\_mention". | Required |
| `url` | String | The URL of the app attachment or entity to reference. | Required |
| `text` | String | Fallback text if attachment not found. | Optional |
| `app_id` | String | The app ID that produced the unfurl. Used to fetch the app's profile for use in rendering. | Optional |
| `entity_id` | String | The Work Object entity ID when the attachment is a Work Object. | Optional |
| `icon_url` | String | An optional override of the icon URL. This will typically be used for adding the Work Object product icon, which can be different from the app's icon. | Optional |
| `channel_id` | String | The encoded channel ID where this attachment lives. | Optional |
| `ts` | String | The encoded message timestamp where this attachment lives. | Optional |
| `full_size_preview_enabled` | Boolean | Whether the work object supports full size preview. | Optional |
| `icon_name` | String | An optional icon name identifier for the attachment (e.g., sf-account, sf-record, sf-list for Salesforce attachments). | Optional |
| `reference_object_type` | String | An optional type identifier for the referenced object (e.g., list\_view, record for Salesforce attachments). | Optional |
| `product_name` | String | The product name for the Work Object (e.g., Google Docs, Google Sheets). Used to determine per-product click behavior preferences when the attachment is not available. | Optional |
| `style` | Object | An object of optional boolean properties that dictate style: `bold`, `italic`, `strike`, `highlight`, `client_highlight`, `underline`, and `unlink`. | Optional |