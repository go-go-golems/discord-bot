---
Title: "rendered composition"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/composition-objects/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

Composition objects can be used inside of [block elements](https://docs.slack.dev/reference/block-kit/block-elements) and certain [message payload](https://docs.slack.dev/messaging#payloads) fields. They are common JSON object patterns that you'll encounter frequently when [building blocks](https://docs.slack.dev/block-kit) or [composing messages](https://docs.slack.dev/messaging).

The list of fields and values for each object describe the JSON that apps can use to generate each object.

Name

Description

[`Confirmation dialog object`](https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object)

Provides a dialog that adds a confirmation step to interactive elements.

[`Dispatch action configuration object`](https://docs.slack.dev/reference/block-kit/composition-objects/dispatch-action-configuration-object)

Defines when a plain-text input element will return a `block_actions` interaction payload.

[`Option object`](https://docs.slack.dev/reference/block-kit/composition-objects/option-object)

Represents a single item in a number of item selection elements.

[`Slack file object`](https://docs.slack.dev/reference/block-kit/composition-objects/slack-file-object)

Defines an object containing Slack file information to be used in an image block or image element.

[`Slack icon object`](https://docs.slack.dev/reference/block-kit/composition-objects/slack-icon-object)

Defines an object containing Slack icon information to be used in a card block.

[`Text object`](https://docs.slack.dev/reference/block-kit/composition-objects/text-object)

Defines text for many different blocks and elements.

[`Trigger object`](https://docs.slack.dev/reference/block-kit/composition-objects/trigger-object)

Defines an object containing trigger information.

[`Workflow object`](https://docs.slack.dev/reference/block-kit/composition-objects/workflow-object)

Defines an object containing workflow information.