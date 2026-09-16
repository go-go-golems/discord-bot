---
Title: "07 block kit"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/block-kit/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

The Block Kit UI framework is built with *blocks*, *block elements* and *composition objects*.

*Blocks* are visual components that can be arranged to create app layouts. Apps can add blocks to *surfaces* like [the Home tab](https://docs.slack.dev/surfaces/app-home), [messages](https://docs.slack.dev/messaging) and [modals](https://docs.slack.dev/surfaces/modals). You can include up to 50 blocks in each message, and 100 blocks in modals or Home tabs.

Blocks may also contain *block elements*. Block elements are usually interactive components, such as buttons and menus.

Blocks and block elements are built with *composition objects*. Composition objects define text, options, or other interactive features within certain blocks and block elements.

![](https://docs.slack.dev/assets/images/bk_landing_bkb-e64c290c97543b50e0b09c0b291c7c78.png)

Whether you're composing layouts for modals, messages, or tabs, the Block Kit building experience is the same — your app will be sculpting specially-structured JSON to express itself. The result is clear, interactive communication between your app and its users.

Read on to learn how you can construct the stacks of blocks that app surfaces love to consume.

Eager to see Block Kit in action? Take a peek in Block Kit Builder.

---

## Placing blocks within surfaces

Blocks are used within all [app surfaces](https://docs.slack.dev/surfaces): [Home tabs](https://docs.slack.dev/surfaces/app-home), [messages](https://docs.slack.dev/messaging) and [modals](https://docs.slack.dev/surfaces/modals) can all be designed using blocks.

Each of them uses a `blocks` array that you prepare by [stacking individual blocks together](#stack_of_blocks).

Check out [app surfaces](https://docs.slack.dev/surfaces) to learn more about using these different surfaces, and how to add blocks to your app's [Home tab](https://docs.slack.dev/surfaces/app-home#composing), [messages](https://docs.slack.dev/messaging), and [modals](https://docs.slack.dev/surfaces/modals#composing_modal).

## Building blocks

There's no special setup needed to start using blocks in [app surfaces](https://docs.slack.dev/surfaces). However, just as when you open a pack of generic, colorful, interlocking plastic bricks, you should read the instructions first.

### Defining a single block

Each block is represented in our APIs as a JSON object. Here's an example of a [`section`](https://docs.slack.dev/reference/block-kit/blocks/section-block) block:

```json
{
  "type": "section",
  "text": {
    "type": "mrkdwn",
    "text": "New Paid Time Off request from <example.com|Fred Enriquez>\n\n<https://example.com|View request>"
  }
}
```

Every block contains a `type` field — specifying which of the [available blocks](https://docs.slack.dev/reference/block-kit/blocks) to use — along with other fields that describe the content of the block.

[Block Kit Builder](https://api.slack.com/tools/block-kit-builder) is a visual prototyping sandbox that will let you choose from, configure, and preview all the available blocks.

If you want to skip the builder, the [block reference guide](https://docs.slack.dev/reference/block-kit/blocks) contains the specifications of every block, and the JSON fields required for each of them.

### Stacking multiple blocks

Individual blocks can be stacked together to create complex visual layouts.

When you've chosen each of the blocks you want in your layout, place each of them in an array, in visual order, like this:

[Block Kit Builder](https://api.slack.com/tools/block-kit-builder) will allow you to drag, drop, and rearrange blocks to design and preview Block Kit layouts.

Alternatively you can use the [block reference guide](https://docs.slack.dev/reference/block-kit/blocks) to manually generate a complete `blocks` array, like the one shown above.

Your newly created array of blocks can be used [with a range of different app surfaces](#adding_blocks).

### Accessibility considerations

When posting messages, it is expected behavior that screen readers will default to the top-level `text` field of the post, and will not read the content of any interior `blocks` in the underlying structure of the message. Therefore, to make an accessible app, you must either:

- include all necessary content for screen reader users in the top-level `text` field of your message, or
- do not include a top-level `text` field if the message has `blocks`, and allow Slack attempt to build it for you by appending content from supported `blocks` to be read by the screen reader.

---

## Adding interactivity to blocks with block elements

Blocks can be made to interact with users via Block Kit *elements*. Elements include interactive components such as buttons, menus and text inputs.

Here's an example of a [`button`](https://docs.slack.dev/reference/block-kit/block-elements/button-element) element within a [`section`](https://docs.slack.dev/reference/block-kit/blocks/section-block) block.

```json
{
  "blocks": [
    {
      "type": "section",
      "text": {
        "type": "mrkdwn",
        "text": "This is a section block with a button."
      },
      "accessory": {
        "type": "button",
        "text": {
          "type": "plain_text",
          "text": "Click Me",
          "emoji": true
        },
        "value": "click_me_123",
        "action_id": "button-action"
      }
    }
  ]
}
```

When you add an interactive component to a surface in your Slack app, you've opened the door to user interaction. People will push your app's buttons and expecting a helpful and prompt reaction.

Apps need to handle the requests that start to flow their way, and respond appropriately. Follow our [guide to handling user interaction](https://docs.slack.dev/interactivity/handling-user-interaction) to prepare your app for the interactivity that Block Kit will inspire.

Block Kit builder allows you to add elements to blocks as well. Give it a try! Alternatively, read the [Block Kit element reference guide](https://docs.slack.dev/reference/block-kit/block-elements) for all the info you'll need for manually implementing individual elements.

---

## Onward

Check out the following guides for everything blocks:

- [Blocks](https://docs.slack.dev/reference/block-kit/blocks)
- [Block elements, including interactive components](https://docs.slack.dev/reference/block-kit/block-elements)
- [Composition objects](https://docs.slack.dev/reference/block-kit/composition-objects)
- [View objects](https://docs.slack.dev/reference/views)