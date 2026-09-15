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

func TestInstallAppSavesRuntimeTokens(t *testing.T) {
	dir := t.TempDir()
	store := slackconfig.New(dir)
	cfg := slackconfig.Config{
		DefaultProfile: "dev",
		Profiles:       map[string]slackconfig.Profile{"dev": {Management: "owner", App: "ping"}},
		Apps:           map[string]slackconfig.App{"ping": {AppID: "A_INSTALL"}},
	}
	creds := slackconfig.Credentials{Management: map[string]slackconfig.ManagementCredential{"owner": {AccessToken: "management-secret"}}}
	require.NoError(t, store.Save(cfg, creds))
	calls := 0
	client := &http.Client{Transport: appRoundTrip(func(r *http.Request) (*http.Response, error) {
		calls++
		require.Equal(t, developerInstallURL, r.URL.String())
		require.Equal(t, "Bearer management-secret", r.Header.Get("Authorization"))
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var payload struct {
			AppID           string   `json:"app_id"`
			BotScopes       []string `json:"bot_scopes"`
			OutgoingDomains []string `json:"outgoing_domains"`
			TeamID          string   `json:"team_id"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "A_INSTALL", payload.AppID)
		require.Equal(t, []string{"chat:write", "commands", "app_mentions:read"}, payload.BotScopes)
		require.Empty(t, payload.OutgoingDomains)
		require.Equal(t, "T_INSTALL", payload.TeamID)
		return appResponse(http.StatusOK, `{"ok":true,"app_id":"A_INSTALL","api_access_tokens":{"bot":"xoxb-secret","app_level":"xapp-secret"}}`), nil
	})}
	root, err := newBotsCommand(zerolog.Nop(), client)
	require.NoError(t, err)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"install", "ping", "--profile", "dev", "--team-id", "T_INSTALL", "--config-dir", dir, "--bot-repository", "../../examples/slack-bots"})
	require.NoError(t, root.ExecuteContext(context.Background()))
	require.Equal(t, 1, calls)
	require.Contains(t, out.String(), "dev-T_INSTALL")
	require.NotContains(t, out.String(), "xoxb-secret")
	require.NotContains(t, out.String(), "xapp-secret")
	after, afterCreds, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, "dev-T_INSTALL", after.Profiles["dev"].Installation)
	require.Equal(t, slackconfig.Installation{App: "ping", TeamID: "T_INSTALL"}, after.Installations["dev-T_INSTALL"])
	require.Equal(t, "xapp-secret", afterCreds.Apps["ping"].AppToken)
	require.Equal(t, "xoxb-secret", afterCreds.Installations["dev-T_INSTALL"].BotToken)
}

func TestInstallAppRejectsInvalidSetupAndResponse(t *testing.T) {
	for _, tc := range []struct {
		name     string
		args     []string
		body     string
		status   int
		expected string
		calls    int
	}{
		{name: "missing profile", args: []string{"install", "ping", "--team-id", "T1", "--config-dir", "DIR"}, expected: "no profile selected", calls: 0},
		{name: "missing team", args: []string{"install", "ping", "--profile", "dev", "--config-dir", "DIR"}, expected: "--team-id is required", calls: 0},
		{name: "slack error", args: []string{"install", "ping", "--profile", "dev", "--team-id", "T1", "--config-dir", "DIR"}, body: `{"ok":false,"error":"not_allowed"}`, status: 200, expected: "not_allowed", calls: 1},
		{name: "mismatch", args: []string{"install", "ping", "--profile", "dev", "--team-id", "T1", "--config-dir", "DIR"}, body: `{"ok":true,"app_id":"A_OTHER","api_access_tokens":{"bot":"xoxb-secret","app_level":"xapp-secret"}}`, status: 200, expected: "expected", calls: 1},
		{name: "missing tokens", args: []string{"install", "ping", "--profile", "dev", "--team-id", "T1", "--config-dir", "DIR"}, body: `{"ok":true,"app_id":"A_INSTALL"}`, status: 200, expected: "no bot and app-level", calls: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			store := slackconfig.New(dir)
			require.NoError(t, store.Save(slackconfig.Config{Profiles: map[string]slackconfig.Profile{"dev": {Management: "owner", App: "ping"}}, Apps: map[string]slackconfig.App{"ping": {AppID: "A_INSTALL"}}}, slackconfig.Credentials{Management: map[string]slackconfig.ManagementCredential{"owner": {AccessToken: "management-secret"}}}))
			calls := 0
			client := &http.Client{Transport: appRoundTrip(func(*http.Request) (*http.Response, error) {
				calls++
				return appResponse(tc.status, tc.body), nil
			})}
			root, err := newBotsCommand(zerolog.Nop(), client)
			require.NoError(t, err)
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&out)
			root.SilenceErrors = true
			root.SilenceUsage = true
			args := append([]string{}, tc.args...)
			for i := range args {
				if args[i] == "DIR" {
					args[i] = dir
				}
			}
			args = append(args, "--bot-repository", "../../examples/slack-bots")
			root.SetArgs(args)
			err = root.ExecuteContext(context.Background())
			require.ErrorContains(t, err, tc.expected)
			require.Equal(t, tc.calls, calls)
			require.NotContains(t, out.String()+err.Error(), "management-secret")
			require.NotContains(t, out.String()+err.Error(), "xoxb-secret")
		})
	}
}

func TestInstallAppDoesNotOverwriteMismatchedInstallation(t *testing.T) {
	dir := t.TempDir()
	store := slackconfig.New(dir)
	cfg := slackconfig.Config{Profiles: map[string]slackconfig.Profile{"dev": {Management: "owner", App: "ping"}}, Apps: map[string]slackconfig.App{"ping": {AppID: "A_INSTALL"}}, Installations: map[string]slackconfig.Installation{"dev-T1": {App: "other", TeamID: "T1"}}}
	creds := slackconfig.Credentials{Management: map[string]slackconfig.ManagementCredential{"owner": {AccessToken: "management-secret"}}}
	require.NoError(t, store.Save(cfg, creds))
	calls := 0
	client := &http.Client{Transport: appRoundTrip(func(*http.Request) (*http.Response, error) { calls++; return appResponse(200, `{"ok":true}`), nil })}
	root, err := newBotsCommand(zerolog.Nop(), client)
	require.NoError(t, err)
	root.SetArgs([]string{"install", "ping", "--profile", "dev", "--team-id", "T1", "--config-dir", dir, "--bot-repository", "../../examples/slack-bots"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	err = root.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "another app")
	require.Equal(t, 0, calls)
}
