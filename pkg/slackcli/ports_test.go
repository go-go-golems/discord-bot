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
			defer func() {
				if err := host.Close(context.Background()); err != nil {
					t.Error(err)
				}
			}()
			require.NoError(t, host.Dispatch(context.Background(), slackbot.Invocation{TeamID: "T", ChannelID: "C", UserID: "U", Command: tc.command, Text: tc.text}, recorder))
			ops := recorder.Operations()
			require.Len(t, ops, 1)
			require.NotNil(t, ops[0].Reply)
			require.Contains(t, ops[0].Reply.Text, tc.expected)
		})
	}
}

func TestPokerRoundAndIsolation(t *testing.T) {
	d, err := Resolve(context.Background(), "../../examples/slack-bots", "poker", 0)
	require.NoError(t, err)
	rec := &slackbot.Recorder{}
	h, err := slackhost.Load(context.Background(), d.ScriptPath, slackhost.Options{Messages: rec, Views: rec, Operations: rec})
	require.NoError(t, err)
	defer func() {
		if err := h.Close(context.Background()); err != nil {
			t.Error(err)
		}
	}()
	dispatch := func(user, cmd, text string) string {
		require.NoError(t, h.Dispatch(context.Background(), slackbot.Invocation{TeamID: "T", ChannelID: "C", UserID: user, Command: cmd, Text: text}, rec))
		ops := rec.Operations()
		return ops[len(ops)-1].Reply.Text
	}
	require.Contains(t, dispatch("U", "/poker-deal", ""), "Round 1")
	require.Contains(t, dispatch("OTHER", "/poker-score", ""), "No hand")
	require.Contains(t, dispatch("U", "/poker-draw", "1,3,5"), "draw used: true")
	require.Contains(t, dispatch("U", "/poker-draw", "1"), "already used")
	require.Contains(t, dispatch("U", "/poker-rank", "As Ks Qs Js Ts"), "Straight")
	require.Contains(t, dispatch("U", "/poker-rank", "As As Qs Js Ts"), "Duplicate")
	require.Contains(t, dispatch("U", "/poker-reset", ""), "cleared")
	require.Contains(t, dispatch("U", "/poker-score", ""), "No hand")
}

func TestSupportDraftAndThreadWorkflows(t *testing.T) {
	d, err := Resolve(context.Background(), "../../examples/slack-bots", "support", 0)
	require.NoError(t, err)
	rec := &slackbot.Recorder{}
	h, err := slackhost.Load(context.Background(), d.ScriptPath, slackhost.Options{Messages: rec, Views: rec, Operations: rec})
	require.NoError(t, err)
	defer func() {
		if err := h.Close(context.Background()); err != nil {
			t.Error(err)
		}
	}()
	i := slackbot.Invocation{TeamID: "T", ChannelID: "C", UserID: "U", Command: "/support-ticket", Text: "Printer"}
	require.NoError(t, h.Dispatch(context.Background(), i, rec))
	require.Equal(t, "ephemeral_reply", rec.Operations()[0].Kind)
	require.Equal(t, "messages.ephemeral", rec.Operations()[1].Kind)
	i.Command = "/support-start-thread"
	require.NoError(t, h.Dispatch(context.Background(), i, rec))
	ops := rec.Operations()
	require.Equal(t, ops[2].Ref.TS, ops[3].Message.ThreadTS)
}
