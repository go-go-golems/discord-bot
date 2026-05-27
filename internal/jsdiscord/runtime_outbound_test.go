package jsdiscord

import (
	"context"
	"testing"
)

func TestDiscordTopLevelChannelsListUsesOutboundOps(t *testing.T) {
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

	result, err := handle.DispatchCommand(context.Background(), DispatchRequest{
		Name: "channels",
		Discord: &DiscordOps{
			ChannelList: func(_ context.Context, guildID string) ([]map[string]any, error) {
				if guildID != "guild-1" {
					t.Fatalf("guildID = %q", guildID)
				}
				return []map[string]any{{"id": "chan-1", "name": "general"}, {"id": "chan-2", "name": "bots"}}, nil
			},
		},
	})
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	m, ok := result.(map[string]any)
	if !ok || m["content"] != "general,bots" {
		t.Fatalf("result = %#v", result)
	}
}

func TestDiscordTopLevelChannelsSendUsesOutboundOps(t *testing.T) {
	handle := loadTestBot(t, writeBotScript(t, `
const discord = require("discord")
const { defineBot } = discord
module.exports = defineBot(({ command }) => {
  command("say", async (ctx) => {
    await discord.channels.send(ctx.options.channelId, { content: ctx.options.message })
    return { content: "sent" }
  })
})
`))

	var gotChannel string
	var gotPayload any
	result, err := handle.DispatchCommand(context.Background(), DispatchRequest{
		Name: "say",
		Args: map[string]any{
			"channelId": "chan-1",
			"message":   "hello from top level",
		},
		Discord: &DiscordOps{
			ChannelSend: func(_ context.Context, channelID string, payload any) error {
				gotChannel = channelID
				gotPayload = payload
				return nil
			},
		},
	})
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if gotChannel != "chan-1" {
		t.Fatalf("channel = %q", gotChannel)
	}
	payload, ok := gotPayload.(map[string]any)
	if !ok || payload["content"] != "hello from top level" {
		t.Fatalf("payload = %#v", gotPayload)
	}
	m, ok := result.(map[string]any)
	if !ok || m["content"] != "sent" {
		t.Fatalf("result = %#v", result)
	}
}
