Source: https://docs.slack.dev/reference/methods/pins.list/
Retrieved: 2026-09-15

## Facts

## Arguments

## Usage info

This method lists the items pinned to a channel.

---

## Response

#### Typical success response

```json
{
  "items": [
    {
      "channel": "C123ABC456",
      "created": 1508881078,
      "created_by": "U123ABC456",
      "message": {
        "permalink": "https://hitchhikers.slack.com/archives/C2U86NC6H/p1508197641000151",
        "pinned_to": [
          "C2U86NC6H"
        ],
        "text": "What is the meaning of life?",
        "ts": "1508197641.000151",
        "type": "message",
        "user": "U123ABC456"
      },
      "type": "message"
    },
    {
      "channel": "C123ABC456",
      "created": 1508880991,
      "created_by": "U123ABC456",
      "message": {
        "permalink": "https://hitchhikers.slack.com/archives/C2U86NC6H/p1508284197000015",
        "pinned_to": [
          "C123ABC456"
        ],
        "text": "The meaning of life, the universe, and everything is 42.",
        "ts": "1503289197.000015",
        "type": "message",
        "user": "U123ABC456"
      },
      "type": "message"
    }
  ],
  "ok": true
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "invalid_auth"
}
```

The response contains a list of pinned items in a channel. Different item types can be pinned. Every item in the list has a `type` property, and the other properties depend on the type of item. The possible types are:

- **`message`**: the item will have a `message` property containing a [message object](/messaging) and a `channel` property containing the channel ID for the message.
- **`file`**: this item will have a `file` property containing a [file object](/reference/objects/file-object).
- **`file_comment`**: the item will have a `file` property containing the [file object](/reference/objects/file-object) and a `comment` property containing the file comment.

The `created` property on each item is a Unix timestamp representing when the item was pinned. The `created_by` property on each item is a string representing the encoded user id of the user who pinned the item.

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

`restricted_action`

The user does not have permission to view the channel.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.