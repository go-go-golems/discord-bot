## Facts

## Arguments

### Optional arguments

**`client_id`** `string` Optional

Issued when you created your application. If possible, avoid sending `client_id` and `client_secret` as parameters in your request and instead supply the Client ID and Client Secret using the HTTP Basic authentication scheme.

*Example:* `2141029472.691202649728`

**`client_secret`** `string` Optional

Issued when you created your application. If possible, avoid sending `client_id` and `client_secret` as parameters in your request and instead supply the Client ID and Client Secret using the HTTP Basic authentication scheme.

*Example:* `e1b9e11dfcd19c1982d5de12921e17e8c`

**`code`** `string` Optional

The `code` param returned via the OAuth callback.

*Example:* `4724469134.4644010092847.232b4e6d82c333b475fc30f5f5a341d294feb1a94392c2fd791f7ab7731a443d1a`

**`code_verifier`** `string` Optional

The code\_verifier param used to generate the code\_challenge originally. Used for PKCE.

*Example:* `secret12345`

**`redirect_uri`** `string` Optional

This must match the originally submitted URI (if one was sent).

*Example:* `http://example.com`

**`grant_type`** `string` Optional

The `grant_type` param as described in the OAuth spec.

*Example:* `authorization_code`

**`refresh_token`** `string` Optional

The `refresh_token` param as described in the OAuth spec.

*Example:* `xoxe-1-abcdefg`

**`assertion`** `string` Optional

Identity assertion JWT authorization grant

*Example:* `eyJhbGciOiJSUzI1NiIsInR…`

## Usage info

This is a step in the [V2 OAuth authentication flow](/authentication/installing-with-oauth). Check out our [guide to new Slack apps](/app-management/quickstart-app-settings) for more information.

If at all possible, avoid sending `client_id` and `client_secret` as parameters in your request. We strongly recommend supplying the Client ID and Client Secret using the HTTP Basic authentication scheme, as discussed [in RFC 6749](https://tools.ietf.org/html/rfc6749#section-2.3.1).

A potential gotcha: while `redirect_uri` is optional, it is *required* if your app passed it as a parameter to `oauth/authorization` in the first step of the OAuth flow.

Look for the `is_enterprise_install` boolean if your app is an [org-ready app](/enterprise/organization-ready-apps) and you need to determine whether it has been installed on an entire organization or a single workspace.

If you're developing a [GovSlack](/govslack) app for use by public sector customers, make your API calls to the `slack-gov.com` domain instead of the `slack.com` domain.

### Rate limits

The rate limit for this method is 600 per minute.

### Refresh tokens

When you use `grant_type=refresh_token` and pass your `refresh_token` as an argument, this method *refreshes* an access token, whether bot or user. Read our [guide to token rotation](/authentication/using-token-rotation) for more information.

---

## Response

#### Successful token request with scopes for both a bot user and a user token.

```json
{
  "ok": true,
  "access_token": "SLACK_REDACTED_EXAMPLE",
  "token_type": "bot",
  "scope": "commands,incoming-webhook",
  "bot_user_id": "U0KRQLJ9H",
  "app_id": "A0KRD7HC3",
  "team": {
    "name": "Slack Softball Team",
    "id": "T9TK3CUKW"
  },
  "enterprise": {
    "name": "slack-sports",
    "id": "E12345678"
  },
  "authed_user": {
    "id": "U1234",
    "scope": "chat:write",
    "access_token": "xoxp-1234",
    "token_type": "user"
  }
}
```

#### Successful token request with scopes for both a bot user and a user token, and token rotation enabled.

```json
{
  "ok": true,
  "access_token": "xoxe.xoxb-1-..",
  "token_type": "bot",
  "scope": "commands,incoming-webhook",
  "bot_user_id": "U0KRQLJ9H",
  "app_id": "A0KRD7HC3",
  "expires_in": 43200,
  "refresh_token": "xoxe-1-...",
  "team": {
    "name": "Slack Softball Team",
    "id": "T9TK3CUKW"
  },
  "enterprise": {
    "name": "slack-sports",
    "id": "E12345678"
  },
  "authed_user": {
    "id": "U1234",
    "scope": "chat:write",
    "access_token": "xoxe.xoxp-1234",
    "expires_in": 43200,
    "refresh_token": "xoxe-1-...",
    "token_type": "user"
  }
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "invalid_client_id"
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

`bad_client_secret`

Value passed for `client_secret` was invalid.

`bad_redirect_uri`

Value passed for `redirect_uri` did not match the `redirect_uri` in the original request.

`cannot_install_an_org_installed_app`

Returned when the the org-installed app cannot be installed on a workspace.

`deprecated_endpoint`

The endpoint has been deprecated.

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

`internal_error`

The server could not complete your operation(s) without encountering an error, likely due to a transient issue on our end. It's possible some aspect of the operation succeeded before the error was raised.

`invalid_arg_name`

The method was passed an argument whose name falls outside the bounds of accepted or expected values. This includes very long names and names with non-alphanumeric characters other than `_`. If you get this error, it is typically an indication that you have made a *very* malformed API call.

`invalid_arguments`

The method was called with invalid arguments.

`invalid_array_arg`

The method was passed an array as an argument. Please only input valid strings.

`invalid_auth`

Some aspect of authentication cannot be validated. Either the provided token is invalid or the request originates from an IP address disallowed from making the request.

`invalid_client_id`

Value passed for `client_id` was invalid.

`invalid_code`

Value passed for `code` was invalid.

`invalid_code_verifier`

The code\_verifier is invalid.

`invalid_form_data`

The method was called via a `POST` request with `Content-Type` `application/x-www-form-urlencoded` or `multipart/form-data`, but the form data was either missing or syntactically invalid.

`invalid_grant_type`

Value passed for `grant_type` was invalid.

`invalid_post_type`

The method was called via a `POST` request, but the specified `Content-Type` was invalid. Valid types are: `application/json` `application/x-www-form-urlencoded` `multipart/form-data` `text/plain`.

`invalid_refresh_token`

The given refresh token is invalid.

`method_deprecated`

The method has been deprecated.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`no_scopes`

Missing `scope` in the request.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`oauth_authorization_url_mismatch`

The OAuth flow was initiated on an incorrect version of the authorization url. The flow must be initiated via /oauth/v2/authorize.

`pkce_not_allowed`

The app is not allowed to use the PKCE flow.

`preview_feature_not_available`

Returned when the API method is not yet available on the team in context.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.

`user_email_unverified`

The users email is unverified.