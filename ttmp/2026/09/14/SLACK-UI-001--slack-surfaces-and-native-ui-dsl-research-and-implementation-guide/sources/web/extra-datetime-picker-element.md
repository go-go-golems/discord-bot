---
Title: "extra datetime picker element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/datetime-picker-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Fields | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of element. In this case `type` is always `datetimepicker`. | Required |
| `action_id` | String | An identifier for the input value when the parent modal is submitted. You can use this when you receive a `view_submission` payload [to identify the value of the input element](https://docs.slack.dev/surfaces/modals#interactions). Should be unique among all other `action_id` s in the containing block. Maximum length is 255 characters. | Optional |
| `initial_date_time` | Integer | The initial date and time that is selected when the element is loaded, represented as a UNIX timestamp in seconds. This should be in the format of 10 digits, for example `1628633820` represents the date and time August 10th, 2021 at 03:17pm PST. | Optional |
| `confirm` | Object | A [confirm object](https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object) that defines an optional confirmation dialog that appears after a time is selected. | Optional |
| `focus_on_load` | Boolean | Indicates whether the element will be set to auto focus within the [`view object`](https://docs.slack.dev/reference/views). Only one element can be set to `true`. Defaults to `false`. | Optional |

## Usage Info

*Interactive component* - see our [guide to enabling interactivity](https://docs.slack.dev/interactivity/handling-user-interaction).

On desktop clients, the time picker will take the form of a dropdown list and the date picker will take the form of a dropdown calendar. Both options will have free-text entry for precise choices. On mobile clients, the time picker and date picker will use native UIs.

## Example

The datetime picker element must be used inside the [actions](https://docs.slack.dev/reference/block-kit/blocks/actions-block) block or [input](https://docs.slack.dev/reference/block-kit/blocks/input-block) block. This example shows an input block containing a datetime picker element: