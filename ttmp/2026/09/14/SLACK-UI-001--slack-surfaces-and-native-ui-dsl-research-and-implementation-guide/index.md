---
Title: Slack surfaces and native UI DSL research and implementation guide
Ticket: SLACK-UI-001
Status: active
Topics:
    - slack
    - architecture
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-14T22:50:30.006245744-04:00
WhatFor: ""
WhenToUse: ""
---

# Slack surfaces and native UI DSL

## Overview

Research and design for extending the existing local Slack host with a small native UI DSL inspired by Discord. The recommended core is rich messages, action routing, and one complete modal editor. The broader Slack surface catalog is reference material and optional follow-up work. No runtime UI functionality has been implemented in this ticket.

## Read first

- [Intern design and implementation guide](design-doc/01-slack-surfaces-and-ui-dsl-intern-guide.md)
- [Detailed investigation diary](reference/01-investigation-diary.md)
- [Source catalog](sources/README.md)
- [Implementation and research tasks](tasks.md)
- [SDK probe](scripts/02-run-sdk-probe.sh) and [observed output](sources/sdk-probe-output.txt)

## Findings

The current Slack API is text-only and drops interactive envelopes. Slack modal submissions need semantic acknowledgment before the interaction deadline. The pinned SDK supports the core slice, but loses fields when decoding newer unknown blocks. The design keeps existing local state and execution machinery with a documented busy/retry limitation.

## Evidence and scope

Baseline discord-bot revision: 5c700b2. Pinned Slack SDK: v0.17.3. Official Slack references were captured on 2026-09-14 local time; some catalogs require browser rendering. Source provenance and extraction recovery are recorded in sources/catalog-complete.json. The live ping process and private credentials were not changed.

## Delivery

The guide is being rendered for reMarkable. The final delivery receipt will be recorded in the diary and artifacts. Implementation tasks remain open.
