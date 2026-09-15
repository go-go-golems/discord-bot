---
Title: "37 reference block kit blocks container block"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/container-block/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of block. For a container block, `type` is always `container`. | Required |
| `title` | [Plain text object](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) | Title of the container, using `plain_text` formatting. Maximum of 150 characters. | Optional (one of `title` or `rich_text_title` is required) |
| `rich_text_title` | [Rich text block](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block) | The title for the container as a [`rich_text`](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block) block. If both `title` and `rich_text_title` are provided, `rich_text_title` takes precedence. | Optional (one of `title` or `rich_text_title` is required) |
| `subtitle` | [Text object](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) | Subtitle of the container, using `plain_text` or `mrkdwn` formatting. Maximum of 150 characters. | Optional |
| `child_blocks` | Array | List of included blocks. Maximum of 10 blocks. | Required |
| `block_id` | String | A unique identifier for a block. If not specified, a `block_id` will be generated. | Optional |
| `width` | String | Sets the width of the container block. The `narrow`, `standard`, and `wide` options use a platform-determined constrained width. The `full` option expands to fill the available space. Default is `standard`. | Optional |
| `icon` | [Image element](https://docs.slack.dev/reference/block-kit/block-elements/image-element) | Link to the small image used next to the card's title and subtitle. Maximum length of 3000 characters. The `alt_text` property has a maximum length of 2000 characters. | Optional |
| `is_collapsible` | Boolean | When `true`, the block can be collapsed to show only the title. Defaults to `false`. | Optional |
| `default_collapsed` | Boolean | When `true` and `is_collapsible` are both `true`, the block initially renders in a collapsed state. Defaults to `false`. | Optional |
| `has_header_divider` | Boolean | When true, a visible border is rendered below the header to visually separate it from the content. Only applies when the block is not collapsible. Defaults to `false`. | Optional |

## Supported child blocks

- [actions](https://docs.slack.dev/reference/block-kit/blocks/actions-block)
- [context](https://docs.slack.dev/reference/block-kit/blocks/context-block)
- [divider](https://docs.slack.dev/reference/block-kit/blocks/divider-block)
- [file](https://docs.slack.dev/reference/block-kit/blocks/file-block)
- [header](https://docs.slack.dev/reference/block-kit/blocks/header-block)
- [image](https://docs.slack.dev/reference/block-kit/blocks/image-block)
- [input](https://docs.slack.dev/reference/block-kit/blocks/input-block)
- [rich\_text](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block)
- [section](https://docs.slack.dev/reference/block-kit/blocks/section-block)
- [table](https://docs.slack.dev/reference/block-kit/blocks/table-block)
- [video](https://docs.slack.dev/reference/block-kit/blocks/video-block)

## Example

A sample container block: