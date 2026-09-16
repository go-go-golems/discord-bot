package slackcli

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-go-golems/discord-bot/internal/slackconfig"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestUpdateManifest(t *testing.T) {
	d, err := Resolve(context.Background(), "../../examples/slack-bots", "ui-showcase", 0)
	require.NoError(t, err)
	for _, tc := range []struct {
		name, body, errorText string
		permissions           bool
	}{
		{"success", `{"ok":true,"app_id":"A1","permissions_updated":false}`, "", false},
		{"reinstall", `{"ok":true,"app_id":"A1","permissions_updated":true}`, "", true},
		{"expired", `{"ok":false,"error":"token_expired"}`, "credentials refresh", false},
		{"rejected", `{"ok":false,"error":"invalid_manifest"}`, "invalid_manifest", false},
		{"wrong app", `{"ok":true,"app_id":"A2"}`, "unexpected app ID", false},
		{"invalid response", `not json`, "invalid JSON", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &http.Client{Transport: appRoundTrip(func(r *http.Request) (*http.Response, error) {
				require.Equal(t, "https://slack.com/api/apps.manifest.update", r.URL.String())
				require.Equal(t, "Bearer test-management-secret", r.Header.Get("Authorization"))
				var args map[string]string
				require.NoError(t, json.NewDecoder(r.Body).Decode(&args))
				require.Equal(t, "A1", args["app_id"])
				expected, err := json.Marshal(Manifest(d))
				require.NoError(t, err)
				require.JSONEq(t, string(expected), args["manifest"])
				require.Contains(t, args["manifest"], "/ui-showcase")
				return appResponse(200, tc.body), nil
			})}
			c := &command{appClient: client}
			permissions, err := c.updateManifest(context.Background(), d, settings{TimeoutMS: 5000}, "A1", "test-management-secret")
			if tc.errorText != "" {
				require.ErrorContains(t, err, tc.errorText)
				require.NotContains(t, err.Error(), "test-management-secret")
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.permissions, permissions)
			}
		})
	}
}

func TestRunUpdatesManifestByDefaultBeforeConnecting(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, slackconfig.New(dir).Save(slackconfig.Config{
		DefaultProfile: "dev",
		Profiles:       map[string]slackconfig.Profile{"dev": {Management: "owner", App: "app", Installation: "installed"}},
		Apps:           map[string]slackconfig.App{"app": {AppID: "A1"}},
		Installations:  map[string]slackconfig.Installation{"installed": {App: "app", TeamID: "T1"}},
	}, slackconfig.Credentials{
		Management:    map[string]slackconfig.ManagementCredential{"owner": {AccessToken: "management-secret"}},
		Apps:          map[string]slackconfig.AppCredential{"app": {AppToken: "xapp-test"}},
		Installations: map[string]slackconfig.InstallationCredential{"installed": {BotToken: "xoxb-test"}},
	}))
	calls := 0
	client := &http.Client{Transport: appRoundTrip(func(r *http.Request) (*http.Response, error) {
		calls++
		require.Equal(t, "/api/apps.manifest.update", r.URL.Path)
		return appResponse(200, `{"ok":true,"app_id":"A1","permissions_updated":true}`), nil
	})}
	root, err := newBotsCommand(zerolog.Nop(), client)
	require.NoError(t, err)
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"run", "ui-showcase", "--config-dir", dir, "--bot-repository", "../../examples/slack-bots"})
	err = root.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "bots install ui-showcase --profile dev --team-id T1")
	require.Equal(t, 1, calls)
	run, _, err := root.Find([]string{"run"})
	require.NoError(t, err)
	require.Equal(t, "false", run.Flags().Lookup("skip-manifest-update").DefValue)

}
