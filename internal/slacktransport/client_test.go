package slacktransport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/stretchr/testify/require"
)

func localClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := NewLocal(LocalOptions{APIURL: server.URL + "/api/", BotToken: "synthetic-bot", AppToken: "synthetic-app", TeamID: "T", AppID: "A"})
	require.NoError(t, err)
	t.Cleanup(client.Close)
	return client
}
func TestPostWireAndValidation(t *testing.T) {
	var count atomic.Int32
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		require.Equal(t, "/api/chat.postMessage", r.URL.Path)
		require.NoError(t, r.ParseForm())
		require.Equal(t, "synthetic-bot", r.Form.Get("token"))
		require.Equal(t, "C", r.Form.Get("channel"))
		require.Equal(t, "1234567890.000001", r.Form.Get("thread_ts"))
		require.Equal(t, "hello", r.Form.Get("text"))
		_, _ = w.Write([]byte(`{"ok":true,"channel":"C","ts":"1234567890.000002"}`))
	})
	_, err := c.Post(context.Background(), slackbot.PostMessage{ChannelID: "C", Text: " "})
	require.Error(t, err)
	require.Zero(t, count.Load())
	ref, err := c.Post(context.Background(), slackbot.PostMessage{ChannelID: "C", Text: "hello", ThreadTS: "1234567890.000001"})
	require.NoError(t, err)
	require.Equal(t, "1234567890.000002", ref.TS)
	require.EqualValues(t, 1, count.Load())
}
func TestPostErrorsDoNotRetryOrLeak(t *testing.T) {
	for _, test := range []struct {
		name, body, code string
		status           int
	}{
		{"rejection", `{"ok":false,"error":"not_in_channel"}`, "not_in_channel", 200},
		{"unknown", `{"ok":false,"error":"synthetic-secret-url"}`, "delivery_unknown", 200},
		{"malformed", `not-json-synthetic-secret`, "delivery_unknown", 200},
		{"rate", ``, "rate_limited", 429},
	} {
		t.Run(test.name, func(t *testing.T) {
			var count atomic.Int32
			c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
				count.Add(1)
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.body))
			})
			_, err := c.Post(context.Background(), slackbot.PostMessage{ChannelID: "C", Text: "hello"})
			var domain *slackbot.Error
			require.ErrorAs(t, err, &domain)
			require.Equal(t, test.code, domain.Code)
			require.NotContains(t, err.Error(), "synthetic-secret")
			require.EqualValues(t, 1, count.Load())
		})
	}
}
func TestPrivateResponseAndEndpointGuard(t *testing.T) {
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/response", r.URL.Path)
		require.Empty(t, r.Header.Get("Authorization"))
		var body map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, map[string]string{"text": "private", "response_type": "ephemeral"}, body)
		w.WriteHeader(http.StatusOK)
	})
	for _, target := range []string{"https://slack.com/response", "http://127.0.0.1:1/response", "http://user@" + c.origin.Host + "/response"} {
		_, err := c.responseCapability(target)
		require.Error(t, err)
	}
	reply, err := c.responseCapability("http://" + c.origin.Host + "/response")
	require.NoError(t, err)
	require.NoError(t, reply.Reply(context.Background(), slackbot.Text{Text: "private"}))
}
func TestConnectionTokenRouting(t *testing.T) {
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/apps.connections.open", r.URL.Path)
		require.Equal(t, "Bearer synthetic-app", r.Header.Get("Authorization"))
		_, _ = w.Write([]byte(`{"ok":true,"url":"ws://127.0.0.1:1/socket"}`))
	})
	_, target, err := c.socket.OpenContext(context.Background())
	require.NoError(t, err)
	// The actual configured dialer must reject a returned URL targeting another port.
	require.Equal(t, "ws://127.0.0.1:1/socket", target)
	_, err = c.transport.DialContext(context.Background(), "tcp", "127.0.0.1:1")
	require.ErrorContains(t, err, "unexpected destination")
}

func TestAcceptedThenLostResponseIsNotRetried(t *testing.T) {
	var count atomic.Int32
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		conn, _, err := w.(http.Hijacker).Hijack()
		require.NoError(t, err)
		_ = conn.Close()
	})
	_, err := c.Post(context.Background(), slackbot.PostMessage{ChannelID: "C", Text: "accepted remotely"})
	var domain *slackbot.Error
	require.ErrorAs(t, err, &domain)
	require.Equal(t, "delivery_unknown", domain.Code)
	require.EqualValues(t, 1, count.Load())
}
