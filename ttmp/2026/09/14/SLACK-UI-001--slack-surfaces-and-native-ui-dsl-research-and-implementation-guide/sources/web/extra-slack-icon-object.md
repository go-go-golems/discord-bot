---
Title: "extra slack icon object"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/composition-objects/slack-icon-object"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

**Defines an object containing Slack icon information to be used in a card block.**

The Slack icon object must be used within the [card](https://docs.slack.dev/reference/block-kit/blocks/card-block) block.

#### Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `type` | string | Always `icon`. | Required |
| `name` | string | One of the following: `archive`, `book`, `bookmark`, `bot`, `bug`, `calendar`, `call`, `caret-left`, `caret-right`, `check`, `clipboard`, `code`, `comment`, `compass`, `copy`, `cube`, `download`, `edit`, `email`, `eye-closed`, `eye-open`, `file`, `flag`, `folder`, `gear`, `globe`, `heart`, `help`, `image`, `info`, `key`, `lightbulb`, `link`, `map`, `mobile`, `new-window`, `pin`, `plus`, `refine`, `refresh`, `rocket`, `save`, `screen`, `share`, `sparkle`, `star`, `star-filled`, `tag`, `thumbs-down`, `thumbs-up`, `trash`, `upload`, `user`, `warning` | Required |