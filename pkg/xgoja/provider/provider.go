package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/discord-bot/internal/jsdiscord"
	"github.com/go-go-golems/discord-bot/pkg/botcli"
	glazedcli "github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/go-go-goja/engine"
	"github.com/go-go-golems/go-go-goja/pkg/jsverbs"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
)

const PackageID = "discord-bot"

type commandProviderConfig struct {
	Repositories     []string `json:"repositories,omitempty"`
	WorkingDirectory string   `json:"workingDirectory,omitempty"`
	RuntimeProfile   string   `json:"runtimeProfile,omitempty"`
}

type xgojaRuntimeFactory interface {
	NewRuntime(ctx context.Context, profile string, opts ...require.Option) (*engine.Runtime, error)
}

func Register(registry *providerapi.Registry) error {
	return registry.Package(PackageID,
		providerapi.Module{
			Name:        "discord",
			DefaultAs:   "discord",
			Description: "Discord bot definition module for JavaScript bot scripts",
			New: func(ctx providerapi.ModuleContext) (require.ModuleLoader, error) {
				moduleName := strings.TrimSpace(ctx.As)
				if moduleName == "" {
					moduleName = "discord"
				}
				return jsdiscord.NewLoader(jsdiscord.Config{ModuleName: moduleName}), nil
			},
		},
		providerapi.Module{
			Name:        "ui",
			DefaultAs:   "ui",
			Description: "Discord UI builder helper module",
			New: func(providerapi.ModuleContext) (require.ModuleLoader, error) {
				return jsdiscord.NewUILoader(), nil
			},
		},
		providerapi.CommandSetProvider{
			Name:         "bots",
			DefaultMount: "bots",
			Description:  "List, inspect, and run JavaScript Discord bots",
			New:          newBotsCommandSet,
		},
	)
}

func newBotsCommandSet(ctx providerapi.CommandSetContext) (*providerapi.CommandSet, error) {
	cfg := commandProviderConfig{}
	if len(ctx.Config) > 0 {
		if err := json.Unmarshal(ctx.Config, &cfg); err != nil {
			return nil, fmt.Errorf("decode discord-bot command provider config: %w", err)
		}
	}
	bootstrap, err := bootstrapFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	profile := strings.TrimSpace(cfg.RuntimeProfile)
	if profile == "" {
		profile = strings.TrimSpace(ctx.RuntimeProfile)
	}
	if profile == "" {
		profile = "main"
	}
	opts := []botcli.CommandOption{}
	if factory, ok := ctx.RuntimeFactory.(xgojaRuntimeFactory); ok {
		opts = append(opts, botcli.WithRuntimeFactory(xgojaBotRuntimeFactory{factory: factory, profile: profile}))
	}
	commands, err := botcli.NewBotsCommands(bootstrap, opts...)
	if err != nil {
		return nil, err
	}
	return &providerapi.CommandSet{
		Commands: commands,
		ParserConfig: &glazedcli.CobraParserConfig{
			AppName:           "discord",
			ShortHelpSections: []string{schema.DefaultSlug, schema.GlobalDefaultSlug},
		},
	}, nil
}

func bootstrapFromConfig(cfg commandProviderConfig) (botcli.Bootstrap, error) {
	args := []string{}
	for _, repo := range cfg.Repositories {
		repo = strings.TrimSpace(repo)
		if repo == "" {
			continue
		}
		args = append(args, "--"+botcli.BotRepositoryFlag, repo)
	}
	opts := []botcli.BuildOption{}
	if strings.TrimSpace(cfg.WorkingDirectory) != "" {
		abs, err := filepath.Abs(cfg.WorkingDirectory)
		if err != nil {
			return botcli.Bootstrap{}, fmt.Errorf("resolve workingDirectory: %w", err)
		}
		opts = append(opts, botcli.WithWorkingDirectory(abs))
	}
	return botcli.BuildBootstrap(args, opts...)
}

type xgojaBotRuntimeFactory struct {
	factory xgojaRuntimeFactory
	profile string
}

func (f xgojaBotRuntimeFactory) HostOptions() []jsdiscord.HostOption {
	return []jsdiscord.HostOption{jsdiscord.WithRuntimeFactory(xgojaHostRuntimeFactory(f))}
}

func (f xgojaBotRuntimeFactory) NewRuntimeForVerb(ctx context.Context, registry *jsverbs.Registry, verb *jsverbs.VerbSpec) (*engine.Runtime, error) {
	if f.factory == nil {
		return nil, fmt.Errorf("xgoja runtime factory is nil")
	}
	opts := []require.Option{}
	if registry != nil {
		opts = append(opts, require.WithLoader(registry.RequireLoader()))
	}
	if verb != nil && verb.File != nil && strings.TrimSpace(verb.File.AbsPath) != "" {
		absScript := strings.TrimSpace(verb.File.AbsPath)
		moduleRootsOpt, err := engine.RequireOptionWithModuleRootsFromScript(absScript, engine.DefaultModuleRootsOptions())
		if err != nil {
			return nil, fmt.Errorf("resolve module roots from script: %w", err)
		}
		opts = append(opts,
			moduleRootsOpt,
			require.WithGlobalFolders(filepath.Dir(absScript), filepath.Join(filepath.Dir(absScript), "node_modules")),
		)
	}
	return f.factory.NewRuntime(ctx, f.profile, opts...)
}

type xgojaHostRuntimeFactory struct {
	factory xgojaRuntimeFactory
	profile string
}

func (f xgojaHostRuntimeFactory) NewRuntime(ctx context.Context, opts ...require.Option) (*engine.Runtime, error) {
	if f.factory == nil {
		return nil, fmt.Errorf("xgoja runtime factory is nil")
	}
	return f.factory.NewRuntime(ctx, f.profile, opts...)
}
