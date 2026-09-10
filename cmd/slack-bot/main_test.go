package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	require.Len(t, all, 1)
	require.Equal(t, "ping", all[0].Name)
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
