package slackbot

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfigProjectsDeclaredFields(t *testing.T) {
	d := Descriptor{Name: "ping", Run: RunSchema{Fields: map[string]Field{"greeting": {Type: "string", Default: "pong"}}}}
	require.NoError(t, d.Validate())
	c, err := d.Config(map[string]any{"botToken": "secret-marker", "appToken": "another-marker"})
	require.NoError(t, err)
	require.Equal(t, map[string]any{"greeting": "pong"}, c)
	_, err = d.Config(map[string]any{"greeting": true})
	require.ErrorContains(t, err, "must be string")
}
func TestInvocationValidation(t *testing.T) {
	i := Invocation{TeamID: "T", ChannelID: "C", UserID: "U", Event: "app_mention", TS: "1741234567.000123"}
	require.NoError(t, i.Validate())
	i.Command = "/ping"
	require.Error(t, i.Validate())
}
