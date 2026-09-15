package slackbot

import (
	"context"
	"fmt"
	"sync"
)

// RecordedOperation is an offline side effect. It contains no tokens or response URLs.
type RecordedOperation struct {
	Kind    string               `json:"kind"`
	Message *PostMessage         `json:"message,omitempty"`
	Reply   *MessagePayload      `json:"reply,omitempty"`
	Ref     *MessageRef          `json:"ref,omitempty"`
	Ack     *InteractionResponse `json:"ack,omitempty"`
	View    *ModalView           `json:"view,omitempty"`
	Trigger string               `json:"triggerId,omitempty"`
}
type Recorder struct {
	mu         sync.Mutex
	operations []RecordedOperation
}

var _ MessageService = (*Recorder)(nil)
var _ Responder = (*Recorder)(nil)
var _ InteractionAcknowledger = (*Recorder)(nil)
var _ ViewService = (*Recorder)(nil)

func (r *Recorder) Post(ctx context.Context, m PostMessage) (MessageRef, error) {
	if err := ctx.Err(); err != nil {
		return MessageRef{}, err
	}
	if err := m.Validate(); err != nil {
		return MessageRef{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	ref := MessageRef{m.ChannelID, fmt.Sprintf("offline.%06d", len(r.operations)+1)}
	r.operations = append(r.operations, RecordedOperation{Kind: "post", Message: &m, Ref: &ref})
	return ref, nil
}
func (r *Recorder) Reply(ctx context.Context, m Text) error {
	return r.ReplyMessage(ctx, MessagePayload{Text: m.Text})
}

func (r *Recorder) ReplyMessage(ctx context.Context, m MessagePayload) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := m.Validate("reply"); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.operations = append(r.operations, RecordedOperation{Kind: "ephemeral_reply", Reply: &m})
	return nil
}

func (r *Recorder) Respond(ctx context.Context, response InteractionResponse) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if response.Kind != "accept" && response.Kind != "errors" {
		return Fail("invalid_argument", "ack", "response kind must be accept or errors")
	}
	if response.Kind == "errors" && len(response.Errors) == 0 {
		return Fail("invalid_argument", "ack", "errors response requires at least one field")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	copyResponse := InteractionResponse{Kind: response.Kind}
	if response.Errors != nil {
		copyResponse.Errors = map[string]string{}
		for key, value := range response.Errors {
			copyResponse.Errors[key] = value
		}
	}
	r.operations = append(r.operations, RecordedOperation{Kind: "ack", Ack: &copyResponse})
	return nil
}

func (r *Recorder) Open(ctx context.Context, triggerID string, view ModalView) (ViewRef, error) {
	if err := ctx.Err(); err != nil {
		return ViewRef{}, err
	}
	if err := view.Validate(); err != nil {
		return ViewRef{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	copyView := view
	copyView.Blocks = append([]Block(nil), view.Blocks...)
	ref := ViewRef{ID: fmt.Sprintf("offline-view.%06d", len(r.operations)+1), Hash: "offline-hash"}
	r.operations = append(r.operations, RecordedOperation{Kind: "open_view", View: &copyView, Trigger: triggerID})
	return ref, nil
}
func (r *Recorder) Operations() []RecordedOperation {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]RecordedOperation, len(r.operations))
	for i, op := range r.operations {
		out[i] = op
		if op.Message != nil {
			v := *op.Message
			out[i].Message = &v
		}
		if op.Reply != nil {
			v := *op.Reply
			if op.Reply.Blocks != nil {
				v.Blocks = append([]Block(nil), op.Reply.Blocks...)
			}
			out[i].Reply = &v
		}
		if op.Ref != nil {
			v := *op.Ref
			out[i].Ref = &v
		}
		if op.Ack != nil {
			v := InteractionResponse{Kind: op.Ack.Kind}
			if op.Ack.Errors != nil {
				v.Errors = map[string]string{}
				for key, value := range op.Ack.Errors {
					v.Errors[key] = value
				}
			}
			out[i].Ack = &v
		}
		if op.View != nil {
			v := *op.View
			v.Blocks = append([]Block(nil), op.View.Blocks...)
			out[i].View = &v
		}
	}
	return out
}
