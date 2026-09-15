---
Title: "rendered blocks"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

Blocks are a series of components that can be combined to create visually rich and compellingly interactive messages. You can include up to 50 blocks in each message, and 100 blocks in modals or Home tabs.

This page describes the JSON payloads that your app can use to generate each block. Select the block you'd like to build:

Name

Description

Surfaces

[`Actions`](https://docs.slack.dev/reference/block-kit/blocks/actions-block)

Holds multiple interactive elements.

Modals Messages Home tabs

[`Alert`](https://docs.slack.dev/reference/block-kit/blocks/alert-block)

Displays alerts, warnings, and informational messages.

Modals

[`Card`](https://docs.slack.dev/reference/block-kit/blocks/card-block)

Displays content in a card.

Modals Messages Home tabs

[`Container`](https://docs.slack.dev/reference/block-kit/blocks/container-block)

A general-purpose wrapper for grouping child blocks together, with a configurable size.

Messages Home tabs

[`Context`](https://docs.slack.dev/reference/block-kit/blocks/context-block)

Provides contextual info, which can include both images and text.

Modals Messages Home tabs

[`Context actions`](https://docs.slack.dev/reference/block-kit/blocks/context-actions-block)

Displays actions as contextual info, which can include both feedback buttons and icon buttons.

Messages

[`Data table`](https://docs.slack.dev/reference/block-kit/blocks/data-table-block)

Displays rich tables that support pagination, sorting, filtering, and interactivity.

Messages Home tabs

[`Data visualization`](https://docs.slack.dev/reference/block-kit/blocks/data-visualization-block)

Displays data visually in pie, bar, area, or line chart formats.

Messages Home tabs

[`Divider`](https://docs.slack.dev/reference/block-kit/blocks/divider-block)

Visually separates pieces of info inside of a message.

Modals Messages Home tabs

[`File`](https://docs.slack.dev/reference/block-kit/blocks/file-block)

Displays info about remote files.

Messages

[`Image`](https://docs.slack.dev/reference/block-kit/blocks/image-block)

Displays an image.

Modals Messages Home tabs

[`Input`](https://docs.slack.dev/reference/block-kit/blocks/input-block)

Collects information from users via elements.

Modals Messages Home tabs

[`Markdown`](https://docs.slack.dev/reference/block-kit/blocks/markdown-block)

Displays formatted markdown.

Messages

[`Rich text`](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block)

Displays formatted, structured representation of text.

Modals Messages Home tabs

[`Section`](https://docs.slack.dev/reference/block-kit/blocks/section-block)

Displays text, possibly alongside elements.

Modals Messages Home tabs

[`Table`](https://docs.slack.dev/reference/block-kit/blocks/table-block)

Displays structured information in a table.

Messages Home tabs

[`Task card`](https://docs.slack.dev/reference/block-kit/blocks/task-card-block)

Displays a single task, representing a single action.

Messages

[`Video`](https://docs.slack.dev/reference/block-kit/blocks/video-block)

Displays an embedded video player.

Modals Messages Home tabs