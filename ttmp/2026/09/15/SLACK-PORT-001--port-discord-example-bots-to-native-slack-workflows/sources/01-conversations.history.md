Source: https://docs.slack.dev/reference/methods/conversations.history/
Retrieved: 2026-09-15

## Facts

## Arguments

### Optional arguments

**`cursor`** `string` Optional

Paginate through collections of data by setting the `cursor` parameter to a `next_cursor` attribute returned by a previous request's `response_metadata`. Default value fetches the first "page" of the collection. See [pagination](/apis/web-api/pagination) for more detail.

*Example:* `dXNlcjpVMDYxTkZUVDI=`

**`include_all_metadata`** `boolean` Optional

Return all metadata associated with this message.

0

*Example:* `true`

**`inclusive`** `boolean` Optional

Include messages with `oldest` or `latest` timestamps in results. Ignored unless either timestamp is specified.

0

*Example:* `true`

**`latest`** `string` Optional

Only messages before this Unix timestamp will be included in results. Default is the current time.

**`limit`** `number` Optional

The maximum number of items to return. Fewer than the requested number of items may be returned, even if the end of the conversation history hasn't been reached. Maximum of 999.

*Default:* `100`

*Example:* `20`

**`oldest`** `string` Optional

Only messages after this Unix timestamp will be included in results.

*Default:* `0`

## Usage info

This method returns a portion of [message events](/reference/events/message) from the specified conversation. Call the method with no `oldest` or `latest` arguments to read the entire history for a conversation.

Each type of token can be used to access a set of conversations given the proper scopes.

| Token type | Required scopes | Accessible conversations |
| --- | --- | --- |
| [App-level access token](/authentication/tokens#app-level) | Relevant [`*:history`](/reference/scopes?query=history) scope | Any conversation the relevant app is a member of. |
| [Bot token](/authentication/tokens#bot) | Relevant [`*:history`](/reference/scopes?query=history) scope | Any conversation the relevant bot is a member of. |
| [User token](/authentication/tokens#user) | Relevant [`*:history`](/reference/scopes?query=history) scope | Any private conversation the user is a member of, and all public conversations. |
| [Legacy bot user tokens](/authentication/tokens#legacy_types) |  | Direct message and multi-party direct message conversations. |

---

## Response

#### Typical success response containing a channel's messages

```json
{
  "ok": true,
  "messages": [
    {
      "type": "message",
      "user": "U123ABC456",
      "text": "I find you punny and would like to smell your nose letter",
      "ts": "1512085950.000216"
    },
    {
      "type": "message",
      "user": "U222BBB222",
      "text": "What, you want to smell my shoes better?",
      "ts": "1512104434.000490"
    }
  ],
  "has_more": true,
  "pin_count": 0,
  "response_metadata": {
    "next_cursor": "bmV4dF90czoxNTEyMDg1ODYxMDAwNTQz"
  }
}
```

#### Typical success response included formatted messages from bots and incoming webhooks

```json
{
  "ok": true,
  "messages": [
    {
      "type": "message",
      "user": "U123ABC456",
      "text": "I find you punny and would like to smell your nose letter",
      "ts": "1512085950.000216"
    },
    {
      "type": "message",
      "user": "U222BBB222",
      "text": "Isn't this whether dreadful? <https://badpuns.example.com/puns/123>",
      "attachments": [
        {
          "service_name": "Leg end nary a laugh, Ink.",
          "text": "This is likely a pun about the weather.",
          "fallback": "We're withholding a pun from you",
          "thumb_url": "https://badpuns.example.com/puns/123.png",
          "thumb_width": 1920,
          "thumb_height": 700,
          "id": 1
        }
      ],
      "ts": "1512085950.218404"
    }
  ],
  "has_more": true,
  "pin_count": 0,
  "response_metadata": {
    "next_cursor": "bmV4dF90czoxNTEyMTU0NDA5MDAwMjU2"
  }
}
```

#### Typical success response with latest timestamp and inclusive parameters specified

```json
{
  "ok": true,
  "latest": "1512085950.000216",
  "messages": [
    {
      "type": "message",
      "user": "U123ABC456",
      "text": "I find you punny and would like to smell your nose letter",
      "ts": "1512085950.000216"
    }
  ],
  "has_more": true,
  "pin_count": 0,
  "response_metadata": {
    "next_cursor": "bmV4dF90czoxNTEyMzU2NTI2MDAwMTMw"
  }
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "channel_not_found"
}
```

### Pagination

This method uses cursor-based pagination to make it easier to incrementally collect information. To begin pagination, specify a `limit` value under `1000`. We recommend no more than `200` results at a time.

Responses will include a top-level `response_metadata` attribute containing a `next_cursor` value. By using this value as a `cursor` parameter in a subsequent request, along with `limit`, you may navigate through the collection page by virtual page.

See [pagination](/apis/web-api/pagination) for more information.

### Pagination by time

This form of pagination can be used in conjunction with cursors.

The `messages` array contains up to 100 messages between the `oldest` and `latest` timestamps. The most recent messages in the time range are returned first.

If there were more than 100 messages between `oldest` and `latest`, then `has_more` will be `true` in the response. In an additional call, set the `ts` value of the final message as `latest` to get the next page of messages.

If a message has the same timestamp as `oldest` or `latest` it will not be included in the list. This functionality allows you to use the timestamps of specific messages as boundaries for the results. You can, however, have both timestamps be included in the time range by setting `inclusive` to `true`. The `inclusive` parameter is ignored when `oldest` or `latest` is not specified.

### Retrieving a single message

`conversations.history` can also be used to find a single message from the archive.

You'll need a message's `ts` value, uniquely identifying it within a conversation. You'll also need that conversation's ID.

If you know the `ts` of a specific message:

1. Set `oldest` to the `ts`
2. Set `inclusive` to `true`
3. Set `limit` to 1

If you know the `ts` of the message that is before or after of the specific message you're looking for; set `inclusive` to `false` and use the `oldest` or `latest` value respectively.

Provide another message's `ts` value *as* the `latest` parameter. Set `limit` to `1`. If it exists, you'll receive the queried message in return.

Finally, use `inclusive=true` because otherwise we'll never retrieve the message we're actually after, just the ones that come after it.

```markdown
GET /api/conversations.history?channel=C123ABC456&latest=1476909142.000007&inclusive=true&limit=1
Authorization: Bearer TOKEN_WITH_CHANNELS_HISTORY_SCOPE
```

To retrieve a message from a thread, check out [`conversations.replies`](/messaging/retrieving-messages#pulling_threads).

You can easily generate a permalink URL for any specific message using [`chat.getPermalink`](/reference/methods/chat.getPermalink).

---

## Retrieving message history from a direct message

`conversations.history` can also be used to export messages from a direct message.

In addition to a [bot token](/authentication/tokens#legacy_types), you'll need to input the direct message ID (similar to a conversation ID and beginning with `D`) in order to retrieve the message history. You can find the direct message ID by using the [`conversations.list`](/reference/methods/conversations.list) method.

## Message types

Messages of type `"message"` are user-entered text messages sent to the channel, while other types are events that happened within the channel. All messages have both a `type` and a sortable `ts`, but the other fields depend on the `type`. For a list of all possible events, see the [channel messages](/reference/events/message) documentation.

Messages that have been reacted to by team members will have a [reactions](/reference/events/message#stars__pins__and_reactions) array delightfully included. If you need a full list of reactions for a message, use the [`reactions.get`](/reference/methods/reactions.get) method.

If a message has been starred by the calling user, the `is_starred` property will be present and true. This property is only added for starred items, so is not present in the majority of messages.

The `is_limited` boolean property is only included for free teams that have reached the free message limit. If true, there are messages before the current result set, but they are beyond the message limit.

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

The user has no access to the channel. This is only applicable to private Salesforce record channels.

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

`invalid_cursor`

Value passed for `cursor` was not valid or is no longer valid.

`invalid_form_data`

The method was called via a `POST` request with `Content-Type` `application/x-www-form-urlencoded` or `multipart/form-data`, but the form data was either missing or syntactically invalid.

`invalid_metadata_filter_keys`

Value passed for `metadata_keys_to_include` was invalid. Must be valid json array of strings.

`invalid_post_type`

The method was called via a `POST` request, but the specified `Content-Type` was invalid. Valid types are: `application/json` `application/x-www-form-urlencoded` `multipart/form-data` `text/plain`.

`invalid_ts_latest`

Value passed for `latest` was invalid

`invalid_ts_oldest`

Value passed for `oldest` was invalid

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

`not_in_channel`

The token used does not have access to the proper channel. Only user tokens can access public channels they are not in.

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