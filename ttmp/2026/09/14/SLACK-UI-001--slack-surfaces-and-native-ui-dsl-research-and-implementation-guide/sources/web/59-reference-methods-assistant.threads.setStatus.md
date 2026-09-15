---
Title: "59 reference methods assistant.threads.setStatus"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/methods/assistant.threads.setStatus/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Facts

## Arguments

### Required arguments

**`channel_id`** `string` Required

Channel ID containing the assistant thread.

**`thread_ts`** `string` Required

Message timestamp of the thread of where to set the status.

**`status`** `string` Required

Status of the specified bot user, e.g., 'is thinking...'. A two minute timeout applies, which will cause the status to be removed if no message has been sent.

### Optional arguments

**`loading_messages`** `array` Optional

The list of messages to rotate through as a loading indicator. Maximum of 10 messages.

**`icon_emoji`** `string` Optional

Emoji to use as the icon for this message. Overrides `icon_url`.

*Example:* `:chart_with_upwards_trend:`

**`icon_url`** `string` Optional

Image URL to use as the icon for this message.

*Example:* `http://lorempixel.com/48/48`

**`username`** `string` Optional

The bot's username to display.

*Example:* `My Bot`

## Usage info

Use this method to push status updates to users in apps using AI features.

This method helps apps set expectations for potentially slow responses; e.g. "app is thinking...". The status will be automatically cleared when the app sends a reply. Sending an empty string in the `status` field will also clear the status indicator. This can be handy if you want to clear the status indicator without sending a new message.

The status is displayed as `<App Name> <status>`, and Slack automatically inserts the `<App Name>`. So for the example below, if the app's name is `YourAssistantJeeves`, the status would render as `YourAssistantJeeves is working on your request...`.

You can also set a custom loading state via the `loading_messages` parameter. The `loading_messages` parameter is an array of strings that Slack will rotate through showing to serve as a loading indicator.

Example request:

```json
{
   "status": "is working on your request...",
   "channel_id": "D324567865",
   "thread_ts": "1724264405.531769"
}
```

### Rate limits

Rate limiting conditions for this method include a default rate limit that is sufficient for most apps, with custom overrides determined on a per-team basis. The default limit is 600 requests per minute (per app per team).

---

## Response

#### Typical success response

```json
{
  "ok": true
}
```

#### Typical error response for invalid channel

```json
{
  "ok": false,
  "error": "channel_not_found",
  "detail": "Invalid channel_id"
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

Error returned when given an invalid channel\_id

`deprecated_endpoint`

The endpoint has been deprecated.

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

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

`invalid_thread_ts`

Error returned when given an invalid thread\_ts

`method_deprecated`

The method has been deprecated.

`method_not_supported_for_channel_type`

This type of conversation cannot be used with this method. Use `agents.sessions.setStatus` to set status in a session channel.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`reserved_username`

Reserved usernames are not allowed to be used.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.