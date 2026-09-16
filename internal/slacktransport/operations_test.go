package slacktransport

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestOperationsKeepPaginationAndErrors(t *testing.T) {
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		require.Equal(t, "Bearer synthetic-bot", r.Header.Get("Authorization"))
		switch r.URL.Path {
		case "/api/conversations.history":
			require.Equal(t, "cursor-1", r.Form.Get("cursor"))
			fmt.Fprint(w, `{"ok":true,"messages":[{"ts":"1.000001","text":"hello"}],"response_metadata":{"next_cursor":"cursor-2"}}`)
		case "/api/chat.delete":
			fmt.Fprint(w, `{"ok":false,"error":"cant_delete_message"}`)
		default:
			t.Fatalf("unexpected method %s", r.URL.Path)
		}
	})
	result, err := c.Call(context.Background(), "conversations.history", map[string]any{"channel": "C", "cursor": "cursor-1"})
	require.NoError(t, err)
	require.Equal(t, "cursor-2", result["response_metadata"].(map[string]any)["next_cursor"])
	_, err = c.Call(context.Background(), "messages.delete", map[string]any{"channel": "C", "ts": "1.000001"})
	require.ErrorContains(t, err, "cant_delete_message")
	_, err = c.Call(context.Background(), "conversations.history", map[string]any{"token": "script-token"})
	require.ErrorContains(t, err, "host-owned")
	_, err = c.Call(context.Background(), "arbitrary.api", map[string]any{})
	require.ErrorContains(t, err, "unknown Slack operation")
}

func TestExternalUploadSequence(t *testing.T) {
	var base string
	calls := []string{}
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.URL.Path)
		if r.URL.Path == "/upload" {
			require.Empty(t, r.Header.Get("Authorization"))
			fmt.Fprint(w, "ok")
			return
		}
		require.NoError(t, r.ParseForm())
		switch r.URL.Path {
		case "/api/files.getUploadURLExternal":
			fmt.Fprintf(w, `{"ok":true,"file_id":"F1","upload_url":%q}`, base+"/upload")
		case "/api/files.completeUploadExternal":
			require.Equal(t, "C", r.Form.Get("channel_id"))
			require.Contains(t, r.Form.Get("files"), "F1")
			fmt.Fprint(w, `{"ok":true,"files":[{"id":"F1"}]}`)
		default:
			t.Fatalf("unexpected %s", r.URL.Path)
		}
	})
	base = c.origin.Scheme + "://" + c.origin.Host
	_, err := c.Call(context.Background(), "files.upload", map[string]any{"filename": "report.md", "channel_id": "C", "content": "Report"})
	require.NoError(t, err)
	require.Equal(t, []string{"/api/files.getUploadURLExternal", "/upload", "/api/files.completeUploadExternal"}, calls)
}
