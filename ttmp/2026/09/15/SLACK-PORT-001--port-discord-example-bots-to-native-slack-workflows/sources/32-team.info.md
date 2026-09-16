Source: https://docs.slack.dev/reference/methods/team.info/
Retrieved: 2026-09-15

## Facts

## Arguments

### Optional arguments

**`domain`** `string` Optional

Query by domain instead of team (only when team is null). This only works for domains in the same enterprise as the querying team token. This also expects the domain to belong to a team and not the enterprise itself. This is the value set up for the 'Joining This Workspace' workspace setting. If it contains more than one domain, the field will contain multiple comma-separated domain values. If no domain is set, the field is empty.

**`team`** `string` Optional

Team to get info about; if omitted, will return information about the current team.

## Usage info

This method provides information about your team.

---

## Response

#### Typical success response

```json
{
  "ok": true,
  "team": {
    "id": "T12345",
    "name": "My Team",
    "domain": "example",
    "email_domain": "example.com",
    "icon": {
      "image_34": "https://...",
      "image_44": "https://...",
      "image_68": "https://...",
      "image_88": "https://...",
      "image_102": "https://...",
      "image_132": "https://...",
      "image_default": true
    },
    "enterprise_id": "E1234A12AB",
    "enterprise_name": "Umbrella Corporation"
  }
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "invalid_auth"
}
```

## Team Icon

If a team has not yet set a custom icon, the value of `team.icon.image_default` will be `true`.

If the team belongs to an [Enterprise organization](/enterprise), the `enterprise_id` and `enterprise_name` fields will indicate the owning enterprise organization.

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

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`enterprise_not_found`

The `enterprise` was not found.

`fail_to_get_teams_for_restricted_user`

Failed to get teams for restricted user.

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

`invalid_form_data`

The method was called via a `POST` request with `Content-Type` `application/x-www-form-urlencoded` or `multipart/form-data`, but the form data was either missing or syntactically invalid.

`invalid_post_type`

The method was called via a `POST` request, but the specified `Content-Type` was invalid. Valid types are: `application/json` `application/x-www-form-urlencoded` `multipart/form-data` `text/plain`.

`method_deprecated`

The method has been deprecated.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`missing_user`

The `user` was not found.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`org_not_found`

The `org` was not found.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`team_not_found`

The `team` was not found.

`team_not_on_enterprise`

Cannot query team by domain because team is not on an enterprise.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.

`user_not_found`

The `user` was not found.