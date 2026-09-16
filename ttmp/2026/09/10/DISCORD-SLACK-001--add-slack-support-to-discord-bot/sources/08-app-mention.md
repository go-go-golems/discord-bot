## Facts

**Required Scopes**

[`app_mentions.read`](/reference/scopes/app_mentions.read)

**Compatible APIs**

[`Events`](/apis/events-api)

## Usage info

This example includes both the [event wrapper](/apis/events-api/#callback-field) and the [event structure](/apis/events-api/#event-type-structure) for the `app_mention` event.

```json
{
    "token": "XXYYZZ",
    "team_id": "T123ABC456",
    "api_app_id": "A123ABC456",
    "event": {
        "type": "app_mention",
        "user": "U061F7AUR",
        "text": "<@U0LAN0Z89> is it everything a river should be?",
        "ts": "1515449522.000016",
        "channel": "C123ABC456",
        "event_ts": "1515449522000016"
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

This [app event](/apis/events-api/#app_events) allows your app to subscribe to [`message`](/reference/events/message) events that directly mention your bot users.

Your Slack app must have a bot user configured and installed to use this event.

If a user mentions your app in a conversation that your app isn't party to, the user has the option to invite your app to the conversation. In that case, you'll receive an `app_mention` [event for the invite](#invite).

You'll receive *only* the messages pertinent to your app with this event. If your app is mentioned but not part of a conversation (and not invited to join), you won't receive an event.

The `app_mention` event is not a [`message`](/reference/events/message) like the `message.*` event types. However, your app can treat its contents similarly.

Messages sent to your app in *direct message* conversations are **not** dispatched via `app_mention`, including messages sent from other apps, regardless of whether your app is explicitly mentioned or otherwise. Subscribe to [`message.im`](/reference/events/message.im) events to receive messages directed to your bot user in direct message conversations.

## Example event payloads

`app_mention` events are just like other `message` events sent over the Events API, but their `type` indicates `app_mention`.

Whether your app has just been invited to a channel *or* is receiving a mention after it's already in a channel, the format and shape is the same.

Read on for more detail on the two scenarios with full Events API envelopes.

- [App mention that invites your app to a channel](#invite)
- [App mention when your app is already in channel](#mention)

### App mention that invites your app to a channel

When your app is directly invited to a channel as a result of being mentioned, you'll receive that inciting message as an `app_mention`:

```json
{
    "token": "ZZZZZZWSxiZZZ2yIvs3peJ",
    "team_id": "T061EG9R6",
    "api_app_id": "A0MDYCDME",
    "event": {
        "type": "app_mention",
        "user": "U0123ABCD",
        "text": "You can count on <@U0LAN0Z89> for an honorable mention.",
        "ts": "1515449483.000108",
        "channel": "C123ABC456",
        "event_ts": "1515449483000108"
    },
    "type": "event_callback",
    "event_id": "Ev0MDYHUEL",
    "event_time": 1515449483000108,
    "authorizations": [
        {
            "user_id": "U0123ABCD",
            "is_bot": false,
            "is_enterprise_install": false,
        }
    ],
}
```

You'll need to monitor your app's channel membership separately from this event. Consider subscribing to [`member_joined_channel`](/reference/events/member_joined_channel) events too.

### Standard app mention when your app is already in channel

Here's the happiest path: you're already part of a conversation and your app is mentioned. It can respond freely after receiving an event like this:

```json
{
    "token": "ZZZZZZWSxiZZZ2yIvs3peJ",
    "team_id": "T123ABC456",
    "api_app_id": "A123ABC456",
    "event": {
        "type": "app_mention",
        "user": "U123ABC456",
        "text": "What is the hour of the pearl, <@U0LAN0Z89>?",
        "ts": "1515449522.000016",
        "channel": "C123ABC456",
        "event_ts": "1515449522000016"
    },
    "type": "event_callback",
    "event_id": "Ev123ABC456",
    "event_time": 1515449522000016,
    "authorizations": [
        {
            "team_id": "T123ABC456",
            "user_id": "U123ABC456",
            "is_bot": false,
            "is_enterprise_install": false,
        }
    ],
}
```