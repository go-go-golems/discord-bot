When you create a Slack app, it resides in one workspace. You can also distribute Slack apps in either listed or unlisted fashion. In this guide, we'll take a look at the different types of distribution available and explain how to set up installation flows and authenticate users for each.

---

## Undistributed apps

When you [create a Slack app](/app-management/quickstart-app-settings#creating), you associate it with a workspace: to do this, [open the app dashboard](https://api.slack.com/apps) for your app and click **Install App to Workspace** within the **Install App** section.

After installing your app, you'll get a single [access token](/authentication/tokens) that can be used to [authenticate](/authentication/installing-with-oauth#using) API method calls on behalf of the app.

Undistributed apps exist on a single workspace and can use the full range of app capabilities, but they can't be distributed to other workspaces. In addition, if your single-workspace app needs to take action on behalf of other users, you'll need to [build an OAuth 2.0 flow](/authentication/installing-with-oauth).

---

## Unlisted distributed apps

By default, a newly built Slack app can only be installed in its associated workspace as mentioned above. Installing your app on other workspaces means that you'll need to set up an [OAuth 2.0](/authentication/installing-with-oauth) flow. Once you've done that, your app will be able to generate [access tokens](/authentication/tokens) for each workspace and user.

Distributing your unlisted app is perfect for when you want to test out your app by running a pilot for early customers; however, apps intended for commercial distribution should be submitted and approved for listing in the Slack Marketplace. Unlike unlisted distributed apps, apps listed in the Slack Marketplace are reviewed against our [requirements & guidelines](/slack-marketplace/slack-marketplace-app-guidelines-and-requirements) to ensure a quality experience for end users.

### Preparing your app for distribution

Before you can distribute your app, there are a few steps to complete that are discussed in the following sections.

#### Handling installations

Once created, your app can be installed to its [associated workspace](/app-management/distribution) without any code for handling authorization. The [one-click install](/app-management/distribution) for an unlisted single-workspace app generates an [access token](/authentication/tokens), which can then be used to [authenticate](//authentication#using_tokens) API requests, but only within that associated workspace.

Therefore, when you're planning to distribute your app to other workspaces, you need to handle authorization. This will allow your app to be installed to any workspace, generating and using an access token programmatically.

✨ Read our [OAuth 2.0](/authentication) guide to understand the flow required for your app to request permissions and to generate access tokens.

In addition, once distributed, your app could potentially be installed to a Slack organization in an Enterprise org.

✨ Read our guide to [supporting Enterprise orgs in apps](/enterprise/developing-for-enterprise-orgs#support) to understand how an Enterprise org environment can affect the way your app works.

#### Creating your onboarding flow

[Onboarding users](/concepts/app-design#onboarding) is an important consideration for any app. For distributed apps, it is especially crucial as your user base could grow from a single workspace to potentially hundreds. Providing direct, hands-on support may be possible with a single workspace, but won't scale once your app is shared with multiple workspaces.

✨ Read our guide to [creating a helpful onboarding flow](/concepts/app-design#onboarding).

#### Enabling SSL

Slack apps open to installation by other workspaces have additional security requirements. Your app must support SSL for all of the following URLs:

- [OAuth redirect URLs](//authentication#redirect_urls)
- [Request URLs for interaction payloads](/interactivity/handling-user-interaction#payloads) from [interactive components](/block-kit#making-things-interactive), [actions](/interactivity/implementing-shortcuts#enabling_components) and [slash commands](/interactivity/implementing-slash-commands#app_command_handling)
- [Block Kit interactive component options load URLs](/reference/block-kit/block-elements/select-menu-element#external_selectt)
- [Events API request URLs](/apis/events-api/#events_api_request_urls)

#### Enabling admin approval

Some workspaces will restrict app installation so that only workspace administrators can provide authorization. Other workspaces may only allow the installation of apps officially listed in the Slack Marketplace.

Workspaces can also enable [admin approval](/admins/managing-app-approvals) requirements so that users can't directly install apps, but can request installation via a guided interface. Administrators can then screen these requests and selectively approve apps for installation. The list of permissions an app requests is also important to the decisions admins make, so ensure that your app only requests the permissions it absolutely needs.

### Enabling unlisted distribution

If you think your app is ready to be distributed, follow these steps:

1. Go to [your app's dashboard](https://api.slack.com/apps).
2. Scroll to *Share Your App with Other Workspaces* within the *Manage Distribution* section.
3. Ensure all items in the supplied checklist are completed.
4. Click **Activate Public Distribution**.

Once distribution is enabled, you'll get access to an embeddable **Add to Slack** button, a shareable URL that kicks off the installation process when clicked, as well as an HTML meta tag to enable [Slack app suggestions](/slack-marketplace/distributing-your-app-in-the-slack-marketplace#suggestions).

You can also go a step further and submit your app to the [Slack Marketplace](#slack_marketplace_apps). This is recommended for apps intended for commercial distribution.

### Disabling unlisted distribution

If distribution is enabled, you can turn it off by clicking **Deactivate Public Distribution** in the same location you [enabled it](#enabling). This will remove the app from any workspaces that have installed it, aside from the original associated workspace.

Note that you cannot disable distribution in this way once your app is published within the Slack Marketplace; instead, you must discontinue your listing there. Follow the instructions in the [Slack Marketplace guide](/slack-marketplace/distributing-your-app-in-the-slack-marketplace#discontinuing) to discontinue listing your app.

### Uninstalling apps

Your app can be uninstalled at any time from a workspace. Subscribe to the [`app_uninstalled`](/reference/events/app_uninstalled) event if you want your app to receive a notification when this happens.

In addition, apps will be automatically uninstalled from workspaces under the following circumstances:

- Apps only using the following [scopes](/reference/scopes) — [`bot`](/reference/scopes/bot), [`incoming-webhook`](/reference/scopes/incoming-webhook), [`commands`](/reference/scopes/commands), and [`identify`](/reference/scopes/identify) —will generally not be automatically uninstalled. However, there is one exception: in the app's associated workspace, the app will be uninstalled when the last creator or [collaborator](https://slack.com/help/articles/7887743766291-Add-internal-app-collaborators) leaves the workspace or becomes a guest user.
- Apps that use scopes beyond those listed above will be automatically uninstalled if the user who installed it leaves the workspace or becomes a guest user.

---

## Listed distributed apps (Slack Marketplace apps)

The Slack Marketplace contains distributed Slack apps that have been reviewed to ensure they meet our standards of quality. Distributing to the Slack Marketplace is the best path for safe and trustworthy third-party apps, and is perfect for when your app has been previously tested with early customers and is ready to be distributed more broadly or at commercial scale.

You can also take advantage of [direct installs](/slack-marketplace/slack-marketplace-review-guide#direct_install), which allow any user to install your app straight from the Slack Marketplace.

### Creating your app's listing page

If your app is approved for listing in the Slack Marketplace, it will have its very own listing page. Use this space to tell your app's story, to give a preview of your app's capabilities and user experience, and to attract users.

✨ Read our [Slack Marketplace guide](/slack-marketplace/slack-marketplace-review-guide) to help you build a great Slack Marketplace listing page, as well as to learn about the review and submission process.

---

## Automated deployment (CI/CD)

The preceding sections focused on the different types of app distribution and the necessary steps to make your app ready for installation by new workspaces. Once your app is distributed, maintaining and updating it becomes a critical task, especially across multiple environments (development, staging, production). To ensure consistency and reliability for all your users, the final step in a professional distribution pipeline is automating your deployment process using Continuous Integration/Continuous Delivery (CI/CD).

The [Slack CLI](/tools/slack-cli/)) is the best tool for managing apps at scale, allowing you to move from manual UI configuration to professional lifecycle management via your terminal or automation platform. The following practices outline how to integrate the [Slack CLI](/tools/slack-cli) and its deploy hooks into a robust CI/CD workflow, ensuring that your app settings and code are always synced before and after deployment.

Ensure your app settings are synced and your app is reinstalled before pushing your code changes to production to ensure a smooth deployment process. In your template’s `.slack/hooks.json` file, add a deploy hook for deploying your application code to your provider of choice:

```json
{
  "hooks": {
    "get-hooks": "npx -q --no-install -p @slack/cli-hooks slack-cli-get-hooks",
    "deploy": "git push heroku main"
  }
}
```

This hook tells the Slack CLI to execute the Heroku deployment command whenever the deploy process is triggered. The `slack deploy` command will:

- Update your app settings if any changes were made to your `manifest.json` file
- Reinstall the app if needed based on app settings changes
- Run the deploy hook to push the application code to Heroku

You can take this a step further by adding a corresponding GitHub Action to your template that runs `slack deploy` on merges to the `main` branch, providing a clear, consistent CI/CD process for every Slack app in your organization. Here’s an example `.github/workflows/deploy.yml` file:

```markdown
name: Deploy Slack app

on:
  push:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest
    timeout-minutes: 5

    steps:
    - uses: actions/checkout@v4
      
      - name: Install Slack CLI
        run: |
          curl -fsSL https://downloads.slack-edge.com/slack-cli/install.sh | bash
      
      - name: Install Heroku CLI
        run: |
          curl https://cli-assets.heroku.com/install.sh | sh
      
      - name: Deploy to Slack and Heroku
        env:
          SLACK_SERVICE_TOKEN: ${{ secrets.SLACK_SERVICE_TOKEN }}
          HEROKU_API_KEY: ${{ secrets.HEROKU_API_KEY }}
        run: slack deploy -s --token $SLACK_SERVICE_TOKEN
```

---

## When is a reinstall required?

When you update your app manifest (via the App Config UI or the [`apps.manifest.update`](/reference/methods/apps.manifest.update) API method), changes fall into two categories: those that take effect immediately for existing installs and those that only take effect after workspaces reinstall your app.

Updating the manifest changes your app's configuration; it does not modify OAuth tokens that have already been issued to workspaces. If a change requires new permissions on existing tokens, the affected workspaces must reinstall.

A reinstallation is necessary when you change permissions or fundamental capabilities, as existing users must approve the new terms. Key changes that require a reinstall include:

- [Enabling org-wide deployment](/enterprise/migrating-to-organization-wide-deployment): Setting `org_deploy_enabled` to `true`.
- Changing scopes: Adding or removing scopes (e.g., adding `canvases:create`).
- Adding features (if those features require additional scopes): Enabling [AI features](/ai/developing-agents) for your app.
- Adding or modifying event subscriptions (if those events require additional scopes): Subscribing to new events like `member_joined_channel`.

### Changes that require reinstallation

#### Scopes

Scopes live under the `oauth_config.scopes` property of the manifest. Changing them updates what your app requests, but existing app installs keep their current tokens with the old scopes. Whenever scopes change, the `apps.manifest.update` API method response includes a `permissions_updated: true` flag. This is your signal that existing installs need to re-authorize to get the new permissions.

| Manifest property | Reinstall needed when... |
| --- | --- |
| `oauth_config.scopes.bot` | New bot scopes are added. |
| `oauth_config.scopes.user` | New user scopes are added. |
| `oauth_config.scopes.bot_optional` | New optional bot scopes are added. |
| `oauth_config.scopes.user_optional` | New optional user scopes are added. |

If you remove a scope, the app configuration updates immediately (new installs won't have that scope), but existing tokens retain the removed scope until the token is revoked or the workspaces uninstall the app. Scope removal may also trigger a breaking change notification to affected workspaces.

### Changes that do not require reinstallation

All manifest properties outside of `oauth_config.scopes` write directly to your app's configuration and take effect for existing installs right away (with an exception for `event_subscriptions`. See [event subscriptions](#event-subscriptions) above).

#### Display and identity

| Manifest property | Propagation behavior |
| --- | --- |
| `display_information.name` | App name updates immediately. Propagates to every installed team to update the app listing and the bot user profile. Note: changing the app name also triggers a bot profile update, which is one-way coupling (see [Edge cases](#edge-cases)). |
| `display_information.description` | Short description written to app config; notifies installed teams automatically. |
| `display_information.long_description` | Long description written to app config; notifies installed teams automatically. |
| `display_information.background_color` | Background color written to app config; notifies installed teams automatically. |

#### Features

| Manifest property | Propagation behavior |
| --- | --- |
| `features.bot_user.display_name` | Propagates to every installed team to update the bot user profile. Works independently of app name changes. Note: changing `display_information.name` also triggers bot profile updates (one-way coupling from app name to bot), but changing bot name does not trigger app name updates. |
| `features.bot_user.always_online` | Online presence flag (`always_active`) propagates to every installed team right away. |
| `features.app_home` | Home tab / messages tab settings written to app config; notifies installed teams automatically. |
| `features.slash_commands` | Config updates immediately. Commands only function on teams where the bot service is installed. |
| `features.shortcuts` | Global and message shortcuts update automatically on all installed teams. |
| `features.unfurl_domains` | Link unfurling domain list updates in app config; notifies teams automatically. |
| `features.workflow_steps` | Workflow step definitions written to app config; propagates to every installed team automatically. |
| `features.assistant_view` | Assistant surface config updates in app config; propagates to every installed team automatically. |
| `features.rich_previews` | Rich preview config updates in app config; propagates to every installed team automatically. |
| `features.search` | Search config updates in app config; propagates to every installed team automatically. |

#### Settings

| Manifest property | Propagation behavior |
| --- | --- |
| `settings.event_subscriptions.request_url` | Events are sent to the new URL immediately; both bot and user event subscription records update automatically on all installed teams. |
| `settings.event_subscriptions.bot_events` | Subscription config writes immediately; propagates to all installed teams automatically. However, at delivery time, bot events are filtered by the team's installed bot scopes. Adding an event that requires a scope the team hasn't installed yet will silently have no effect until the team reinstalls with that scope. |
| `settings.event_subscriptions.user_events` | Same as bot events: config writes immediately, but delivery is gated by the user token's installed scopes. Events requiring uninstalled scopes won't fire. |
| `settings.interactivity` | Request URL and message menu options URL update in app config; propagates to installed teams automatically. |
| `settings.socket_mode_enabled` | Delivery mechanism switches in app config. Takes effect for new event deliveries. |
| `settings.allowed_ip_address_ranges` | IP allowlist written to app config and enforced immediately on incoming API calls. |
| `settings.org_deploy_enabled` | Org-wide deployment eligibility flag updates in app config. |
| `settings.siws_links` | Sign in with Slack link config updates in app config. |

### Edge cases

There are a few edge cases to these general rules, outlined here.

| Scenario | Behavior |
| --- | --- |
| `oauth_config.redirect_urls` | Takes effect immediately for new OAuth flows. No reinstall needed, but the API response reports `permissions_update: true`. (The flag also fires for `pkce_enabled` and `token_management_enabled` changes, not just scopes.) |
| `oauth_config.token_management_enabled` | App-level config, takes effect for new token issuance. Cannot be turned off once enabled. |
| `oauth_config.pkce_enabled` | App-level config, takes effect for new OAuth flows. Cannot be turned off once enabled. |
| `settings.token_rotation_enabled` | Takes effect for newly issued tokens. Existing tokens don't automatically start rotating. |
| Adding `features.bot_user` for the first time | Requires an install to provision a bot token (no bot token exists yet). |
| Bot events and scope gating | Bot event subscriptions write immediately, but at delivery time each event type is filtered against the team's installed bot scopes. Adding an event that requires an uninstalled scope has no effect until that scope is granted via reinstall. |
| User events and scope gating | Same pattern as bot events: config writes immediately, delivery is filtered by the user token's scopes. Events requiring uninstalled user scopes won't fire. |
| App name to bot name coupling | Whenever the app name changes, the bot profile updates across all teams. This is a one-way coupling: changing the bot display name does not trigger app name updates. |
| `permissions_updated` flag scope | Returns true when any OAuth field changes, including `redirect_urls`, `pkce_enabled`, and `token_management_enabled`. This is not limited to scope changes, and callers should not assume this flag means that scopes changed. |