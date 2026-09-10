package slackbot

import (
	"context"
	"fmt"
	"sync"
)

// RecordedOperation is an offline side effect. It contains no tokens or response URLs.
type RecordedOperation struct {
	Kind    string       `json:"kind"`
	Message *PostMessage `json:"message,omitempty"`
	Reply   *Text        `json:"reply,omitempty"`
	Ref     *MessageRef  `json:"ref,omitempty"`
}
type Recorder struct {
	mu         sync.Mutex
	operations []RecordedOperation
}

var _ MessageService = (*Recorder)(nil)
var _ Responder = (*Recorder)(nil)

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
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := m.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.operations = append(r.operations, RecordedOperation{Kind: "ephemeral_reply", Reply: &m})
	return nil
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
			out[i].Reply = &v
		}
		if op.Ref != nil {
			v := *op.Ref
			out[i].Ref = &v
		}
	}
	return out
}
