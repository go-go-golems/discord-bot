---
Title: "rendered elements"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

This page lists the JSON payloads that your app can use to generate each element. Select the block element you'd like to build:

Name

Description

Blocks

Surfaces

[`Attachment mention`](https://docs.slack.dev/reference/block-kit/block-elements/attachment-mention-element)

Renders as a rich app attachment or entity reference.

Rich text

Modals Messages Home tabs

[`Broadcast`](https://docs.slack.dev/reference/block-kit/block-elements/broadcast-element)

Displays a broadcast mention such as here, channel, or everyone.

Rich text

Modals Messages Home tabs

[`Button`](https://docs.slack.dev/reference/block-kit/block-elements/button-element)

Allows users a direct path to performing basic actions.

Section Actions

Modals Messages Home tabs

[`Canvas`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-element)

Renders as a link to a canvas.

Rich text

Modals Messages Home tabs

[`Canvas message unfurl`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-message-unfurl-element)

Renders as an inline preview of a message inside a canvas.

Rich text

Modals Messages Home tabs

[`Canvas user mention`](https://docs.slack.dev/reference/block-kit/block-elements/canvas-user-mention-element)

Renders as a user mention in canvas content.

Rich text

Modals Messages Home tabs

[`Channel`](https://docs.slack.dev/reference/block-kit/block-elements/channel-element)

Renders as a mention of a channel.

Rich text

Modals Messages Home tabs

[`Checkboxes`](https://docs.slack.dev/reference/block-kit/block-elements/checkboxes-element)

Allows users to choose multiple items from a list of options.

Section Actions Input

Modals Messages Home tabs

[`Citation`](https://docs.slack.dev/reference/block-kit/block-elements/citation-element)

Renders as an AI citation (file, external, web, message, or memory).

Rich text

Modals Messages Home tabs

[`Color`](https://docs.slack.dev/reference/block-kit/block-elements/color-element)

Displays a color swatch from a hex value.

Rich text

Modals Messages Home tabs

[`Date`](https://docs.slack.dev/reference/block-kit/block-elements/date-element)

Displays a formatted, localized date.

Rich text

Modals Messages Home tabs

[`Date picker`](https://docs.slack.dev/reference/block-kit/block-elements/date-picker-element)

Allows users to select a date from a calendar style UI.

Section Actions Input

Modals Messages Home tabs

[`Datetime picker`](https://docs.slack.dev/reference/block-kit/block-elements/datetime-picker-element)

Allows users to select both a date and a time of day.

Actions Input

Modals Messages

[`Email input`](https://docs.slack.dev/reference/block-kit/block-elements/email-input-element)

Allows user to enter an email into a single-line field.

Input

Modals

[`Emoji`](https://docs.slack.dev/reference/block-kit/block-elements/emoji-element)

Displays an emoji.

Rich text

Modals Messages Home tabs

[`Feedback buttons`](https://docs.slack.dev/reference/block-kit/block-elements/feedback-buttons-element)

Buttons to indicate positive or negative feedback.

Context actions

Messages

[`File`](https://docs.slack.dev/reference/block-kit/block-elements/file-element)

Renders as a link to a Slack file.

Rich text

Modals Messages Home tabs

[`File input`](https://docs.slack.dev/reference/block-kit/block-elements/file-input-element)

Allows user to upload files.

Input

Modals

[`Icon button`](https://docs.slack.dev/reference/block-kit/block-elements/icon-button-element)

An icon button to perform actions.

Context actions

Messages

[`Image`](https://docs.slack.dev/reference/block-kit/block-elements/image-element)

Displays an image as part of a larger block of content.

Section Context

Modals Messages Home tabs

[`Link`](https://docs.slack.dev/reference/block-kit/block-elements/link-element)

Displays a hyperlink.

Rich text

Modals Messages Home tabs

[`List record`](https://docs.slack.dev/reference/block-kit/block-elements/list-record-element)

Renders as a link to a Slack list record.

Rich text

Modals Messages Home tabs

[`Message mention`](https://docs.slack.dev/reference/block-kit/block-elements/message-mention-element)

Renders as a link to a Slack message.

Rich text

Modals Messages Home tabs

[`Number input`](https://docs.slack.dev/reference/block-kit/block-elements/number-input-element)

Allows user to enter a number into a single-line field.

Input

Modals

[`Plain-text input`](https://docs.slack.dev/reference/block-kit/block-elements/plain-text-input-element)

Allows users to enter freeform text data into a single-line or multi-line field.

Input

Modals Messages Home tabs

[`Radio button group`](https://docs.slack.dev/reference/block-kit/block-elements/radio-button-group-element)

Allows users to choose one item from a list of possible options.

Section Actions Input

Modals Messages Home tabs

[`Rich text input`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-input-element)

Allows users to enter formatted text in a WYSIWYG composer, offering the same messaging writing experience as in Slack.

Input Table

Modals Home tabs

[`Rich text list`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element)

Displays a list of rich text items.

Rich text

Modals Messages Home tabs

[`Rich text preformatted`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-preformatted-element)

Displays a preformatted rich text element.

Rich text

Modals Messages Home tabs

[`Rich text quote`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-quote-element)

Displays a rich text quote block.

Rich text

Modals Messages Home tabs

[`Rich text section`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element)

A section element that holds rich text elements.

Rich text

Modals Messages Home tabs

[`Salesforce data field`](https://docs.slack.dev/reference/block-kit/block-elements/salesforce-data-field-element)

Renders as a Salesforce data field reference.

Rich text

Modals Messages Home tabs

[`Tag`](https://docs.slack.dev/reference/block-kit/block-elements/tag-element)

Renders as a colored tag or pill.

Rich text

Modals Messages Home tabs

[`Team`](https://docs.slack.dev/reference/block-kit/block-elements/team-element)

Renders as a mention of a workspace or team.

Rich text

Modals Messages Home tabs

[`Text`](https://docs.slack.dev/reference/block-kit/block-elements/text-element)

Displays text, optionally with styling.

Rich text

Modals Messages Home tabs

[`Time picker`](https://docs.slack.dev/reference/block-kit/block-elements/time-picker-element)

Allows users to enter numerical data into a single-line field.

Section Actions Input

Modals Messages Home tabs

[`URL input`](https://docs.slack.dev/reference/block-kit/block-elements/url-input-element)

Allows user to enter a URL into a single-line field.

Input

Modals

[`URL source`](https://docs.slack.dev/reference/block-kit/block-elements/url-source-element)

Displays a URL source for referencing within a task card block.

Task card

Messages

[`User`](https://docs.slack.dev/reference/block-kit/block-elements/user-element)

Renders as a mention of a user.

Rich text

Modals Messages Home tabs

[`Usergroup`](https://docs.slack.dev/reference/block-kit/block-elements/usergroup-element)

Renders as a mention of a user group.

Rich text

Modals Messages Home tabs

[`Work object mention`](https://docs.slack.dev/reference/block-kit/block-elements/work-object-mention-element)

Renders as a Work Object reference.

Rich text

Modals Messages Home tabs

[`Workflow button`](https://docs.slack.dev/reference/block-kit/block-elements/workflow-button-element)

Allows users to run a link trigger with customizable inputs.

Section Actions

Messages

[`Workflow mention`](https://docs.slack.dev/reference/block-kit/block-elements/workflow-mention-element)

Renders as a link to a workflow.

Rich text

Modals Messages Home tabs