Source: https://docs.slack.dev/reference/methods/conversations.replies/
Retrieved: 2026-09-15

## Facts

## Arguments

### Required arguments

**`channel`** `string` Required

Conversation ID to fetch thread from.

**`ts`** `string` Required

Unique identifier of either a thread’s parent message or a message in the thread. `ts` must be the timestamp of an existing message with 0 or more replies. If there are no replies then just the single message referenced by `ts` will return - it is just an ordinary, unthreaded message.

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

Only messages before this Unix timestamp will be included in results.

*Default:* `now`

**`limit`** `number` Optional

The maximum number of items to return. Fewer than the requested number of items may be returned, even if the end of the users list hasn't been reached.

*Default:* `1000`

*Example:* `20`

**`oldest`** `string` Optional

Only messages after this Unix timestamp will be included in results.

*Default:* `0`

## Usage info

This [Conversations API](/apis/web-api/using-the-conversations-api) method returns a cursor-paginated thread of messages posted to a conversation.

The `channel` and `ts` arguments are always required. `ts` must be the timestamp of an existing message. If no replies, just the single message referenced by `ts` will return.

The `reply_users` field returned by this method sometimes contains bot IDs rather than user IDs. This is essentially a fallback that can occur depending on how the message was posted. It is recommended to check for this either on the prefix or by implementing a retry mechanism if user lookup fails.

The `thread_not_found` error shown in the example error response can also apply to the [`channel_leave`](/reference/events/message/channel_leave) and [`channel_join`](/reference/events/message/channel_join) message subtypes, as these message subtypes cannot be threaded.

---

## Response

#### Typical success response

#### Typical error response

```json
{
  "ok": false,
  "error": "thread_not_found"
}
```

### Pagination

This method uses cursor-based pagination to make it easier to incrementally collect information. To begin pagination, specify a `limit` value under `1000`. We recommend no more than `200` results at a time.

Responses will include a top-level `response_metadata` attribute containing a `next_cursor` value. By using this value as a `cursor` parameter in a subsequent request, along with `limit`, you may navigate through the collection page by virtual page.

See [pagination](/apis/web-api/pagination) for more information.

### Pagination by time

This form of pagination can be used in conjunction with cursors.

The `messages` array contains up to 1000 messages between the `oldest` and `latest` timestamps. The earliest messages in the time range are returned first.

If there were more than 1000 messages between `oldest` and `latest`, then `has_more` will be `true` in the response. In an additional call, set the `ts` value of the final message as `latest` to get the next page of messages.

If a message has the same timestamp as `oldest` or `latest` it will not be included in the list. This functionality allows you to use the timestamps of specific messages as boundaries for the results. You can, however, have both timestamps be included in the time range by setting `inclusive` to `true`. The `inclusive` parameter is ignored when `oldest` or `latest` is not specified.

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

Value for `channel` was missing or invalid.

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

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`thread_not_found`

Value for `ts` was missing or invalid.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.