---
Title: "29 reference block kit blocks markdown block"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/blocks/markdown-block/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | String | The type of block. For a markdown block, `type` is always `markdown`. | Required |
| `text` | String | The standard markdown-formatted text. The cumulative limit for all `markdown` blocks in a single payload is 12,000 characters. | Required |
| `block_id` | String | The `block_id` is ignored in markdown blocks and will not be retained. | Optional |

The following are supported `markdown` types.

| Markdown Type | Example | Result |
| --- | --- | --- |
| Bold | `**this is bold**` or   `__this is bold__` | **this is bold** |
| Italic | `*this text is italicized*` or   `_this text is italicized_` | *this text is italicized* |
| Bold and nested italic | `**this text is _extremely_ important**` | **this text is *extremely* important** |
| All bold and italic | `***all of this is important***` | ***all of this is important*** |
| Links | `[my text](https://www.google.com)` | [my text](https://google.com) |
| Lists (unordered) | ```markdown - first item in list - second item in a list - third item in a list ``` | - first item in a list - second item in a list - third item in a list |
| Lists (ordered) | ```markdown 1. first item in a list 2. second item in a list 3. third item in a list ``` | 1. first item in a list 2. second item in a list 3. third item in a list |
| Strikethrough | `~~this is strikethrough text~~` | ~~this is strikethrough text~~ |
| Headers (level 1) | `# Header 1` | Renders as a header. Note that all header levels are rendered at the same size. |
| Headers (level 2+) | ```markdown ## Header 2 ### Header 3  etc. ``` | Renders as a header. Note that all header levels are rendered at the same size. |
| In-line code | `` `this is my code` `` | `this is my code` |
| Block quote | `> this is a block quote` | > this is a block quote |
| Code blocks | ```markdown \`\`\` this is a code block \`\`\` ``` | ```markdown this is a code block ``` |
| Code blocks with syntax highlighting | ```markdown \`\`\`python print("hello") \`\`\` ``` | Renders as a code block with syntax highlighting for the specified language. |
| Dividers (horizontal rules) | `---` | Renders as a horizontal divider line. |
| Tables | ```markdown \| Col 1 \| Col 2 \| \| ----- \| ----- \| \| A      \| B      \| ``` | Renders as a formatted table. |
| Task lists | ```markdown - [ ] incomplete task - [x] completed task ``` | Renders as a task list with checkboxes. |
| Images | `![Logo](https://example.com/logo.png)` | This is translated as hyperlink text, i.e. [Logo](https://example.com/logo.png) |
| Escaping special characters:   - \\ backslash - \` backtick - \* asterisk - \_ underscore - curly braces - \[\] square brackets - () parentheses - \# hash mark - \+ plus sign - \- minus sign (hypen) - . dot - ! exclamation mark - ampersand | `\*This is special text\*` | \*This is special text\* |

## Usage info

This block can be used with [apps that use platform AI features](https://docs.slack.dev/ai/) when you expect a markdown response from an LLM that can get lost in translation rendering in Slack. Providing it in a markdown block leaves the translating to Slack to ensure your message appears as intended. Note that passing a single block may result in multiple blocks after translation.

![An example of an markdown block](data:image/png;base64,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)