Source: https://docs.slack.dev/reference/methods/conversations.list/
Retrieved: 2026-09-15

## Facts

## Arguments

### Optional arguments

**`cursor`** `string` Optional

Paginate through collections of data by setting the `cursor` parameter to a `next_cursor` attribute returned by a previous request's `response_metadata`. Default value fetches the first "page" of the collection. See [pagination](/apis/web-api/pagination) for more detail.

*Example:* `dXNlcjpVMDYxTkZUVDI=`

**`exclude_archived`** `boolean` Optional

Set to `true` to exclude archived channels from the list.

*Default:* `false`

*Example:* `true`

**`limit`** `number` Optional

The maximum number of items to return. Fewer than the requested number of items may be returned, even if the end of the list hasn't been reached. Must be an integer under 1000.

*Default:* `100`

*Example:* `20`

**`team_id`** `string` Optional

encoded team id to list channels in, required if token belongs to org-wide app

**`types`** `string` Optional

Mix and match channel types by providing a comma-separated list of any combination of `public_channel`, `private_channel`, `mpim`, `im`

*Default:* `public_channel`

*Example:* `public_channel,private_channel`

## Usage info

This [Conversations API](/apis/web-api/using-the-conversations-api) method returns a list of all [channel-like conversations](/reference/objects/conversation-object) in a workspace. The "channels" returned depend on what the calling token has access to and the directives placed in the `types` parameter.

The `team_id` is only relevant when using an org-level token. This field will be ignored if the API call is sent using a workspace-level token. When paginating, any filters used in the request are applied *after* retrieving a virtual page's `limit`. For example, using `exclude_archived=true` when `limit=20` on a virtual page that would contain 15 archived channels will return you the virtual page with only `5` results. Additional results are available from the next `cursor` value.

---

## Response

#### Typical success response with only public channels

```json
{
  "ok": true,
  "channels": [
    {
      "id": "C012AB3CD",
      "name": "general",
      "is_channel": true,
      "is_group": false,
      "is_im": false,
      "created": 1449252889,
      "creator": "U012A3CDE",
      "is_archived": false,
      "is_general": true,
      "unlinked": 0,
      "name_normalized": "general",
      "is_member": true,
      "is_private": false,
      "is_mpim": false,
      "updated": 1678229664302,
      "topic": {
        "value": "Company-wide announcements and work-based matters",
        "creator": "",
        "last_set": 0
      },
      "purpose": {
        "value": "This channel is for team-wide communication and announcements. All team members are in this channel.",
        "creator": "",
        "last_set": 0
      },
      "previous_names": [],
      "num_members": 4
    },
    {
      "id": "C061EG9T2",
      "name": "random",
      "is_channel": true,
      "is_group": false,
      "is_im": false,
      "created": 1449252889,
      "creator": "U061F7AUR",
      "is_archived": false,
      "is_general": false,
      "unlinked": 0,
      "name_normalized": "random",
      "is_member": true,
      "is_private": false,
      "is_mpim": false,
      "updated": 1678229664302,
      "topic": {
        "value": "Non-work banter and water cooler conversation",
        "creator": "",
        "last_set": 0
      },
      "purpose": {
        "creator": "",
        "last_set": 0
      },
      "previous_names": [],
      "num_members": 4
    }
  ],
  "response_metadata": {
    "next_cursor": "dGVhbTpDMDYxRkE1UEI="
  }
}
```

#### Example response when mixing different conversation types together, like im and mpim

```json
{
  "ok": true,
  "channels": [
    {
      "id": "G0AKFJBEU",
      "name": "mpdm-mr.banks--slactions-jackson--beforebot-1",
      "is_channel": false,
      "is_group": true,
      "is_im": false,
      "created": 1493657761,
      "creator": "U061F7AUR",
      "is_archived": false,
      "is_general": false,
      "unlinked": 0,
      "name_normalized": "mpdm-mr.banks--slactions-jackson--beforebot-1",
      "is_member": true,
      "is_private": true,
      "is_mpim": true,
      "is_open": true,
      "updated": 1678229664302,
      "topic": {
        "value": "Group messaging",
        "creator": "U061F7AUR",
        "last_set": 1493657761
      },
      "purpose": {
        "value": "Group messaging with: @mr.banks @slactions-jackson @beforebot",
        "creator": "U061F7AUR",
        "last_set": 1493657761
      },
      "priority": 0
    },
    {
      "id": "D0C0F7S8Y",
      "created": 1498500348,
      "is_im": true,
      "user": "U0BS9U4SV",
      "is_user_deleted": false,
      "priority": 0
    },
    {
      "id": "D0BSHH4AD",
      "created": 1498511030,
      "is_im": true,
      "user": "U0C0NS9HN",
      "is_user_deleted": false,
      "priority": 0
    }
  ],
  "response_metadata": {
    "next_cursor": "aW1faWQ6RDBCSDk1RExI"
  }
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "invalid_auth"
}
```

Returns a list of limited channel-like [conversation objects](/reference/objects/conversation-object). To get a full [conversation object](/reference/objects/conversation-object), call the [`conversations.info`](/reference/methods/conversations.info) method.

Use [`conversations.members`](/reference/methods/conversations.members) to retrieve and traverse membership.

See [conversation object](/reference/objects/conversation-object) for more detail on returned fields.

Some fields in the response, like `unread_count` and `unread_count_display`, are included for DM conversations only.

### Pagination

This method uses cursor-based pagination to make it easier to incrementally collect information. To begin pagination, specify a `limit` value under `1000`. We recommend no more than `200` results at a time.

Responses will include a top-level `response_metadata` attribute containing a `next_cursor` value. By using this value as a `cursor` parameter in a subsequent request, along with `limit`, you may navigate through the collection page by virtual page.

See [pagination](/apis/web-api/pagination) for more information.

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

The server could not complete your operation(s) without encountering an error, likely due to a transient issue on our end. It's possible some aspect of the operation succeeded before the error was raised.

`invalid_arg_name`

The method was passed an argument whose name falls outside the bounds of accepted or expected values. This includes very long names and names with non-alphanumeric characters other than `_`. If you get this error, it is typically an indication that you have made a *very* malformed API call.

`invalid_arguments`

The method was called with invalid arguments.

`invalid_array_arg`

The method was passed an array as an argument. Please only input valid strings.

`invalid_auth`

Some aspect of authentication cannot be validated. Either the provided token is invalid or the request originates from an IP address disallowed from making the request.

`invalid_cursor`

Value passed for `cursor` was not valid or is no longer valid.

`invalid_form_data`

The method was called via a `POST` request with `Content-Type` `application/x-www-form-urlencoded` or `multipart/form-data`, but the form data was either missing or syntactically invalid.

`invalid_limit`

Value passed for `limit` is not understood.

`invalid_post_type`

The method was called via a `POST` request, but the specified `Content-Type` was invalid. Valid types are: `application/json` `application/x-www-form-urlencoded` `multipart/form-data` `text/plain`.

`invalid_types`

Value passed for `type` could not be used based on the method's capabilities or the permission scopes granted to the used token.

`method_deprecated`

The method has been deprecated.

`method_not_supported_for_channel_type`

This type of conversation cannot be used with this method.

`missing_argument`

A required argument is missing.

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