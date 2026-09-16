## Facts

## Arguments

### Required arguments

**`refresh_token`** `string` Required

The `xoxe` refresh token that was issued along with the old app configuration token.

*Example:* `xoxe-1-abcdefg`

## Usage info

## Rotating configuration tokens

Each app configuration token will expire 12 hours after it has been generated. In order to continually rotate your config tokens, you are also provided with a **refresh token**.

It's strongly suggested that you refresh your token before it expires, rather than waiting for it to expire and checking for an error from the Slack API.

In order to refresh config tokens, make a call to [`tooling.tokens.rotate`](/reference/methods/tooling.tokens.rotate), using the refresh token in the `refresh_token` argument. In response you'll receive something like this:

```json
{
    "ok": true,
    "token": "xoxe.xoxp-...",
    "refresh_token": "xoxe-...",
    "team_id": "...",
    "user_id": "...",
    "iat": 1633095660,
    "exp": 1633138860
}
```

The `token` field contains your new config access token, which you can then store and use for Manifest API calls. The `refresh_token` field contains a *new* refresh token.

The remainder of the response above contains fields which identify the source workspace and user of each token, as well as timestamps which indicate when the token was issued and when it will expire.

To learn how to use app configuration tokens, read [our guide to the App Manifest APIs](/app-manifests/configuring-apps-with-app-manifests#manifest_apis).

---

## Response

#### Typical success response

```json
{
  "ok": true,
  "token": "xoxe.xoxp-...",
  "refresh_token": "xoxe-...",
  "team_id": "...",
  "user_id": "...",
  "iat": 1633095660,
  "exp": 1633138860
}
```

#### Typical error response if incorrect refresh token used

```json
{
  "ok": false,
  "error": "invalid_refresh_token"
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

`deprecated_endpoint`

The endpoint has been deprecated.

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

`internal_error`

Internal error

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

`invalid_refresh_token`

The given refresh token is invalid.

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

`unknown_error`

Temporary error for dev only restriction