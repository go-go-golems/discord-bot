// Package slackhost composes the Slack JavaScript host with injected Go services.
// It performs no networking of its own.
package slackhost

import (
	"context"
	"github.com/go-go-golems/discord-bot/internal/jsslack"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"time"
)

type Options = jsslack.Options

// Host exposes lifecycle and dispatch without exposing the VM.
type Host struct{ host *jsslack.Host }

var _ slackbot.Dispatcher = (*Host)(nil)

func Load(ctx context.Context, script string, opts Options) (*Host, error) {
	h, err := jsslack.Load(ctx, script, opts)
	if err != nil {
		return nil, err
	}
	return &Host{host: h}, nil
}
func (h *Host) Descriptor() slackbot.Descriptor { return h.host.Descriptor() }
func (h *Host) Dispatch(ctx context.Context, input slackbot.Invocation, responder slackbot.Responder) error {
	return h.host.Dispatch(ctx, input, responder)
}
func (h *Host) Close(ctx context.Context) error { return h.host.Close(ctx) }

func Inspect(ctx context.Context, script string, timeout time.Duration) (slackbot.Descriptor, error) {
	return jsslack.Inspect(ctx, script, timeout)
}
