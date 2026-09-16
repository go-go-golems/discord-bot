Source: https://docs.slack.dev/reference/methods/users.info/
Retrieved: 2026-09-15

## Facts

## Arguments

### Optional arguments

**`include_locale`** `boolean` Optional

Set this to `true` to receive the locale for this user. Defaults to `false`

**`user`** `string` Optional

User to get info on

## Usage info

This method returns information about a member of a workspace.

---

## Response

#### Typical success response

```json
{
  "ok": true,
  "user": {
    "id": "W012A3CDE",
    "team_id": "T012AB3C4",
    "name": "spengler",
    "deleted": false,
    "color": "9f69e7",
    "real_name": "Egon Spengler",
    "tz": "America/Los_Angeles",
    "tz_label": "Pacific Daylight Time",
    "tz_offset": -25200,
    "profile": {
      "avatar_hash": "ge3b51ca72de",
      "status_text": "Print is dead",
      "status_emoji": ":books:",
      "real_name": "Egon Spengler",
      "display_name": "spengler",
      "real_name_normalized": "Egon Spengler",
      "display_name_normalized": "spengler",
      "email": "spengler@ghostbusters.example.com",
      "image_original": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
      "image_24": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
      "image_32": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
      "image_48": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
      "image_72": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
      "image_192": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
      "image_512": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
      "team": "T012AB3C4"
    },
    "is_admin": true,
    "is_owner": false,
    "is_primary_owner": false,
    "is_restricted": false,
    "is_ultra_restricted": false,
    "is_bot": false,
    "updated": 1502138686,
    "is_app_user": false,
    "has_2fa": false
  }
}
```

#### Typical error response

```json
{
  "ok": false,
  "error": "user_not_found"
}
```

Returns a [user object](/reference/objects/user-object).

The profile hash contains as much information as the user has supplied in the default profile fields: `first_name`, `last_name`, `real_name`, `display_name`, `skype`, and the `image_*` fields. Only the `image_*` fields are guaranteed to be included. Data that has not been supplied may not be present at all, may be null or may contain the empty string ("").

The composition of user objects can vary greatly depending on the API being used, or the context of each Slack workspace. One example where the profile information varies is getting data for a user connected through Slack Connect; an example of that type of user object can be found in the [External members section](/apis/slack-connect/#members) of the Slack Connect page.

A user's custom profile fields may be discovered using [`users.profile.get`](/reference/methods/users.profile.get).

### Email addresses

Apps created after January 4th, 2017 must request *both* the `users:read` and `users:read.email` OAuth permission scopes simultaneously when using the [OAuth app installation flow](/authentication/installing-with-oauth) to enable access to the `email` field of user objects returned by this method.

### Presence

A user's presence is found using [`users.getPresence`](/reference/methods/users.getPresence).

### Shared channels

When looking up a user belonging to a foreign workspace party to a shared channel, you'll find an `is_stranger` boolean attribute.

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

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`request_timeout`

The method was called via a `POST` request, but the `POST` data was either missing or truncated.

`team_access_not_granted`

The token used is not granted the specific workspace access required to complete this request.

`token_expired`

Authentication token has expired

`token_revoked`

Authentication token is for a deleted user or workspace or the app has been removed when using a `user` token.

`too_many_users`

Too many users.

`two_factor_setup_required`

Two factor setup is required.

`user_not_found`

Value passed for `user` was invalid.

`user_not_visible`

The requested user is not visible to the calling user