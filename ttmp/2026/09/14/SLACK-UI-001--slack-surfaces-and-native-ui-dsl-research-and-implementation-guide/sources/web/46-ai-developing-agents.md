---
Title: "46 ai developing agents"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/ai/developing-agents/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

This guide takes you through developing the response loop of an agent. The response loop is a cycle of receive input → reason → call tools → stream/render output → repeat if needed.

## Enabling the agent feature

First things first: follow the Quickstart guide to [create an app](https://docs.slack.dev/quickstart#cli-and-bolt). Once completed, open the [app settings](https://api.slack.com/apps), find the **Agents** feature in the sidebar, and enable it.

The [`assistant:write`](https://docs.slack.dev/reference/scopes/assistant.write) scope is needed for this, and thus is automatically added to your app when you enable the feature in the app settings. It also allows your agent to take advantage of suggested prompts and thread title customization. You'll also want to add the [`chat:write`](https://docs.slack.dev/reference/scopes/chat.write) scope; do so under the **OAuth & Permissions** section of the navigation. Remember to provide an agent overview in the setup too!

### Subscribe to events

In the **Event Subscriptions** menu tab of the app settings, toggle on the **Enable Events** setting. Under **Subscribe to bot events**, subscribe to the [`app_home_opened`](https://docs.slack.dev/reference/events/app_home_opened), [`app_context_changed`](https://docs.slack.dev/reference/events/app_context_changed), [`message.im`](https://docs.slack.dev/reference/events/message.im), [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped), and [`agent_session_title_changed`](https://docs.slack.dev/reference/events/agent_session_title_changed) events.

This section is split between the Agent and Assistant messaging experiences. The Assistant messaging experience is the legacy implementation of developing apps using soon-to-be-deprecated methods. Follow the Agent messaging experience below for new apps. To migrate an old app to the new messaging experience, check out [Migrating to the Agent messaging experience](https://docs.slack.dev/ai/migrating-to-agent-messaging).

### Onboarding and welcome message

Send a call to action or suggest next steps when a user interacts with an agent for the very first time. Once this requirement is completed, optimize for repeat use and avoid repetitive prompts and 'getting started' types of messaging. This is important especially when it is necessary for the user to sign in, connect an account, agree to terms of service, or review a code of conduct.

## Setting suggested prompts

Present the user with suggested prompts using the [`assistant.threads.setSuggestedPrompts`](https://docs.slack.dev/reference/methods/assistant.threads.setSuggestedPrompts) API method. Suggested prompts live at the top of the Messages tab. We recommend using the Bolt framework to handle the details for you.

![suggested prompts](https://docs.slack.dev/img/guides/ai_container/suggestedprompts.png)

## Listening for the message.im event

The user then will type a message or click on a prompt which triggers a [`message.im`](https://docs.slack.dev/reference/events/message.im) event. The event is the same whether the user clicked the suggested prompt or typed it manually. Users can message your app via the container or through your app's Messages tab.

After the user sends a new message, calling the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method with `status: "processing"` on that thread opens the thread to keep the conversation going. Only do this if you intend to reply in thread. (The legacy [`assistant.threads.setStatus`](https://docs.slack.dev/reference/methods/assistant.threads.setStatus) method behaves the same way through the [compatibility bridge](https://docs.slack.dev/ai/migrating-to-agent-messaging#compatibility-bridge).)

Your app can respond to the user directly or it can pass back the `thread_ts` parameter to continue in the same thread. In most situations, you will want to call the [`chat.postMessage`](https://docs.slack.dev/reference/methods/chat.postMessage) method with the `thread_ts` parameter.

When your app receives the `thread_ts` parameter, you can retrieve the conversation by using `thread_ts` as the unique identifier. This is useful if your app stores the long-lived context or the state of a thread.

You can also fetch previous thread messages using the [`conversations.replies`](https://docs.slack.dev/reference/methods/conversations.replies) method and choose which other messages from the conversation to include in the LLM prompt or your app logic.

*Note: @-mentions in channels can happen like they do today; whether you support this or not is up to you. You can engage with the user or ask them to use the container to converse with your app.*

Your app should then call the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method with `status: "processing"` to display the loading indicator in the container. We recommend doing so immediately for the user's benefit.

Loading states indicate to your user that the app is working on a response. While a session is in `processing`, Slack shows a standard loading UX, along with a stop button if your app subscribes to the [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped) event. Custom loading messages are not supported by the `agents.sessions.setStatus` method.

![loading state](https://docs.slack.dev/img/guides/ai_container/loadingstates.png)

We recommend using the Bolt framework to handle the details for you.

## Responding to the user

Once your app finishes its work, call the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method with `status: "active"` to clear the loading indicator and mark the session ready for the next prompt.

Formulate and send a response, using [text streaming](https://docs.slack.dev/ai/developing-agents#streaming).

### Text streaming

Text streaming is handled by three different API methods: [`chat.startStream`](https://docs.slack.dev/reference/methods/chat.startStream), [`chat.appendStream`](https://docs.slack.dev/reference/methods/chat.appendStream), and [`chat.stopStream`](https://docs.slack.dev/reference/methods/chat.stopStream). These allow the user to see the response from the LLM as a text stream, rather than a single block of text sent all at once, providing closer alignment with expected behavior from other major LLM tools.

<video><source src="/img/guides/ai_container/textstreaming.mp4" type="video/mp4"></video>

When using text streaming, there are a couple of caveats to keep in mind. [Blocks](https://docs.slack.dev/block-kit) may be used in the `chat.stopStream` method, but not the `chat.startStream` or `chat.appendStream` method, in order to prevent having them broken up. Also, unfurling is disabled in streaming messages.

If you're a Python or JavaScript fan, our Bolt frameworks in those languages have a streamer utility to allow you to quickly implement the functionality of these API methods into your apps.

#### Display modes for streaming text

Use blocks from Block Kit to help visualize the response. Tasks can then be displayed using [task card](https://docs.slack.dev/reference/block-kit/blocks/task-card-block) blocks along with the comprehensive [plan](https://docs.slack.dev/reference/block-kit/blocks/plan-block) display. Task cards display individual steps your agent is taking; a plan groups those tasks together.

Apps can display a [task update](https://docs.slack.dev/reference/methods/chat.appendStream#task_update-chunks) view for users to better understand what the app is doing. The task update display mode is best suited for short tasks with narration text. It can be in one of three different states: `in_progress`, `completed`, and `error`.

![Task update display mode](https://docs.slack.dev/img/guides/timeline.png)

The plan display mode uses the [plan block](https://docs.slack.dev/reference/block-kit/blocks/plan-block/) to present a list of tasks all together. It can be in one of four different states: `pending`, `in_progress`, `completed`, and `error`.

![Plan display mode](https://a.slack-edge.com/bf0a72a/img/api/partner_docs/thinking_steps/combined_plan.png)

With every message, provide an opportunity for feedback on the response with:

- [interactive elements](https://docs.slack.dev/messaging/creating-interactive-messages),
- [context actions block](https://docs.slack.dev/reference/block-kit/blocks/context-actions-block/),
- [icon button block](https://docs.slack.dev/reference/block-kit/block-elements/icon-button-element/), and
- [feedback button block](https://docs.slack.dev/reference/block-kit/block-elements/feedback-buttons-element/).

You can also subscribe to [`reaction_added`](https://docs.slack.dev/reference/events/reaction_added) events to collect feedback based on reactions.

![Image of context block](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAG0AAAAmCAMAAAD5ueXLAAAAD1BMVEX///9qaWq4uLjc3NySkZJWgBFNAAAACXBIWXMAAAsTAAALEwEAmpwYAAABQUlEQVR42u1VwZLFIAhLwP//4kL20J231Yrdmb6jOdWiICEgsLGxsfEeHH+0AMBcn3JGvcUFpEMAYzQNa0sHKdc6nKHc0hIAJQCiqY7m/JhJahErRcztlpTLSZKEWlbRTLi6MJVUUa7iNoRgYkiSZNlR0C7f6mrRUFTGEwkh/yOLZCeM1jN0QfhcAEhYQJQjnsN5uWrD6YpJC7hAiROlEIBdaM4qt8OkP7psdnMmGg4AtOS9fSYH5FEweZKVMD5xFMiVji4bOzHZaJZqCh+qMhrv1vZmECnn3XgyzCA8usLZ90ehFEg/RwR6Hb3JzcRZRU95ewiMHNN+g1JMhxPhvJM8dr8BFvWEfEjtU9LwyWU4UMOPNGchWxD6lWKs36UGSVq+b+z6eFoqgICFA1o2pefEB6sGitfq9C/42NjY2JjgB7Pii0Q9wJPWAAAAAElFTkSuQmCC)

A simple thumbs up/down reaction emoji will work, but consider opening a modal to collect more information when the response was graded poorly so that you can learn more about what the issue was.

### App threads

Slack will automatically group your app conversations into threads, shown in a timeline above the composer in the Messages tab. You can set the title of these threads using the [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) API method, or use the Bolt framework utility to handle the details. The title shows in the reply bar of an individual message. When viewing the thread, the title shows in the header.

When a user renames a session themselves, your app receives an [`agent_session_title_changed`](https://docs.slack.dev/reference/events/agent_session_title_changed) event so you can keep titles in sync. The legacy [`assistant.threads.setTitle`](https://docs.slack.dev/reference/methods/assistant.threads.setTitle) method still works through the [compatibility bridge](https://docs.slack.dev/ai/migrating-to-agent-messaging#compatibility-bridge), and the Bolt `setTitle` utilities below wrap it.

### Messaging guidelines

#### Block Kit and interactivity

Provide interactive [Block Kit](https://docs.slack.dev/block-kit) elements, such as drop-down menus and buttons, to allow your user to interact with the app. Block Kit is not required, however; you can forgo interactivity and message the user via plain text and Slack markdown.

When updating longer messages sent to a user, only call the [`chat.update`](https://docs.slack.dev/reference/methods/chat.update) method once every 3 seconds with new content, otherwise your calls may hit the rate limit.

You can also set a [section block element's](https://docs.slack.dev/reference/block-kit/blocks/section-block) `expand` property to `true` to allow your app to post long messages without the user needing to click 'see more' to see the full text of the message.

When not using Block Kit, use [Slack `mrkdwn`](https://docs.slack.dev/messaging/formatting-message-text#basic-formatting) for sending rich text. The formatting system of Slack is different from the common markup language used elsewhere on the web. It is typical for LLMs to default to use this common markdown syntax unless prompted otherwise, which will not render correctly when posted in Slack. Using the [Markdown Block](https://docs.slack.dev/reference/block-kit/blocks/markdown-block) with standard markdown entered as the `text` input will ensure Slack translates the formatting correctly.

Add a disclaimer at the footer of app messages indicating that the response was generated by an LLM (large language model) and provide any disclaimers that may be appropriate or applicable to communicate with the user. An example of this might sound something like:

> This content was generated by an LLM. Check generated content for vulnerabilities, and do not use to generate code that is visible outside of Slack. Review carefully before acting on the response; it may contain bias or hallucinations.

#### Media support

Where applicable, build apps that can handle a wide array of media types to provide the best user experience. Make it clear which media the app supports and gracefully fail when necessary. Refer to [working with files](https://docs.slack.dev/messaging/working-with-files) for this guidance.

For example, if a user sends an image file and your app does not support receiving images, reply with a message like this:

> It looks like the image didn't come through! 📸 Feel free to describe what you need help with, and I'll do my best to assist you. 😊

#### Sending notifications

When your app has the Agents feature toggled on, every DM with the user is a thread.

When sending a notification to a user outside of an existing thread:

1. Use the [`chat.postMessage`](https://docs.slack.dev/reference/methods/chat.postMessage) method as normal, but look for the `ts` parameter in the response.
2. Call the [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) method, sending the new `ts` parameter as the `thread_ts` parameter to set the title of the thread. This allows a user to see the new notification with a titled thread when they view the app's DM.

Slack has the **Activity** side rail tab to show new activity in a workspace. This area is optimized for users to quickly see and respond to notifications from your app.

#### References, citations, and annotations

Sources and attribution should be used and displayed consistently. There should be a concise way to reference internal messages and files from external sources. We recommend including these in each response that cited a source or used knowledge from a message or document to generate the output. Doing so builds trust with your app's users.

Graceful failure means the agent treats its own partial progress as something worth preserving.

Agents make mistakes, they make things up, they omit important details and get stuck in endless loops of thinking. These situations can happen fairly often, and designing for them is a critical part of building a good agent experience.

When an agent experiences an error, it should:

- Save what it's accomplished
- Explain where it got stuck and why
- Give the user a clear set of options, including:
	- provide the missing information
	- skip the blocked step
	- or take over manually from a known state.

As a last resort, clear the status so the app is not stuck 'thinking' indefinitely. It's important for the user to know when there are errors beyond their control and they need to try again or report a bug.

Something as simple as sending a message like this can go a long way for a user understanding why something isn't working:

> Ope sorry! TeamworkDreamwork App isn't enabled for you.

## Full example

Here is a full code example of the response loop.

## Additional guidelines

### Data retention

Do not store any Slack data you obtain. Instead, store metadata and pull in data in real time if needed.

### Members only

Workspace guests are not permitted to access apps with the Agents feature enabled.

### Beware of prompt injection

Integrating with AI carries an inherent risk of prompt injection. Read more about the risk of data exfiltration and how to prevent it in the [security](https://docs.slack.dev/concepts/security#prompt-injection) docs.

✨ Integrate the [Slack MCP server](https://docs.slack.dev/ai/slack-mcp-server) to access Slack data and perform user-authorized actions in your app.

✨ Discover the various [interaction surfaces](https://docs.slack.dev/interactivity) your app can employ to interact with its users.

✨ Explore all of the [surfaces](https://docs.slack.dev/surfaces) on which your app can exist.