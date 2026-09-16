package slackcli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-go-golems/discord-bot/internal/slacktransport"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/go-go-golems/discord-bot/pkg/slackhost"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/sources"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type settings struct {
	SkipManifestUpdate  bool   `glazed:"skip-manifest-update"`
	ConfigTokenFile     string `glazed:"config-token-file"`
	CredentialsFile     string `glazed:"credentials-file"`
	ConfigDir           string `glazed:"config-dir"`
	Profile             string `glazed:"profile"`
	TeamID              string `glazed:"team-id"`
	LocalConnectionFile string `glazed:"local-connection-file"`
	LogLevel            string `glazed:"log-level"`
	Repository          string `glazed:"bot-repository"`
	Name                string `glazed:"name"`
	EventFile           string `glazed:"event-file"`
	ConfigFile          string `glazed:"bot-config-file"`
	TimeoutMS           int    `glazed:"timeout-ms"`
}
type command struct {
	*cmds.CommandDescription
	operation string
	logger    zerolog.Logger
	appClient *http.Client
}

var _ cmds.WriterCommand = (*command)(nil)

// Offline commands produce one JSON document; run-local waits for cancellation. These are writer commands, not
// structured-row commands, and do not depend on Glazed's evolving output flags.
func NewBotsCommand(logger zerolog.Logger) (*cobra.Command, error) {
	return newBotsCommand(logger, appHTTPClient())
}

func newBotsCommand(logger zerolog.Logger, appClient *http.Client) (*cobra.Command, error) {
	root := &cobra.Command{Use: "bots", Short: "Inspect, create, install and run Slack bots"}
	for _, op := range []string{"list", "inspect", "manifest", "create-app", "install", "simulate", "run", "run-local"} {
		short := op + " a Slack bot (offline)"
		if op == "create-app" {
			short = "Create a Slack app from the bot's manifest using the Slack API"
		}
		if op == "install" {
			short = "Install a local Slack app and save runtime tokens"
		}
		if op == "run" {
			short = "Update the app manifest and run a Slack bot using stored credentials"
		}
		desc := cmds.NewCommandDescription(op, cmds.WithShort(short), cmds.WithFlags(
			fields.New("log-level", fields.TypeString, fields.WithDefault("info"), fields.WithHelp("Log level (debug, info, warn, error)")),
			fields.New("bot-repository", fields.TypeString, fields.WithDefault("examples/slack-bots"), fields.WithHelp("Repository of Slack bot entrypoints")),
			fields.New("timeout-ms", fields.TypeInteger, fields.WithDefault(5000), fields.WithHelp("Inspection, invocation and API request deadline in milliseconds")),
		))
		if op != "list" {
			cmds.WithArguments(fields.New("name", fields.TypeString, fields.WithIsArgument(true), fields.WithRequired(true), fields.WithHelp("Bot name")))(desc)
		}
		if op == "run-local" {
			cmds.WithFlags(fields.New("local-connection-file", fields.TypeString, fields.WithRequired(true), fields.WithHelp("JSON local mock connection file; never exposed to JavaScript")), fields.New("bot-config-file", fields.TypeString, fields.WithHelp("Optional declared bot configuration JSON")))(desc)
		}
		if op == "create-app" {
			cmds.WithFlags(
				fields.New("config-token-file", fields.TypeString, fields.WithHelp("File containing an app-configuration access token from api.slack.com/apps")),
				fields.New("credentials-file", fields.TypeString, fields.WithHelp("Optional new private file for returned app credentials (0600; refuses overwrite)")),
				fields.New("config-dir", fields.TypeString, fields.WithHelp("Local Slack profile and credentials directory (default: user config directory)")),
				fields.New("profile", fields.TypeString, fields.WithHelp("Named local profile (alternative to --config-token-file)")),
			)(desc)
		}
		if op == "install" {
			cmds.WithFlags(
				fields.New("config-dir", fields.TypeString, fields.WithHelp("Local Slack profile and credentials directory (default: user config directory)")),
				fields.New("profile", fields.TypeString, fields.WithHelp("Named local profile")),
				fields.New("team-id", fields.TypeString, fields.WithHelp("Slack workspace/team ID")),
			)(desc)
		}
		if op == "run" {
			cmds.WithFlags(
				fields.New("skip-manifest-update", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Connect without updating the installed app manifest (no management token required)")),
				fields.New("config-dir", fields.TypeString, fields.WithHelp("Local Slack profile and credentials directory (default: user config directory)")),
				fields.New("profile", fields.TypeString, fields.WithHelp("Named local profile")),
				fields.New("bot-config-file", fields.TypeString, fields.WithHelp("Optional declared bot configuration JSON")),
			)(desc)
		}
		if op == "simulate" {
			cmds.WithFlags(fields.New("event-file", fields.TypeString, fields.WithRequired(true), fields.WithHelp("JSON invocation fixture")), fields.New("bot-config-file", fields.TypeString, fields.WithHelp("Optional JSON object of declared bot configuration")))(desc)
		}
		// Use the pinned Glazed parser directly: its high-level builder calls
		// os.Exit on domain errors. The application owns errors and output here.
		domain := &command{CommandDescription: desc, operation: op, logger: logger, appClient: appClient}
		c := cli.NewCobraCommandFromCommandDescription(desc)
		parser, err := cli.NewCobraParserFromSections(desc.Schema, &cli.CobraParserConfig{
			SkipCommandSettingsSection: true,
			MiddlewaresFunc: func(_ *values.Values, cmd *cobra.Command, args []string) ([]sources.Middleware, error) {
				return []sources.Middleware{sources.FromCobra(cmd), sources.FromArgs(args), sources.FromDefaults()}, nil
			},
		})
		if err != nil {
			return nil, err
		}
		if err := parser.AddToCobraCommand(c); err != nil {
			return nil, err
		}
		c.RunE = func(cmd *cobra.Command, args []string) error {
			vals, err := parser.Parse(cmd, args)
			if err != nil {
				return err
			}
			return domain.RunIntoWriter(cmd.Context(), vals, cmd.OutOrStdout())
		}
		root.AddCommand(c)
	}
	return root, nil
}
func (c *command) RunIntoWriter(ctx context.Context, vals *values.Values, w io.Writer) error {
	var s settings
	if err := vals.DecodeSectionInto(schema.DefaultSlug, &s); err != nil {
		return err
	}
	if s.TimeoutMS <= 0 || s.TimeoutMS > 60000 {
		return errors.New("timeout-ms must be between 1 and 60000")
	}
	level, err := zerolog.ParseLevel(s.LogLevel)
	if err != nil {
		return errors.Wrap(err, "log-level")
	}
	timeout := time.Duration(s.TimeoutMS) * time.Millisecond
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if c.operation == "list" {
		all, err := Discover(ctx, s.Repository, timeout)
		if err != nil {
			return err
		}
		return enc.Encode(all)
	}
	d, err := Resolve(ctx, s.Repository, s.Name, timeout)
	if err != nil {
		return err
	}
	switch c.operation {
	case "create-app":
		return c.createApp(ctx, d, s, w)
	case "install":
		return c.installApp(ctx, d, s, w)
	case "run":
		return c.runRemote(ctx, d, s)
	case "run-local":
		var connection struct {
			APIURL, BotToken, AppToken, TeamID, AppID, UserID string
			AllowedChannels                                   []string
		}
		if err := readJSON(s.LocalConnectionFile, &connection); err != nil {
			return errors.New("invalid local connection file")
		}
		client, err := slacktransport.NewLocal(slacktransport.LocalOptions{APIURL: connection.APIURL, BotToken: connection.BotToken, AppToken: connection.AppToken, TeamID: connection.TeamID, AppID: connection.AppID, AllowedChannels: connection.AllowedChannels, Logger: c.logger.Level(level)})
		if err != nil {
			return err
		}
		defer client.Close()
		config := map[string]any{}
		if s.ConfigFile != "" {
			if err := readJSON(s.ConfigFile, &config); err != nil {
				return err
			}
		}
		host, err := slackhost.Load(ctx, d.ScriptPath, slackhost.Options{Messages: client, Views: client, Config: config, Timeout: timeout, Logger: c.logger.Level(level)})
		if err != nil {
			return err
		}
		defer func() { _ = host.Close(context.Background()) }()
		return client.Run(ctx, host)
	case "inspect":
		return enc.Encode(d)
	case "manifest":
		return enc.Encode(Manifest(d))
	case "simulate":
		var invocation slackbot.Invocation
		if err := readJSON(s.EventFile, &invocation); err != nil {
			return err
		}
		config := map[string]any{}
		if s.ConfigFile != "" {
			if err := readJSON(s.ConfigFile, &config); err != nil {
				return err
			}
		}
		recorder := &slackbot.Recorder{}
		if invocation.Interaction != nil {
			invocation.Interaction.Ack = recorder
		}
		h, err := slackhost.Load(ctx, d.ScriptPath, slackhost.Options{Messages: recorder, Views: recorder, Config: config, Timeout: timeout, Logger: c.logger.Level(level)})
		if err != nil {
			return err
		}
		defer func() { _ = h.Close(context.Background()) }()
		if err := h.Dispatch(ctx, invocation, recorder); err != nil {
			return err
		}
		return enc.Encode(recorder.Operations())
	}
	return errors.New("unknown operation")
}
func readJSON(path string, target any) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "open fixture")
	}
	defer func() { _ = f.Close() }()
	data, err := io.ReadAll(io.LimitReader(f, 1024*1024+1))
	if err != nil {
		return errors.Wrap(err, "read fixture")
	}
	if len(data) > 1024*1024 {
		return errors.New("fixture exceeds 1 MiB")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return errors.Wrap(err, "decode fixture")
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return errors.New("fixture must contain exactly one JSON value (maximum 1 MiB)")
	}
	return nil
}

// Manifest emits a reviewable Slack app manifest; it never contacts Slack.
func Manifest(d slackbot.Descriptor) map[string]any {
	commands := []map[string]any{}
	for _, c := range d.Commands {
		commands = append(commands, map[string]any{"command": c.Name, "description": c.Description, "should_escape": true})
	}
	scopes := []string{"chat:write"}
	if len(commands) > 0 || len(d.Shortcuts) > 0 {
		scopes = append(scopes, "commands")
	}
	if len(d.Events) > 0 {
		scopes = append(scopes, "app_mentions:read")
	}
	for _, scope := range d.Scopes {
		found := false
		for _, existing := range scopes {
			if existing == scope {
				found = true
				break
			}
		}
		if !found {
			scopes = append(scopes, scope)
		}
	}
	features := map[string]any{"bot_user": map[string]any{"display_name": d.Name, "always_online": false}}
	if len(commands) > 0 {
		features["slash_commands"] = commands
	}
	if len(d.Shortcuts) > 0 {
		shortcuts := []map[string]any{}
		for _, shortcut := range d.Shortcuts {
			shortcuts = append(shortcuts, map[string]any{"name": shortcut.Name, "description": shortcut.Description, "callback_id": shortcut.CallbackID, "type": shortcut.Type})
		}
		features["shortcuts"] = shortcuts
	}
	settings := map[string]any{"socket_mode_enabled": true, "interactivity": map[string]any{"is_enabled": true}}
	if len(d.Events) > 0 {
		settings["event_subscriptions"] = map[string]any{"bot_events": d.Events}
	}
	return map[string]any{"display_information": map[string]any{"name": d.Name}, "features": features, "oauth_config": map[string]any{"scopes": map[string]any{"bot": scopes}}, "settings": settings}
}
