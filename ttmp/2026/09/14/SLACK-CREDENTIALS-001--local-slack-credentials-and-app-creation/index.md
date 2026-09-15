---
Title: Local Slack credentials and app creation
Ticket: SLACK-CREDENTIALS-001
Status: active
Topics:
    - discord-bot
    - api-design
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://pkg/slackcli/commands.go
      Note: Manifest and CLI integration
    - Path: repo://pkg/slackcli/create_app.go
      Note: Existing creation foundation
    - Path: repo://pkg/slackcli/create_app_test.go
      Note: Focused request tests
ExternalSources: []
Summary: Implement named local profiles, private credential files, explicit refresh, and profile-based app creation. Manual recovery is sufficient.
LastUpdated: 2026-09-14T21:31:18.814835481-04:00
WhatFor: ""
WhenToUse: ""
---

# Local Slack credentials and app creation

## Overview

Make it convenient to create Slack apps using saved local credentials and select the correct app/workspace. Use one metadata file and one private credentials file. This is a single-user development tool; manual token re-import and dashboard recovery are acceptable.

This ticket replaces the implementation scope proposed in DISCORD-SLACK-001 design document 03. The earlier document remains historical research, not the implementation specification. Existing manifest and create-app code should be reused.

Planning is complete; the credential manager is not implemented. The existing create-app changes are present in the working tree and need to be finished as the first implementation task.

## Key Links

- [Implementation plan](design-doc/01-pragmatic-implementation-plan.md)
- [Tasks](tasks.md)
- [Diary](reference/01-implementation-diary.md)
- [Archived API sources](sources/README.md)

## Status

Current status: **active**

## Topics

- discord-bot
- api-design

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.
