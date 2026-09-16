package slackcli

import (
	"context"
	"testing"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/go-go-golems/discord-bot/pkg/slackhost"
	"github.com/stretchr/testify/require"
)

func TestInitialNativeSlackPorts(t *testing.T) {
	for _, tc := range []struct{ bot, command, text, expected string }{
		{"announcements", "/announce-preview", "Release", "Preview ready for Release"},
		{"unified-demo", "/unified-ping", "", "apiKey=(unset)"},
		{"interaction-types", "/fun", "roll -1", "Sides must be an integer"},
		{"hater", "/apology-status", "", "No apology on record"},
	} {
		t.Run(tc.bot, func(t *testing.T) {
			d, err := Resolve(context.Background(), "../../examples/slack-bots", tc.bot, 0)
			require.NoError(t, err)
			recorder := &slackbot.Recorder{}
			host, err := slackhost.Load(context.Background(), d.ScriptPath, slackhost.Options{Messages: recorder, Views: recorder})
			require.NoError(t, err)
			defer host.Close(context.Background())
			require.NoError(t, host.Dispatch(context.Background(), slackbot.Invocation{TeamID: "T", ChannelID: "C", UserID: "U", Command: tc.command, Text: tc.text}, recorder))
			ops := recorder.Operations()
			require.Len(t, ops, 1)
			require.NotNil(t, ops[0].Reply)
			require.Contains(t, ops[0].Reply.Text, tc.expected)
		})
	}
}
