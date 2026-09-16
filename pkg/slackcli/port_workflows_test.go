package slackcli

import (
	"context"
	"path/filepath"
	"strings"
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
	id := strings.TrimPrefix(ops[4].Params["text"].(string), "Posted and pinned show ")
	require.NoError(t, portCommand(h, rec, "T", "MANAGER", "/cancel-show", id))
	ops = rec.Operations()
	require.Equal(t, "messages.update", ops[5].Kind)
	require.Equal(t, "pins.remove", ops[6].Kind)
	require.NoError(t, portCommand(h, rec, "T", "OTHER", "/show", id))
	ops = rec.Operations()
	require.Contains(t, ops[len(ops)-1].Reply.Text, "[cancelled]")

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
	require.Equal(t, "files.upload", ops[4].Kind)
	content := ops[4].Params["content"].(string)
	require.Contains(t, content, "older")
	require.Contains(t, content, "newer")
	require.Contains(t, content, "https://example.test/file")
	require.Contains(t, ops[5].Reply.Text, "Archived 2")
}

func TestModerationDeniesBeforeMutation(t *testing.T) {
	rec := &slackbot.Recorder{}
	h := portHost(t, "moderation", map[string]any{"moderatorIds": "MODERATOR", "enableWorkspaceRemoval": true}, rec, rec)
	require.NoError(t, portCommand(h, rec, "T", "OTHER", "/mod-remove-workspace-user", "TARGET"))
	require.Len(t, rec.Operations(), 1)
	require.Contains(t, rec.Operations()[0].Reply.Text, "permission required")
	require.NoError(t, portCommand(h, rec, "T", "MODERATOR", "/mod-remove-workspace-user", "TARGET"))
	require.Equal(t, "admin.removeUser", rec.Operations()[1].Kind)
}

func TestShowcaseStateAndModalValidation(t *testing.T) {
	rec := &slackbot.Recorder{}
	h := portHost(t, "ui-showcase", nil, rec, rec)
	require.NoError(t, portCommand(h, rec, "T", "U", "/demo-search", "UI"))
	require.Contains(t, rec.Operations()[0].Reply.Text, "Articles")
	i := slackbot.Invocation{TeamID: "T", UserID: "U", ChannelID: "C", Action: &slackbot.Action{ActionID: "demo.article.verify", Value: "art-3"}}
	require.NoError(t, h.Dispatch(context.Background(), i, rec))
	require.Equal(t, "replace_original", rec.Operations()[1].Kind)
	require.Contains(t, rec.Operations()[1].Reply.Text, "[verified]")
	i.Action = &slackbot.Action{ActionID: "demo.pager.next", Value: "0"}
	require.NoError(t, h.Dispatch(context.Background(), i, rec))
	require.Contains(t, rec.Operations()[2].Reply.Text, "page 2")
	i.Action = nil
	i.Interaction = &slackbot.Interaction{Type: "view_submission", CallbackID: "demo.form", Values: map[string]map[string]any{"title": {"value": map[string]any{"value": "x"}}}, Ack: rec}
	require.NoError(t, h.Dispatch(context.Background(), i, nil))
	require.Equal(t, "errors", rec.Operations()[3].Ack.Kind)
	i.Interaction.Values["title"]["value"] = map[string]any{"value": "Valid title"}
	require.NoError(t, h.Dispatch(context.Background(), i, nil))
	require.Equal(t, "update", rec.Operations()[4].Ack.Kind)
}

func TestPingFeedbackAndDynamicSearch(t *testing.T) {
	rec := &slackbot.Recorder{}
	h := portHost(t, "ping", nil, rec, rec)
	i := slackbot.Invocation{TeamID: "T", UserID: "U", ChannelID: "C", Command: "/golem-feedback", TriggerID: "trigger"}
	require.NoError(t, h.Dispatch(context.Background(), i, rec))
	require.Equal(t, "open_view", rec.Operations()[0].Kind)
	i.Command = ""
	i.Interaction = &slackbot.Interaction{Type: "block_suggestion", CallbackID: "ping.search.query", Query: "arch", Ack: rec}
	require.NoError(t, h.Dispatch(context.Background(), i, nil))
	require.Len(t, rec.Operations()[1].Ack.Options, 2)
	require.Equal(t, "architecture", rec.Operations()[1].Ack.Options[0]["value"])
}

func TestHaterTargetAndApologyResult(t *testing.T) {
	rec := &slackbot.Recorder{}
	h := portHost(t, "hater", nil, rec, rec)
	require.NoError(t, portCommand(h, rec, "T", "U", "/roast", "<@TARGET|target>"))
	require.Contains(t, rec.Operations()[0].Reply.Text, "TARGET")
	i := slackbot.Invocation{TeamID: "T", UserID: "U", Interaction: &slackbot.Interaction{Type: "view_submission", CallbackID: "hater.apology.submit", Ack: rec, Values: map[string]map[string]any{"subject": {"value": map[string]any{"value": "Sorry"}}, "details": {"value": map[string]any{"value": "A sufficiently detailed apology"}}}}}
	require.NoError(t, h.Dispatch(context.Background(), i, nil))
	require.Equal(t, "update", rec.Operations()[1].Ack.Kind)
	require.NoError(t, portCommand(h, rec, "T", "U", "/apology-status", ""))
	require.Contains(t, rec.Operations()[2].Reply.Text, "Apology received")
}

type operationFixture func(context.Context, string, map[string]any) (map[string]any, error)

func (f operationFixture) Call(ctx context.Context, op string, params map[string]any) (map[string]any, error) {
	return f(ctx, op, params)
}

func TestInteractionAvatarAndQuote(t *testing.T) {
	rec := &slackbot.Recorder{}
	services := operationFixture(func(ctx context.Context, op string, params map[string]any) (map[string]any, error) {
		require.Equal(t, "users.info", op)
		require.Equal(t, "TARGET", params["user"])
		return map[string]any{"user": map[string]any{"name": "Target", "profile": map[string]any{"image_512": "https://example.test/avatar.png"}}}, nil
	})
	h := portHost(t, "interaction-types", nil, rec, services)
	i := slackbot.Invocation{TeamID: "T", UserID: "U", TriggerID: "trigger", Shortcut: &slackbot.Shortcut{Type: "message_action", CallbackID: "quote-message", Message: map[string]any{"text": "Quoted source"}}}
	require.NoError(t, h.Dispatch(context.Background(), i, nil))
	require.Contains(t, rec.Operations()[0].View.Blocks[0]["text"].(map[string]any)["text"], "Quoted source")
	i.Shortcut = nil
	i.Interaction = &slackbot.Interaction{Type: "view_submission", CallbackID: "avatar.submit", Ack: rec, Values: map[string]map[string]any{"user": {"selected": map[string]any{"selected_user": "TARGET"}}}}
	require.NoError(t, h.Dispatch(context.Background(), i, nil))
	require.Equal(t, "update", rec.Operations()[1].Ack.Kind)
	require.Equal(t, "https://example.test/avatar.png", rec.Operations()[1].Ack.View.Blocks[0]["accessory"].(map[string]any)["image_url"])
}

func TestArchiveAbortsWithoutPartialUpload(t *testing.T) {
	for _, tc := range []struct {
		name      string
		cursors   []string
		failPage  int
		wantPages int
		wantError string
	}{
		{name: "repeated cursor", cursors: []string{"A", "A"}, wantPages: 2, wantError: "repeated pagination cursor"},
		{name: "cursor cycle", cursors: []string{"A", "B", "A"}, wantPages: 3, wantError: "repeated pagination cursor"},
		{name: "page failure", cursors: []string{"A"}, failPage: 2, wantPages: 2, wantError: "rate_limited"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := &slackbot.Recorder{}
			pages := 0
			services := operationFixture(func(ctx context.Context, op string, params map[string]any) (map[string]any, error) {
				_, err := rec.Call(ctx, op, params)
				if err != nil {
					return nil, err
				}
				if op == "conversations.info" {
					return map[string]any{"channel": map[string]any{"name": "test"}}, nil
				}
				if op != "conversations.history" {
					return map[string]any{"ok": true}, nil
				}
				pages++
				if pages == tc.failPage {
					return nil, slackbot.Fail("rate_limited", op, "retry later")
				}
				if pages > len(tc.cursors) {
					return nil, slackbot.Fail("unexpected_page", op, "cursor cycle was not rejected")
				}
				return map[string]any{
					"messages":          []any{map[string]any{"ts": "1.000001", "text": "partial page"}},
					"response_metadata": map[string]any{"next_cursor": tc.cursors[pages-1]},
				}, nil
			})
			h := portHost(t, "archive-helper", nil, rec, services)
			err := portCommand(h, rec, "T", "U", "/archive-channel", "500")
			require.ErrorContains(t, err, tc.wantError)
			require.Equal(t, tc.wantPages, pages)
			for _, op := range rec.Operations() {
				require.NotEqual(t, "files.upload", op.Kind, "failed archive must not publish partial history")
				require.NotEqual(t, "reply", op.Kind, "failed archive must not report success")
			}
		})
	}
}
