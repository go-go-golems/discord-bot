---
Title: "33 reference block kit blocks data table block"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/data-table-block/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

The data table block is a rich table that supports pagination, sorting, filtering, and rich interactivity, such as opening a [Work Object flexpane](https://docs.slack.dev/messaging/work-objects-overview#flexpane) or clickable links in cells. This is different from the existing [table block](https://docs.slack.dev/reference/block-kit/blocks/table-block), which only supports filtering and basic interactivity.

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of block. For a data table block, `type` is always `data_table`. | Required |
| `rows` | Array | An array consisting of table rows. | Required |
| `block_id` | String | A unique identifier for a block. If not specified, a `block_id` will be generated. | Optional |
| `page_size` | Integer | Number of rows per page. Min `1`, Max `100`. Defaults to `5` if omitted. | Optional |
| `caption` | String | A caption for the table; used as the value for the HTML caption element. | Required |
| `row_header_column_index` | Integer | The 0-based index of the column that uniquely identifies each row (the row header). This column is treated as the row's primary identifier for screen readers. Defaults to 0 if omitted. | Optional |

## Usage info

You can use `rich_text`, `raw_text` (simple text), or `raw_number` (numeric values) for cell content. The first row of the table is a header, and `rich_text` cannot be used for header cells. You can have a minimum of 2 rows (1 regular row plus the header) and a maximum of 201 rows (200 regular rows plus the header). All rows must have the same number of values.

A single table's character count across all cells cannot exceed 20,000 characters. Additionally, the aggregate character count across all table cells for a single message cannot exceed 20,000 characters. Large tables should be broken up into separate messages.

Sorting rows by column is done alphabetically by default. If a column contains cells all of type `raw_number`, a numeric sort will be performed instead. You can have a minimum of 1 column and a maximum of 20 columns.

### Schema for raw\_text

```json
"properties": {
    "type": {
        "type": "string",
        "enum": ["raw_text"]
    },
    "text": {
        "type": "string",
        "minLength": 1
    }
}
```

### Schema for raw\_number

```json
"properties": {
    "type": {
        "type": "string",
        "enum": ["raw_number"]
    },
    "value": {
        "type": "number"
    },
    "text": {
        "type": "string",
        "minLength": 1
    }
}
```

## Example

A sample data table block: