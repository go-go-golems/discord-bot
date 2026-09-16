package slackcli

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/go-go-golems/discord-bot/pkg/slackhost"
	"github.com/stretchr/testify/require"
)

func portHost(t *testing.T, name string, config map[string]any, rec *slackbot.Recorder, ops slackbot.OperationService) *slackhost.Host {
	t.Helper()
	d, err := Resolve(context.Background(), "../../examples/slack-bots", name, 0)
	require.NoError(t, err)
	h, err := slackhost.Load(context.Background(), d.ScriptPath, slackhost.Options{Messages: rec, Views: rec, Operations: ops, Config: config})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, h.Close(context.Background())) })
	return h
}
func portCommand(h *slackhost.Host, rec *slackbot.Recorder, team, user, name, text string) error {
	return h.Dispatch(context.Background(), slackbot.Invocation{TeamID: team, ChannelID: "C", UserID: user, Command: name, Text: text}, rec)
}
func TestCustomKBPersistsAndSeparatesWorkspaces(t *testing.T) {
	config := map[string]any{"dbPath": filepath.Join(t.TempDir(), "links.sqlite")}
	rec := &slackbot.Recorder{}
	h := portHost(t, "custom-kb", config, rec, rec)
	require.NoError(t, portCommand(h, rec, "T", "U", "/kb-link", "https://example.test/docs | Example docs | Useful documentation | docs"))
	require.NoError(t, h.Close(context.Background()))
	h = portHost(t, "custom-kb", config, rec, rec)
	require.NoError(t, portCommand(h, rec, "T", "U", "/kb-list", ""))
	require.Contains(t, rec.Operations()[1].Reply.Text, "Found 1")
	require.NoError(t, portCommand(h, rec, "OTHER", "U", "/kb-list", ""))
	require.Contains(t, rec.Operations()[2].Reply.Text, "No links")
	require.NoError(t, portCommand(h, rec, "T", "U", "/kb-link", "https://example.test/docs | Updated title | Updated summary | docs"))
	require.NoError(t, portCommand(h, rec, "T", "U", "/kb-list", ""))
	require.Contains(t, rec.Operations()[4].Reply.Text, "Found 1")
}
func TestKnowledgeReviewAuthorization(t *testing.T) {
	rec := &slackbot.Recorder{}
	h := portHost(t, "knowledge-base", map[string]any{"dbPath": filepath.Join(t.TempDir(), "knowledge.sqlite"), "trustedReviewerIds": "REVIEWER"}, rec, rec)
	require.NoError(t, portCommand(h, rec, "T", "U", "/kb-reject", "getting-started"))
	require.Contains(t, rec.Operations()[0].Reply.Text, "not an authorized")
	require.NoError(t, portCommand(h, rec, "T", "REVIEWER", "/kb-stale", "getting-started"))
	require.Contains(t, rec.Operations()[1].Reply.Text, "[stale]")
	require.NoError(t, h.Close(context.Background()))
}
func TestShowManagementAndPinReference(t *testing.T) {
	rec := &slackbot.Recorder{}
	h := portHost(t, "show-space", map[string]any{"dbPath": filepath.Join(t.TempDir(), "shows.sqlite"), "managerIds": "MANAGER", "seedShows": false}, rec, rec)
	require.NoError(t, portCommand(h, rec, "T", "OTHER", "/add-show", ""))
	require.Len(t, rec.Operations(), 1)
	require.Contains(t, rec.Operations()[0].Reply.Text, "authorize")
	values := map[string]map[string]any{}
	for k, v := range map[string]string{"artist": "Band", "date": "2027-01-02", "doors": "7pm", "age": "All ages", "price": "10", "notes": ""} {
		values[k] = map[string]any{"value": map[string]any{"value": v}}
	}
	require.NoError(t, h.Dispatch(context.Background(), slackbot.Invocation{TeamID: "T", UserID: "MANAGER", Interaction: &slackbot.Interaction{Type: "view_submission", CallbackID: "show.save", PrivateMetadata: "C", Values: values, Ack: rec}}, nil))
	ops := rec.Operations()
	require.Equal(t, "ack", ops[1].Kind)
	require.Equal(t, "post", ops[2].Kind)
	require.Equal(t, "pins.add", ops[3].Kind)
	require.Equal(t, ops[2].Ref.TS, ops[3].Params["timestamp"])
	require.Equal(t, ops[2].Ref.ChannelID, ops[3].Params["channel"])
}

type archiveFixture struct {
	*slackbot.Recorder
	pages int
}

func (f *archiveFixture) Call(ctx context.Context, operation string, params map[string]any) (map[string]any, error) {
	_, err := f.Recorder.Call(ctx, operation, params)
	if err != nil {
		return nil, err
	}
	if operation == "conversations.info" {
		return map[string]any{"channel": map[string]any{"name": "test"}}, nil
	}
	if operation == "conversations.history" {
		f.pages++
		if f.pages == 1 {
			return map[string]any{"messages": []any{map[string]any{"ts": "2.000001", "text": "newer", "files": []any{map[string]any{"name": "attachment.txt", "permalink": "https://example.test/file"}}}}, "response_metadata": map[string]any{"next_cursor": "next"}}, nil
		}
		return map[string]any{"messages": []any{map[string]any{"ts": "1.000001", "text": "older"}}, "response_metadata": map[string]any{"next_cursor": ""}}, nil
	}
	return map[string]any{"ok": true}, nil
}
func TestArchivePaginatesAndPreservesAttachments(t *testing.T) {
	rec := &slackbot.Recorder{}
	fixture := &archiveFixture{Recorder: rec}
	h := portHost(t, "archive-helper", nil, rec, fixture)
	require.NoError(t, portCommand(h, rec, "T", "U", "/archive-channel", "500"))
	require.Equal(t, 2, fixture.pages)
	ops := rec.Operations()
	require.Equal(t, "next", ops[2].Params["cursor"])
	require.Equal(t, "files.upload", ops[3].Kind)
	content := ops[3].Params["content"].(string)
	require.Contains(t, content, "older")
	require.Contains(t, content, "newer")
	require.Contains(t, content, "https://example.test/file")
	require.Contains(t, ops[4].Reply.Text, "Archived 2")
}
