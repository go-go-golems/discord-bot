package slackhost_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/go-go-golems/discord-bot/pkg/slackhost"
	"github.com/stretchr/testify/require"
)

type acknowledgments struct{ count int }

var _ slackbot.Acknowledger = (*acknowledgments)(nil)

func (a *acknowledgments) Ack(context.Context, string, bool) error { a.count++; return nil }
func TestIngressThroughRealJavaScriptHost(t *testing.T) {
	path, err := filepath.Abs("../../examples/slack-bots/ping/index.js")
	require.NoError(t, err)
	recorder := &slackbot.Recorder{}
	h, err := slackhost.Load(context.Background(), path, slackhost.Options{Messages: recorder})
	require.NoError(t, err)
	defer func() { require.NoError(t, h.Close(context.Background())) }()
	ingress, err := slackbot.NewIngress(context.Background(), slackbot.IngressOptions{TeamID: "T", AppID: "A", Capacity: 2, DedupeCapacity: 16, DedupeTTL: time.Minute}, h)
	require.NoError(t, err)
	defer func() { require.NoError(t, ingress.Close(context.Background())) }()
	ack := &acknowledgments{}
	envelope := slackbot.Envelope{ID: "delivery1", EventID: "event1", AppID: "A", Invocation: slackbot.Invocation{TeamID: "T", ChannelID: "C", UserID: "U", Event: "app_mention", TS: "1234567890.000123"}}
	ackCtx, cancel := context.WithCancel(context.Background())
	status, err := ingress.Admit(ackCtx, envelope, ack)
	require.NoError(t, err)
	require.Equal(t, slackbot.Accepted, status)
	cancel() // Receipt lifetime must not own the admitted handler.
	envelope.ID = "delivery2"
	status, err = ingress.Admit(context.Background(), envelope, ack)
	require.NoError(t, err)
	require.Equal(t, slackbot.Duplicate, status)
	require.Eventually(t, func() bool { return len(recorder.Operations()) == 1 }, time.Second, time.Millisecond)
	operations := recorder.Operations()
	require.Equal(t, "1234567890.000123", operations[0].Message.ThreadTS)
	require.Equal(t, 2, ack.count)
	require.NoError(t, ingress.Close(context.Background()))
	for err := range ingress.Errors() {
		require.NoError(t, err)
	}
}
