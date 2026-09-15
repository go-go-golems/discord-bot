package slacktransport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-go-golems/discord-bot/internal/slacktransport"
	"github.com/go-go-golems/discord-bot/pkg/slackcli"
	"github.com/go-go-golems/discord-bot/pkg/slackhost"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestShowcaseWireDoesNotReplyPong(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Inspection runs ping before showcase, as the actual CLI does.
	d, err := slackcli.Resolve(ctx, "../../examples/slack-bots", "ui-showcase", time.Second)
	require.NoError(t, err)
	var mu sync.Mutex
	var replies []map[string]any
	completed := make(chan struct{})
	mux := http.NewServeMux()
	var server *httptest.Server
	mux.HandleFunc("/api/auth.test", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true,"team_id":"T","user_id":"BOT"}`))
	})
	mux.HandleFunc("/api/apps.connections.open", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "url": "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"})
	})
	mux.HandleFunc("/reply", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			return
		}
		mu.Lock()
		replies = append(replies, body)
		mu.Unlock()
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}).Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		_ = conn.SetReadDeadline(time.Now().Add(4 * time.Second))
		_ = conn.WriteJSON(map[string]any{"type": "hello"})
		for _, command := range []string{"/golem-ping", "/ui-showcase"} {
			if err := conn.WriteJSON(map[string]any{
				"type": "slash_commands", "envelope_id": command, "accepts_response_payload": true,
				"payload": map[string]any{"is_enterprise_install": false, "team_id": "T", "api_app_id": "A", "channel_id": "C", "user_id": "U", "command": command, "response_url": server.URL + "/reply"},
			}); err != nil {
				t.Error(err)
				return
			}
			var ack map[string]any
			if err := conn.ReadJSON(&ack); err != nil {
				t.Error(err)
				return
			}
			if ack["envelope_id"] != command || ack["payload"] != nil {
				t.Errorf("unexpected ACK: %v", ack)
				return
			}
		}
		close(completed)
		_, _, _ = conn.ReadMessage()
	})
	server = httptest.NewServer(mux)
	defer server.Close()
	c, err := slacktransport.NewLocal(slacktransport.LocalOptions{APIURL: server.URL + "/api/", BotToken: "test-bot", AppToken: "test-app", TeamID: "T", AppID: "A"})
	require.NoError(t, err)
	defer c.Close()
	h, err := slackhost.Load(ctx, d.ScriptPath, slackhost.Options{Messages: c, Views: c})
	require.NoError(t, err)
	defer h.Close(context.Background())
	var group errgroup.Group
	group.Go(func() error { return c.Run(ctx, h) })
	defer func() { cancel(); require.NoError(t, group.Wait()) }()
	select {
	case <-completed:
	case <-ctx.Done():
		t.Fatal("missing socket acknowledgments")
	}
	require.Eventually(t, func() bool { mu.Lock(); defer mu.Unlock(); return len(replies) > 0 }, time.Second, time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	require.Len(t, replies, 1)
	require.Equal(t, "UI showcase: choose an action", replies[0]["text"])
	require.NotEmpty(t, replies[0]["blocks"])
}
