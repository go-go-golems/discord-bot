## Facts

## Arguments

## Usage info

This method checks authentication and tells "you" who you are, even if you might be a bot.

You can also use this method to test whether Slack API authentication is functional.

---

## Response

#### Standard success response when used with a user token

```json
{
  "ok": true,
  "url": "https://subarachnoid.slack.com/",
  "team": "Subarachnoid Workspace",
  "user": "grace",
  "team_id": "T12345678",
  "user_id": "W12345678"
}
```

#### Standard failure response when used with an invalid token

```json
{
  "ok": false,
  "error": "invalid_auth"
}
```

#### Success response when using a bot user token

```json
{
  "ok": true,
  "url": "https://subarachnoid.slack.com/",
  "team": "Subarachnoid Workspace",
  "user": "bot",
  "team_id": "T0G9PQBBK",
  "user_id": "W23456789",
  "bot_id": "BZYBOTHED"
}
```

#### Error response when omitting a token

```json
{
  "ok": false,
  "error": "not_authed"
}
```

When working against a team within an [Enterprise organization](/enterprise), you'll also find their `enterprise_id` here.

## Rate limiting

This method allows hundreds of requests per minute. Use it as often as is reasonably required. Please consult [rate limits](/apis/web-api/rate-limits) for more information.

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

`deprecated_endpoint`

The endpoint has been deprecated.

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

`internal_error`

An internal error has been found.

`internal_error`

The server could not complete your operation(s) without encountering an error, likely due to a transient issue on our end. It's possible some aspect of the operation succeeded before the error was raised.

`invalid_arg_name`

The method was passed an argument whose name falls outside the bounds of accepted or expected values. This includes very long names and names with non-alphanumeric characters other than `_`. If you get this error, it is typically an indication that you have made a *very* malformed API call.

`invalid_arguments`

The method was called with invalid arguments.

`invalid_array_arg`

The method was passed an array as an argument. Please only input valid strings.

`invalid_auth`

Method was called with invalid credentials

`invalid_auth`

Some aspect of authentication cannot be validated. Either the provided token is invalid or the request originates from an IP address disallowed from making the request.

`invalid_form_data`

The method was called via a `POST` request with `Content-Type` `application/x-www-form-urlencoded` or `multipart/form-data`, but the form data was either missing or syntactically invalid.

`invalid_post_type`

The method was called via a `POST` request, but the specified `Content-Type` was invalid. Valid types are: `application/json` `application/x-www-form-urlencoded` `multipart/form-data` `text/plain`.

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