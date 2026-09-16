## Facts

**Required Scopes**

[`im:history`](/reference/scopes/im.history)

**Compatible APIs**

[`Events`](/apis/events-api)

## Usage info

This example includes both the [event wrapper](/apis/events-api/#callback-field) and the [event structure](/apis/events-api/#event-type-structure) for the `message.im` event.

```json
{
    "token": "one-long-verification-token",
    "team_id": "T123ABC456",
    "api_app_id": "A0PNCHHK2",
    "event": {
        "type": "message",
        "channel": "D024BE91L",
        "user": "U2147483697",
        "text": "Hello hello can you hear me?",
        "ts": "1355517523.000005",
        "event_ts": "1355517523.000005",
        "channel_type": "im"
    },
    "type": "event_callback",
    "authorizations": [
        {
            "team_id": "T123ABC456",
            "user_id": "U2147483697",
            "is_bot": false,
            "is_enterprise_install": false,
        }
    ],
    "event_id": "Ev0PV52K21",
    "event_time": 1355517523
}
```

The semantics for this message event type are similar to the core [message](/reference/events/message) event sent through the RTM API. Messages often include many more fields than those shown above.

Differentiate direct messages from other `message.*` events by looking for the `event` 's `channel_type` field set to `"im"`.