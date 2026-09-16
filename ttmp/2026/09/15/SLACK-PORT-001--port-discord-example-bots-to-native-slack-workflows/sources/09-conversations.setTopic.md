Source: https://docs.slack.dev/reference/methods/conversations.setTopic/
Retrieved: 2026-09-15

## Facts

## Arguments

### Required arguments

**`channel`** `string` Required

Conversation to set the topic of

**`topic`** `string` Required

The new topic string. Does not support formatting or linkification.

*Example:* `Apply topically for best effects`

## Usage info

This method is used to change the topic of a conversation. The calling user must be a member of the conversation. Not all conversation types support a new topic.

---

## Response

#### A conversation object is returned:

```json
{
  "ok": true,
  "channel": {
    "id": "C12345678",
    "name": "tips-and-tricks",
    "is_channel": true,
    "is_group": false,
    "is_im": false,
    "is_mpim": false,
    "is_private": false,
    "created": 1649195947,
    "is_archived": false,
    "is_general": false,
    "unlinked": 0,
    "name_normalized": "tips-and-tricks",
    "is_frozen": false,
    "parent_conversation": null,
    "creator": "U12345678",
      "T12345678"
    ],
    "pending_connected_team_ids": [],
    "is_member": true,
    "last_read": "1649869848.627809",
    "latest": {
      "type": "message",
      "subtype": "channel_topic",
      "ts": "1649952691.429799",
      "user": "U12345678",
      "text": "set the channel topic: Apply topically for best effects",
      "topic": "Apply topically for best effects"
    },
    "unread_count": 1,
    "unread_count_display": 0,
    "topic": {
      "value": "Apply topically for best effects",
      "creator": "U12345678",
      "last_set": 1649952691
    },
    "purpose": {
      "value": "",
      "creator": "",
      "last_set": 0
    },
    "previous_names": []
  }
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "invalid_arguments",
  "response_metadata": {
    "messages": [
      "[ERROR] missing required field: topic"
    ]
  }
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

`is_archived`

Channel has been archived.

`method_deprecated`

The method has been deprecated.

`method_not_supported_for_channel_type`

This type of conversation cannot be used with this method.

`missing_scope`

The calling token is not granted the necessary scopes to complete this operation.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`not_in_channel`

Authenticated user is not in the channel.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`too_long`

Topic was longer than 250 characters.

`two_factor_setup_required`

Two factor setup is required.

`user_is_restricted`

Setting the topic is a restricted action.