package slackcli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

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
	LogLevel   string `glazed:"log-level"`
	Repository string `glazed:"bot-repository"`
	Name       string `glazed:"name"`
	EventFile  string `glazed:"event-file"`
	ConfigFile string `glazed:"bot-config-file"`
	TimeoutMS  int    `glazed:"timeout-ms"`
}
type command struct {
	*cmds.CommandDescription
	operation string
	logger    zerolog.Logger
}

var _ cmds.WriterCommand = (*command)(nil)

// Commands intentionally produce one JSON document. They are writer commands, not
// structured-row commands, and do not depend on Glazed's evolving output flags.
func NewBotsCommand(logger zerolog.Logger) (*cobra.Command, error) {
	root := &cobra.Command{Use: "bots", Short: "Inspect and simulate Slack bots offline"}
	for _, op := range []string{"list", "inspect", "manifest", "simulate"} {
		desc := cmds.NewCommandDescription(op, cmds.WithShort(op+" a Slack bot (offline)"), cmds.WithFlags(
			fields.New("log-level", fields.TypeString, fields.WithDefault("info"), fields.WithHelp("Log level (debug, info, warn, error)")),
			fields.New("bot-repository", fields.TypeString, fields.WithDefault("examples/slack-bots"), fields.WithHelp("Repository of Slack bot entrypoints")),
			fields.New("timeout-ms", fields.TypeInteger, fields.WithDefault(5000), fields.WithHelp("Inspection and invocation deadline in milliseconds")),
		))
		if op != "list" {
			cmds.WithArguments(fields.New("name", fields.TypeString, fields.WithIsArgument(true), fields.WithRequired(true), fields.WithHelp("Bot name")))(desc)
		}
		if op == "simulate" {
			cmds.WithFlags(fields.New("event-file", fields.TypeString, fields.WithRequired(true), fields.WithHelp("JSON invocation fixture")), fields.New("bot-config-file", fields.TypeString, fields.WithHelp("Optional JSON object of declared bot configuration")))(desc)
		}
		// Use the pinned Glazed parser directly: its high-level builder calls
		// os.Exit on domain errors. The application owns errors and output here.
		domain := &command{desc, op, logger}
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
		h, err := slackhost.Load(ctx, d.ScriptPath, slackhost.Options{Messages: recorder, Config: config, Timeout: timeout, Logger: c.logger.Level(level)})
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
	if len(commands) > 0 {
		scopes = append(scopes, "commands")
	}
	if len(d.Events) > 0 {
		scopes = append(scopes, "app_mentions:read")
	}
	features := map[string]any{"bot_user": map[string]any{"display_name": d.Name, "always_online": false}}
	if len(commands) > 0 {
		features["slash_commands"] = commands
	}
	settings := map[string]any{"socket_mode_enabled": true, "interactivity": map[string]any{"is_enabled": true}}
	if len(d.Events) > 0 {
		settings["event_subscriptions"] = map[string]any{"bot_events": d.Events}
	}
	return map[string]any{"display_information": map[string]any{"name": d.Name}, "features": features, "oauth_config": map[string]any{"scopes": map[string]any{"bot": scopes}}, "settings": settings}
}
