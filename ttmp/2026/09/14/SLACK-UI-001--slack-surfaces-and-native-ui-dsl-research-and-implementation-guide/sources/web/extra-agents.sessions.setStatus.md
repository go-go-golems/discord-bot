---
Title: "extra agents.sessions.setStatus"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/methods/agents.sessions.setStatus/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Facts

## Arguments

### Required arguments

**`token`** `string` Required

Authentication token bearing the `chat:write` scope. Must be a granular bot token.

*Example:* `xxxx-xxxxxxxxx-xxxx`

**`status`** `string` Required

The lifecycle status to set.

*Acceptable values:* `active` `processing` `suspended` `closed`

### Optional arguments

**`channel_id`** `string` Optional

ID of the channel containing the agent session. Required for public channels.

**`thread_ts`** `string` Optional

Timestamp of the thread root message the session is scoped to. Required for thread-based sessions in regular channels and DMs. Must be omitted for session channels.

*Example:* `1234567890.123456`

**`title`** `string` Optional

Title for the agent session (max 200 characters). Only used when creating a new session; ignored if the session already exists. To rename an existing session, use `agents.sessions.rename`.

**`initiator_user_id`** `string` Optional

The user who initiated the session. Only used when creating a new session; ignored if the session already exists. Must be a member of the channel.

**`icon_emoji`** `string` Optional

Emoji to use as the agent's icon. Takes priority over `icon_url`. Remains in effect until you clear it (pass `null`) or set a new value. Requires the `chat:write.customize` scope.

**`icon_url`** `string` Optional

URL to an image to use as the agent's icon. Remains in effect until you clear it (pass `null`) or set a new value. Requires the `chat:write.customize` scope.

**`username`** `string` Optional

Display name override for the agent (max 200 characters). Remains in effect until you clear it (pass `null`) or set a new value. Requires the `chat:write.customize` scope.

## Usage info

This method sets the lifecycle status of an agent session, creating the session if it does not already exist. This is the primary API method for managing an agent session's status.

To rename a session, use the [`agents.sessions.rename`](https://docs.slack.dev/reference/methods/agents.sessions.rename) API method.

Your app must be a member of the channel specified in `channel_id`.

Setting status on a session thread (regular channel or DM):

```json
{
  "channel_id": "C123ABC",
  "thread_ts": "1717171717.123456",
  "status": "processing",
  "title": "Scuba diving research"
}
```

The `status` field in the response is the session-level status, derived from the statuses of all agents in the session (`suspended` > `processing` > `active` > `closed`). The `agent_status` field is the calling agent's own status, the value your call just wrote.

#### Status values

| Status | Meaning |
| --- | --- |
| `suspended` | The agent cannot make progress until the user intervenes, for example when the agent needs user clarification or a tool approval. |
| `processing` | The agent is working on a user's task. A stop button is shown to the user if your app subscribes to the [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped) event. |
| `active` | The agent is alive and ready for the next prompt or task. |
| `closed` | The agent has closed the session and will no longer respond on it. |

In a session with more than one agent, these can differ: if another agent is still `processing` when you set `active`, the response is `{ "status": "processing", "agent_status": "active" }`. The session is still busy, but your write went through.

### Starting work

Call this when your agent begins working on a task. While in `processing`, Slack shows a loading UX to the user, along with a stop button if your app subscribes to the [`agent_session_stopped`](https://docs.slack.dev/reference/events/agent_session_stopped) event.

```json
{
  "channel_id": "C123ABC",
  "thread_ts": "1717171717.123456",
  "status": "processing"
}
```

### Creating a session with a title and initiator

The `title` and `initiator_user_id` arguments are applied only when the session is first created:

```json
{
  "channel_id": "C123ABC",
  "thread_ts": "1717171717.123456",
  "status": "processing",
  "title": "Scuba diving research",
  "initiator_user_id": "U123ABC456"
}
```

### Customizing agent identity

You can override the agent's icon and display name with the `icon_emoji`, `icon_url`, and `username` parameter. These require the [`chat:write.customize`](https://docs.slack.dev/reference/scopes/chat.write.customize) scope.

Overrides persist across status transitions. A call that omits all three parameters (or passes `null` for all three) leaves the existing overrides in place. Each call that sets at least one of the three replaces the whole set, so any of the three you leave out of that call is cleared.

```json
{
  "channel_id": "C123ABC",
  "thread_ts": "1717171717.123456",
  "status": "processing",
  "icon_emoji": ":robot_face:",
  "username": "Custom Agent Name"
}
```

### Suspending (awaiting user input)

Call this when your agent needs clarification or approval before continuing:

```json
{
  "channel_id": "C123ABC",
  "thread_ts": "1717171717.123456",
  "status": "suspended"
}
```

### Closing a session

Call this when the conversation is complete:

```json
{
  "channel_id": "C123ABC",
  "thread_ts": "1717171717.123456",
  "status": "closed"
}
```

---

## Response

#### Typical success response

```json
{
  "ok": true,
  "status": "processing",
  "agent_status": "processing",
  "title": "Scuba diving research"
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "thread_ts_required"
}
```

## Errors

This table lists the expected errors that this method could return. However, other errors can be returned in the case where the service is down or other unexpected factors affect processing. Callers should always check the value of the `ok` parameter in the response.

Error

Description

`access_denied`

Access to a resource specified in the request is denied.

`accesslimited`

Access to this method is limited on the current network

`account_inactive`

Authentication token is for a deleted user or workspace when using a `bot` token.

`channel_not_found`

The specified channel does not exist or is not accessible.

`deprecated_endpoint`

The endpoint has been deprecated.

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

`feature_disabled`

The agent tasks feature is not enabled for this workspace.

`internal_error`

An internal error occurred while updating the session.

`internal_error`

The server could not complete your operation(s) without encountering an error, likely due to a transient issue on our end. It's possible some aspect of the operation succeeded before the error was raised.

`invalid_arg_name`

The method was passed an argument whose name falls outside the bounds of accepted or expected values. This includes very long names and names with non-alphanumeric characters other than `_`. If you get this error, it is typically an indication that you have made a *very* malformed API call.

`invalid_arguments`

The method was called with invalid arguments.

`invalid_array_arg`

The method was passed an array as an argument. Please only input valid strings.

`invalid_auth`

Some aspect of authentication cannot be validated. Either the provided token is invalid or the request originates from an IP address disallowed from making the request.

`invalid_form_data`

The method was called via a `POST` request with `Content-Type` `application/x-www-form-urlencoded` or `multipart/form-data`, but the form data was either missing or syntactically invalid.

`invalid_post_type`

The method was called via a `POST` request, but the specified `Content-Type` was invalid. Valid types are: `application/json` `application/x-www-form-urlencoded` `multipart/form-data` `text/plain`.

`invalid_status`

The `status` value is not a valid agent session status. Must be one of: `active`, `processing`, `suspended`, `closed`.

`method_deprecated`

The method has been deprecated.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`not_authorized`

The caller is not a member of the specified channel.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`thread_ts_not_allowed`

`thread_ts` must not be provided for session channels. The channel-level session is updated automatically.

`thread_ts_required`

`thread_ts` is required for thread-based sessions in regular channels and DMs.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.

`user_not_found`

The specified `initiator_user_id` does not exist or is not accessible.