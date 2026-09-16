package slackcli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-go-golems/discord-bot/internal/slackconfig"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type appRoundTrip func(*http.Request) (*http.Response, error)

var _ http.RoundTripper = appRoundTrip(nil)

func (f appRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func appResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: http.Header{}}
}

func TestCreateAppFromProfileSavesApp(t *testing.T) {
	dir := t.TempDir()
	store := slackconfig.New(dir)
	cfg := slackconfig.Config{DefaultProfile: "dev", Profiles: map[string]slackconfig.Profile{"dev": {Management: "owner"}}}
	creds := slackconfig.Credentials{Management: map[string]slackconfig.ManagementCredential{"owner": {AccessToken: "profile-access", RefreshToken: "profile-refresh"}}}
	require.NoError(t, store.Save(cfg, creds))
	client := &http.Client{Transport: appRoundTrip(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "Bearer profile-access", r.Header.Get("Authorization"))
		return appResponse(200, `{"ok":true,"app_id":"A_PROFILE","credentials":{"client_id":"C1","client_secret":"S1"}}`), nil
	})}
	root, err := newBotsCommand(zerolog.Nop(), client)
	require.NoError(t, err)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"create-app", "ping", "--bot-repository", "../../examples/slack-bots", "--profile", "dev", "--config-dir", dir})
	require.NoError(t, root.ExecuteContext(context.Background()))
	require.Contains(t, out.String(), "A_PROFILE")
	after, appCreds, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, "dev", after.Profiles["dev"].App)
	require.Equal(t, "A_PROFILE", after.Apps["dev"].AppID)
	require.Equal(t, "S1", appCreds.Apps["dev"].ClientSecret)
}

func executeApp(t *testing.T, client *http.Client, tokenPath string, extra ...string) (string, error) {
	t.Helper()
	root, err := newBotsCommand(zerolog.Nop(), client)
	require.NoError(t, err)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SilenceErrors = true
	root.SilenceUsage = true
	args := []string{"create-app", "ping", "--bot-repository", "../../examples/slack-bots", "--config-token-file", tokenPath}
	root.SetArgs(append(args, extra...))
	err = root.ExecuteContext(context.Background())
	return out.String(), err
}

func testToken(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "token")
	require.NoError(t, os.WriteFile(path, []byte("synthetic-config-secret\n"), 0600))
	return path
}

func TestCreateAppWireAndPrivateCredentials(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: appRoundTrip(func(r *http.Request) (*http.Response, error) {
		calls++
		require.Equal(t, createAppURL, r.URL.String())
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "Bearer synthetic-config-secret", r.Header.Get("Authorization"))
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var payload map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload, 1)
		d, err := Resolve(context.Background(), "../../examples/slack-bots", "ping", 0)
		require.NoError(t, err)
		expected, err := json.Marshal(Manifest(d))
		require.NoError(t, err)
		require.JSONEq(t, string(expected), payload["manifest"])
		return appResponse(200, `{"ok":true,"app_id":"A123","oauth_authorize_url":"https://slack.com/oauth/v2/authorize?client_id=123","credentials":{"client_secret":"returned-secret","signing_secret":"signing-secret"}}`), nil
	})}
	path := filepath.Join(t.TempDir(), "app.json")
	out, err := executeApp(t, client, testToken(t), "--credentials-file", path)
	require.NoError(t, err)
	require.Equal(t, 1, calls)
	require.Contains(t, out, `"app_id": "A123"`)
	require.NotContains(t, out, "returned-secret")
	require.NotContains(t, out, "synthetic-config-secret")
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Contains(t, string(b), "returned-secret")
	require.NotContains(t, string(b), "synthetic-config-secret")
	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0600), info.Mode().Perm())
	_, err = executeApp(t, client, testToken(t), "--credentials-file", path)
	require.ErrorContains(t, err, "create credentials-file")
	require.Equal(t, 1, calls, "existing credential destination must prevent another app creation")
}

func TestCreateAppFailureAndNoRetry(t *testing.T) {
	for _, tc := range []struct {
		name, body, expected string
		status               int
	}{
		{"auth", `{"ok":false,"error":"token_expired"}`, "token_expired", 200},
		{"manifest", `{"ok":false,"error":"invalid_manifest","errors":[{"pointer":"/settings","message":"invalid synthetic-config-secret"}]}`, "/settings: invalid [redacted]", 200},
		{"malformed", `synthetic-config-secret`, "outcome unknown", 200},
		{"missing-id", `{"ok":true}`, "missing app_id", 200},
		{"rate", ``, "rate limited", 429},
		{"server", `synthetic-config-secret`, "HTTP 500", 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := &http.Client{Transport: appRoundTrip(func(*http.Request) (*http.Response, error) {
				calls++
				return appResponse(tc.status, tc.body), nil
			})}
			path := filepath.Join(t.TempDir(), "credentials.json")
			out, err := executeApp(t, client, testToken(t), "--credentials-file", path)
			require.NoFileExists(t, path)
			require.ErrorContains(t, err, tc.expected)
			require.NotContains(t, out+err.Error(), "synthetic-config-secret")
			require.Equal(t, 1, calls)
		})
	}
}

func TestCreateAppCancellationAndRedirect(t *testing.T) {
	t.Run("deadline", func(t *testing.T) {
		client := &http.Client{Transport: appRoundTrip(func(r *http.Request) (*http.Response, error) {
			<-r.Context().Done()
			return nil, errors.New("synthetic-config-secret")
		})}
		_, err := executeApp(t, client, testToken(t), "--timeout-ms", "100")
		require.ErrorContains(t, err, "outcome unknown")
		require.NotContains(t, err.Error(), "synthetic-config-secret")
	})
	t.Run("redirect", func(t *testing.T) {
		client := appHTTPClient()
		calls := 0
		client.Transport = appRoundTrip(func(*http.Request) (*http.Response, error) {
			calls++
			response := appResponse(307, "")
			response.Header.Set("Location", "https://example.com/steal")
			return response, nil
		})
		_, err := executeApp(t, client, testToken(t))
		require.Error(t, err)
		require.Equal(t, 1, calls)
	})
}

func TestCreateAppRequiresToken(t *testing.T) {
	client := &http.Client{Transport: appRoundTrip(func(*http.Request) (*http.Response, error) {
		t.Fatal("missing token must not contact Slack")
		return nil, errors.New("unreachable")
	})}
	_, err := executeApp(t, client, filepath.Join(t.TempDir(), "missing"))
	require.ErrorContains(t, err, "config-token-file")
}

func TestCreateAppCanRetryWithSameCredentialsPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	calls := 0
	client := &http.Client{Transport: appRoundTrip(func(*http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return appResponse(200, `{"ok":false,"error":"invalid_manifest"}`), nil
		}
		return appResponse(200, `{"ok":true,"app_id":"A_RETRY","credentials":{"client_secret":"synthetic-secret"}}`), nil
	})}
	token := testToken(t)
	_, err := executeApp(t, client, token, "--credentials-file", path)
	require.ErrorContains(t, err, "invalid_manifest")
	require.NoFileExists(t, path)
	_, err = executeApp(t, client, token, "--credentials-file", path)
	require.NoError(t, err)
	require.FileExists(t, path)
	require.Equal(t, 2, calls)
}

type failingAppOutput struct{}

var _ io.Writer = failingAppOutput{}

func (failingAppOutput) Write([]byte) (int, error) { return 0, errors.New("output unavailable") }

func TestCreateAppKeepsCredentialsAfterOutputFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	client := &http.Client{Transport: appRoundTrip(func(*http.Request) (*http.Response, error) {
		return appResponse(200, `{"ok":true,"app_id":"A_SAVED","credentials":{"client_secret":"synthetic-secret"}}`), nil
	})}
	c := &command{appClient: client}
	err := c.createApp(context.Background(), slackbot.Descriptor{Name: "test"}, settings{ConfigTokenFile: testToken(t), CredentialsFile: path}, failingAppOutput{})
	require.ErrorContains(t, err, "output failed; do not recreate")
	content, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	require.Contains(t, string(content), "A_SAVED")
	require.Contains(t, string(content), "synthetic-secret")
}
