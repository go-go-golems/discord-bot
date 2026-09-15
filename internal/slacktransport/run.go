package slacktransport

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/pkg/errors"
	"github.com/slack-go/slack/socketmode"
	"golang.org/x/sync/errgroup"
)

type receipt struct {
	socket          *socketmode.Client
	id              string
	acceptsResponse bool
	deadline        time.Time
	remember        func(any)
	mu              sync.Mutex
	used            bool
}

var _ slackbot.Acknowledger = (*receipt)(nil)

func (r *receipt) Ack(ctx context.Context, id string, busy bool) error {
	var payload any
	if busy && r.acceptsResponse {
		payload = map[string]string{"response_type": "ephemeral", "text": "Bot is busy. Please try again."}
	}
	return r.send(ctx, id, payload, slackbot.InteractionResponse{Kind: "auto"})
}

func (r *receipt) Respond(ctx context.Context, response slackbot.InteractionResponse) error {
	if response.Kind != "accept" && response.Kind != "errors" {
		return slackbot.Fail("invalid_argument", "ack", "response kind must be accept or errors")
	}
	if response.Kind == "errors" && len(response.Errors) == 0 {
		return slackbot.Fail("invalid_argument", "ack", "errors response requires at least one field")
	}
	var payload any
	if response.Kind == "errors" {
		payload = map[string]any{"response_action": "errors", "errors": response.Errors}
	}
	return r.send(ctx, r.id, payload, response)
}

func (r *receipt) send(ctx context.Context, id string, payload any, response slackbot.InteractionResponse) error {
	r.mu.Lock()
	if r.used {
		r.mu.Unlock()
		return slackbot.Fail("ack_already_sent", "ack", "interaction acknowledgment was already chosen")
	}
	if !r.deadline.IsZero() && time.Now().After(r.deadline) {
		r.mu.Unlock()
		return slackbot.Fail("ack_expired", "ack", "interaction acknowledgment deadline has passed")
	}
	if response.Kind == "errors" && !r.acceptsResponse {
		r.mu.Unlock()
		return slackbot.Fail("ack_payload_unsupported", "ack", "Slack did not accept an acknowledgment payload")
	}
	r.used = true
	r.mu.Unlock()
	if !r.deadline.IsZero() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, r.deadline)
		defer cancel()
	}
	if err := r.socket.AckCtx(ctx, id, payload); err != nil {
		return slackbot.Fail("ack_failed", "ack", "could not schedule acknowledgment")
	}
	if r.remember != nil {
		r.remember(payload)
	}
	return nil
}

func (r *receipt) replay(ctx context.Context, payload any) error {
	if !r.deadline.IsZero() && time.Now().After(r.deadline) {
		return slackbot.Fail("ack_expired", "ack", "interaction acknowledgment deadline has passed")
	}
	if !r.deadline.IsZero() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, r.deadline)
		defer cancel()
	}
	if err := r.socket.AckCtx(ctx, r.id, payload); err != nil {
		return slackbot.Fail("ack_failed", "ack", "could not schedule acknowledgment")
	}
	return nil
}

// Run verifies the expected workspace, then owns ingress and socket lifetimes.
// One Client is intended for one Run call. SDK reconnects remain within the same ingress lifetime.
func (c *Client) Run(ctx context.Context, dispatcher slackbot.Dispatcher) error {
	identity, err := c.api.AuthTestContext(ctx)
	if err != nil {
		return safeError(ctx, err, "auth.test")
	}
	if identity.TeamID != c.opts.TeamID {
		return errors.New("authenticated workspace does not match configured team")
	}
	group, workerCtx := errgroup.WithContext(ctx)
	ingress, err := slackbot.NewIngress(workerCtx, slackbot.IngressOptions{
		TeamID: c.opts.TeamID, AppID: c.opts.AppID, SelfUserID: identity.UserID, AllowedChannels: c.opts.AllowedChannels,
		Capacity: 32, DedupeCapacity: 4096, DedupeTTL: 5 * time.Minute,
	}, dispatcher)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = ingress.Close(closeCtx)
	}()
	group.Go(func() error {
		err := c.socket.RunContext(workerCtx)
		if workerCtx.Err() != nil {
			return nil
		}
		if err != nil {
			return errors.New("Socket Mode connection stopped")
		}
		return nil
	})
	group.Go(func() error {
		for {
			select {
			case <-workerCtx.Done():
				return nil
			case failure, ok := <-ingress.Errors():
				if !ok {
					return nil
				}
				// Handler errors may contain script data. Avoid exposing raw payloads or credentials.
				if failure != nil {
					c.opts.Logger.Warn().Msg("Slack handler failed")
				}
			case event, ok := <-c.socket.Events:
				if !ok {
					return errors.New("Socket Mode event stream closed")
				}
				if event.Type == socketmode.EventTypeInvalidAuth {
					return errors.New("Socket Mode authentication failed")
				}
				if event.Request == nil || event.Request.EnvelopeID == "" {
					continue
				}
				ackCtx, cancel := context.WithTimeout(workerCtx, 2*time.Second)
				envelope, decodeErr := c.decode(*event.Request)
				ack := &receipt{socket: c.socket, id: event.Request.EnvelopeID, acceptsResponse: event.Request.AcceptsResponsePayload && (event.Request.Type == socketmode.RequestTypeSlashCommands || event.Request.Type == socketmode.RequestTypeInteractive), deadline: time.Now().Add(3 * time.Second)}
				if event.Request.Type == socketmode.RequestTypeInteractive {
					ack.remember = func(payload any) { c.rememberAck(event.Request.EnvelopeID, payload) }
				}
				if decodeErr != nil {
					err = ack.Ack(ackCtx, event.Request.EnvelopeID, false)
					c.opts.Logger.Debug().Msg("Dropped unsupported or malformed Slack envelope")
				} else {
					if envelope.Invocation.Interaction != nil {
						envelope.Invocation.Interaction.Ack = ack
					}
					var decision slackbot.Admission
					decision, err = ingress.Admit(ackCtx, envelope, ack)
					if err == nil && decision == slackbot.Duplicate && envelope.Invocation.Interaction != nil {
						if payload, ok := c.replayAck(event.Request.EnvelopeID); ok {
							err = ack.replay(ackCtx, payload)
						}
					}
					c.opts.Logger.Debug().Str("admission", string(decision)).Msg("Slack receipt")
				}
				cancel()
				if err != nil {
					c.opts.Logger.Warn().Msg("Slack acknowledgment failed")
				}
			}
		}
	})
	return group.Wait()
}
func (c *Client) decode(r socketmode.Request) (slackbot.Envelope, error) {
	e := slackbot.Envelope{ID: r.EnvelopeID}
	switch r.Type {
	case socketmode.RequestTypeEventsAPI:
		var p struct {
			Type    string `json:"type"`
			TeamID  string `json:"team_id"`
			AppID   string `json:"api_app_id"`
			EventID string `json:"event_id"`
			Event   struct {
				Type, Channel, User, Text, TS string
				ThreadTS                      string `json:"thread_ts"`
				BotID                         string `json:"bot_id"`
				Subtype                       string `json:"subtype"`
			}
		}
		if err := json.Unmarshal(r.Payload, &p); err != nil {
			return e, errors.New("invalid event payload")
		}
		if p.Type != "event_callback" || p.Event.Type != "app_mention" {
			return e, errors.New("unsupported event")
		}
		e.AppID = p.AppID
		e.EventID = p.EventID
		e.Bot = p.Event.BotID != "" || p.Event.Subtype != ""
		e.Invocation = slackbot.Invocation{ID: r.EnvelopeID, TeamID: p.TeamID, ChannelID: p.Event.Channel, UserID: p.Event.User, Event: p.Event.Type, Text: p.Event.Text, TS: p.Event.TS, ThreadTS: p.Event.ThreadTS}
	case socketmode.RequestTypeSlashCommands:
		var p struct {
			TeamID        string `json:"team_id"`
			AppID         string `json:"api_app_id"`
			ChannelID     string `json:"channel_id"`
			UserID        string `json:"user_id"`
			ResponseURL   string `json:"response_url"`
			Command, Text string
		}
		if err := json.Unmarshal(r.Payload, &p); err != nil {
			return e, errors.New("invalid command payload")
		}
		var err error
		e.Responder, err = c.responseCapability(p.ResponseURL)
		if err != nil {
			return e, err
		}
		e.AppID = p.AppID
		e.Invocation = slackbot.Invocation{ID: r.EnvelopeID, TeamID: p.TeamID, ChannelID: p.ChannelID, UserID: p.UserID, Command: p.Command, Text: p.Text}
	case socketmode.RequestTypeInteractive:
		var p struct {
			Type string `json:"type"`
			Team struct {
				ID string `json:"id"`
			} `json:"team"`
			User struct {
				ID string `json:"id"`
			} `json:"user"`
			Channel struct {
				ID string `json:"id"`
			} `json:"channel"`
			APIAppID    string `json:"api_app_id"`
			ResponseURL string `json:"response_url"`
			CallbackID  string `json:"callback_id"`
			TriggerID   string `json:"trigger_id"`
			Message     struct {
				TS string `json:"ts"`
			} `json:"message"`
			Container struct {
				ChannelID string `json:"channel_id"`
				MessageTS string `json:"message_ts"`
				ThreadTS  string `json:"thread_ts"`
			} `json:"container"`
			View struct {
				ID              string `json:"id"`
				Hash            string `json:"hash"`
				CallbackID      string `json:"callback_id"`
				PrivateMetadata string `json:"private_metadata"`
				State           struct {
					Values map[string]map[string]any `json:"values"`
				} `json:"state"`
			} `json:"view"`
			Actions []struct {
				ActionID        string           `json:"action_id"`
				BlockID         string           `json:"block_id"`
				Type            string           `json:"type"`
				Value           string           `json:"value"`
				SelectedOption  map[string]any   `json:"selected_option"`
				SelectedOptions []map[string]any `json:"selected_options"`
			} `json:"actions"`
		}
		if err := json.Unmarshal(r.Payload, &p); err != nil || p.Type == "" {
			return e, errors.New("unsupported interactive payload")
		}
		if p.Type == "view_submission" {
			e.AppID = p.APIAppID
			e.Invocation = slackbot.Invocation{
				ID: r.EnvelopeID, TeamID: p.Team.ID, UserID: p.User.ID,
				Interaction: &slackbot.Interaction{Type: p.Type, CallbackID: p.View.CallbackID, PrivateMetadata: p.View.PrivateMetadata, ViewID: p.View.ID, ViewHash: p.View.Hash, Values: p.View.State.Values},
			}
			return e, e.Invocation.Validate()
		}
		if p.Type != "block_actions" || len(p.Actions) == 0 {
			return e, errors.New("unsupported interactive payload")
		}
		channelID := p.Channel.ID
		if channelID == "" {
			channelID = p.Container.ChannelID
		}
		threadTS := p.Container.ThreadTS
		messageTS := p.Message.TS
		if messageTS == "" {
			messageTS = p.Container.MessageTS
		}
		var responder slackbot.Responder
		if p.ResponseURL != "" {
			var capabilityErr error
			responder, capabilityErr = c.responseCapability(p.ResponseURL)
			if capabilityErr != nil {
				return e, capabilityErr
			}
		}
		a := p.Actions[0]
		selected := make([]any, len(a.SelectedOptions))
		for i := range a.SelectedOptions {
			selected[i] = a.SelectedOptions[i]
		}
		e.AppID = p.APIAppID
		e.Responder = responder
		e.Invocation = slackbot.Invocation{
			ID: r.EnvelopeID, TeamID: p.Team.ID, ChannelID: channelID, UserID: p.User.ID,
			Action: &slackbot.Action{Type: a.Type, ActionID: a.ActionID, BlockID: a.BlockID, Value: a.Value, SelectedOption: a.SelectedOption, SelectedOptions: selected, MessageTS: messageTS, ThreadTS: threadTS, ResponseURL: p.ResponseURL, TriggerID: p.TriggerID, CallbackID: p.CallbackID},
		}
	default:
		return e, errors.New("unsupported envelope")
	}
	return e, e.Invocation.Validate()
}
