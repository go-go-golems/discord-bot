---
Title: Investigation diary
Ticket: DISCORD-SLACK-001
Status: active
Topics:
    - discord-bot
    - architecture
    - api-design
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-09-10T18:20:32.32301754-04:00
WhatFor: ""
WhenToUse: ""
---

# Investigation diary

## Step 1: Repository investigation and Slack design

Mapped the existing Discord host, discovery, configuration, runtime owner, response handling, and xgoja integration before designing Slack support. The user clarified that a separate implementation following the same concepts is welcome, so the guide uses new Slack packages and a new binary in the existing module.

### Prompt Context

**User prompt (verbatim):**

> Create a new docmgr ticket to add slack support to discord-bot . Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.
>
> Then tell me what tokens and stuff you need to test / develop yourself

**Follow-up (verbatim):**

> save all useful resources to the ticket's sources folder, with defuddle for html/markdown

**Follow-up (verbatim):**

> yeah, use inspiration from discord-bot but feel free to make a separate implementation following the same concepts. You can also look at ~/code/wesen/go-go-golems/go-go-parc for more information about go-go-goja JS api design and also the work behind discord-bot

**Assistant interpretation:** Deliver design, sources, testing prerequisites, and reMarkable upload; no runtime implementation or live Slack messaging requested.

**Inferred user intent:** Give an intern enough architectural and operational context to build Slack support independently.

**Code revision:** `2ea219e6d0a37bb27916dc29a2ba88c1f03b6182`; initial git status clean.

### What I did

- Read README, configuration, lifecycle, JS dispatch, discovery, provider, and runtime ownership source.
- Created ticket and guide through docmgr.
- Researched official Slack APIs and the slack-go repository.
- Archived 14 remote references with defuddle and 5 vault notes with original and defuddle-extracted Markdown copies.
- Archived 19 current code/document inputs with original paths and hashes.
- Wrote the implementation guide with typed contracts, diagrams, pseudocode, phases, test cases, and credential setup.

### Why

Slack acknowledgment, command registration, identity, threading, and rendering differ from Discord. A separate implementation makes those contracts explicit without a compatibility layer. The vault reinforces Go ownership of lifecycle and domain state, and JS ownership of bot composition.

### What worked

`docmgr ticket create-ticket` and document creation succeeded. Public source extraction succeeded after enabling network access for the bounded source collection command. Local source archiving preserved original Markdown as well as normalized reading copies.

### What did not work

- `cat .docmgr.yaml`: `No such file or directory`. Repository configuration is `.ttmp.yaml`; used the installed docmgr configuration successfully.
- Initial `defuddle parse` calls failed with `Error loading content: getaddrinfo ENOTFOUND docs.slack.dev` and `getaddrinfo ENOTFOUND github.com`. Re-ran the collector with approved network access; all 14 sources succeeded.
- Large combined code/note reads exceeded output budgets; used targeted symbol searches and section reads for review-critical details.

### What I learned

Discovery and descriptors are coupled to Discord, not only transport. Current runtime config uses a host-field blacklist, so the Slack design uses a declared-field allowlist. Inspection executes scripts. Adjacent go-go-goja source and vault APIs may differ from the pinned dependency.

### What was tricky to design

Acknowledging before arbitrary JS work is necessary, but in-memory admission is not durable. Modal validation cannot use the same unconditional early ACK flow. Promise settlement must return to the VM owner while preserving invocation cancellation.

### What warrants a second pair of eyes

Review ingress dedupe/reservation ordering, owner-thread promise settlement, credential projection, and the selected SDK's retry/ACK behavior before implementation.

### What should be done in the future

Implement phases 0–4 before expanding to buttons, modals, xgoja integration, or persistence. Obtain the credentials and explicit test-message authorization described in guide section 12.

### Code review instructions

Start with the guide and source inventory. Follow the current file map and compare it to the archived revision. Use `docmgr doctor --ticket DISCORD-SLACK-001 --stale-after 30` for documentation hygiene. Runtime tests are future implementation gates, not claimed research results.

### Technical details

Source scripts are `scripts/01-collect-sources.py` and `scripts/02-archive-local.py`. Remote URL/time/hash records are in `sources/manifest.json`; vault and repository paths/hashes are in `sources/local-manifest.json`. No environment credentials were read and no Slack messages were sent.

## Step 2: Document validation and print layout

The guide contains approximately 6,500 words and the print edition has 21 pages. The first render exposed overlapping long code paths in table columns and page breaks inside diagrams. A generated print derivative converts tables to labeled records and keeps fenced examples together; the canonical guide retains its Markdown tables.

- `docmgr doctor --ticket DISCORD-SLACK-001 --stale-after 30`: all checks passed.
- Source SHA-256 verification, canonical guide relative links, and fenced-block balance passed.
- Initial print-wrapper attempt failed with `pandoc failed: Error producing PDF. ! Missing $ inserted.` at `event("app_`. Raw LaTeX environments consumed Markdown code; corrected by isolating wrappers in Pandoc raw-LaTeX fences.
- Final `remarquee upload bundle ... --pdf-only`: `OK: generated /tmp/slack-ticket-pdf-final/DISCORD-SLACK-001 Intern Implementation Guide.pdf`.
- Inspected rendered pages containing both architecture diagrams, file references, and test records. The overlaps and split diagrams were resolved.
- Saved the reviewed PDF to `artifacts/slack-intern-guide.pdf`. Print derivation script: `scripts/03-prepare-print.py`.
- Code snapshots use `.txt` extensions so the archive cannot create Go packages or nested modules. Added two adjacent go-go-goja runtime snapshots, each with revision provenance.
- No application code changed; Go build/test results are not claimed. Runtime validation belongs to the implementation tasks.

## Step 3: reMarkable delivery

Uploaded the completed guide using remarquee to `/ai/2026/09/10/DISCORD-SLACK-001`. The command returned `OK: uploaded DISCORD-SLACK-001 Intern Implementation Guide.pdf -> /ai/2026/09/10/DISCORD-SLACK-001`. A transient `remote tree has changed, refresh the file tree` warning preceded the successful upload. This is evidence of cloud delivery, not a physical-device synchronization check.

The final doctor pass initially reported `missing_numeric_prefix` for the generated print Markdown. Renamed that derivative with a `01-` prefix and updated the generator. The delivered content is unchanged. Implementation remains an unchecked future task.
