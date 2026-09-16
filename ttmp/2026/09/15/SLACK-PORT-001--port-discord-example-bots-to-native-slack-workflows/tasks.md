# Tasks

## Scope

Every feature with a Slack equivalent across all 13 Discord examples. Initial
scripts are not evidence of complete parity. Follow design-doc/01 and
reference/02-source-handler-acceptance-matrix.md.

## Implementation

- [x] Inventory source examples, archive API references, and write the intern guide.
- [x] Upload the intern guide to reMarkable.
- [x] T1: Complete feature mappings and acceptance tests for all source workflows.
- [x] T2: Fix optional inputs and add native UI controls needed by the ports.
- [x] T3: Complete shortcuts, options, modal-result and reply-update interactions.
- [x] T4: Implement operational Slack services and explicit capability scopes.
- [x] T5: Add lifecycle-owned persistent database capability and reopen tests.
- [x] T6: Finish announcements, ping, hater, interaction-types, unified-demo and poker.
- [x] T7: Port custom-kb and knowledge-base with persistence and review/search/export.
- [x] T8: Port support, show-space, archive-helper and moderation with authorization.
- [x] T9: Complete UI showcase, documentation and local parity qualification.

## Qualification boundary

Implementation tasks above are complete against the local source inventory and
fixture suite. The matrix contains 167 literal Discord registrations and 181
native registration expectations, including computed registrations. Registration
coverage proves handlers exist; it does not prove every possible input or Slack
plan/token combination. Representative workflow tests cover each bot family.

Live follow-up: install the desired bot with its generated manifest, reinstall
when scopes change, configure authorized actors, and exercise commands, shortcuts,
modals and event subscriptions in the intended workspace. Conditional user-token
and Enterprise operations require separately granted permissions. No live
moderation or workspace mutation was performed during this porting work.

The full Go tests, build and vet passed. Repository-wide lint and vulnerability
checks retain existing findings; see the Step 7 diary and validation artifacts.
- [x] Harden archive pagination against cursor cycles and verify no partial upload on API failures <!-- t:dj85 -->
- [x] Address PR 19 review: failed credentials-file cleanup and documented profiles invocation <!-- t:qkhd -->
- [ ] Align lefthook and CI checks with pinned Glazed tooling; fix lint/security jobs <!-- t:sl8h -->
