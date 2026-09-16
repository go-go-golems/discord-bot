package jsslack

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
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
	ackMu     sync.Mutex
	ackChosen bool
}
type outcome struct {
	pending bool
	message *slackbot.MessagePayload
	err     error
}

// Dispatch serializes complete invocations, while network work and promise settlement
// run outside the owner. A timeout is bounded even for CPU-bound JavaScript.
func (h *Host) Dispatch(ctx context.Context, input slackbot.Invocation, responder slackbot.Responder) (dispatchErr error) {
	started := time.Now()
	logger := h.logger.With().Str("bot", h.descriptor.Name).Str("invocation", input.ID).
		Str("command", input.Command).Str("event", input.Event).Logger()
	logger.Debug().Msg("Slack dispatch started")
	defer func() {
		if dispatchErr == nil {
			logger.Debug().Dur("duration_ms", time.Since(started)).Msg("Slack dispatch completed")
			return
		}
		failure := logger.Warn().Dur("duration_ms", time.Since(started))
		var domainErr *slackbot.Error
		if errors.As(dispatchErr, &domainErr) {
			failure = failure.Str("error_code", domainErr.Code)
		} else if errors.Is(dispatchErr, context.DeadlineExceeded) {
			failure = failure.Str("error_code", "deadline_exceeded")
		} else {
			failure = failure.Str("error_code", "handler_error")
		}
		failure.Msg("Slack dispatch failed")
	}()
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
		} else if input.Action != nil {
			key = "action:" + input.Action.ID()
		} else if input.Shortcut != nil {
			key = "shortcut:" + input.Shortcut.CallbackID
		} else if input.Interaction != nil {
			prefix := "view:"
			if input.Interaction.Type == "block_suggestion" {
				prefix = "options:"
			}
			key = prefix + input.Interaction.CallbackID
		}
		fn, ok := h.handlers[key]
		if !ok {
			h.logger.Warn().Str("bot", h.descriptor.Name).Str("invocation", input.ID).Str("handler", key).Msg("No registered Slack handler")
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
				case goja.PromiseStateFulfilled:
					value = p.Result()
				}
			}
			if goja.IsNull(value) || goja.IsUndefined(value) {
				return outcome{}, nil
			}
			var message slackbot.MessagePayload
			if err := decode(vm, value, &message); err != nil {
				return outcome{err: err}, nil
			}
			return outcome{message: &message, err: message.Validate("reply")}, nil
		})
		if err != nil {
			return err
		}
		o := raw.(outcome)
		if o.err != nil {
			return o.err
		}
		if !o.pending {
			if s.input.Interaction != nil && !s.hasAckChosen() {
				return slackbot.Fail("ack_required", "ack", "view handler must choose an acknowledgment")
			}
			if o.message != nil {
				// Claim the reply slot on the owner; perform network I/O outside it.
				_, err = h.call(ctx, "slack.auto-reply", func(*goja.Runtime) (any, error) { return nil, claimReply(s) })
				if err != nil {
					return err
				}
				if _, err = h.sendReply(s, *o.message); err != nil {
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
func (h *Host) sendReply(s *invocationState, message slackbot.MessagePayload) (result any, replyErr error) {
	logger := h.logger.With().Str("bot", h.descriptor.Name).Str("invocation", s.input.ID).
		Str("command", s.input.Command).Str("text_sha256", fmt.Sprintf("%x", sha256.Sum256([]byte(message.Text)))).
		Int("text_bytes", len(message.Text)).Int("blocks", len(message.Blocks)).Logger()
	logger.Debug().Msg("Slack reply sending")
	defer func() {
		logger.Debug().Bool("delivered", replyErr == nil).Msg("Slack reply finished")
	}()
	if s.input.Command != "" || s.input.Action != nil {
		if s.responder == nil {
			return nil, slackbot.Fail("unavailable", "reply", "no response capability")
		}
		var err error
		if len(message.Blocks) > 0 {
			rich, ok := s.responder.(slackbot.RichResponder)
			if !ok {
				return nil, slackbot.Fail("unavailable", "reply", "response capability does not support blocks")
			}
			err = rich.ReplyMessage(s.ctx, message)
		} else {
			err = s.responder.Reply(s.ctx, slackbot.Text{Text: message.Text})
		}
		if err != nil {
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
	ref, err := h.messages.Post(s.ctx, slackbot.PostMessage{ChannelID: s.input.ChannelID, Text: message.Text, Blocks: message.Blocks, ThreadTS: thread})
	return map[string]any{"channelId": ref.ChannelID, "ts": ref.TS}, err
}

func (s *invocationState) chooseAck() bool {
	s.ackMu.Lock()
	defer s.ackMu.Unlock()
	if s.ackChosen {
		return false
	}
	s.ackChosen = true
	return true
}

func (s *invocationState) hasAckChosen() bool {
	s.ackMu.Lock()
	defer s.ackMu.Unlock()
	return s.ackChosen
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
	must(vm, c.Set("event", map[string]any{"data": s.input.EventData, "type": s.input.Event, "text": s.input.Text, "ts": s.input.TS, "threadTs": s.input.ThreadTS, "channelId": s.input.ChannelID, "userId": s.input.UserID}))
	if s.input.Action != nil {
		action := map[string]any{
			"type": s.input.Action.Type, "actionId": s.input.Action.ActionID,
			"blockId": s.input.Action.BlockID, "value": s.input.Action.Value,
			"selectedOption":  s.input.Action.SelectedOption,
			"selectedOptions": s.input.Action.SelectedOptions,
			"selection":       s.input.Action.Selection,
			"messageTs":       s.input.Action.MessageTS, "threadTs": s.input.Action.ThreadTS,
		}
		must(vm, c.Set("action", action))
	} else {
		must(vm, c.Set("action", goja.Undefined()))
	}
	if s.input.Shortcut != nil {
		must(vm, c.Set("shortcut", map[string]any{"type": s.input.Shortcut.Type, "callbackId": s.input.Shortcut.CallbackID, "message": s.input.Shortcut.Message}))
	}
	if s.input.Interaction != nil {
		interaction := s.input.Interaction
		must(vm, c.Set("view", map[string]any{"type": interaction.Type, "callbackId": interaction.CallbackID, "privateMetadata": interaction.PrivateMetadata, "id": interaction.ViewID, "hash": interaction.ViewHash}))
		must(vm, c.Set("query", interaction.Query))
		values := vm.NewObject()
		all := map[string]map[string]any{}
		for blockID, actions := range interaction.Values {
			all[blockID] = actions
		}
		must(vm, values.Set("all", all))
		must(vm, values.Set("text", func(call goja.FunctionCall) goja.Value {
			blockID, actionID := call.Argument(0).String(), call.Argument(1).String()
			field, ok := interaction.Values[blockID][actionID]
			if !ok {
				return goja.Undefined()
			}
			fieldMap, ok := field.(map[string]any)
			if !ok {
				return goja.Undefined()
			}
			if value, ok := fieldMap["value"].(string); ok {
				return vm.ToValue(value)
			}
			return goja.Undefined()
		}))
		must(vm, c.Set("values", values))
		ack := vm.NewObject()
		must(vm, ack.Set("accept", func(goja.FunctionCall) goja.Value {
			if !s.chooseAck() {
				panic(jsError(vm, slackbot.Fail("ack_already_sent", "ack", "interaction acknowledgment was already chosen")))
			}
			return h.async(vm, s, "ack", func() (any, error) {
				if interaction.Ack == nil {
					return nil, slackbot.Fail("unavailable", "ack", "no interaction acknowledgment capability")
				}
				return nil, interaction.Ack.Respond(s.ctx, slackbot.InteractionResponse{Kind: "accept"})
			})
		}))
		must(vm, ack.Set("errors", func(call goja.FunctionCall) goja.Value {
			if !s.chooseAck() {
				panic(jsError(vm, slackbot.Fail("ack_already_sent", "ack", "interaction acknowledgment was already chosen")))
			}
			var errs map[string]string
			must(vm, decode(vm, call.Argument(0), &errs))
			return h.async(vm, s, "ack", func() (any, error) {
				if interaction.Ack == nil {
					return nil, slackbot.Fail("unavailable", "ack", "no interaction acknowledgment capability")
				}
				return nil, interaction.Ack.Respond(s.ctx, slackbot.InteractionResponse{Kind: "errors", Errors: errs})
			})
		}))
		for _, kind := range []string{"update", "options"} {
			kind := kind
			must(vm, ack.Set(kind, func(call goja.FunctionCall) goja.Value {
				response := slackbot.InteractionResponse{Kind: kind}
				if kind == "update" {
					if interaction.Type != "view_submission" {
						panic(vm.NewTypeError("ack.update requires view_submission"))
					}
					response.View = &slackbot.ModalView{}
					must(vm, decode(vm, call.Argument(0), response.View))
				} else {
					if interaction.Type != "block_suggestion" {
						panic(vm.NewTypeError("ack.options requires block_suggestion"))
					}
					raw, err := call.Argument(0).ToObject(vm).MarshalJSON()
					must(vm, err)
					must(vm, json.Unmarshal(raw, &response.Options))
					if response.Options == nil {
						response.Options = []slackbot.Block{}
					}
				}
				must(vm, response.Validate())
				if !s.chooseAck() {
					panic(jsError(vm, slackbot.Fail("ack_already_sent", "ack", "interaction acknowledgment was already chosen")))
				}
				return h.async(vm, s, "ack", func() (any, error) {
					if interaction.Ack == nil {
						return nil, slackbot.Fail("unavailable", "ack", "no interaction acknowledgment capability")
					}
					return nil, interaction.Ack.Respond(s.ctx, response)
				})
			}))
		}
		must(vm, c.Set("ack", ack))
	} else {
		must(vm, c.Set("view", goja.Undefined()))
		must(vm, c.Set("values", goja.Undefined()))
		must(vm, c.Set("ack", goja.Undefined()))
	}
	must(vm, c.Set("openModal", func(call goja.FunctionCall) goja.Value {
		trigger := s.input.TriggerID
		if trigger == "" && s.input.Action != nil {
			trigger = s.input.Action.TriggerID
		}
		if trigger == "" {
			panic(jsError(vm, slackbot.Fail("unavailable", "views.open", "no modal trigger is available")))
		}
		if h.views == nil {
			panic(jsError(vm, slackbot.Fail("unavailable", "views.open", "no view service")))
		}
		var view slackbot.ModalView
		must(vm, decode(vm, call.Argument(0), &view))
		must(vm, view.Validate())
		return h.async(vm, s, "views.open", func() (any, error) {
			ref, err := h.views.Open(s.ctx, trigger, view)
			return map[string]any{"id": ref.ID, "hash": ref.Hash}, err
		})
	}))
	must(vm, c.Set("reply", func(call goja.FunctionCall) goja.Value {
		var message slackbot.MessagePayload
		must(vm, decode(vm, call.Argument(0), &message))
		must(vm, message.Validate("reply"))
		if err := claimReply(s); err != nil {
			panic(jsError(vm, err))
		}
		return h.async(vm, s, "reply", func() (any, error) { return h.sendReply(s, message) })
	}))
	must(vm, c.Set("replaceOriginal", func(call goja.FunctionCall) goja.Value {
		if s.input.Action == nil {
			panic(vm.NewTypeError("replaceOriginal requires a message action"))
		}
		updater, ok := s.responder.(slackbot.UpdatingResponder)
		if !ok {
			panic(jsError(vm, slackbot.Fail("unavailable", "replaceOriginal", "no response URL update capability")))
		}
		var message slackbot.MessagePayload
		must(vm, decode(vm, call.Argument(0), &message))
		must(vm, message.Validate("replaceOriginal"))
		must(vm, claimReply(s))
		return h.async(vm, s, "replaceOriginal", func() (any, error) { return nil, updater.Replace(s.ctx, message) })
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
	h.addOperations(vm, s, slack, messages)
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
