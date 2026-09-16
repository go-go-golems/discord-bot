---
Title: "56 reference block kit composition objects confirmation dialog object"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

**Defines a dialog that adds a confirmation step to interactive elements.**

An object that defines a dialog that provides a confirmation step to any interactive element. This dialog will ask the user to confirm their action by offering a confirm and deny buttons.

#### Fields

| Field | Type | Description | Required? |
| --- | --- | --- | --- |
| `title` | Object | A [`plain_text`](#text) text object that defines the dialog's title. Maximum length for this field is 100 characters. | Required |
| `text` | Object | A [`plain_text`](#text) text object that defines the explanatory text that appears in the confirm dialog. Maximum length for the `text` in this field is 300 characters. | Required |
| `confirm` | Object | A [`plain_text`](#text) text object to define the text of the button that confirms the action. Maximum length for the `text` in this field is 30 characters. | Required |
| `deny` | Object | A [`plain_text`](#text) text object to define the text of the button that cancels the action. Maximum length for the `text` in this field is 30 characters. | Required |
| `style` | String | Defines the color scheme applied to the `confirm` button. A value of `danger` will display the button with a red background on desktop, or red text on mobile. A value of `primary` will display the button with a green background on desktop, or blue text on mobile. If this field is not provided, the default value will be `primary`. | Optional |

#### Example

The confirmation dialog object must be used within an interactive element. It is shown here within the [button](https://docs.slack.dev/reference/block-kit/block-elements/button-element) element.

---