package slacktransport

import (
	"encoding/json"
	"testing"

	"github.com/slack-go/slack/socketmode"
	"github.com/stretchr/testify/require"
)

func TestDecodePortInteractions(t *testing.T) {
	c := &Client{}
	for _, tc := range []struct{ kind, body string }{
		{"shortcut", `"callback_id":"capture","trigger_id":"trigger"`},
		{"message_action", `"callback_id":"quote","trigger_id":"trigger","channel":{"id":"C"},"message":{"text":"Source","ts":"123.000001"}`},
		{"block_suggestion", `"action_id":"search","value":"query"`},
		{"block_actions", `"actions":[{"type":"users_select","action_id":"owner","selected_user":"U2"}]`},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			payload := `{"type":"` + tc.kind + `","team":{"id":"T"},"user":{"id":"U"},"api_app_id":"A",` + tc.body + `}`
			e, err := c.decode(socketmode.Request{Type: socketmode.RequestTypeInteractive, Payload: json.RawMessage(payload), EnvelopeID: "envelope"})
			require.NoError(t, err)
			if tc.kind == "block_suggestion" {
				require.Equal(t, "query", e.Invocation.Interaction.Query)
			}
			if tc.kind == "block_actions" {
				require.Equal(t, "U2", e.Invocation.Action.Selection["selected_user"])
			}
			if tc.kind == "message_action" {
				require.Equal(t, "Source", e.Invocation.Shortcut.Message["text"])
			}
		})
	}
}
