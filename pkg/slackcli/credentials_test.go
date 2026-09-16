package slackcli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/discord-bot/internal/slackconfig"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestImportStatusAndRefresh(t *testing.T) {
	dir := t.TempDir()
	access := filepath.Join(dir, "a")
	refresh := filepath.Join(dir, "r")
	require.NoError(t, os.WriteFile(access, []byte("old-access\n"), 0600))
	require.NoError(t, os.WriteFile(refresh, []byte("old-refresh\n"), 0600))
	root := NewCredentialsCommand(zerolog.Nop(), &http.Client{Transport: appRoundTrip(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "https://slack.com/api/tooling.tokens.rotate", r.URL.String())
		b, _ := io.ReadAll(r.Body)
		vals, _ := url.ParseQuery(string(b))
		require.Equal(t, "old-refresh", vals.Get("refresh_token"))
		return appResponse(200, `{"ok":true,"token":"new-access","refresh_token":"new-refresh","team_id":"T1","user_id":"U1","exp":2000000000}`), nil
	})})
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"import-management", "--config-dir", dir, "--profile", "dev", "--management", "owner", "--access-token-file", access, "--refresh-token-file", refresh})
	require.NoError(t, root.ExecuteContext(context.Background()))
	root.SetArgs([]string{"refresh", "--config-dir", dir, "--profile", "dev"})
	require.NoError(t, root.ExecuteContext(context.Background()))
	require.NotContains(t, out.String(), "new-access")
	_, cr, err := slackconfig.New(dir).Load()
	require.NoError(t, err)
	require.Equal(t, "new-access", cr.Management["owner"].AccessToken)
	require.Equal(t, "new-refresh", cr.Management["owner"].RefreshToken)
	var status bytes.Buffer
	root.SetOut(&status)
	root.SetArgs([]string{"status", "--config-dir", dir, "--profile", "dev"})
	require.NoError(t, root.ExecuteContext(context.Background()))
	require.NotContains(t, status.String(), "new-access")
	var safe map[string]any
	require.NoError(t, json.Unmarshal(status.Bytes(), &safe))
	require.Equal(t, true, safe["has_access_token"])
}

func TestImportRuntimeTokens(t *testing.T) {
	dir := t.TempDir()
	bot := filepath.Join(dir, "bot")
	app := filepath.Join(dir, "app")
	user := filepath.Join(dir, "user")
	require.NoError(t, os.WriteFile(user, []byte("user-secret"), 0600))
	require.NoError(t, os.WriteFile(bot, []byte("bot-secret\n"), 0600))
	require.NoError(t, os.WriteFile(app, []byte("app-secret\n"), 0600))
	store := slackconfig.New(dir)
	require.NoError(t, store.Save(slackconfig.Config{DefaultProfile: "dev", Profiles: map[string]slackconfig.Profile{"dev": {App: "app"}}}, slackconfig.Credentials{}))
	root := NewCredentialsCommand(zerolog.Nop(), nil)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"import-runtime", "--config-dir", dir, "--profile", "dev", "--installation", "test", "--team-id", "T1", "--bot-token-file", bot, "--app-token-file", app, "--user-token-file", user})
	require.NoError(t, root.ExecuteContext(context.Background()))
	_, cr, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, "bot-secret", cr.Installations["test"].BotToken)
	require.Equal(t, "user-secret", cr.Installations["test"].UserToken)
	require.Equal(t, "app-secret", cr.Apps["app"].AppToken)
	require.NotContains(t, out.String(), "secret")
}
