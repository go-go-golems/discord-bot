Source: https://docs.slack.dev/reference/methods/pins.add/
Retrieved: 2026-09-15

## Facts

## Arguments

### Required arguments

**`channel`** `string` Required

Channel to pin the messsage to. You must also include a `timestamp` when pinning messages.

### Optional arguments

**`timestamp`** `string` Optional

Timestamp of the message to pin. You must also include the `channel`.

## Usage info

This method pins a message to a particular conversation or channel.

Both the `channel` and `timestamp` arguments are functionally required. The provided `channel` should be the ID of a public or private channel you want to pin to and the `timestamp` should be the `ts` value of a message within that conversation.

In the past, files and file comments could be pinned to a channel as well.

---

## Response

#### Typical success response

```json
{
  "ok": true
}
```

#### Typical error response

```json
{
  "error": "channel_not_found",
  "ok": false
}
```

After processing, a [`pin_added`](/reference/events/pin_added) event is broadcast via the [Events](/apis/events-api/) and [RTM](/legacy/legacy-rtm-api) APIs.

Some objects cannot be pinned: channel join messages, files, and file comments. A `not_pinnable` error is thrown when we " *no can do*."

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

`already_pinned`

The specified item is already pinned to the channel.

`bad_timestamp`

Value passed for `timestamp` was invalid.

`channel_not_found`

The `channel` argument was not specified or was invalid

`deprecated_endpoint`

The endpoint has been deprecated.

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`external_channel_migrating`

Channel is undergoing an active migration.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

`file_not_found`

File not found.

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

Message specified by `channel` and `timestamp` does not exist.

`method_deprecated`

The method has been deprecated.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`no_item_specified`

One of `file`, `file_comment`, or `timestamp` was not specified.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`not_in_channel`

Item is not in channel.

`not_pinnable`

This message type is not pinnable.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`restricted_action`

The user does not have permission to add pins to the channel.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`too_many_pins`

Too many pins in channel.

`two_factor_setup_required`

Two factor setup is required.