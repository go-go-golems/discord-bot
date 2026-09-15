## Facts

## Arguments

### Required arguments

**`token`** `string` Required

An [app configuration access token](/authentication/tokens#config).

*Example:* `xxxx-xxxxxxxxx-xxxx`

**`manifest`** `string` Required

A JSON app manifest encoded as a string. This manifest **must** use a valid [app manifest schema - read our guide to creating one](/app-manifests/configuring-apps-with-app-manifests#fields).

### Optional arguments

**`team_id`** `string` Optional

When called with an org token, which specific team to create app on

## Usage info

### Using manifests

This method accepts an app manifest as an argument. Read our [guide to app manifests](/app-manifests/configuring-apps-with-app-manifests) to learn how create or reuse them.

If you receive an `invalid_manifest` response when trying to use any App Manifest API, it indicates that the manifest you supplied didn't match the [correct schema](/reference/app-manifest#fields).

To better locate the problem with your manifest, the `invalid_manifest` error should be accompanied by an `errors` array:

```json
{
    "ok": false,
    "error": "invalid_manifest",
    "errors": [
        {
            "message": "Event Subscription requires either Request URL or Socket Mode Enabled",
            "pointer": "/settings/event_subscriptions"
        },
        {
            "message": "Interactivity requires a Request URL",
            "pointer": "/settings/interactivity"
        },
        {
            "message": "Interactivity requires Socket Mode enabled",
            "pointer": "/settings/interactivity"
        }
    ]
}
```

Each of the items in this array contain a `message` which describes the problem, and a `pointer` which indicates the problem's location within your supplied manifest. Use these two pieces of info to correct your manifest and try again.

---

## Response

#### Typical success response

```json
{
  "ok": true,
  "app_id": "A012ABCD0A0",
  "credentials": {
    "client_id": "...",
    "client_secret": "...",
    "verification_token": "...",
    "signing_secret": "..."
  },
  "oauth_authorize_url": "https://slack.com/oauth/v2/authorize?client_id=...&scope=commands,workflow.steps:execute"
}
```

## Errors

This table lists the expected errors that this method could return. However, other errors can be returned in the case where the service is down or other unexpected factors affect processing. Callers should always check the value of the `ok` parameter in the response.

Error

Description

`access_denied`

Access to a resource specified in the request is denied.

`accesslimited`

Access to this method is limited on the current network

`account_inactive`

Authentication token is for a deleted user or workspace when using a `bot` token.

`deprecated_endpoint`

The endpoint has been deprecated.

`dynamic_client_registration_failed`

Dynamic client registration failed for an MCP server that requires OAuth authentication.

`dynamic_client_registration_invalid_redirect_uri`

The remote server rejected Slack's redirect URI during Dynamic Client Registration. The server must allowlist https://oauth2.slack.com/external/auth/callback as a valid redirect URI.

`dynamic_client_registration_invalid_response`

The MCP server returned an invalid response during Dynamic Client Registration.

`dynamic_client_registration_missing_oauth_endpoints`

The MCP server's OAuth metadata is missing a required authorization or token endpoint.

`dynamic_client_registration_not_supported`

The MCP server's OAuth metadata does not include a registration endpoint. The server may not support Dynamic Client Registration.

`dynamic_client_registration_rejected`

The MCP server rejected the Dynamic Client Registration request.

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

Org-level tokens are not allowed.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`failed_adding_collaborator`

Failed writing a collaborator record for this new app

`failed_creating_app`

Failed to create the app model

`failed_datastore_operation`

Failed while managing datastore infrastructure

`failed_generating_app_token`

App level token failed to generate.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

`feature_not_enabled`

The required feature is not enabled for this workspace.

`internal_error`

The server could not complete your operation(s) without encountering an error, likely due to a transient issue on our end. It's possible some aspect of the operation succeeded before the error was raised.

`invalid_app`

An app created from the provided manifest would not be valid.

`invalid_arg_name`

The method was passed an argument whose name falls outside the bounds of accepted or expected values. This includes very long names and names with non-alphanumeric characters other than `_`. If you get this error, it is typically an indication that you have made a *very* malformed API call.

`invalid_arguments`

Invalid API arguments provided.

`invalid_arguments`

The method was called with invalid arguments.

`invalid_array_arg`

The method was passed an array as an argument. Please only input valid strings.

`invalid_auth`

Some aspect of authentication cannot be validated. Either the provided token is invalid or the request originates from an IP address disallowed from making the request.

`invalid_form_data`

The method was called via a `POST` request with `Content-Type` `application/x-www-form-urlencoded` or `multipart/form-data`, but the form data was either missing or syntactically invalid.

`invalid_manager_app`

The calling app is not enrolled as a manager app. Ensure the app has `app_configurations:write` in its configured scopes and its home team has manager app support enabled.

`invalid_manifest`

The provided manifest file does not validate against schema. Consult the additional errors field to locate specific issues.

`invalid_post_type`

The method was called via a `POST` request, but the specified `Content-Type` was invalid. Valid types are: `application/json` `application/x-www-form-urlencoded` `multipart/form-data` `text/plain`.

`invalid_team_id`

The provided team ID is not valid

`managed_app_limit_reached`

The manager app has reached the maximum number of managed apps it can create.

`method_deprecated`

The method has been deprecated.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`not_in_team`

Cannot create an app in a team that user is not a member of

`oauth_metadata_not_found`

Could not discover OAuth metadata for the MCP server. Ensure the server exposes a.well-known/oauth-authorization-server or.well-known/openid-configuration endpoint.

`ratelimited`

Too many calls in succession to create endpoint during a short period of time.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`socket_mode_not_enabled`

Socket mode is not enabled in manifest.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.

`unknown_method`

Unknown method