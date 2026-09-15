---
Title: "22 reference methods views.publish"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/reference/methods/views.publish/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

## Facts

## Arguments

### Required arguments

**`user_id`** `string` Required

`id` of the user you want publish a view to.

*Example:* `U0BPQUNTA`

**`view`** Required

A [view payload](https://docs.slack.dev/reference/views). This must be a JSON-encoded string.

### Optional arguments

**`hash`** `string` Optional

A string that represents view state to protect against possible race conditions.

*Example:* `156772938.1827394`

**`interactivity_pointer`** `string` Optional

## Usage info

Create or update the view that comprises [an app's Home tab](https://docs.slack.dev/surfaces/app-home) for a specific user.

---

## Response

#### Typical success response includes the published view payload.

```json
{
  "ok": true,
  "view": {
    "id": "VMHU10V25",
    "team_id": "T8N4K1JN",
    "type": "home",
    "close": null,
    "submit": null,
    "blocks": [
      {
        "type": "section",
        "block_id": "2WGp9",
        "text": {
          "type": "mrkdwn",
          "text": "A simple section with some sample sentence.",
          "verbatim": false
        }
      }
    ],
    "private_metadata": "Shh it is a secret",
    "callback_id": "identify_your_home_tab",
    "state": {
      "values": {}
    },
    "hash": "156772938.1827394",
    "clear_on_close": false,
    "notify_on_close": false,
    "root_view_id": "VMHU10V25",
    "previous_view_id": null,
    "app_id": "AA4928AQ",
    "external_id": "",
    "bot_id": "BA13894H"
  }
}
```

#### Typical error response, before getting to any possible validation errors.

```json
{
  "ok": false,
  "error": "invalid_arguments",
  "response_metadata": {
    "messages": [
      "invalid \`user_id\`"
    ]
  }
}
```

Assuming your view object was properly formatted, valid, and the `user_id` was viable, you will receive a success response.

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

`duplicate_external_id`

Error returned when the given `external_id` has already be used.

`ekm_access_denied`

Administrators have suspended the ability to post a message.

`enterprise_is_restricted`

The method cannot be called from an Enterprise.

`fatal_error`

The server could not complete your operation(s) without encountering a catastrophic error. It's possible some aspect of the operation succeeded before the error was raised.

`hash_conflict`

Error returned when the provided `hash` doesn't match the current stored value.

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

`missing_profile_id`

A profile id was not provided when trying to publish a view of type profile.

`missing_scope`

The token used is not granted the specific scope permissions required to complete this request.

`no_permission`

The workspace token used in this request does not have the permissions necessary to complete the request. Make sure your app is a member of the conversation it's attempting to post a message to.

`not_allowed_token_type`

The type of token your app used when requesting this method is not allowed.

`not_allowed_token_type`

The token type used in this request is not allowed.

`not_authed`

No authentication token provided.

`not_enabled`

Error returned if a `home` view is published but the Home tab isn't enabled for the app.

`not_implemented`

The profile view experiment is not enabled for this user.

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

`view_too_large`

Error returned if the provided view is greater than 250kb.