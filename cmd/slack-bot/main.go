package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-go-golems/discord-bot/pkg/slackcli"
	"github.com/go-go-golems/discord-bot/pkg/slackdoc"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func newRoot() (*cobra.Command, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	root := &cobra.Command{Use: "slack-bot", Short: "Inspect, create and run Go-hosted Slack bots", SilenceUsage: true, SilenceErrors: true, Long: "Inspect, generate manifests, and simulate Slack bot behavior offline. The create-app command creates a Slack app using an explicit configuration token file. The install command installs a local app and saves runtime credentials. The run command connects to Slack using a selected profile; run-local connects only to an explicitly configured loopback Socket Mode mock."}
	bots, err := slackcli.NewBotsCommand(logger)
	if err != nil {
		return nil, err
	}
	root.AddCommand(bots)
	root.AddCommand(slackcli.NewCredentialsCommand(logger, nil))
	root.AddCommand(slackcli.NewProfilesCommand())
	hs := help.NewHelpSystem()
	if err := slackdoc.AddTo(hs); err != nil {
		return nil, err
	}
	help_cmd.SetupCobraRootCommand(hs, root)
	return root, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
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
