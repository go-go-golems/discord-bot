Source: https://docs.slack.dev/reference/methods/chat.delete/
Retrieved: 2026-09-15

## Facts

## Arguments

### Required arguments

**`channel`** `string` Required

Channel containing the message to be deleted.

**`ts`** `string` Required

Timestamp of the message to be deleted.

*Example:* `"1405894322.002768"`

### Optional arguments

**`as_user`** `boolean` Optional

(Legacy) Pass true to delete the message as the authed user with `chat:write:user` scope. Bot users in this context are considered authed users. See [legacy `as_user` parameter](#legacy_as_user) below.

*Example:* `true`

## Usage info

This method deletes a message from a conversation.

When used with a user token, this method may only [delete messages that user themselves can delete in Slack](https://slack.com/intl/en-ie/help/articles/202395258-Edit-or-delete-messages).

When used with a bot token, this method may delete only messages posted by that bot.

This method has no accommodations for impersonation. If a message is sent impersonating another user, you will not be able to call `chat.delete` to delete that same message.

---

## Response

#### Typical success response

```json
{
  "ok": true,
  "channel": "C123ABC456",
  "ts": "1401383885.000061"
}
```

#### Typical error response

```json
{
  "error": "message_not_found",
  "ok": false
}
```

The response includes the `channel` and `timestamp` properties of the deleted message.

---

## Legacy concerns

This feature works differently for classic apps.

### Legacy authorship

Classic apps using the umbrella `bot` scope can't request additional scopes to adjust message authorship.

#### Legacy as\_user parameter

For classic apps, the best way to control the authorship of a message was to be explicit with the legacy `as_user` parameter. If you didn't use the `as_user` parameter, `chat.delete` would guess the most appropriate `as_user` interpretation based on the kind of token you were using. If `as_user` was not provided at all, the value was inferred based on the scopes granted to the caller: If the caller *could* delete with `as_user` passed as `false`, then that was how the method behaved; otherwise, the method behaved as if `as_user` were passed as `true`. When the `as_user` parameter was set to `false`, messages were deleted as " [`bot_messages`](/reference/events/message/bot_message) ", with message authorship attributed to the user name and icons associated with the classic app.

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

`cant_delete_message`

Authenticated user does not have permission to delete this message.

`channel_not_found`

Value passed for `channel` was invalid.

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

`message_not_found`

No message exists with the requested timestamp.

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

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.