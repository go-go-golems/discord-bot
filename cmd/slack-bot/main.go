package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-go-golems/discord-bot/pkg/slackcli"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func newRoot() (*cobra.Command, error) {
	level := "info"
	// A shared level writer lets child commands use the final parsed log level.
	writer := &levelWriter{writer: zerolog.ConsoleWriter{Out: os.Stderr}, level: zerolog.InfoLevel}
	logger := zerolog.New(writer).With().Timestamp().Logger()
	root := &cobra.Command{Use: "slack-bot", Short: "Go-hosted Slack bots: offline inspection and simulation", SilenceUsage: true, SilenceErrors: true, Long: "Inspect, generate manifests, and simulate Slack bot behavior without tokens or network access. Socket Mode is not implemented yet.", PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		parsed, err := zerolog.ParseLevel(level)
		if err != nil {
			return errors.Wrap(err, "log-level")
		}
		writer.level = parsed
		return nil
	}}
	root.PersistentFlags().StringVar(&level, "log-level", "info", "Log level (debug, info, warn, error)")
	bots, err := slackcli.NewBotsCommand(logger)
	if err != nil {
		return nil, err
	}
	root.AddCommand(bots)
	return root, nil
}

type levelWriter struct {
	writer zerolog.ConsoleWriter
	level  zerolog.Level
}

var _ zerolog.LevelWriter = (*levelWriter)(nil)

func (w *levelWriter) Write(p []byte) (int, error) { return w.writer.Write(p) }
func (w *levelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level < w.level {
		return len(p), nil
	}
	return w.writer.Write(p)
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
