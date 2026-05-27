.DEFAULT_GOAL := smoke

EXAMPLE_DIR := $(CURDIR)
WORKSPACE_ROOT := $(abspath ../../../..)
XGOJA_REPO := $(WORKSPACE_ROOT)/go-go-goja
BIN := $(EXAMPLE_DIR)/dist/xdiscord
XGOJA := cd $(XGOJA_REPO) && GOWORK=off go run ./cmd/xgoja
TMUX_SESSION := xgoja-discord-bot

.PHONY: smoke doctor build eval bots-list bots-help tmux-run tmux-stop clean

smoke: doctor build eval bots-list bots-help

doctor:
	$(XGOJA) doctor -f $(EXAMPLE_DIR)/xgoja.yaml

build:
	$(XGOJA) build -f $(EXAMPLE_DIR)/xgoja.yaml --output $(BIN) --xgoja-replace $(XGOJA_REPO) --keep-work

eval:
	cd $(EXAMPLE_DIR) && $(BIN) eval 'require("discord"); require("fs").readFileSync("./bot-data/message.txt", "utf8").trim()'

bots-list:
	cd $(EXAMPLE_DIR) && $(BIN) bots list --http-enabled=false --output json

bots-help:
	cd $(EXAMPLE_DIR) && $(BIN) bots help fs-express-smoke --http-enabled=false --output json

tmux-run: build
	@test -n "$(DISCORD_BOT_TOKEN)" || (echo "DISCORD_BOT_TOKEN is required" && exit 1)
	@test -n "$(DISCORD_APPLICATION_ID)" || (echo "DISCORD_APPLICATION_ID is required" && exit 1)
	@test -n "$(DISCORD_GUILD_ID)" || (echo "DISCORD_GUILD_ID is required" && exit 1)
	@tmux has-session -t $(TMUX_SESSION) 2>/dev/null && tmux kill-session -t $(TMUX_SESSION) || true
	tmux new-session -d -s $(TMUX_SESSION) -- bash -lc 'cd "$(EXAMPLE_DIR)" && source "$(WORKSPACE_ROOT)/.envrc" && exec "$(BIN)" bots fs-express-smoke run --sync-on-start --http-listen 127.0.0.1:8787'
	@echo "tmux session $(TMUX_SESSION) is running"
	@echo "attach: tmux attach -t $(TMUX_SESSION)"
	@echo "http: curl http://127.0.0.1:8787/"
	@echo "say: curl -X POST http://127.0.0.1:8787/say -H 'Content-Type: application/json' -d '{\"channelId\":\"<channel-id>\",\"content\":\"hello from xgoja express\"}'"

tmux-stop:
	tmux kill-session -t $(TMUX_SESSION)

clean:
	rm -rf $(EXAMPLE_DIR)/dist
