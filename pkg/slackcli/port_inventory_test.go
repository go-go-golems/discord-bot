package slackcli

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/stretchr/testify/require"
)

// Expectations are derived from the Discord source inventory and explicit native
// mappings, not from the Slack implementation being tested.
func TestSourceRegistrationParity(t *testing.T) {
	raw, err := os.ReadFile("testdata/port-registration-expectations.json")
	require.NoError(t, err)
	var expected []struct{ Bot, Kind, Name, Source string }
	require.NoError(t, json.Unmarshal(raw, &expected))
	descriptors, err := Discover(context.Background(), "../../examples/slack-bots", 0)
	require.NoError(t, err)
	byName := map[string]slackbot.Descriptor{}
	for _, d := range descriptors {
		byName[d.Name] = d
	}
	require.Len(t, byName, 13)
	for _, e := range expected {
		d, ok := byName[e.Bot]
		require.True(t, ok, e.Bot)
		var actual []string
		switch e.Kind {
		case "command":
			for _, c := range d.Commands {
				actual = append(actual, c.Name)
			}
		case "event":
			actual = d.Events
		case "action":
			actual = d.Actions
		case "view":
			actual = d.Views
		case "options":
			actual = d.Options
		case "shortcut":
			for _, s := range d.Shortcuts {
				actual = append(actual, s.CallbackID)
			}
		case "verb":
			for _, v := range d.Verbs {
				actual = append(actual, v.Name)
			}
		default:
			t.Fatalf("unknown expectation kind %s", e.Kind)
		}
		require.Contains(t, actual, e.Name, "%s from %s", e.Bot, e.Source)
	}
}
