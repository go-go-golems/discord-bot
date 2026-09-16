Source: https://docs.slack.dev/reference/events/message/
Retrieved: 2026-09-15

## Usage info

A `message` can be delivered from multiple sources:

- They are returned via calls to [`conversations.history`](/reference/methods/conversations.history).
- They are sent via the [Real Time Messaging API](/legacy/legacy-rtm-api) when a message is sent to a channel to which you subscribe. This message should immediately be displayed in the client.

This example includes both the [event wrapper](/apis/events-api/#callback-field) and the [event structure](/apis/events-api/#event-type-structure) for the `message` event.

```json
{
    "token": "XXYYZZ",
    "team_id": "T123ABC456",
    "api_app_id": "A123ABC456",
    "event": {
        "type": "message",
        "channel": "C123ABC456",
        "user": "U123ABC456",
        "text": "Hello world",
        "ts": "1355517523.000005",
        "channel_type": "channel"
    },
    "type": "event_callback",
    "authorizations": [
        {
            "team_id": "T123ABC456",
            "user_id": "U123ABC456",
            "is_bot": false,
            "is_enterprise_install": false,
        }
    ],
    "event_id": "Ev123ABC456",
    "event_time": 123456789
}
```

The `channel` property is the ID of the channel, private group, or DM channel this message is posted in. The `user` property is the ID of the user, the `text` property is the message, and `ts` is the unique (per-channel) timestamp. Messages can also include an `attachments` property, containing a list of [attachment objects](/messaging/formatting-message-text).

If the message has been edited after posting, it will include an `edited` property, which includes the user ID of the editor and the timestamp for when the edit occurred. The original text of the message is not available.

Example:

```json
{
    "type": "message",
    "channel": "C123ABC456",
    "user": "U123ABC456",
    "text": "Hello, world!",
    "ts": "1355517523.000005",
    "edited": {
        "user": "U123ABC456",
        "ts": "1355517536.000001"
    }
}
```

If the app is subscribed to the [`app_context_changed`](/reference/events/app_context_changed) event, Slack will also send the `app_context` to these events. This context provides more information about what a user is currently viewing when they send a message.

Example:

```json
{
  "type": "message",
  "user": "U012345ABCDEF",
  "text": "What's the latest on agents?",
  "app_context": {
    "entities": [
      {
        "type": "slack#/types/channel_id",
        "value": "C0123456ABC",
        "team_id": "T0123ABCDEF"
      }
    ]
  },
  "channel": "D01ABCDEFG9",
  "channel_type": "im"
}
```

### Message events in the Events API

In the [Events API](/apis/events-api/), the exact message event is defined differently, and the specific event you subscribe to also dictates what scopes your app needs: for example, private channels require the `groups:history` scope. Depending on your associated OAuth scope or the channels your bot user belongs to, you must subscribe to one of the following events:

- [`message.channels`](/reference/events/message.channels) - for messages appearing within channels. Will only send message events from public channels to your app.
- [`message.im`](/reference/events/message.im) - for messages appearing within direct messages.
- [`message.groups`](/reference/events/message.groups) - subscribe to this event if you want message events from private channels.
- [`message.mpim`](/reference/events/message.mpim) - for messages appearing within multiparty direct messages.

The semantics for the events are the same — you're still receiving an event of type `message`.

### Message subtypes

Unlike other event types, `message` events can have a subtype. For example, this is a channel join message:

```json
{
    "type": "message",
    "subtype": "channel_join",
    "text": "<@U123ABC456|bobby> has joined the channel",
    "ts": "1403051575.000407",
    "user": "U123ABC456"
}
```

They are structured in this way so that clients can either support them fully by having distinct behavior for each different subtype, or they can fallback to just displaying the text of the message. The full list of message subtypes is below.

| Subtype | Description |
| --- | --- |
| [`assistant_app_thread`](/reference/events/message/assistant_app_thread) | The message sent is an app assistant thread |
| [`bot_message`](/reference/events/message/bot_message) | A message was posted by an integration |
| [`channel_archive`](/reference/events/message/channel_archive) | A channel was archived. |
| [`channel_convert_to_private`](/reference/events/message/channel_convert_to_private) | This channel was made private. Now, it can only be viewed or joined by invitation |
| [`channel_convert_to_public`](/reference/events/message/channel_convert_to_public) | This channel was made public. Any member in this workspace can see and join it |
| [`channel_join`](/reference/events/message/channel_join) | A member joined a channel |
| [`channel_leave`](/reference/events/message/channel_leave) | A member left a channel |
| [`channel_name`](/reference/events/message/channel_name) | A channel was renamed |
| [`channel_posting_permissions`](/reference/events/message/channel_posting_permissions) | The posting permissions for a channel changed |
| [`channel_purpose`](/reference/events/message/channel_purpose) | A channel purpose was updated |
| [`channel_topic`](/reference/events/message/channel_topic) | A channel topic was updated |
| [`channel_unarchive`](/reference/events/message/channel_unarchive) | A channel was unarchived |
| [`document_mention`](/reference/events/message/document_mention) | A bot is mentioned in the body of a canvas |
| [`ekm_access_denied`](/reference/events/message/ekm_access_denied) | Message content redacted due to Enterprise Key Management (EKM) |
| [`file_comment`](/reference/events/message/file_comment) | A comment was added to a file |
| [`file_mention`](/reference/events/message/file_mention) | A file was mentioned in a channel |
| [`file_share`](/reference/events/message/file_share) | A file was shared into a channel |
| [`group_archive`](/reference/events/message/group_archive) | A group was archived |
| [`group_join`](/reference/events/message/group_join) | A member joined a group |
| [`group_leave`](/reference/events/message/group_leave) | A member left a group |
| [`group_name`](/reference/events/message/group_name) | A group was renamed |
| [`group_purpose`](/reference/events/message/group_purpose) | A group purpose was updated |
| [`group_topic`](/reference/events/message/group_topic) | A group topic was updated |
| [`group_unarchive`](/reference/events/message/group_unarchive) | A group was unarchived |
| [`me_message`](/reference/events/message/me_message) | A /me message was sent |
| [`message_changed`](/reference/events/message/message_changed) | A message was changed |
| [`message_deleted`](/reference/events/message/message_deleted) | A message was deleted |
| [`message_replied`](/reference/events/message/message_replied) | A message thread received a reply |
| [`pinned_item`](/reference/events/message/pinned_item) | An item was pinned in a channel |
| [`reminder_add`](/reference/events/message/reminder_add) | A reminder was added to the channel |
| [`reply_broadcast`](/reference/events/message/reply_broadcast) | (No longer served) A message thread's reply was broadcast to a channel |
| [`thread_broadcast`](/reference/events/message/thread_broadcast) | A message thread's reply was broadcast to a channel |
| [`unpinned_item`](/reference/events/message/unpinned_item) | An item was unpinned from a channel |

### Hidden subtypes

Some subtypes have a special hidden property. These indicate messages that are part of the history of a channel, but should not be displayed to users. Examples include records of message edits or deletes:

```json
{
    "type": "message",
    "subtype": "message_deleted",
    "hidden": true,
    "channel": "C123ABC456",
    "ts": "1358878755.000001",
    "deleted_ts": "1358878749.000002",
    "event_ts": "1358878755.000002"
}
```

Hidden messages will not appear as the `latest` property on [`channel`](/reference/objects/channel-object), [`group`](/reference/objects/group-object) or [`im`](/reference/objects/im-object) objects. They will also not return in calls to [`conversations.history`](/reference/methods/conversations.history).

### Stars, pins, and reactions

Messages can have several extra properties depending on whether or not they have been starred, pinned, or reacted to:

```json
{
    "type": "message",
    "channel": "C123ABC456",
    "user": "U123ABC456",
    "text": "Hello world",
    "ts": "1355517523.000005",
    "is_starred": true,
    "pinned_to": ["C024BE7LT", ...],
    "reactions": [
        {
            "name": "astonished",
            "count": 3,
            "users": [ "U1", "U2", "U3" ]
        },
        {
            "name": "facepalm",
            "count": 1034,
            "users": [ "U1", "U2", "U3", "U4", "U5" ]
        }
    ]
}
```

The `is_starred` property is present and true if the calling user has starred the message, else it is omitted.

The `pinned_to` array, if present, contains the IDs of any channels to which the message is currently pinned.

The `reactions` property, if present, contains any reactions that have been added to the message and gives information about the type of reaction, the total number of users who added that reaction, and a (possibly incomplete) list of users who have added that reaction to the message. The users array in the `reactions` property might not always contain all users that have reacted (we limit it to X users, and X might change) — however, `count` will always represent the count of all users who made that reaction (i.e. it may be greater than `users.length`). If the authenticated user has a given reaction then they are guaranteed to appear in the `users` array, regardless of whether or not `count` is greater than `users.length`.