---
Title: Implementation diary
Ticket: SLACK-CREDENTIALS-001
Status: active
Topics:
    - discord-bot
    - api-design
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-14T21:31:19.295803347-04:00
WhatFor: ""
WhenToUse: ""
---

# Implementation diary

## Step 1: Reset scope for a local tool

Created a separate implementation ticket after the user rejected the earlier reliability-heavy design. The replacement uses two files, named maps, explicit refresh, and manual installation/recovery. The credential manager is not implemented in this planning step.

### Prompt Context

**User prompt (verbatim):** "ok, start again, make a pragmatic implementation ticket."

**Follow-up (verbatim):** "mark the previous one deprecated"

**Assistant interpretation:** Replace the implementation scope with a small local workflow and retire the old credentials design.

**Inferred user intent:** Make progress on a useful local bot without production reliability infrastructure.

### What I did

- Created a plan, ordered tasks, command forms, and bounded acceptance tests.
- Linked the existing defuddle source archive.
- Deprecated the old credentials design while keeping the broader Slack-support ticket active.

### Why

Detailed explanation does not imply automatic recovery from every rare failure. Manual re-import is acceptable here.

### What worked

Existing manifest/create-app source provides a concrete starting point. A single secrets document keeps replacement tokens together without generations or locks.

### What didn't work

The preceding design exceeded the user's needs. Also, the first combined documentation patch failed with `apply_patch verification failed: Failed to find expected lines` because a diary template comment differed; reread the actual template and applied the corrected patch. No partial changes from that patch were present.

### What I learned

Documentation depth should explain a small implementation, not expand its reliability requirements.

### What was tricky to build

App secrets and workspace installation tokens still need distinct map entries. Shared management identities require references, but not a generalized storage framework.

### What warrants a second pair of eyes

Ensure implementation does not add automatic refresh/retries implicitly. The current create-app source is an uncommitted foundation; profile commands remain proposed.

### What should be done in the future

Implement the ordered tasks. Keep production transport separate.

### Code review instructions

Read the plan/tasks, then `pkg/slackcli/create_app.go` and `commands.go`. Validate ticket metadata and relative links before committing.

### Technical details

Storage is `config.yaml` plus `credentials.json`, with private permissions, temporary-file rename, and one writer at a time. This planning milestone makes no API requests and reads no credentials.
