---
Title: "47 ai migrating to agent messaging"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/ai/migrating-to-agent-messaging/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

In the [Agent messaging experience](https://docs.slack.dev/changelog/2026/06/30/agent-messages-tab), agent conversations happen in the classic Messages tab of the app, with threads shown in a timeline above the composer. Apps that already use the Assistant messaging experience can continue to use it for now, but `assistant_view` will eventually be deprecated, and we'll ask existing apps to migrate.

To move an existing app from the Assistant messaging experience to the Agent messaging experience, work through the following steps.

1. **Switch the manifest from `assistant_view` to `agent_view`.**
	This change enables the Agent messaging experience for all users of your app. See the [app manifest reference](https://docs.slack.dev/reference/app-manifest#features) for the field definitions. To configure in the [app settings](https://api.slack.com/apps), navigate to the **Agent** tab in the left nav, then select the button to update the app. When switching to `agent_view`, the nested `assistant_description` switches to `agent_description`.
2. **Detect when a user opens a DM with the [`app_home_opened`](https://docs.slack.dev/reference/events/app_home_opened) event.**
	Subscribe to the `app_home_opened` event to be notified when a user has opened a DM with your app. Specifically, check the value of the `tab` property of the event to verify when its value is `"messages"`.
	JavaScript example:
	```js
	if (event.tab === 'messages')
	```
	Python example:
	```markdown
	if event.get("tab") == "messages"
	```
	This replaces the [`assistant_thread_started`](https://docs.slack.dev/reference/events/assistant_thread_started) event. Without the `assistant_thread_started` events, messages received in the DM channel will not have a root message with `subtype = "assistant_app_thread"`; rather, the root message will appear from the user.
3. **Revisit how you set suggested prompts.**
	[Suggested prompts](https://docs.slack.dev/reference/methods/assistant.threads.setSuggestedPrompts) now live at the top of the Messages tab instead of within a thread. If your app set prompts contextually for each thread, review that logic.
4. **Move from `assistant.threads.*` methods to `agents.sessions.*`.**
	The [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) and [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) methods replace [`assistant.threads.setStatus`](https://docs.slack.dev/reference/methods/assistant.threads.setStatus) and [`assistant.threads.setTitle`](https://docs.slack.dev/reference/methods/assistant.threads.setTitle). Your existing calls keep working through a [compatibility bridge](#compatibility-bridge), so this step is not yet required to switch experiences, but we recommend migrating as part of this move. See [Migrating from `assistant.threads.*` methods](#migration) below.
5. **Implement the stop button.**
	Slack shows a native stop button while a session is in `processing` if your app subscribes to the [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped) event. Subscribe to it so users can stop your agent, and transition the session out of `processing` when you receive the event. See [how to implement stop](#migration-stop) below.

## Migrating from assistant.threads.\* methods

The [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) and [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) methods replace [`assistant.threads.setStatus`](https://docs.slack.dev/reference/methods/assistant.threads.setStatus) and [`assistant.threads.setTitle`](https://docs.slack.dev/reference/methods/assistant.threads.setTitle). The sessions API works with channels and DMs and supports richer lifecycle states. For the full lifecycle and concepts, see [Agent sessions](https://docs.slack.dev/ai/agent-sessions).

### Compatibility bridge

Existing apps using the legacy APIs will continue to work. Slack bridges calls to the sessions system automatically:

- [`assistant.threads.setStatus`](https://docs.slack.dev/reference/methods/assistant.threads.setStatus) with a non-empty `status` string sets the session to `processing`. An empty `status` string sets the session to `active`.
- [`assistant.threads.setTitle`](https://docs.slack.dev/reference/methods/assistant.threads.setTitle) applies the provided title to the thread's agent session.

For apps using the streaming API methods ([`chat.startStream`](https://docs.slack.dev/reference/methods/chat.startStream) and [`chat.stopStream`](https://docs.slack.dev/reference/methods/chat.stopStream)), the bridge ensures your app participates in the sessions UX without code changes by automatically creating sessions. We plan to deprecate the `assistant.threads.setStatus` and `assistant.threads.setTitle` methods in favor of the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method, and we recommend migrating as part of this move.

### Replace setStatus and setTitle

Replace calls to the legacy API methods with the sessions methods:

| Before | After |
| --- | --- |
| `assistant.threads.setStatus` with non-empty `status` | `agents.sessions.setStatus` with `status: "processing"` |
| `assistant.threads.setStatus` with empty `status` | `agents.sessions.setStatus` with `status: "active"` |
| `assistant.threads.setTitle` | `agents.sessions.rename` with `title: "..."` |

Status and title are managed by separate methods. You can also use the `suspended` and `closed` statuses for richer lifecycle management.

```json
// Before
POST /api/assistant.threads.setStatus
{ "channel_id": "C123", "thread_ts": "1234123.232342", "status": "is typing..." }

POST /api/assistant.threads.setTitle
{ "channel_id": "C123", "thread_ts": "1234123.232342", "title": "Deep sea diving research" }

// After
POST /api/agents.sessions.setStatus
{
  "channel_id": "C123",
  "thread_ts": "1234123.232342",
  "status": "processing",
  "title": "Deep sea diving research"
}

POST /api/agents.sessions.rename
{
  "channel_id": "C123",
  "thread_ts": "1234123.232342",
  "title": "Deep sea diving research"
}
```

### Implement stop

The native stop button only appears if your app subscribes to the [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped) event; while your app is not subscribed, the user sees a non-interactive loading indicator instead. Subscribe to the event and transition the session out of `processing` when you receive it. See [Stopping a session](https://docs.slack.dev/ai/agent-sessions#stopping) for details.