# xgoja discord-bot provider example

This example builds a generated `xdiscord` binary that imports the real `discord-bot` xgoja provider.

It mounts provider-owned bot commands under `bots` and runs bot scripts with the selected xgoja runtime profile. The profile includes:

- `discord` from `discord-bot/pkg/xgoja/provider`
- `ui` from `discord-bot/pkg/xgoja/provider`
- `fs` from `go-go-goja/pkg/xgoja/providers/host`
- `express` from `go-go-goja/pkg/xgoja/providers/http`

The sample bot is `fs-express-smoke`:

- `/ping` returns a static pong.
- `/read-config` reads `./bot-data/message.txt` through `require("fs")`.
- `/express-status` reports the xgoja-owned HTTP route status.
- `GET /` returns JSON status from the Express provider.
- `POST /say` sends a Discord message through `require("discord").channels.send(...)`.

## Smoke without Discord

```bash
make smoke
```

This validates the spec, builds the generated binary, verifies `eval` can require `discord` and `fs`, and verifies `bots list` / `bots help fs-express-smoke`. The static list/help smoke disables the HTTP server with `--http-enabled=false` so repeated smoke subcommands do not compete for the same listen port.

## Run against Discord and HTTP

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

Test the HTTP endpoint:

```bash
curl http://127.0.0.1:8787/
```

Send a Discord message through HTTP by replacing `<channel-id>` with a channel ID the bot can write to:

```bash
curl -X POST http://127.0.0.1:8787/say \
  -H 'Content-Type: application/json' \
  -d '{"channelId":"<channel-id>","content":"hello from xgoja express"}'
```

Stop with:

```bash
tmux kill-session -t xgoja-discord-bot
```
