This quickstart guide covers how to create a Slack app using the [Slack CLI](/tools/slack-cli) and the Bolt frameworks, in either [Python](/tools/bolt-python) or [JavaScript](/tools/bolt-js).

By the end, you'll have a local environment configured with a running customized app to modify and make your own.

## Prerequisites

Before you begin building, grab the required tools by following the three steps below. If you don't have a Slack workspace already, we recommend also setting up a [developer sandbox](/tools/developer-sandboxes).

1

### Install the Slack CLI

Download the command-line tool for developing Slack apps.

`curl -fsSL https://downloads.slack-edge.com/slack-cli/install.sh | bash`

3

### Use our preconfigured template

Create a project by selecting "Starter App", then either "Bolt for Python" or "Bolt for JavaScript".

`slack create first-slack-app`

  

Now you've got everything you need to start developing on the Slack platform!

## Creating a project

The [`slack create`](/tools/slack-cli/reference/commands/slack_create) Slack CLI command created a new Bolt project on your system. This contains the code that handles logic for your app. Navigate into that project:

```markdown
$ cd first-slack-app
```

and open the project in an editor of your choice.

## Running the app

The starter app is ready to run from the start. By using the [`slack run`](/tools/slack-cli/reference/commands/slack_run) command, and selecting "Create a new app", you'll install the app to the team you choose. The details of the app shown are found in the `manifest.json` file.

```sh
$ slack run
...
INFO:slack_bolt.App:⚡️ Bolt app is running!
...
```

With the app running, test it out with the following steps in the Slack workspace your app is running in:

1. Open a direct message with your app or invite the bot `@first-slack-app (local)` to a public channel.
2. Send "hello" to the current conversation and wait for a response.
3. Click the attached button labeled "Click Me" to post another reply.

After confirming the app responds, celebrate! But now it's time to add on.

## Updating the app

You may have a running Slack app, but the defaults included leave opportunities abound. Let's personalize this app by editing this code to respond with a kind farewell.

Once the file is updated, save the changes. You should notice the app restarting in your terminal.

Open Slack to perform these steps:

1. Return to the direct message or public channel with your bot.
2. Send "goodbye" to the conversation.
3. Receive a parting response from before and repeat "goodbye" to find another one.

Your app can be stopped by pressing `CTRL+C` in the terminal to end these chats.

#### Customizing app settings

The created app will have some placeholder values and a small set of [scopes](/reference/scopes) to start, but we recommend exploring the customizations possible on app settings.

Open app settings for your app with the following command:

```sh
$ slack app settings
```

This will open the following page in a web browser:

![Basic Information page](/assets/images/basic-information-page-e7d531fe4721830376d61a91de5d933e.png "Basic Information page")

On these pages you're free to make changes such as updating your app icon, configuring app features, and perhaps even distributing your app!

You can now continue customizing your app with various features to make it right for whatever job's at hand. Here are some ideas about what to explore next:

- Check out the [Agent quickstart](/ai/agent-quickstart) to get up and running with an agent.
- Browse our [curated catalog of samples](/samples) for more apps to use as a starting point for development.
- Try using the [Slack MCP Plugin](/ai/slack-skills-plugin) to aid in agent-assisted development.