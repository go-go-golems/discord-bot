package slacktransport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

type blockedDispatch struct {
	started chan struct{}
	count   atomic.Int32
}

var _ slackbot.Dispatcher = (*blockedDispatch)(nil)

func (d *blockedDispatch) Dispatch(ctx context.Context, _ slackbot.Invocation, _ slackbot.Responder) error {
	if d.count.Add(1) == 1 {
		close(d.started)
	}
	<-ctx.Done()
	return ctx.Err()
}
func TestSocketAckWhileHandlerBlockedAndEventRedelivered(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dispatch := &blockedDispatch{started: make(chan struct{})}
	observed := make(chan struct{})
	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth.test", func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		require.Equal(t, "synthetic-bot", r.Form.Get("token"))
		_, _ = w.Write([]byte(`{"ok":true,"team_id":"T","user_id":"BOT"}`))
	})
	mux.HandleFunc("/api/apps.connections.open", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer synthetic-app", r.Header.Get("Authorization"))
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "url": "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"})
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return r.Header.Get("Origin") == "https://api.slack.com" }}).Upgrade(w, r, nil)
		require.NoError(t, err)
		defer func() { _ = conn.Close() }()
		require.NoError(t, conn.SetReadDeadline(time.Now().Add(4*time.Second)))
		require.NoError(t, conn.WriteJSON(map[string]any{"type": "hello", "num_connections": 1}))
		for _, id := range []string{"delivery-1", "delivery-2"} {
			require.NoError(t, conn.WriteJSON(map[string]any{"type": "events_api", "envelope_id": id, "payload": map[string]any{
				"type": "event_callback", "team_id": "T", "api_app_id": "A", "event_id": "same-event",
				"event": map[string]any{"type": "app_mention", "channel": "C", "user": "HUMAN", "ts": "1234567890.000001", "text": "hello"},
			}}))
			var ack struct {
				EnvelopeID string `json:"envelope_id"`
			}
			require.NoError(t, conn.ReadJSON(&ack))
			require.Equal(t, id, ack.EnvelopeID)
		}
		close(observed)
		_, _, _ = conn.ReadMessage() // Client cancellation closes the connection.
	})
	server = httptest.NewServer(mux)
	defer server.Close()
	c, err := NewLocal(LocalOptions{APIURL: server.URL + "/api/", BotToken: "synthetic-bot", AppToken: "synthetic-app", TeamID: "T", AppID: "A"})
	require.NoError(t, err)
	defer c.Close()
	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error { return c.Run(ctx, dispatch) })
	defer func() { cancel(); require.NoError(t, group.Wait()) }()
	select {
	case <-observed:
	case <-ctx.Done():
		t.Fatal("ACKs did not arrive while handler blocked")
	}
	select {
	case <-dispatch.started:
	case <-ctx.Done():
		t.Fatal("handler did not start")
	}
	require.EqualValues(t, 1, dispatch.count.Load())
}
