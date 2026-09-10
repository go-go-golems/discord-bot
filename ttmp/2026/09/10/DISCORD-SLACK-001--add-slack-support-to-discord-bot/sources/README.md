# Archived research sources

HTML pages extracted with `defuddle parse URL --md -o FILE`. Retrieval metadata and SHA-256 hashes are in `manifest.json`. These are research snapshots; linked upstream pages remain authoritative.

| Snapshot | Original URL | Bytes |
|---|---|---|
| [01-socket-mode](01-socket-mode.md) | https://docs.slack.dev/apis/events-api/using-socket-mode/ | 16788 |
| [02-connections-scope](02-connections-scope.md) | https://docs.slack.dev/reference/scopes/connections.write/ | 537 |
| [03-slash-commands](03-slash-commands.md) | https://docs.slack.dev/interactivity/implementing-slash-commands/ | 38647 |
| [04-post-message](04-post-message.md) | https://docs.slack.dev/reference/methods/chat.postMessage/ | 22048 |
| [05-interactions](05-interactions.md) | https://docs.slack.dev/interactivity/handling-user-interaction/ | 14165 |
| [06-rate-limits](06-rate-limits.md) | https://docs.slack.dev/apis/web-api/rate-limits/ | 9810 |
| [07-app-manifest](07-app-manifest.md) | https://docs.slack.dev/reference/app-manifest/ | 32816 |
| [08-app-mention](08-app-mention.md) | https://docs.slack.dev/reference/events/app_mention/ | 4535 |
| [09-request-signing](09-request-signing.md) | https://docs.slack.dev/authentication/verifying-requests-from-slack/ | 11084 |
| [10-slack-go](10-slack-go.md) | https://github.com/slack-go/slack | 27493 |
| [11-auth-test](11-auth-test.md) | https://docs.slack.dev/reference/methods/auth.test/ | 4650 |
| [12-chat-update](12-chat-update.md) | https://docs.slack.dev/reference/methods/chat.update/ | 12144 |
| [13-views-open](13-views-open.md) | https://docs.slack.dev/reference/methods/views.open/ | 5972 |
| [14-message-im](14-message-im.md) | https://docs.slack.dev/reference/events/message.im/ | 1318 |

## Local evidence

Vault Markdown is preserved byte-for-byte as `.original.md`; a Pandoc HTML conversion was passed through defuddle for the companion extracted `.md`. Originals preserve code fences and Obsidian metadata. `local-manifest.json` records paths and hashes. Code snapshots are verbatim, not transformed.

- [20-dsl-design](20-dsl-design.md) — [original](20-dsl-design.original.md)
- [21-discord-framework](21-discord-framework.md) — [original](21-discord-framework.original.md)
- [22-discord-ui-dsl](22-discord-ui-dsl.md) — [original](22-discord-ui-dsl.original.md)
- [23-xgoja-discord](23-xgoja-discord.md) — [original](23-xgoja-discord.original.md)
- [24-context-management](24-context-management.md) — [original](24-context-management.original.md)

Discord code snapshots: `code/`, revision `2ea219e6d0a37bb27916dc29a2ba88c1f03b6182`.

Source-code snapshots use `.txt` suffixes to avoid entering Go package discovery or creating nested Go modules. Adjacent go-go-goja runtime/bridge snapshots include their own revision in the local manifest and are not asserted to match the pinned dependency.
