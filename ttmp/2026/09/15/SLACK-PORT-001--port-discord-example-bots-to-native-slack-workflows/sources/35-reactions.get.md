Source: https://docs.slack.dev/reference/methods/reactions.get/
Retrieved: 2026-09-15

## Facts

## Arguments

### Optional arguments

**`channel`** `string` Optional

Channel where the message to get reactions for was posted.

*Example:* `C0NF841BK`

**`file`** `string` Optional

File to get reactions for.

*Example:* `F1234567890`

**`full`** `boolean` Optional

If true always return the complete reaction list.

**`timestamp`** `string` Optional

Timestamp of the message to get reactions for.

*Example:* `1524523204.000192`

## Usage info

This method returns a list of all reactions for a single item (file, file comment, channel message, group message, or direct message).

---

## Response

#### The response contains the item with reactions.

```json
{
  "ok": true,
  "type": "message",
  "message": {
    "type": "message",
    "text": "Hi there!",
    "user": "W123456",
    "ts": "1648602352.215969",
    "team": "T123456",
    "reactions": [
      {
        "name": "grinning",
        "users": [
          "W222222"
        ],
        "count": 1
      },
      {
        "name": "question",
        "users": [
          "W333333"
        ],
        "count": 1
      }
    ],
    "permalink": "https://xxx.slack.com/archives/C123456/p1648602352215969"
  },
  "channel": "C123ABC456"
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "invalid_auth"
}
```

An item will always have a `type` property, while the other properties depend on the type of item. The possible types are:

- **`message`**: the item will have a `message` property containing a [message object](/messaging) and a `channel` property containing the channel ID for the message
- **`file`**: this item will have a `file` property containing a [file object](/reference/objects/file-object)
- **`file_comment`**: the item will have a `file` property containing the [file object](/reference/objects/file-object) and a `comment` property containing the file comment

The `users` array in the `reactions` property will always contain the authenticated user, but might not always contain all users that have reacted. The value of `count`, however, will always represent the count of all users who made that reaction (i.e. it may be greater than `users.length`).

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

`bad_timestamp`

Value passed for `timestamp` was invalid.

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

`file_comment_not_found`

File comment specified by `file_comment` does not exist.

`file_not_found`

File specified by `file` does not exist.

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

`file`, `file_comment`, or combination of `channel` and `timestamp` was not specified.

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