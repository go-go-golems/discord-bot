package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-go-golems/discord-bot/pkg/slackcli"
	"github.com/go-go-golems/discord-bot/pkg/slackdoc"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func newRoot() (*cobra.Command, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	root := &cobra.Command{Use: "slack-bot", Short: "Go-hosted Slack bots: offline inspection and simulation", SilenceUsage: true, SilenceErrors: true, Long: "Inspect, generate manifests, and simulate Slack bot behavior without tokens or network access. Socket Mode is not implemented yet."}
	bots, err := slackcli.NewBotsCommand(logger)
	if err != nil {
		return nil, err
	}
	root.AddCommand(bots)
	hs := help.NewHelpSystem()
	if err := slackdoc.AddTo(hs); err != nil {
		return nil, err
	}
	help_cmd.SetupCobraRootCommand(hs, root)
	return root, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	root, err := newRoot()
	if err == nil {
		err = root.ExecuteContext(ctx)
	}
	if err != nil {
		logger := zerolog.New(os.Stderr)
		logger.Error().Err(err).Msg("Slack command failed")
		os.Exit(1)
	}
}
