Source: https://docs.slack.dev/reference/methods/admin.users.assign/
Retrieved: 2026-09-15

## Facts

## Arguments

### Required arguments

**`team_id`** Required

The ID (`T1234`) of the workspace.

**`user_id`** Required

The ID of the user to add to the workspace.

### Optional arguments

**`is_restricted`** `boolean` Optional

True if user should be added to the workspace as a guest.

*Example:* `true`

**`is_ultra_restricted`** `boolean` Optional

True if user should be added to the workspace as a single-channel guest.

*Example:* `true`

**`channel_ids`** `string` Optional

Comma separated values of channel IDs to add user in the new workspace.

*Example:* `C123,C3456`

## Usage info

This Admin API assigns a user to a workspace. If the user has never been a member of the workspace, they will be added. If they've previously been removed or left the workspace, the user will be reinstated as a member. If this method is used on a workspace-level user who is deactivated on the org level, the user will be reactivated on both the workspace and the org level.

---

## Response

#### Typical success response

```json
{
  "ok": true
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

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

`feature_not_enabled`

The Admin APIs feature is not enabled for this team

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

`invalid_role_for_user`

The requested status is incompatible with the user's role on another workspace.

`invited_user_not_created`

The invited user could not be created.

`invited_user_not_reactivated`

The invited user could not be reactivated.

`invitor_cannot_see_channel`

The invitor is not part of one or many channels that the user was requested to be added to

`method_deprecated`

The method has been deprecated.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_an_admin`

This method wasn't called by an admin.

`not_authed`

No authentication token provided.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`team_not_found`

\`team\_id was not found.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`two_factor_setup_required`

Two factor setup is required.

`user_already_team_member`

The given user is already active on the given team.

`user_cannot_be_added_to_workspace`

The user can't be added to this workspace because of their account type or the workspace's configuration.

`user_is_bot`

The given user is a bot from an app. This api only works with humans.

`user_not_found`

`user_id` was not found.