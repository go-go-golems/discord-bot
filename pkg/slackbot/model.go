// Package slackbot defines Slack-specific contracts without a JavaScript or SDK dependency.
package slackbot

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// Error is a stable, credential-free error exposed to bot scripts.
type Error struct {
	Code      string `json:"code"`
	Operation string `json:"operation"`
	Message   string `json:"message"`
}

var _ error = (*Error)(nil)

func (e *Error) Error() string            { return e.Operation + ": " + e.Code + ": " + e.Message }
func Fail(code, op, message string) error { return &Error{code, op, message} }

type Text struct {
	Text string `json:"text"`
}

func (m Text) Validate() error {
	if strings.TrimSpace(m.Text) == "" || utf8.RuneCountInString(m.Text) > 4000 {
		return Fail("invalid_argument", "message", "text must contain 1–4000 characters")
	}
	return nil
}

// Block is a detached Block Kit object. The framework keeps the wire shape as
// JSON so that supported helpers and the low-level message escape hatch share
// one transport boundary without exposing Slack SDK values to JavaScript.
type Block map[string]any

func (b Block) Type() string {
	if b == nil {
		return ""
	}
	v, _ := b["type"].(string)
	return strings.TrimSpace(v)
}

func (b Block) Validate() error {
	if b.Type() == "" {
		return Fail("invalid_argument", "message.blocks", "block type is required")
	}
	return nil
}

// MessagePayload is the common rich-message shape used by replies and posts.
// Text remains mandatory as the notification and accessibility fallback.
type MessagePayload struct {
	Text   string  `json:"text"`
	Blocks []Block `json:"blocks,omitempty"`
}

func (m MessagePayload) Validate(operation string) error {
	if err := (Text{Text: m.Text}).Validate(); err != nil {
		return Fail("invalid_argument", operation, "text must contain 1–4000 characters")
	}
	if len(m.Blocks) > 50 {
		return Fail("invalid_argument", operation, "at most 50 blocks are supported")
	}
	for i, block := range m.Blocks {
		if err := block.Validate(); err != nil {
			return Fail("invalid_argument", operation, "blocks["+strconv.Itoa(i)+"] is invalid")
		}
	}
	return nil
}

func (m MessagePayload) JSON() ([]byte, error) { return json.Marshal(m) }

type MessageRef struct {
	ChannelID string `json:"channelId"`
	TS        string `json:"ts"`
}
type PostMessage struct {
	ChannelID string  `json:"channelId"`
	Text      string  `json:"text"`
	ThreadTS  string  `json:"threadTs,omitempty"`
	Blocks    []Block `json:"blocks,omitempty"`
}

func (m PostMessage) Validate() error {
	if strings.TrimSpace(m.ChannelID) == "" {
		return Fail("invalid_argument", "messages.post", "channelId is required")
	}
	return MessagePayload{Text: m.Text, Blocks: m.Blocks}.Validate("messages.post")
}

type MessageService interface {
	Post(context.Context, PostMessage) (MessageRef, error)
}

// Responder retains the response URL privately in its implementation. It sends an ephemeral reply.
type Responder interface {
	Reply(context.Context, Text) error
}

// RichResponder is implemented by response URL transports that can preserve
// Block Kit payloads. Text-only responders remain valid for simple tests and
// integrations; callers must check this capability before sending blocks.
type RichResponder interface {
	ReplyMessage(context.Context, MessagePayload) error
}
type Invocation struct {
	ID        string  `json:"id"`
	TeamID    string  `json:"teamId"`
	ChannelID string  `json:"channelId"`
	UserID    string  `json:"userId"`
	Command   string  `json:"command,omitempty"`
	Event     string  `json:"event,omitempty"`
	Text      string  `json:"text"`
	TS        string  `json:"ts,omitempty"`
	ThreadTS  string  `json:"threadTs,omitempty"`
	Action    *Action `json:"action,omitempty"`
}

// Action is the normalized subset of a Slack block action needed by a local
// bot. Selection fields remain detached JSON so the transport can preserve
// Slack's different element-specific values without exposing SDK objects.
type Action struct {
	Type            string         `json:"type"`
	ActionID        string         `json:"actionId"`
	BlockID         string         `json:"blockId,omitempty"`
	Value           string         `json:"value,omitempty"`
	SelectedOption  map[string]any `json:"selectedOption,omitempty"`
	SelectedOptions []any          `json:"selectedOptions,omitempty"`
	MessageTS       string         `json:"messageTs,omitempty"`
	ThreadTS        string         `json:"threadTs,omitempty"`
	ResponseURL     string         `json:"-"`
}

func (a Action) ID() string { return a.ActionID }

func (i Invocation) Validate() error {
	if i.TeamID == "" || i.UserID == "" {
		return Fail("invalid_argument", "dispatch", "teamId and userId are required")
	}
	kinds := 0
	if i.Command != "" {
		kinds++
	}
	if i.Event != "" {
		kinds++
	}
	if i.Action != nil {
		kinds++
	}
	if kinds != 1 {
		return Fail("invalid_argument", "dispatch", "provide exactly one command, event or action")
	}
	if i.Command != "" && i.ChannelID == "" {
		return Fail("invalid_argument", "dispatch", "channelId is required for commands")
	}
	if i.Event != "" && (i.Event != "app_mention" || i.TS == "") {
		return Fail("invalid_argument", "dispatch", "only app_mention with a string ts is supported")
	}
	if i.Action != nil && i.Action.ID() == "" {
		return Fail("invalid_argument", "dispatch", "action id is required")
	}
	return nil
}

type Field struct {
	Type     string `json:"type"`
	Default  any    `json:"default,omitempty"`
	Required bool   `json:"required,omitempty"`
	Help     string `json:"help,omitempty"`
}
type RunSchema struct {
	Fields map[string]Field `json:"fields,omitempty"`
}
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Descriptor struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ScriptPath  string    `json:"scriptPath"`
	Run         RunSchema `json:"run"`
	Commands    []Command `json:"commands"`
	Events      []string  `json:"events"`
	Actions     []string  `json:"actions,omitempty"`
}

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var fieldPattern = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)

func (d Descriptor) Validate() error {
	if !namePattern.MatchString(d.Name) {
		return errors.New("bot name must match [a-z][a-z0-9-]*")
	}
	for name, f := range d.Run.Fields {
		if !fieldPattern.MatchString(name) {
			return errors.Errorf("invalid bot field name %q", name)
		}
		if f.Type != "string" && f.Type != "bool" && f.Type != "number" {
			return errors.Errorf("field %s: unsupported type %q", name, f.Type)
		}
		if f.Default != nil && !fieldValueValid(f.Type, f.Default) {
			return errors.Errorf("field %s: default must be %s", name, f.Type)
		}
	}
	return nil
}
func fieldValueValid(kind string, v any) bool {
	switch kind {
	case "string":
		_, ok := v.(string)
		return ok
	case "bool":
		_, ok := v.(bool)
		return ok
	case "number":
		switch v.(type) {
		case float64, int, int64:
			return true
		}
	}
	return false
}

// Config projects ONLY declared bot fields. Host settings never enter this map.
func (d Descriptor) Config(input map[string]any) (map[string]any, error) {
	out := map[string]any{}
	for name, f := range d.Run.Fields {
		value, ok := input[name]
		if !ok {
			value = f.Default
		}
		if value == nil {
			if f.Required {
				return nil, errors.Errorf("field %s is required", name)
			}
			continue
		}
		if !fieldValueValid(f.Type, value) {
			return nil, errors.Errorf("field %s must be %s", name, f.Type)
		}
		out[name] = value
	}
	return out, nil
}
