package botcli

import (
	"context"
	"fmt"

	glazed_cli "github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/spf13/cobra"

	"github.com/go-go-golems/discord-bot/internal/jsdiscord"
)

func botCLIParserConfig(appName string) glazed_cli.CobraParserConfig {
	return glazed_cli.CobraParserConfig{
		AppName:           appName,
		ShortHelpSections: []string{schema.DefaultSlug, schema.GlobalDefaultSlug},
	}
}

// NewBotsCommand builds the public repo-driven bot command tree from a bootstrap.
func NewBotsCommand(bootstrap Bootstrap, opts ...CommandOption) (*cobra.Command, error) {
	cfg, err := applyCommandOptions(opts...)
	if err != nil {
		return nil, err
	}
	parserConfig := botCLIParserConfig(cfg.appName)

	root := &cobra.Command{
		Use:   "bots",
		Short: "List, inspect, and run named JavaScript bot implementations",
	}

	commands, err := NewBotsCommands(bootstrap, opts...)
	if err != nil {
		return nil, err
	}
	if err := glazed_cli.AddCommandsToRootCommand(root, commands, nil, glazed_cli.WithParserConfig(parserConfig)); err != nil {
		return nil, fmt.Errorf("register bot commands: %w", err)
	}
	return root, nil
}

func NewBotsCommands(bootstrap Bootstrap, opts ...CommandOption) ([]cmds.Command, error) {
	cfg, err := applyCommandOptions(opts...)
	if err != nil {
		return nil, err
	}
	hostOpts := cfg.hostOptions()
	ret := []cmds.Command{}
	staticCommands, err := staticBotsCommands(bootstrap, hostOpts)
	if err != nil {
		return nil, err
	}
	ret = append(ret, staticCommands...)
	if len(bootstrap.Repositories) > 0 {
		discoveredCommands, err := discoveredBotsCommands(bootstrap, cfg, hostOpts)
		if err != nil {
			return nil, err
		}
		ret = append(ret, discoveredCommands...)
	}
	return ret, nil
}

func staticBotsCommands(bootstrap Bootstrap, hostOpts []jsdiscord.HostOption) ([]cmds.Command, error) {
	listDesc := cmds.NewCommandDescription(
		"list",
		cmds.WithShort("List discovered bot implementations"),
		cmds.WithLong("Emit all discovered bots as a structured table."),
	)
	listCmd := &listBotsCommand{CommandDescription: listDesc, bootstrap: bootstrap, hostOpts: hostOpts}

	helpDesc := cmds.NewCommandDescription(
		"help",
		cmds.WithShort("Show metadata for one discovered bot implementation"),
		cmds.WithLong("Emit bot metadata (commands, events, components, modals) as structured rows."),
		cmds.WithArguments(
			fields.New("bot-name", fields.TypeString,
				fields.WithIsArgument(true),
				fields.WithRequired(true),
				fields.WithHelp("Bot name or source label"),
			),
		),
	)
	helpCmd := &helpBotsCommand{CommandDescription: helpDesc, bootstrap: bootstrap, hostOpts: hostOpts}

	return []cmds.Command{listCmd, helpCmd}, nil
}

func discoveredBotsCommands(bootstrap Bootstrap, cfg commandOptions, hostOpts []jsdiscord.HostOption) ([]cmds.Command, error) {
	discoveredBots, err := DiscoverBots(context.Background(), bootstrap, hostOpts...)
	if err != nil {
		return nil, fmt.Errorf("discover bots: %w", err)
	}
	botByScriptPath := map[string]DiscoveredBot{}
	for _, discoveredBot := range discoveredBots {
		botByScriptPath[discoveredBot.ScriptPath()] = discoveredBot
	}

	registries, err := ScanBotRepositories(bootstrap.Repositories)
	if err != nil {
		return nil, fmt.Errorf("scan bot repositories: %w", err)
	}

	runCommandByBot := map[string]*cmds.CommandDescription{}
	ret := make([]cmds.Command, 0)
	for _, registry := range registries {
		for _, verb := range registry.Verbs() {
			if verb.Name == "run" {
				baseDesc, err := registry.CommandDescriptionForVerb(verb)
				if err != nil {
					return nil, fmt.Errorf("build run command description for %s: %w", verb.SourceRef(), err)
				}
				desc := ensureRunCommandDefaults(baseDesc)
				ret = append(ret, &botRunCommand{CommandDescription: desc, scriptPath: verb.File.AbsPath, hostOpts: hostOpts})
				if bot, ok := botByScriptPath[verb.File.AbsPath]; ok {
					runCommandByBot[bot.Name()] = desc
				}
				continue
			}

			cmd, err := registry.CommandForVerbWithInvoker(verb, makeBotVerbInvoker(cfg))
			if err != nil {
				return nil, fmt.Errorf("build jsverb command for %s: %w", verb.SourceRef(), err)
			}
			ret = append(ret, cmd)
		}
	}

	for _, discoveredBot := range discoveredBots {
		if _, ok := runCommandByBot[discoveredBot.Name()]; ok {
			continue
		}
		baseDesc := buildSyntheticBotRunDescription(discoveredBot, discoveredBot.Name())
		ret = append(ret, &botRunCommand{CommandDescription: baseDesc, scriptPath: discoveredBot.ScriptPath(), hostOpts: hostOpts})
		runCommandByBot[discoveredBot.Name()] = baseDesc
	}

	return ret, nil
}
