---
Title: "extra time picker element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/time-picker-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Usage info

*Interactive component* - see our [guide to enabling interactivity](https://docs.slack.dev/interactivity/handling-user-interaction).

On desktop clients, this time picker will take the form of a dropdown list with free-text entry for precise choices. On mobile clients, the time picker will use native time picker UIs.

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of element. In this case `type` is always `timepicker`. | Required |
| `action_id` | String | An identifier for the action triggered when a time is selected. You can use this when you receive an interaction payload to [identify the source of the action](https://docs.slack.dev/interactivity/handling-user-interaction#payloads). Should be unique among all other `action_id` s in the containing block. Maximum length is 255 characters. | Optional |
| `initial_time` | String | The initial time that is selected when the element is loaded. This should be in the format `HH:mm`, where `HH` is the 24-hour format of an hour (00 to 23) and `mm` is minutes with leading zeros (00 to 59), for example `22:25` for 10:25pm. | Optional |
| `confirm` | Object | A [confirm object](https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object) that defines an optional confirmation dialog that appears after a time is selected. | Optional |
| `focus_on_load` | Boolean | Indicates whether the element will be set to auto focus within the [`view object`](https://docs.slack.dev/reference/views). Only one element can be set to `true`. Defaults to `false`. | Optional |
| `placeholder` | Object | A [`plain_text`](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) only text object that defines the placeholder text shown on the time picker. Maximum length for the `text` in this field is 150 characters. | Optional |
| `timezone` | String | A string in the IANA format, e.g. "America/Chicago". The timezone is displayed to end users as hint text underneath the time picker. It is also passed to the app upon certain interactions, such as `view_submission`. | Optional |

## Example

The time picker element must be used inside of the [section](https://docs.slack.dev/reference/block-kit/blocks/section-block) block, [actions](https://docs.slack.dev/reference/block-kit/blocks/actions-block) block, or [input](https://docs.slack.dev/reference/block-kit/blocks/input-block) block. This example shows a section block containing a time picker element, with the initial time set to 11:40am: