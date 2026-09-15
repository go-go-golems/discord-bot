---
Title: "extra date element"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/block-elements/date-element"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

This is a rich text element, compatible only with the [`rich_text`](https://docs.slack.dev/reference/block-kit/blocks/rich-text-block) block. It must be used within the [`rich_text_list`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element), [`rich_text_quote`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-quote-element), or [`rich_text_section`](https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element) block element within the `rich_text` block's `elements` array.

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of object; in this case "date". | Required |
| `timestamp` | Number | A Unix timestamp for the date to be displayed in seconds. | Required |
| `format` | String | A template string containing curly-brace-enclosed tokens to substitute your provided `timestamp`. See details below. | Required |
| `timezone` | String | The timezone in which the date is in. | Optional |
| `url` | String | URL to link the entire `format` string to. | Optional |
| `fallback` | String | Text to display in place of the date should parsing, formatting or displaying fail. | Optional |
| `style` | Object | An object of optional boolean properties that dictate style: `bold`, `italic`, `strike`, `highlight`, `client_highlight`, and `underline`. | Optional |

### Date format strings

The following are the template strings allowed by the `format` property of the `date` element type.

- `{day_divider_pretty}`: Shows `today`, `yesterday` or `tomorrow` if applicable. Otherwise, if the date is in current year, uses the `{date_long}` format without the year. Otherwise, falls back to using the `{date_long}` format.
- `{date_num}`: Shows date as YYYY-MM-DD.
- `{date_slash}`: Shows date as DD/MM/YYYY (subject to locale preferences).
- `{date_long}`: Shows date as a long-form sentence including day-of-week, e.g. `Monday, December 23rd, 2013`.
- `{date_long_full}`: Shows date as a long-form sentence without day-of-week, e.g. `August 9, 2020`.
- `{date_long_pretty}`: Shows `yesterday`, `today` or `tomorrow`, otherwise uses the `{date_long}` format.
- `{date}`: Same as `{date_long_full}` but without the year.
- `{date_pretty}`: Shows `today`, `yesterday` or `tomorrow` if applicable, otherwise uses the `{date}` format.
- `{date_short}`: Shows date using short month names without day-of-week, e.g. `Aug 9, 2020`.
- `{date_short_pretty}`: Shows `today`, `yesterday` or `tomorrow` if applicable, otherwise uses the `{date_short}` format.
- `{time}`: Depending on user preferences, shows just the time-of-day portion of the timestamp using either 12 or 24 hour clock formats, e.g. `2:34 PM` or `14:34`.
- `{time_secs}`: Depending on user preferences, shows just the time-of-day portion of the timestamp using either 12 or 24 hour clock formats, including seconds, e.g. `2:34:56 PM` or `14:34:56`.
- `{ago}`: A human-readable period of time, e.g. `3 minutes ago`, `4 hours ago`, `2 days ago`.