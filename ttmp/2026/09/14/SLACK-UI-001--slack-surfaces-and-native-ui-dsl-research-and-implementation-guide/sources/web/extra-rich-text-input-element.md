---
Title: "extra rich text input element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/rich-text-input-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

![An example of a rich text input element](https://docs.slack.dev/assets/images/bk_richtext_modal-524385b01f9c6354fd2988a32cfa529d.png)

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of element. In this case `type` is always `rich_text_input`. | Required |
| `action_id` | String | An identifier for the input value when the parent modal is submitted. You can use this when you receive a `view_submission` payload [to identify the value of the input element](https://docs.slack.dev/surfaces/modals#interactions). Should be unique in the containing block. Maximum length is 255 characters. | Required |
| `initial_value` | [Rich text](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block) | The initial value in the rich text input when it is loaded. | Optional |
| `dispatch_action_config` | Object | A [dispatch configuration object](https://docs.slack.dev/reference/block-kit/composition-objects/dispatch-action-configuration-object) that determines when during text input the element returns a [`block_actions`](https://docs.slack.dev/reference/interaction-payloads/block_actions-payload) payload. | Optional |
| `focus_on_load` | Boolean | Indicates whether the element will be set to auto focus within the [`view object`](https://docs.slack.dev/reference/views). Only one element can be set to `true`. Defaults to `false`. | Optional |
| `placeholder` | Object | A [`plain_text`](https://docs.slack.dev/reference/block-kit/composition-objects/text-object) object that defines the placeholder text shown in the plain-text input. Maximum length for the `text` in this field is 150 characters. | Optional |
| `min_lines` | Integer | The minimum number of visible text lines the input should display before scrolling. Controls the initial height of the input. Must be between 1 and 100. | Optional |
| `max_lines` | Integer | The maximum number of visible text lines the input can grow to before scrolling. Defaults to 8 when omitted. Must be between 1 and 100. | Optional |

## Example

The rich text input element must be used inside of the [input](https://docs.slack.dev/reference/block-kit/blocks/input-block) block. This example shows an input block containing a rich text input element.