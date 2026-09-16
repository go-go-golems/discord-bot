Follow along to get up and running with a Slack app using Bolt for Python. Along the way, we’ll create a new Slack app, set up your local environment, and develop an app that listens and responds to messages from a Slack workspace.

When you're finished, you'll have created the [Getting Started app](https://github.com/slackapi/bolt-python/tree/main/examples/getting_started) to run, modify, and make your own. ⚡️

---

### Create a new app

First thing's first: before you start developing with Bolt, you'll want to [create a Slack app](https://api.slack.com/apps/new).

After you fill out an app name (*you can change it later*) and pick a workspace to install it to, hit the `Create App` button and you'll land on your app's **Basic Information** page.

This page contains an overview of your app in addition to important credentials you'll need later.

![Basic Information page](/assets/images/basic-information-page-e7d531fe4721830376d61a91de5d933e.png "Basic Information page")

Look around, add an app icon and description, and then let's start configuring your app 🔩

---

### Tokens and installing apps

Slack apps use [OAuth to manage access to Slack's APIs](/authentication/installing-with-oauth). When an app is installed, you'll receive a token that the app can use to call API methods.

There are three main token types available to a Slack app: user (`xoxp`), bot (`xoxb`), and app-level (`xapp`) tokens.

- [User tokens](/authentication/tokens#user) allow you to call API methods on behalf of users after they install or authenticate the app. There may be several user tokens for a single workspace.
- [Bot tokens](/authentication/tokens#bot) are associated with bot users, and are only granted once in a workspace where someone installs the app. The bot token your app uses will be the same no matter which user performed the installation. Bot tokens are the token type that *most* apps use.
- [App-level tokens](/authentication/tokens#app-level) represent your app across organizations, including installations by all individual users on all workspaces in a given organization and are commonly used for creating WebSocket connections to your app.

We're going to use bot and app-level tokens for this guide.

1. Navigate to **OAuth & Permissions** on the left sidebar and scroll down to the **Bot Token Scopes** section. Click **Add an OAuth Scope**.
2. For now, we'll just add one scope: [`chat:write`](/reference/scopes/chat.write). This grants your app the permission to post messages in channels it's a member of.
3. Scroll up to the top of the **OAuth & Permissions** page and click **Install App to Workspace**. You'll be led through Slack's OAuth UI, where you should allow your app to be installed to your development workspace.
4. Once you authorize the installation, you'll land on the **OAuth & Permissions** page and see a **Bot User OAuth Access Token**.

![OAuth Tokens](/assets/images/bot-token-3d6c761238c7a66557fd08d00a2a1b0c.png "Bot OAuth Token")

1. Head over to **Basic Information** and scroll down under the App Token section and click **Generate Token and Scopes** to generate an app-level token. Add the `connections:write` scope to this token and save the generated `xapp` token.
2. Navigate to **Socket Mode** on the left side menu and toggle to enable.

---

### Setting up your project

With the initial configuration handled, it's time to set up a new Bolt project. This is where you'll write the code that handles the logic for your app.

If you don’t already have a project, let’s create a new one. Create an empty directory:

```sh
$ mkdir first-bolt-app
$ cd first-bolt-app
```

Next, we recommend using a [Python virtual environment](https://packaging.python.org/guides/installing-using-pip-and-virtual-environments/#creating-a-virtual-environment) to manage your project's dependencies. This is a great way to prevent conflicts with your system's Python packages. Let's create and activate a new virtual environment with [Python 3.7 or later](https://www.python.org/downloads/):

```sh
$ python3 -m venv .venv
$ source .venv/bin/activate
$ pip install -r requirements.txt
```

We can confirm that the virtual environment is active by checking that the path to `python3` is *inside* your project ([a similar command is available on Windows](https://packaging.python.org/guides/installing-using-pip-and-virtual-environments/#activating-a-virtual-environment)):

```sh
$ which python3
# Output: /path/to/first-bolt-app/.venv/bin/python3
```

Before we install the Bolt for Python package to your new project, let's save the **bot token** and **app-level token** that were generated when you configured your app.

1. **Copy your bot (xoxb) token from the OAuth & Permissions page** and store it in a new environment variable. The following example works on Linux and macOS; but [similar commands are available on Windows](https://superuser.com/questions/212150/how-to-set-env-variable-in-windows-cmd-line/212153#212153).
```sh
$ export SLACK_BOT_TOKEN=xoxb-<your-bot-token>
```
1. **Copy your app-level (xapp) token from the Basic Information page** and then store it in a new environment variable.
```sh
$ export SLACK_APP_TOKEN=<your-app-level-token>
```

Now, let's create your app. Install the `slack_bolt` Python package to your virtual environment using the following command:

```sh
$ pip install slack_bolt
```

Create a new file called `app.py` in this directory and add the following code:

```python
import os
from slack_bolt import App
from slack_bolt.adapter.socket_mode import SocketModeHandler

# Initializes your app with your bot token and socket mode handler
app = App(token=os.environ.get("SLACK_BOT_TOKEN"))

# Start your app
if __name__ == "__main__":
    SocketModeHandler(app, os.environ["SLACK_APP_TOKEN"]).start()
```

Your tokens are enough to create your first Bolt app. Save your `app.py` file then on the command line run the following:

```sh
$ python3 app.py
```

Your app should let you know that it's up and running. 🎉

---

### Setting up events

Your app behaves similarly to people on your team — it can post messages, add emoji reactions, and listen and respond to events.

To listen for events happening in a Slack workspace (like when a message is posted or when a reaction is posted to a message) you'll use the [Events API to subscribe to event types](/apis/events-api/).

For those just starting, we recommend using [Socket Mode](/apis/events-api/using-socket-mode). Socket Mode allows your app to use the Events API and interactive features without exposing a public HTTP Request URL. This can be helpful during development, or if you're receiving requests from behind a firewall.

That being said, you're welcome to set up an app with a public HTTP Request URL. HTTP is more useful for apps being deployed to hosting environments to respond within a large corporate Slack workspaces/organization, or apps intended for distribution via the Slack Marketplace.

We've provided instructions for both ways in this guide.

Navigate to **Event Subscriptions** on the left sidebar and toggle to enable. Under **Subscribe to Bot Events**, you can add events for your bot to respond to. There are four events related to messages:

- [`message.channels`](/reference/events/message.channels) listens for messages in public channels that your app is added to.
- [`message.groups`](/reference/events/message.groups) listens for messages in 🔒 private channels that your app is added to.
- [`message.im`](/reference/events/message.im) listens for messages in your app's DMs with users.
- [`message.mpim`](/reference/events/message.mpim) listens for messages in multi-person DMs that your app is added to.

If you want your bot to listen to messages from everywhere it is added to, choose all four message events. After you’ve selected the events you want your bot to listen to, click the green **Save Changes** button.

---

### Listening and responding to a message

Your app is now ready for some logic. Let's start by using the `message()` method to attach a listener for messages.

The following example listens and responds to all messages in channels/DMs where your app has been added that contain the word "hello":

If you restart your app, so long as your bot user has been added to the channel or DM conversation, when you send any message that contains "hello", it will respond.

This is a basic example, but it gives you a place to start customizing your app based on your own goals. Let's try something a little more interactive by sending a button rather than plain text.

---

### Sending and responding to actions

To use features like buttons, select menus, datepickers, modals, and shortcuts, you’ll need to enable interactivity. Head over to **Interactivity & Shortcuts** in your app configuration.

When interactivity is enabled, interactions with shortcuts, modals, or interactive components (such as buttons, select menus, and datepickers) will be sent to your app as events.

Now, let's go back to your app's code and add logic to handle those events:

- First, we'll send a message that contains an interactive component (in this case a button).
- Next, we'll listen for the action of a user clicking the button before responding.

Below, the code from the last section is modified to send a message containing a button rather than just a string:

The value inside of `say()` is now an object that contains an array of `blocks`. Blocks are the building components of a Slack message and can range from text to images to datepickers. In this case, your app will respond with a section block that includes a button as an accessory. Since we're using `blocks`, the `text` is a fallback for notifications and accessibility.

You'll notice in the button `accessory` object, there is an `action_id`. This will act as a unique identifier for the button so your app can specify which action it wants to respond to.

Now, if you restart your app and say "hello" in a channel your app is in, you'll see a message with a button. But if you click the button, nothing happens (*yet!*).

Let's add a handler to send a follow-up message when someone clicks the button:

You can see that we used `app.action()` to listen for the `action_id` that we named `button_click`. If you restart your app and click the button, you'll see a new message from your app that says you clicked the button.

---

You just built your first [Bolt for Python app](https://github.com/slackapi/bolt-python/tree/main/examples/getting_started)! 🎉

Now that you have a basic app up and running, you can start exploring how to make your Bolt app stand out. Here are some ideas about what to explore next:

- Read through the concepts pages to learn about the different methods and features your Bolt app has access to.
- Explore the different events your bot can listen to with the [`app.event()`](/tools/bolt-python/concepts/event-listening) method. View the full events reference docs [here](/reference/events).
- Bolt allows you to [call Web API methods](/tools/bolt-python/concepts/web-api) with the client attached to your app. There are over 200 methods; view them [here](/reference/methods).
- Learn more about the different token types in the [tokens guide](/authentication/tokens). Your app may need different tokens depending on the actions you want it to perform.