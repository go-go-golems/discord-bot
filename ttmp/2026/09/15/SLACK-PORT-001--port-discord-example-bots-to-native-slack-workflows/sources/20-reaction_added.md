Source: https://docs.slack.dev/reference/events/reaction_added/
Retrieved: 2026-09-15

## Facts

**Required Scopes**

[`reactions:read`](/reference/scopes/reactions.read)

**Compatible APIs**

## Usage info

This example includes both the [event wrapper](/apis/events-api/#callback-field) and the [event structure](/apis/events-api/#event-type-structure) for the `reaction_added` event.

```json
{
    "token": "XXYYZZ",
    "team_id": "T123ABC456",
    "api_app_id": "A123ABC456",
    "event": {
        "type": "reaction_added",
        "user": "U123ABC456",
        "reaction": "thumbsup",
        "item_user": "U222222222",
        "item": {
            "type": "message",
            "channel": "C123ABC456",
            "ts": "1360782400.498405",
            "channel_type": "channel"
        },
        "event_ts": "1360782804.083113"
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

When a reaction is added to an item the `reaction_added` event is sent to all connected clients for users who can see the content that was reacted to.

The `user` field indicates the ID of the user who performed this event. The `item_user` field represents the ID of the user that created the original item that has been reacted to.

Some messages aren't authored by "users," like those created by [incoming webhooks](/messaging/sending-messages-using-incoming-webhooks). `reaction_added` events related to these messages will not include an `item_user`.

The `item` field is a brief reference to what was reacted to. The above example describes a reaction to a message.

### Embedded item objects

Embedded `item` nodes are more lightweight than the structures you'll find in [`reactions.list`](/reference/methods/reactions.list).

Here are some examples:

#### Message:

```json
"item": {
    "type": "message",
    "channel": "C123ABC456",
    "ts": "1360782400.498405"
}
```

#### File:

```json
"item": {
    "type": "file",
    "file": "F123ABC456"
}
```

#### File Comment:

```json
"item": {
    "type": "file_comment",
    "file_comment": "Fc123ABC456",
    "file": "F123ABC456"
}
```