package slackcli

import (
	"bytes"
	"context"
	"testing"

	"github.com/go-go-golems/discord-bot/internal/slackconfig"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestRunRemoteRequiresStoredRuntimeCredentials(t *testing.T) {
	dir := t.TempDir()
	store := slackconfig.New(dir)
	require.NoError(t, store.Save(slackconfig.Config{
		DefaultProfile: "dev",
		Profiles:       map[string]slackconfig.Profile{"dev": {Management: "owner", App: "ping", Installation: "dev-T1"}},
		Apps:           map[string]slackconfig.App{"ping": {AppID: "A1"}},
		Installations:  map[string]slackconfig.Installation{"dev-T1": {App: "ping", TeamID: "T1"}},
	}, slackconfig.Credentials{Management: map[string]slackconfig.ManagementCredential{"owner": {AccessToken: "management"}}}))
	root, err := newBotsCommand(zerolog.Nop(), appHTTPClient())
	require.NoError(t, err)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"run", "ping", "--config-dir", dir, "--bot-repository", "../../examples/slack-bots"})
	root.SilenceErrors = true
	root.SilenceUsage = true
	err = root.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "no bot token")
}
