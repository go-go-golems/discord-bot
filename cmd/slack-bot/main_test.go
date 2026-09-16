package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/stretchr/testify/require"
)

func execute(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root, err := newRoot()
	require.NoError(t, err)
	var b bytes.Buffer
	root.SetOut(&b)
	root.SetErr(&b)
	root.SetArgs(args)
	err = root.ExecuteContext(context.Background())
	return b.String(), err
}
func TestOfflineCLI(t *testing.T) {
	repo, err := filepath.Abs("../../examples/slack-bots")
	require.NoError(t, err)
	out, err := execute(t, "bots", "list", "--bot-repository", repo)
	require.NoError(t, err)
	var all []slackbot.Descriptor
	require.NoError(t, json.Unmarshal([]byte(out), &all))
	names := make([]string, 0, len(all))
	for _, d := range all {
		names = append(names, d.Name)
	}
	require.Contains(t, names, "ping")
	require.Contains(t, names, "ui-showcase")
	out, err = execute(t, "bots", "inspect", "ping", "--bot-repository", repo)
	require.NoError(t, err)
	require.Contains(t, out, "greeting")
	out, err = execute(t, "bots", "manifest", "ping", "--bot-repository", repo)
	require.NoError(t, err)
	require.Contains(t, out, "socket_mode_enabled")
	require.Contains(t, out, "app_mentions:read")
	out, err = execute(t, "bots", "simulate", "ping", "--bot-repository", repo, "--event-file", "../../examples/slack-bots/fixtures/command.json")
	require.NoError(t, err)
	require.Contains(t, out, "ephemeral_reply")
	require.Contains(t, out, "pong")
	out, err = execute(t, "bots", "simulate", "ping", "--bot-repository", repo, "--event-file", "../../examples/slack-bots/fixtures/mention.json")
	require.NoError(t, err)
	require.Contains(t, out, "1741234567.000001")
	_, err = execute(t, "bots", "inspect", "missing", "--bot-repository", repo)
	require.ErrorContains(t, err, "not found")
	_, err = execute(t, "bots", "list", "--timeout-ms", "0")
	require.ErrorContains(t, err, "timeout-ms")
}
func TestHelpAndParseErrors(t *testing.T) {
	out, err := execute(t, "bots", "simulate", "--help")
	require.NoError(t, err)
	require.Contains(t, out, "event-file")
	require.Contains(t, out, "bot-config-file")
	require.Contains(t, out, "log-level")
	require.NotContains(t, out, "bot-token")
	require.NotContains(t, out, "application-id")
	_, err = execute(t, "bots", "simulate")
	require.Error(t, err)
	_, err = execute(t, "bots", "list", "--log-level", "bogus")
	require.ErrorContains(t, err, "log-level")
}

func TestLocalRunnerRejectsExternalConfiguration(t *testing.T) {
	repo, err := filepath.Abs("../../examples/slack-bots")
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "connection.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"apiURL":"https://slack.com/api/","botToken":"secret-marker","appToken":"secret-marker","teamID":"T","appID":"A"}`), 0600))
	_, err = execute(t, "bots", "run-local", "ping", "--bot-repository", repo, "--local-connection-file", path)
	require.Error(t, err)
	require.NotContains(t, err.Error(), "secret-marker")
	out, err := execute(t, "bots", "run-local", "--help")
	require.NoError(t, err)
	require.Contains(t, out, "local-connection-file")
}

func TestLocalUnifiedVerb(t *testing.T) {
	out, err := execute(t, "bots", "invoke", "unified-demo", "status", "--bot-repository", "../../examples/slack-bots")
	require.NoError(t, err)
	require.Contains(t, out, `"active": true`)
	require.NotContains(t, out, "apiKey")
	_, err = execute(t, "bots", "invoke", "unified-demo", "missing", "--bot-repository", "../../examples/slack-bots")
	require.ErrorContains(t, err, "not found")
}

func TestDocumentedProfilesCommand(t *testing.T) {
	out, err := execute(t, "profiles", "--config-dir", t.TempDir())
	require.NoError(t, err)
	require.JSONEq(t, "[]", out)
}
