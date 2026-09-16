Source: https://docs.slack.dev/reference/events/message/message_deleted/
Retrieved: 2026-09-15

## Usage info

This example includes both the [event wrapper](/apis/events-api/#callback-field) and the [event structure](/apis/events-api/#event-type-structure) for the `message_deleted` message subtype event.

```json
{
    "token": "XXYYZZ",
    "team_id": "T123ABC456",
    "api_app_id": "A123ABC456",
    "event": {
        "type": "message",
        "subtype": "message_deleted",
        "hidden": true,
        "channel": "C123ABC456",
        "ts": "1358878755.000001",
        "deleted_ts": "1358878749.000002"
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

A `message_deleted` message is sent when a message in a channel is deleted, usually via the [`chat.delete` method](/reference/methods/chat.delete).

The `deleted_ts` property gives the timestamp of the message that was deleted.

If clients find an existing message with the same `deleted_ts` and `channel`, the existing message should be removed from the local model and UI. The original message will no longer return in history calls.

All types of messages are eligible for deletion, not just user-sent messages.