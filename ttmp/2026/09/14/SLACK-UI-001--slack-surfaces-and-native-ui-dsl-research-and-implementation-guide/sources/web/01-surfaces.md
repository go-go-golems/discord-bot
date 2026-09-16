---
Title: "01 surfaces"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/surfaces/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

Surfaces are where users engage with your app. Slack offers several surfaces where users can access and interact with your app.

Some of these surfaces serve as entry points: places where a user's action can initially invoke your agent. Most surfaces can be built using [Block Kit](https://docs.slack.dev/block-kit) layout blocks and elements. Our [guide to building block layouts](https://docs.slack.dev/block-kit) will help you learn how. Canvases, on the other hand, use [markdown](https://docs.slack.dev/messaging/formatting-message-text#basic-formatting) for content formatting.

---

## App Home

The App Home is a private, one-to-one space in Slack shared by a user and an app. The Home tab, a specific App Home view, is an optional ever-present space, retaining its content and state until the app chooses to update it.

Present each of your users with a unique Home tab just for them, always found in the same place.

Although not every app needs to have a Home tab, the 'always-on' nature of the space makes it an important surface for many Slack apps.

The App Home is not available for apps created with the Deno Slack SDK.

![](https://docs.slack.dev/assets/images/app_home_abstract-f4c341508b05b02d5fb0aac5a7ad61ba.png)

➡️ **To get started with App Home**, read our [App Home guide](https://docs.slack.dev/surfaces/app-home).

---

## Canvases

Canvases are built-in documents in Slack, existing either tied to a channel or as a standalone space.

Use canvases to store channel guidelines or instructions, welcome new team members with an onboarding flow, or present project updates with canvases. Add links to team resources, helpful videos, and even kick off a workflow from a canvas.

![](https://docs.slack.dev/assets/images/canvas-5260e8e383df4e487e8bfca3c8883b24.png)

➡️ **To get started with canvases**, read our [Canvas guide](https://docs.slack.dev/surfaces/canvases).

---

## Lists

![](https://docs.slack.dev/assets/images/lists-bc547b497580bd64f82e996a24d2bcfc.png)

➡️ **To get started with Lists**, read our [Lists guide](https://docs.slack.dev/surfaces/lists).

---

## Messages

App-published messages are dynamic yet transient spaces. They allow users to complete workflows as Slack conversations.

Apps can [send messages](https://docs.slack.dev/messaging) whenever they want to, as long as they have the relevant permissions and access. Our [guide to formatting text for app surfaces](https://docs.slack.dev/messaging/formatting-message-text) will show you what formatting is possible.

When an app is invoked, it can respond with a message. Further action can flow from that message, forming a conversational interface connected to any of Slack's features.

![](https://docs.slack.dev/assets/images/message-abstract-06be210d128e91ff97e3ca6d791ef7d9.png)

➡️ **To get started with messages**, read our [Messages guide](https://docs.slack.dev/messaging).

✨ **To level up your messages with interactive components such as buttons and select menus**, read our Creating interactive messages guide for traditional [Slack apps](https://docs.slack.dev/messaging/creating-interactive-messages) or [apps created with the Deno Slack SDK](https://docs.slack.dev/tools/deno-slack-sdk/guides/creating-an-interactive-message).

---

## Modals

Modals are prominent and pervasive spaces ideal for requesting and collecting data from users, or temporarily displaying dynamic and interactive information.

Modals appear in front of any other interface element in Slack. As a result, they are short-lived and invoked only when a specific task is to be completed. Apps can *only* create modals in response to [user invocation](https://docs.slack.dev/interactivity#user), such as a [shortcut](https://docs.slack.dev/interactivity/implementing-shortcuts).

Modals contain one to three [**views**](https://docs.slack.dev/surfaces/modals#lifecycle) that can be chained together to create complex, non-linear workflows.

![](https://docs.slack.dev/assets/images/modal-abstract-f84c7b1e74a116b1376d94dd07121db0.png)

➡️ **To get started with modals**, read our [Modals guide](https://docs.slack.dev/surfaces/modals).

---

## Split view

When the **Agents** feature is enabled in your [app settings](https://api.slack.com/apps), the top nav entry point and agent container are available.

The native split pane in the Slack client is accessible from the top bar. Use it for conversational agent interactions: back-and-forth dialogue, multi-turn reasoning, and contextual responses. Interactions here can be implemented using the Bolt `Assistant` class.

The [Bolt](https://docs.slack.dev/tools#bolt) `Assistant` class wraps related API events into handler callbacks, providing a single lifecycle for the agent container. Read more about setting suggested prompts and user interactions in [Developing agents](https://docs.slack.dev/ai/developing-agents#thread-started).

![Image of split view](https://docs.slack.dev/img/guides/ai_container/splitview.png)

➡️ **To get started with split view**, read our [Split view guide](https://docs.slack.dev/surfaces/split-view).

---

## Using app surfaces together

App surfaces can also be used together to create a rich interactive experience for your users. For example, imagine the following Task App, which presents a task dashboard that resides in the app's [Home tab](https://docs.slack.dev/surfaces/app-home):

1. A user can click a [button](https://docs.slack.dev/reference/block-kit/block-elements/button-element) to add a task.
2. The user is then presented with a [modal](https://docs.slack.dev/surfaces/modals) to [enter](https://docs.slack.dev/reference/block-kit/blocks/input-block) some [plain text](https://docs.slack.dev/reference/block-kit/block-elements/plain-text-input-element) and [select from a list of categories](https://docs.slack.dev/reference/block-kit/block-elements/select-menu-element).
3. Upon submitting, a [message](https://docs.slack.dev/messaging) is sent to a triage channel elsewhere in Slack.
4. Finally, a different user in the triage channel can click a [button](https://docs.slack.dev/reference/block-kit/block-elements/button-element) to claim that task.

Explore all the possibilities and get some tips and inspiration by reading our [guides to planning Slack apps](https://docs.slack.dev/concepts/app-design).