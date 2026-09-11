// Package slackprobe verifies the actual SDK against a separately started local mock.
package slackprobe

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

var runDirectory = flag.String("slack-probe-directory", "", "Explicit directory from the separately started Bun probe")

// Dial only literal loopback addresses, including URLs returned by the server.
// No DNS, environment proxy selection, or external fallback is permitted.
func localDial(ctx context.Context, network, address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return nil, errors.New("probe refused non-loopback destination")
	}
	return (&net.Dialer{Timeout: time.Second}).DialContext(ctx, network, address)
}

func TestSDKInteroperability(t *testing.T) {
	if *runDirectory == "" {
		t.Skip("run the explicit local-mock gate documented in testdata/slack/mock/README.md")
	}
	var cfg struct{ APIURL, BotToken, AppToken, TeamID, UserID string }
	data, err := os.ReadFile(filepath.Join(*runDirectory, "config.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &cfg))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log := zerolog.New(os.Stderr).With().Str("test", "slack-sdk-probe").Logger()
	transport := &http.Transport{DialContext: localDial}
	defer transport.CloseIdleConnections()
	httpClient := &http.Client{Transport: transport, Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("probe refuses redirects") }}
	api := slack.New(cfg.BotToken, slack.OptionAppLevelToken(cfg.AppToken), slack.OptionAPIURL(cfg.APIURL), slack.OptionHTTPClient(httpClient))
	identity, err := api.AuthTestContext(ctx)
	require.NoError(t, err)
	require.Equal(t, cfg.TeamID, identity.TeamID)
	require.Equal(t, cfg.UserID, identity.UserID)
	socket := socketmode.New(api, socketmode.OptionDialer(&websocket.Dialer{NetDialContext: localDial, HandshakeTimeout: 5 * time.Second}))
	group, workerCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err := socket.RunContext(workerCtx)
		if workerCtx.Err() != nil {
			return nil
		}
		return err
	})
	defer func() { cancel(); require.NoError(t, group.Wait()) }()
	handled := 0
	for handled < 2 {
		select {
		case <-workerCtx.Done():
			t.Fatal("probe ended before two deliveries")
		case event := <-socket.Events:
			if event.Request == nil || event.Request.EnvelopeID == "" {
				continue
			}
			request := event.Request
			require.NoError(t, socket.AckCtx(workerCtx, request.EnvelopeID, nil))
			switch request.Type {
			case "events_api":
				var payload struct {
					Event struct{ Type, Channel, Text, TS string }
				}
				require.NoError(t, json.Unmarshal(request.Payload, &payload))
				require.Equal(t, "app_mention", payload.Event.Type)
				channel, ts, err := api.PostMessageContext(workerCtx, payload.Event.Channel, slack.MsgOptionText("probe-thread-reply", false), slack.MsgOptionTS(payload.Event.TS))
				require.NoError(t, err)
				require.Equal(t, payload.Event.Channel, channel)
				require.NotEmpty(t, ts)
			case "slash_commands":
				var payload struct {
					Command     string
					ResponseURL string `json:"response_url"`
				}
				require.NoError(t, json.Unmarshal(request.Payload, &payload))
				require.Equal(t, "/golem-ping", payload.Command)
				body := bytes.NewBufferString(`{"response_type":"ephemeral","text":"probe-private-reply"}`)
				req, err := http.NewRequestWithContext(workerCtx, http.MethodPost, payload.ResponseURL, body)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				response, err := httpClient.Do(req)
				require.NoError(t, err)
				_, err = io.Copy(io.Discard, response.Body)
				require.NoError(t, err)
				require.NoError(t, response.Body.Close())
				require.Equal(t, http.StatusOK, response.StatusCode)
			default:
				t.Fatalf("unexpected envelope type %q", request.Type)
			}
			handled++
			log.Info().Str("envelope_type", request.Type).Msg("handled local delivery")
		}
	}
	require.Eventually(t, func() bool { _, err := os.Stat(filepath.Join(*runDirectory, "result.json")); return err == nil }, 5*time.Second, 25*time.Millisecond)
	var result struct {
		ThreadPreserved, EphemeralRecipient, EmptyCommandAck bool
		PublicReplyCount                                     int
		Deliveries                                           []struct {
			Acked    bool
			Attempts int
		}
	}
	data, err = os.ReadFile(filepath.Join(*runDirectory, "result.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &result))
	require.True(t, result.ThreadPreserved)
	require.True(t, result.EphemeralRecipient)
	require.True(t, result.EmptyCommandAck)
	require.Zero(t, result.PublicReplyCount)
	require.Len(t, result.Deliveries, 2)
	for _, delivery := range result.Deliveries {
		require.True(t, delivery.Acked)
		require.Equal(t, 1, delivery.Attempts)
	}
	cancel()
	require.NoError(t, group.Wait())
	require.Eventually(t, func() bool { _, err := os.Stat(filepath.Join(*runDirectory, "shutdown.json")); return err == nil }, 5*time.Second, 25*time.Millisecond)
}

func TestLocalDialRejectsExternalDestinations(t *testing.T) {
	for _, address := range []string{"slack.com:443", "192.0.2.1:443", "[2001:db8::1]:443"} {
		_, err := localDial(context.Background(), "tcp", address)
		require.ErrorContains(t, err, "non-loopback")
	}
}
