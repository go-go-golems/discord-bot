package jsslack

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type invocationState struct {
	ctx       context.Context
	input     slackbot.Invocation
	responder slackbot.Responder
	replied   bool // owner-only
	workers   errgroup.Group
}
type outcome struct {
	pending bool
	text    *slackbot.Text
	err     error
}

// Dispatch serializes complete invocations, while network work and promise settlement
// run outside the owner. A timeout is bounded even for CPU-bound JavaScript.
func (h *Host) Dispatch(ctx context.Context, input slackbot.Invocation, responder slackbot.Responder) error {
	if err := input.Validate(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()
	stop := context.AfterFunc(h.lifetime, cancel)
	defer stop()
	select {
	case h.gate <- struct{}{}:
		defer func() { <-h.gate }()
	case <-ctx.Done():
		return ctx.Err()
	}
	if h.lifetime.Err() != nil {
		return slackbot.Fail("context_closed", "dispatch", "host is closed")
	}
	s := &invocationState{ctx: ctx, input: input, responder: responder}
	s.workers.SetLimit(16)
	// Cancellation releases any outstanding context-aware service work before returning.
	defer func() { cancel(); _ = s.workers.Wait() }()
	var result goja.Value // only inspected on owner
	_, err := h.call(ctx, "slack.dispatch", func(vm *goja.Runtime) (any, error) {
		key := "command:" + input.Command
		if input.Event != "" {
			key = "event:" + input.Event
		}
		fn, ok := h.handlers[key]
		if !ok {
			return nil, slackbot.Fail("not_found", "dispatch", "no handler for "+key)
		}
		var err error
		result, err = fn(goja.Undefined(), h.buildContext(vm, s))
		return nil, err
	})
	if err != nil {
		return errors.Wrap(err, "Slack handler")
	}
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()
	for {
		raw, err := h.call(ctx, "slack.settle", func(vm *goja.Runtime) (any, error) {
			value := result
			if p, ok := value.Export().(*goja.Promise); ok {
				switch p.State() {
				case goja.PromiseStatePending:
					return outcome{pending: true}, nil
				case goja.PromiseStateRejected:
					return outcome{err: errors.New(rejectionText(vm, p.Result()))}, nil
				default:
					value = p.Result()
				}
			}
			if goja.IsNull(value) || goja.IsUndefined(value) {
				return outcome{}, nil
			}
			var text slackbot.Text
			if err := decode(vm, value, &text); err != nil {
				return outcome{err: err}, nil
			}
			return outcome{text: &text, err: text.Validate()}, nil
		})
		if err != nil {
			return err
		}
		o := raw.(outcome)
		if o.err != nil {
			return o.err
		}
		if !o.pending {
			if o.text != nil {
				// Claim the reply slot on the owner; perform network I/O outside it.
				_, err = h.call(ctx, "slack.auto-reply", func(*goja.Runtime) (any, error) { return nil, claimReply(s) })
				if err != nil {
					return err
				}
				if _, err = h.sendReply(s, *o.text); err != nil {
					return publicError(err, "reply")
				}
			}
			return s.workers.Wait()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
		}
	}
}
func rejectionText(vm *goja.Runtime, v goja.Value) string {
	if o, ok := v.(*goja.Object); ok {
		if stack := o.Get("stack"); stack != nil && !goja.IsUndefined(stack) {
			return stack.String()
		}
	}
	return v.String()
}
func claimReply(s *invocationState) error {
	if s.ctx.Err() != nil {
		return slackbot.Fail("context_closed", "reply", "invocation is closed")
	}
	if s.replied {
		return slackbot.Fail("already_replied", "reply", "one implicit reply is allowed")
	}
	s.replied = true
	return nil
}
func (h *Host) sendReply(s *invocationState, text slackbot.Text) (any, error) {
	if s.input.Command != "" {
		if s.responder == nil {
			return nil, slackbot.Fail("unavailable", "reply", "no response capability")
		}
		if err := s.responder.Reply(s.ctx, text); err != nil {
			return nil, err
		}
		return map[string]any{"delivered": true, "via": "response_url"}, nil
	}
	if h.messages == nil {
		return nil, slackbot.Fail("unavailable", "reply", "no message service")
	}
	thread := s.input.ThreadTS
	if thread == "" {
		thread = s.input.TS
	}
	ref, err := h.messages.Post(s.ctx, slackbot.PostMessage{ChannelID: s.input.ChannelID, Text: text.Text, ThreadTS: thread})
	return map[string]any{"channelId": ref.ChannelID, "ts": ref.TS}, err
}
func publicError(err error, op string) error {
	var domain *slackbot.Error
	if errors.As(err, &domain) {
		return domain
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return slackbot.Fail("context_closed", op, "operation canceled")
	}
	return slackbot.Fail("service_error", op, "service operation failed")
}
func jsError(vm *goja.Runtime, err error) goja.Value {
	e := vm.NewGoError(err)
	var d *slackbot.Error
	if errors.As(err, &d) {
		_ = e.Set("code", d.Code)
		_ = e.Set("operation", d.Operation)
	}
	return e
}
func (h *Host) async(vm *goja.Runtime, s *invocationState, op string, work func() (any, error)) goja.Value {
	if s.ctx.Err() != nil {
		panic(jsError(vm, slackbot.Fail("context_closed", op, "invocation is closed")))
	}
	promise, resolve, reject := vm.NewPromise()
	if !s.workers.TryGo(func() error {
		value, err := work()
		_, _ = h.call(s.ctx, "slack.complete", func(vm *goja.Runtime) (any, error) {
			if err != nil {
				return nil, reject(jsError(vm, publicError(err, op)))
			}
			return nil, resolve(value)
		})
		return nil
	}) {
		_ = reject(jsError(vm, slackbot.Fail("busy", op, "too many pending operations")))
	}
	return vm.ToValue(promise)
}
func (h *Host) buildContext(vm *goja.Runtime, s *invocationState) *goja.Object {
	c := vm.NewObject()
	for key, value := range map[string]any{"id": s.input.ID, "teamId": s.input.TeamID, "channelId": s.input.ChannelID, "userId": s.input.UserID, "command": s.input.Command, "text": s.input.Text, "config": h.config} {
		// JSON roundtrip ensures no reflected Go map remains mutable from JavaScript.
		b, _ := json.Marshal(value)
		var detached any
		_ = json.Unmarshal(b, &detached)
		must(vm, c.Set(key, detached))
	}
	must(vm, c.Set("event", map[string]any{"type": s.input.Event, "text": s.input.Text, "ts": s.input.TS, "threadTs": s.input.ThreadTS, "channelId": s.input.ChannelID, "userId": s.input.UserID}))
	must(vm, c.Set("reply", func(call goja.FunctionCall) goja.Value {
		var text slackbot.Text
		must(vm, decode(vm, call.Argument(0), &text))
		must(vm, text.Validate())
		if err := claimReply(s); err != nil {
			panic(jsError(vm, err))
		}
		return h.async(vm, s, "reply", func() (any, error) { return h.sendReply(s, text) })
	}))
	messages := vm.NewObject()
	must(vm, messages.Set("post", func(call goja.FunctionCall) goja.Value {
		var input slackbot.PostMessage
		must(vm, decode(vm, call.Argument(0), &input))
		must(vm, input.Validate())
		return h.async(vm, s, "messages.post", func() (any, error) {
			if h.messages == nil {
				return nil, slackbot.Fail("unavailable", "messages.post", "no message service")
			}
			ref, err := h.messages.Post(s.ctx, input)
			return map[string]any{"channelId": ref.ChannelID, "ts": ref.TS}, err
		})
	}))
	slack := vm.NewObject()
	must(vm, slack.Set("messages", messages))
	must(vm, c.Set("slack", slack))
	logger := vm.NewObject()
	for _, level := range []string{"debug", "info", "warn", "error"} {
		must(vm, logger.Set(level, func(call goja.FunctionCall) goja.Value {
			if s.ctx.Err() != nil {
				panic(jsError(vm, slackbot.Fail("context_closed", "log", "invocation is closed")))
			}
			// Only script-supplied message text is logged; no host settings or transport payloads.
			msg := strings.TrimSpace(call.Argument(0).String())
			event := h.logger.Info()
			switch level {
			case "debug":
				event = h.logger.Debug()
			case "warn":
				event = h.logger.Warn()
			case "error":
				event = h.logger.Error()
			}
			event.Str("bot", h.descriptor.Name).Str("invocation", s.input.ID).Msg(msg)
			return goja.Undefined()
		}))
	}
	must(vm, c.Set("store", h.storeObject(vm, s)))
	must(vm, c.Set("log", logger))
	return c
}
