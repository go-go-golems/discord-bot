Source: https://docs.slack.dev/reference/events/message/message_changed/
Retrieved: 2026-09-15

## Usage info

This example includes both the [event wrapper](/apis/events-api/#callback-field) and the [event structure](/apis/events-api/#event-type-structure) for the `message_changed` message subtype event.

```json
{
    "token": "XXYYZZ",
    "team_id": "T123ABC456",
    "api_app_id": "A123ABC456",
    "event": {
        "type": "message",
        "subtype": "message_changed",
        "hidden": true,
        "channel": "C123ABC456",
        "ts": "1358878755.000001",
        "message": {
            "type": "message",
            "user": "U123ABC456",
            "text": "Hello, world!",
            "ts": "1355517523.000005",
            "edited": {
                "user": "U123ABC456",
                "ts": "1358878755.000001"
            }
        }
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

A `message_changed` message is sent when a message in a channel is edited using the [`chat.update`](/reference/methods/chat.update) API method. The `message` property contains the updated message object.

This event can also sometimes be triggered by our automatic language detection, which can add or update language or locale information to the metadata for the message, prompting the event to be dispatched.

When clients receive this message type, they should look for an existing message with the same `message.ts` in that `channel`. If they find one, the existing message should be replaced with the new one.