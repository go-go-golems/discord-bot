Source: https://docs.slack.dev/apis/events-api/using-socket-mode/
Retrieved: 2026-09-15

Socket Mode allows your app to use [the Events API](/apis/events-api/) and [interactive features](/interactivity) — *without* exposing a public HTTP **Request URL**.

Instead of sending payloads to a public endpoint, Slack will use a [WebSocket URL](https://tools.ietf.org/html/rfc6455) to communicate with your app. WebSockets use a bidirectional stateful protocol with low latency to communicate between two parties—in this case, Slack and your app.

Unlike a public HTTP endpoint, the WebSocket URL you listen to is not static. The URL is created at runtime by calling the [`apps.connections.open`](/reference/methods/apps.connections.open) method, and it refreshes regularly.

Socket Mode helps developers working behind a corporate firewall, or who have other security concerns that don't allow exposing a static HTTP endpoint.

You can still switch between a direct HTTP endpoint and Socket Mode at any time in the [app settings](https://api.slack.com/apps).

We recommend [using our Bolt framework or SDKs for Java, JavaScript, or Python](/apis/events-api/using-socket-mode#sdks) to handle the details of Socket Mode. The process is streamlined, and you get access to all the other pleasant features of our SDKs.

---

## Setting up an app for Socket Mode

App setup for Socket Mode has three steps:

1. [Create an app](#creating)
2. [Toggle on Socket Mode](#toggling)
3. [Generate an app-level token](#token)

### 1\. Create an app

If you don't have one already, you'll need to create a Slack app:

Fill out your **App Name** and select the Development Workspace where you'll play around and build your app.

Socket Mode apps currently can't be listed in the Slack Marketplace. If you'd like to distribute your Socket Mode app to your entire Enterprise organization, you can [make your app deployable organization-wide](/enterprise/organization-ready-apps).

#### Creating an app with a manifest

You can [create an app through a manifest](/app-manifests/configuring-apps-with-app-manifests) and turn on Socket Mode by configuring a few options. Below is a YAML example which includes adding a bot user with the [`app_mentions:read`](/reference/scopes/chat.write.customize) scope and subscribing to the [`app_mention`](/reference/events/app_mention) event subscription:

### 2\. Toggle on Socket Mode

In your [app settings](https://api.slack.com/apps), navigate to the **Socket Mode** section. Toggle the **Enable Socket Mode** button to turn on receiving payloads via WebSockets.

Socket Mode can be toggled on or off whenever you'd like. When you toggle Socket Mode on, you'll **only** receive events and interactive payloads over your WebSocket connections—not over HTTP.

If your app is actively receiving events when you toggle Socket Mode on, you may lose events until you establish a connection to your WebSocket URL.

At this point, you can subscribe to some events in the **Event Subscriptions** section of your app settings. A favorite is the [`app_mention`](/reference/events/app_mention) event.

When using Socket Mode, your app does not need a Request UR to use the [Events API](/apis/events-api/). Your app's connection to the WebSocket replaces the need for Slack to dispatch to a Request URL.

### 3\. Generate an app-level token

When you toggle Socket Mode on, the app settings will guide you through obtaining an app-level token if you haven't created one.

You can also create or generate an app-level token manually. Under **Basic Information**, scroll to the **App-level tokens** section and click the button to generate an [app-level token](/authentication/tokens#app).

Your app-level token allows your app, either directly or with an SDK, to generate a WebSocket URL for communication with Slack via the [`apps.connections.open`](/reference/methods/apps.connections.open) method.

You're finished with setting up your app for Socket Mode. The rest can be handled by [using the Slack SDKs](#sdks). Or, you can [implement the Socket Mode protocol yourself](#implementing).

---

## Using Socket Mode with the Bolt framework

Socket Mode is automatically supported our Bolt framework and Slack SDKs.

The only thing you need to do is set your app-level token as an environment variable:

```markdown
export SLACK_APP_TOKEN='xapp-***'
```

Now you can lean on the [Bolt framework](/tools#bolt) for [JavaScript](/tools/bolt-js), [Java](/tools/java-slack-sdk), or [Python](/tools/bolt-python) to take care of the rest.

Here's some code to get your app capturing events using Socket Mode using Bolt:

If a Bolt framework suits your goals, then you're ready to to continue developing your app as normal! No need to follow the rest of this guide.

---

## Implementing Socket Mode without Bolt

If you prefer to implement the Socket Mode protocol yourself, read on.

You'll need a library or programming language that supports the [WebSocket](https://tools.ietf.org/html/rfc6455) protocol.

Connections refresh regularly. You should be ready to receive and connect to new WebSocket URLs as quickly as possible to maintain service.

### 1\. Call the apps.connections.open endpoint

Call the [`apps.connections.open`](/reference/methods/apps.connections.open) endpoint with your app-level token to receive a WebSocket URL:

```markdown
curl -X POST "https://slack.com/api/apps.connections.open" \
-H "Content-type: application/x-www-form-urlencoded" \
-H "Authorization: Bearer xapp-1-123"
```

Remember to send the token in the `Authorization` header, not as a parameter.

You'll receive a URL in response:

```markdown
{
    "ok": true,
    "url": "wss:\/\/wss.slack.com\/link\/?ticket=1234-5678"
}
```

You'll notice that, with Socket Mode turned on, your app settings doesn't require or even allow you to enter a **Request URL**. That's because this WebSocket URL replaces the public Request URL that Slack would have sent payloads to.

You can turn off Socket Mode **at any time** to return to the direct HTTP protocol for events and interactive features. If you've already set a **Request URL**, it'll be saved for later once you enter Socket Mode. When you turn Socket Mode off, the Request URL will be used once again to receive events.

---

### 2\. Connect to the WebSocket

Use your WebSocket library to connect to the URL specified in the above response.

Here's an example in JavaScript:

```markdown
if (response.ok) {
  let wssUrl = response.url;
  let socket = new WebSocket(wssUrl);

  socket.onopen = function(e) {
    // connection established
  }

  socket.onmessage = function(event) {
    // application received message
  }
}
```

One helpful trick: you can append `&debug_reconnects=true` to your WebSocket URL when you connect to it in order to make the connection time significantly shorter (360 seconds). That way, you can test and debug reconnects without waiting around.

After you connect to the WebSocket, Slack will send a `hello` message:

```markdown
{
    "type": "hello",
    "connection_info": {
        "app_id": "A1234"
    },
    "num_connections": 1,
    "debug_info": {
        "host": "applink-….",
        "started" "2020-10-11 12:12:12.120",
        "build_number": 54,
        "approximate_connection_time": 3600
    }
}
```

The `approximate_connection_time` (in seconds) can be used to estimate how long the connection will persist until Slack refreshes it.

### Using multiple connections

Socket Mode allows your app to maintain *up to 10* open WebSocket connections at the same time.

When multiple connections are active, each payload may be sent to *any* of the connections. It's best not to assume any particular pattern for how payloads will be distributed across multiple open connections.

There are a few reasons to make use of multiple active connections:

- If you'd like to handle a scheduled connection restart gracefully, you can generate an additional connection before the restart occurs.
- If you're having trouble keeping up with a large throughput of events from one connection, multiple connections can allow you to load balance.
- If you'd like to gracefully restart your app's services, you can use multiple connections for temporary active-active redundancy.

---

### 3\. Handle disconnects gracefully

Expect disconnects to your WebSocket connection. These may happen when you toggle off Socket Mode in the app settings, or for other reasons.

Here's a disconnect message you'll receive if you toggle off Socket Mode:

```markdown
{
  "type": "disconnect",
  "reason": "link_disabled",
  "debug_info": {
    "host": "wss-111.slack.com"
  }
}
```

No matter what, you'll need to handle connection refreshes once every few hours.

You may receive a warning about 10 seconds before the disconnect:

```markdown
{
  "type": "disconnect",
  "reason": "warning",
  "debug_info": {
    "host": "wss-111.slack.com"
  }
}
```

Even if you don't receive a warning, you'll still want to expect a `refresh_requested` message:

```markdown
{
  "type" : "disconnect",
  "reason": "refresh_requested",
  "debug_info": {
    "host": "wss-111.slack.com"
  }
}
```

You may want to use [multiple connections](#connections) in order to maintain uptime during a connection restart.

---

### 4\. Receive events

Event payloads sent to your app via Socket Mode are identical to the typical [Events API](/apis/events-api/) payloads, with some additional metadata:

```markdown
{
  "payload": <event_payload>,
  "envelope_id": <unique_identifier_string>,
  "type": <event_type_enum>,
  "accepts_response_payload": <accepts_response_payload_bool>
}
```

You can use the `accepts_response_payload` to determine whether a response can include additional information.

Your app still needs to acknowledge receiving *each event* so that Slack knows whether to retry.

### 5\. Acknowledge events

While acknowledging each event is required, there's no need to verify or validate inbound events, because you're receiving the events over a pre-authenticated WebSocket. Note that this is a different pattern from receiving events directly over HTTP, where validation is required for each event.

Use the `envelope_id` field in the object you receive from your WebSocket to send a response back to Slack acknowledging that you've received the event:

```markdown
{
    "envelope_id": <$unique_identifier_string>,
    "payload": <$payload_shape> // optional
}
```

---

## Using interactive features

[Interactive features](/interactivity) also send payloads as normal to your app—along with additional metadata specific to Socket Mode. Here's the general structure to expect:

```markdown
{
  "payload": <interactive_component_payload>,
  "envelope_id": <unique_identifier_string>,
  "type": <event_type_enum>,
  "accepts_response_payload": <accepts_response_payload_bool>
}
```

See below for specific examples on the following:

- [Slash commands](#command)
- [Block Kit buttons](#button)
- [App Home interactions](#home)
- [Modals](#modal)
- [Dynamic menus](#menu)

### Slash commands

Here's what's sent:

```markdown
{
  "payload": {
    "token": "bHKJ2n9AW6Ju3MjciOHfbA1b",
    "team_id": "T123ABC456",
    "team_domain": "maria",
    "channel_id": "C123ABC456",
    "channel_name": "general",
    "user_id": "U123ABC456",
    "user_name": "rainer",
    "command": "/randorilke",
    "text": "",
    "response_url": "https://rilke.slack.com/commands/T0SNL8S4S/37053613554/YMB2ZESDLNjNLqSFZ1quhNAh",
    "trigger_id": "37053613634.26768298162.440952c06ef4de2653466a48fe495f93"
  },
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545",
  "type": "slash_commands",
  "accepts_response_payload": true
}
```

Here's an example response from your app:

```markdown
{
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545",
  "payload": {
    "blocks": [
      {
        "type": "section",
        "text": {
          "type": "mrkdwn",
          "text": "Book of Hours"
        }
      },
      {
        "type": "section",
        "text": {
          "type": "mrkdwn",
          "text": "Duino Elegies"
        }
      }
    ]
  }
}
```

### Block Kit buttons

Here's what your app receives:

```markdown
{
  "payload": {
    "type": "block_actions",
    "team": {
      "id": "T123ABC456",
      "domain": "Duino"
    },
    "user": {
      "id": "U123ABC456",
      "username": "RMR",
      "team_id": "T123ABC456"
    },
    "api_app_id": "AABA1ABCD",
    "token": "9s8d9as89d8as9d8as989",
    "container": {
      "type": "message_attachment",
      "message_ts": "1548261231.000200",
      "attachment_id": 1,
      "channel_id": "C123ABC456",
      "is_ephemeral": false,
      "is_app_unfurl": false
    },
    "trigger_id": "12321423423.333649436676.d8c1bb837935619ccad0f624c448ffb3",
    "channel": {
      "id": "C123ABC456",
      "name": "review-updates"
    },
    "message": {
      "bot_id": "B123ABC456",
      "type": "message",
      "text": "Who if I cried out would hear me.",
      "user": "U123ABC456",
      "ts": "1548261231.000200",
      ...
    },
    "response_url": "https://hooks.slack.com/actions/AABA1ABCD/1232321423432/D09sSasdasdAS9091209",
    "actions": [{
      "action_id": "WaXA",
      "block_id": "=qXel",
      "text": {
        "type": "plain_text",
        "text": "View",
        "emoji": true
      },
      "value": "click_me_123",
      "type": "button",
      "action_ts": "1548426417.840180"
    }]
  },
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545",
  "type": "interactive",
  "accepts_response_payload": true
}
```

Here's an example response from your app:

```markdown
{
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545"
}
```

### App Home interactions

Here's what your app receives:

```markdown
{
  "payload": {
    "type": "app_home_opened",
    "user": "U123ABC456",
    "channel": "C123ABC456",
    "event_ts": "1515449522000016",
    "tab": "home",
    "view": {
      "id": "V123ABC456",
      "team_id": "T123ABC456",
      "type": "home",
      "blocks": [
        ...
      ],
      "private_metadata": "",
      "callback_id": "",
      "state":{
        ...
      },
      "hash":"1231232323.12321312",
      "clear_on_close": false,
      "notify_on_close": false,
      "root_view_id": "V123ABC456",
      "app_id": "A123ABC456",
      "external_id": "",
      "app_installed_team_id": "T123ABC456",
      "bot_id": "B123ABC456"
    }
  }
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545",
  "type": "events_api",
  "accepts_response_payload": false
}
```

Here's an example response from your app:

```markdown
{
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545"
}
```

### Modals

Here's what your app receives:

```markdown
{
  "payload": {
    "type": "view_submission",
    "team": { ... },
    "user": { ... },
    "view": {
      "id": "V123ABC456",
      "type": "modal",
      "title": { ... },
      "submit": { ... },
      "blocks": [ ... ],
      "private_metadata": "shhh-its-secret",
      "callback_id": "modal-with-inputs",
      "state": {
        "values": {
          "multi-line": {
            "ml-value": {
              "type": "plain_text_input",
              "value": "Archaic torso of Apollo"
            }
          }
        }
      },
      "hash": "156663117.cd33ad1f"
    }
  },
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545",
  "type": "interactive",
  "accepts_response_payload": true
}
```

Here's an example response from your app:

```markdown
{
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545",
  "payload": {
    "response_action": "update",
    "view": {
      "type": "modal",
      "callback_id": "updated-view-id",
      "title": {
        "type": "plain_text",
        "text": "Updated view"
      },
      "blocks": [{
        "type": "section",
        "text": {
          "type": "plain_text",
          "text": "You must change your life."
        }
      }]
    }
  }
}
```

Here's what your app receives:

```markdown
{
  "payload": {
    "type": "block_suggestion",
    "user": {
      "id": "U123ABC456",
      "name": "panther"
    },
    "team": {
      "id": "T123ABC456",
      "domain": "rilke"
    },
    "block_id": "search-block",
    "action_id": "search-action",
    "value": "an",
    "view": {"id": "V111", "type": "modal", "callback_id": "view-id"}
  },
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545",
  "type": "interactive",
  "accepts_response_payload": true
}
```

Here's an example response from your app:

```markdown
{
  "envelope_id": "dbdd0ef3-1543-4f94-bfb4-133d0e6c1545",
  "payload": {
    "options": [
      {
        "text": {
          "type": "plain_text",
          "text": "Give me your hand"
        },
        "value": "AI-2323"
      },
      {
        "text": {
          "type": "plain_text",
          "text": "Beauty and terror"
        },
        "value": "SUPPORT-42"
      }
    ]
  }
}
```