# xgoja discord-bot provider example

This example builds a generated `xdiscord` binary that imports the real `discord-bot` xgoja provider.

It mounts provider-owned bot commands under `bots` and runs bot scripts with the selected xgoja runtime profile. The profile includes:

- `discord` from `discord-bot/pkg/xgoja/provider`
- `ui` from `discord-bot/pkg/xgoja/provider`
- `fs` from `go-go-goja/pkg/xgoja/providers/host`

The sample bot is `fs-express-smoke`:

- `/ping` returns a static pong.
- `/read-config` reads `./bot-data/message.txt` through `require("fs")`.
- `/express-status` reports that express host lifecycle is still planned.

## Smoke without Discord

```bash
make smoke
```

This validates the spec, builds the generated binary, verifies `eval` can require `discord` and `fs`, and verifies `bots list` / `bots help fs-express-smoke`.

## Run against Discord

Only do this when `DISCORD_BOT_TOKEN`, `DISCORD_APPLICATION_ID`, and `DISCORD_GUILD_ID` are set.

```bash
make tmux-run
```

The tmux session is named `xgoja-discord-bot`. Attach with:

```bash
tmux attach -t xgoja-discord-bot
```

Then test these slash commands in the configured guild:

- `/ping`
- `/read-config`
- `/express-status`

Stop with:

```bash
tmux kill-session -t xgoja-discord-bot
```
