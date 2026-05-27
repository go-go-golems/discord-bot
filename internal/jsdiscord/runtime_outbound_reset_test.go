package jsdiscord

import (
	"context"
	"strings"
	"testing"
)

func TestDiscordDispatchClearsOutboundOpsWhenRequestOmitsDiscord(t *testing.T) {
	handle := loadTestBot(t, writeBotScript(t, `
const discord = require("discord")
const { defineBot } = discord
module.exports = defineBot(({ command }) => {
  command("channels", async () => {
    const channels = await discord.channels.list("guild-1")
    return { content: channels.map((channel) => channel.name).join(",") }
  })
})
`))

	first, err := handle.DispatchCommand(context.Background(), DispatchRequest{
		Name: "channels",
		Discord: &DiscordOps{
			ChannelList: func(_ context.Context, guildID string) ([]map[string]any, error) {
				if guildID != "guild-1" {
					t.Fatalf("guildID = %q", guildID)
				}
				return []map[string]any{{"id": "chan-1", "name": "general"}}, nil
			},
		},
	})
	if err != nil {
		t.Fatalf("first dispatch: %v", err)
	}
	if got := responseContent(t, first); got != "general" {
		t.Fatalf("first content = %q, want general", got)
	}

	second, err := handle.DispatchCommand(context.Background(), DispatchRequest{Name: "channels"})
	if err == nil {
		t.Fatalf("second dispatch unexpectedly reused stale outbound ops and returned %#v", second)
	}
	if !strings.Contains(err.Error(), "discord channel list API is not ready") {
		t.Fatalf("second dispatch error = %v, want cleared outbound ops error", err)
	}
}

func responseContent(t *testing.T, result any) string {
	t.Helper()
	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result = %#v, want map", result)
	}
	content, _ := m["content"].(string)
	return content
}
