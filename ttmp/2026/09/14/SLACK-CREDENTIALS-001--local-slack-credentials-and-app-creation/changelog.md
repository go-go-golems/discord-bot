# Changelog

## 2026-09-14

- Initial workspace created

## 2026-09-14

Created pragmatic local implementation scope: two files, named profiles, explicit refresh, manual recovery. Replaces DISCORD-SLACK-001 design 03.

## 2026-09-14

Implemented pragmatic two-file profile store, management import/status/refresh commands, and profile-based app creation. Focused CLI/storage tests pass.

### Related Files

- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/internal/slackconfig/store.go — Local config and secret persistence
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/pkg/slackcli/create_app.go — Profile-based app creation
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/pkg/slackcli/credentials.go — Credential commands

## 2026-09-14

Completed manual runtime-token import and full validation. All repository tests pass with loopback access; build, vet, focused tests, help, and docmgr doctor pass.

## 2026-09-14

Documented and implemented a pragmatic local bots install command using the Slack CLI observed apps.developerInstall method; it persists bot and Socket Mode tokens in the existing profile store and falls back to manual installation if the private method is unavailable.

### Related Files

- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/pkg/slackcli/install_app.go — developer-install request and credential persistence
- /home/manuel/workspaces/2026-09-10/add-slack-support/discord-bot/ttmp/2026/09/14/SLACK-CREDENTIALS-001--local-slack-credentials-and-app-creation/design-doc/02-local-app-installation-with-slack-cli-developerinstall.md — contract and limitations
