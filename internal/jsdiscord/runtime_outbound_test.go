package jsdiscord

import (
	"context"
	"testing"
)

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
