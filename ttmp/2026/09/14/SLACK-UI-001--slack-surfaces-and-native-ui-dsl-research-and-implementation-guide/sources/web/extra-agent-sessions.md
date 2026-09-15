---
Title: "extra agent sessions"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/ai/agent-sessions/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

The `agent_sessions.*` API methods help users manage their conversations with your agent. Sessions surface a status ("processing", "active", "suspended", or "closed"), a title for recall, and a stop button while the agent is working.

Agent sessions are thread-based and scoped to a single thread in a conversation, like a channel or DM.

Only apps declared as agents in the [app settings](https://api.slack.com/apps) can create sessions. Declaring an app as an agent adds the `assistant:write` scope to it. The [`chat:write`](https://docs.slack.dev/reference/scopes/chat.write) scope is required to manage thread-based sessions.

## Available methods

Create sessions when your agent sends messages or signals its intent to reply to a user. The primary API method for managing a session's status is the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method, which creates the session if it doesn't exist and updates its status if it does. Use the [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) method to change a session's title.

| Method | Description |
| --- | --- |
| [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) | Set an agent session's lifecycle status, creating the session if needed. |
| [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) | Rename an agent session. |

Refer to each method's reference page for full details.

## Session lifecycle

1. A user sends a message to your agent in a thread (in a channel or DM).
2. Your app calls the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) API method with `status: "processing"` (and, on creation, an optional `title`).
3. Your app does its work, posting messages with the [`chat.postMessage`](https://docs.slack.dev/reference/methods/chat.postMessage) API method or streaming with the [`chat.startStream`](https://docs.slack.dev/reference/methods/chat.startStream) API method.
4. When the agent finishes and is ready for the next prompt, it calls the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) API method with `status: "active"`.
5. If the agent needs user input to continue, it sets `status: "suspended"`.
6. When the conversation is complete, the agent sets `status: "closed"`.

Users see their subscribed sessions with pinned ones first, then in chronological order, and can pin, rename, or archive them.

### Processing status

When a session is in the `processing` status, Slack shows a "Working..." loading UX to the user, along with a stop button if your app subscribes to the [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped) event. See [Stopping a session](#stopping). The [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method does not accept a custom loading message.

#### Session timeouts

Sessions in the `processing` status time out after one hour and automatically transition to `active`. Setting the status to `processing` again restarts the one-hour timer. Apps can periodically re-send `processing` to reset the timer and keep the session alive during long-running operations.

#### Customizing your agent's appearance

You can override the agent's icon and display name with the `icon_emoji`, `icon_url`, and `username` parameter. These require the [`chat:write.customize`](https://docs.slack.dev/reference/scopes/chat.write.customize) scope.

Overrides persist across status transitions. A call that omits all three parameters (or passes `null` for all three) leaves the existing overrides in place. Each call that sets at least one of the three replaces the whole set, so any of the three you leave out of that call is cleared.

These overrides apply only to the loading UX and session surface, not the messages. To keep a consistent name and icon on the messages your agent sends, set them via the corresponding method, either the [`chat.postMessage`](https://docs.slack.dev/reference/methods/chat.postMessage) method or the [`chat.startStream`](https://docs.slack.dev/reference/methods/chat.startStream) method.

```json
POST /api/agents.sessions.setStatus
{
  "channel_id": "C123ABC",
  "thread_ts": "1717171717.123456",
  "status": "processing",
  "title": "Scuba diving research",
  "icon_emoji": ":robot_face:",
  "username": "Custom Agent Name"
}
```

## Messaging interactions

Agent sessions integrate with the Slack streaming API methods for a streamlined experience:

- The [`chat.startStream`](https://docs.slack.dev/reference/methods/chat.startStream) method creates a session if one doesn't exist, and sets the session status to `processing`. It requires `channel` and `thread_ts` (the thread to stream into). When streaming into a channel rather than a DM with your app, `recipient_user_id` and `recipient_team_id` are also required. When a stream creates a new session, the `initiator_user_id` is set to the author of the thread's root message.
	```json
	POST /api/chat.startStream
	{
	  "channel": "C123ABC",
	  "thread_ts": "1717171717.123456",
	  "recipient_user_id": "U123ABC456",
	  "recipient_team_id": "T123ABC456",
	  "markdown_text": "I'll research scuba diving and create a canvas for you..."
	}
	```
- The [`chat.stopStream`](https://docs.slack.dev/reference/methods/chat.stopStream) method can set the session status via the `session_status` parameter (defaults to `active`). If a stream times out, the session status is set to `active`.
	```json
	POST /api/chat.stopStream
	{
	  "channel": "C123ABC",
	  "ts": "1717171717.654321",
	  "markdown_text": "What is the intended audience? Experienced diver or newbie?",
	  "session_status": "suspended"
	}
	```

Agent sessions do not require using the streaming API methods. You can instead use the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method with the [`chat.postMessage`](https://docs.slack.dev/reference/methods/chat.postMessage) method and the [`chat.update`](https://docs.slack.dev/reference/methods/chat.update) method. However, when not using the streaming API methods, your app must call the [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method to manage the session status, and the [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) method to manage the title, as there is no implicit state management.

## Stopping a session

If your app subscribes to the [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped) event, Slack displays a native stop button to the user while the session is in the `processing` status. If you don't subscribe to this event, Slack has no way to deliver the stop request, and the user sees a non-interactive loading indicator instead.

Subscribe to the event so users can stop your agent. The [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) method returns a `missing_agent_session_stopped_event_subscription` warning while your app is not subscribed.

When the user clicks the stop button, your app will receive an [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped) event. Your app should:

1. Stop any in-progress work for the given channel and thread.
2. Clean up resources and confirm to the user that work has stopped.
3. Set the session status away from `processing` using either [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) or [`chat.stopStream`](https://docs.slack.dev/reference/methods/chat.stopStream).

The session status will not update automatically when the user clicks stop. Your app is responsible for transitioning the status.

```json
// agent_session_stopped event payload
{
  "channel": "C123ABC456",
  "event_ts": "1783536983.783769",
  "streaming_message_ts": ["1782234987.693923"],
  "thread_ts": "1782234671.392669",
  "type": "agent_session_stopped",
  "user": "U123ABC456"
}
```

The `streaming_message_ts` array lists the timestamps of your app's in-progress streaming messages that Slack stopped in response to the click. It is an empty array when no stream was active, so the event still tells you to stop work that isn't a stream. As with all Events API deliveries, `team_id` is on the enclosing `event_callback` envelope rather than in the event itself.

## Renaming a session

Users may change the title of a session at any time, even if the agent previously set it. When this happens, your app receives an [`agent_session_title_changed`](https://docs.slack.dev/reference/events/agent_session_title_changed) event. If your agent keeps sessions in sync between Slack and your own UI, reflect the user's title change accordingly.

```json
// agent_session_title_changed event payload
{
  "channel": "C123ABC456",
  "event_ts": "1783536983.783769",
  "previous_title": "Scuba diving research",
  "team_id": "T123ABC456",
  "thread_ts": "1782234671.392669",
  "title": "Bora Bora trip prep",
  "type": "agent_session_title_changed",
  "user": "U123ABC456"
}
```

The `previous_title` field is omitted when the session had no title before the change, and `enterprise_id` is included for org-level installs.

## Session visibility

Users subscribe to agent sessions by:

- **Initiating them**, if they are listed as the `initiator_user_id` of a new session.
- **Participating in a thread** that contains an agent session.

The sessions UX shows subscribed sessions with pinned ones first, then in chronological order. Users can pin, rename, or archive sessions to manage their list.

All users in a channel can:

- See the status and title of an agent session in that channel.
- Change the title of an agent session in that channel.

## Migrating from assistant.threads.\*

The [`agents.sessions.setStatus`](https://docs.slack.dev/reference/methods/agents.sessions.setStatus) and [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) methods replace the [`assistant.threads.setStatus`](https://docs.slack.dev/reference/methods/assistant.threads.setStatus) and [`assistant.threads.setTitle`](https://docs.slack.dev/reference/methods/assistant.threads.setTitle) methods. Existing apps will keep working for now through a compatibility bridge, but migrating is recommended.

For step-by-step migration guidance, including the compatibility bridge details, the method replacement table, and how to implement the stop button, see [Migrating to the Agent messaging experience](https://docs.slack.dev/ai/migrating-to-agent-messaging#migration).