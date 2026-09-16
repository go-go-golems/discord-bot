Source: https://docs.slack.dev/reference/methods/conversations.info/
Retrieved: 2026-09-15

## Facts

## Arguments

### Optional arguments

**`include_locale`** `boolean` Optional

Set this to `true` to receive the locale for this conversation. Defaults to `false`

**`include_num_members`** `boolean` Optional

Set to `true` to include the member count for the specified conversation. Defaults to `false`

*Default:* `false`

*Example:* `true`

## Usage info

This [Conversations API](/apis/web-api/using-the-conversations-api) method returns information about a workspace [conversation](/reference/objects/conversation-object).

---

## Response

#### Typical success response for a public channel. A response from a private channel and a multi-party IM is very similar to this example. Note that if the properties.tabs parameter is an empty set, it will not be included in the channel object.

```json
{
  "ok": true,
  "channel": {
    "id": "C012AB3CD",
    "name": "general",
    "is_channel": true,
    "is_group": false,
    "is_im": false,
    "is_mpim": false,
    "is_private": false,
    "created": 1654868334,
    "is_archived": false,
    "is_general": true,
    "unlinked": 0,
    "name_normalized": "general",
    "is_frozen": false,
    "context_team_id": "T123ABC456",
    "updated": 1723130875818,
    "parent_conversation": null,
    "creator": "U123ABC456",
      "T123ABC456"
    ],
    "pending_connected_team_ids": [],
    "topic": {
      "value": "For public discussion of generalities",
      "creator": "W012A3BCD",
      "last_set": 1449709364
    },
    "purpose": {
      "value": "This part of the workspace is for fun. Make fun here.",
      "creator": "W012A3BCD",
      "last_set": 1449709364
    },
    "properties": {
      "tabs": [
        {
          "id": "workflows",
          "label": "",
          "type": "workflows"
        },
        {
          "id": "files",
          "label": "",
          "type": "files"
        },
        {
          "id": "bookmarks",
          "label": "",
          "type": "bookmarks"
        }
      ]
    },
    "previous_names": []
  }
}
```

#### Typical success response for a 1:1 direct message

#### When using the method with the include\_num\_members parameter, we return a num\_members field

#### Typical error response when a channel cannot be found

```json
{
  "ok": false,
  "error": "channel_not_found"
}
```

Returns a [conversation object](/reference/objects/conversation-object), which could be a public channel, private channel, direct message, multi-person direct message, depending completely on the `channel` ID and the permissions granted to your token. Some fields in the response, like `unread_count` and `unread_count_display`, are included for DM conversations only.

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

`channel_is_limited_access`

The user has no access to the channel. Only applicable to Salesforce limited access channels.

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

The calling token is not granted the necessary scopes to complete this operation.

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