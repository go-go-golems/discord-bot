---
Title: Port inventory and implementation plan
Ticket: SLACK-PORT-001
Status: active
Topics:
    - slack
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-15T20:06:08.612721797-04:00
WhatFor: ""
WhenToUse: ""
---

# Port inventory and implementation plan

## Executive Summary

<!-- Provide a high-level overview of the design proposal -->

## Problem Statement

<!-- Describe the problem this design addresses -->

## Proposed Solution

<!-- Describe the proposed solution in detail -->

## Design Decisions

<!-- Document key design decisions and rationale -->

## Alternatives Considered

<!-- List alternative approaches that were considered and why they were rejected -->

## Implementation Plan

<!-- Outline the steps to implement this design -->

## Open Questions

<!-- List any unresolved questions or concerns -->

## References

<!-- Link to related documents, RFCs, or external resources -->


## Scope and current inventory

Port 13 examples into examples/slack-bots using native Slack payloads and handlers.
No Discord compatibility runtime is planned. The user selected every feature with a Slack equivalent. The initial command/UI
ports are a checkpoint, not acceptance of reduced functionality. Unsupported
Discord-specific operations must be documented with evidence; missing Slack
services must be implemented instead of silently dropping their workflows.

| Bot | Dependencies and port considerations |
| --- | --- |
| announcements | Preview text and a Block Kit header. |
| ping | Existing baseline; rich Discord API demo needs expanded Slack examples. |
| hater | Text, buttons, modal; use app mentions for explicit interaction. |
| interaction-types | Slash text parsing; user/message context commands need Slack-specific alternatives. |
| unified-demo | Declared configuration and status; Slack runner has no jsverbs discovery integration. |
| poker | Reuse game math; port state scoping and renderers to native Slack. |
| support | Ticket preview/status plus channel-thread operations. |
| custom-kb | Persistent SQLite link storage, search and entry form. |
| knowledge-base | Persistent entries, capture, source links, review and search workflows. |
| show-space | Shows, permissions, posting/pinning, persistent state and seed data. |
| archive-helper | Paginated history/thread retrieval, rendering and file delivery. |
| moderation | Workspace/channel/member read APIs plus permission-aware native actions; Discord roles/bans/timeouts need explicit mapping. |
| ui-showcase | Current basic modal works; broaden to search, pagination, selections, cards and review. |

## Implementation order

1. Commands and existing interactive primitives; executable fixtures.
2. Reusable Slack services needed by operational bots, with scopes declared per bot.
3. Persistent application stores and workflows.
4. Expanded UI examples and a parity matrix documenting platform differences.
5. Offline tests, CLI discovery/manifest validation, documentation and commits.

Porting source is authorized; changing the currently installed Slack app or running
administrative actions in the workspace is not part of this implementation pass.
