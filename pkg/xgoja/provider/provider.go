package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/discord-bot/internal/jsdiscord"
	"github.com/go-go-golems/discord-bot/pkg/botcli"
	glazedcli "github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/go-go-goja/pkg/jsverbs"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerutil"
)

const PackageID = "discord-bot"

type commandProviderConfig struct {
	Repositories     []string `json:"repositories,omitempty"`
	WorkingDirectory string   `json:"workingDirectory,omitempty"`
	RuntimeProfile   string   `json:"runtimeProfile,omitempty"`
}

func Register(registry *providerapi.ProviderRegistry) error {
	return registry.Package(PackageID,
		providerapi.Module{
			Name:        "discord",
			DefaultAs:   "discord",
			Description: "Discord bot definition module for JavaScript bot scripts",
			NewModuleFactory: func(ctx providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
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
			NewModuleFactory: func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
				return jsdiscord.NewUILoader(), nil
			},
		},
		providerapi.CommandSetProvider{
			Name:         "bots",
			DefaultMount: "bots",
			Description:  "List, inspect, and run JavaScript Discord bots",
			NewCommandSet: newBotsCommandSet,
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
		profile = "main"
	}
	sections, err := providerutil.CollectGlazedConfigSections(ctx.SelectedModules, providerapi.SectionRequest{
		CommandProviderID: ctx.Name,
	}, map[string]string{schema.DefaultSlug: "bot command schema"})
	if err != nil {
		return nil, err
	}

	opts := []botcli.CommandOption{}
	var runtimeFactory *xgojaBotRuntimeFactory
	if ctx.RuntimeFactory != nil {
		runtimeFactory = &xgojaBotRuntimeFactory{factory: ctx.RuntimeFactory, profile: profile, selectedModules: ctx.SelectedModules}
		opts = append(opts, botcli.WithRuntimeFactory(runtimeFactory))
	}
	commands, err := botcli.NewBotsCommands(bootstrap, opts...)
	if err != nil {
		return nil, err
	}
	if len(sections) > 0 {
		for _, command := range commands {
			if command == nil || command.Description() == nil {
				continue
			}
			for _, section := range sections {
				command.Description().SetSections(section)
			}
		}
	}
	if runtimeFactory != nil {
		commands = wrapCommandsWithValues(commands, runtimeFactory)
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
	factory         providerapi.RuntimeFactory
	profile         string
	selectedModules []providerapi.ModuleDescriptor
}

func (f *xgojaBotRuntimeFactory) HostOptions() []jsdiscord.HostOption {
	return []jsdiscord.HostOption{jsdiscord.WithRuntimeFactory(xgojaHostRuntimeFactory{parent: f})}
}

func (f *xgojaBotRuntimeFactory) HostOptionsWithValues(vals *values.Values) []jsdiscord.HostOption {
	return []jsdiscord.HostOption{jsdiscord.WithRuntimeFactory(xgojaHostRuntimeFactory{parent: f, values: vals})}
}

type commandValuesContextKey struct{}

func contextWithCommandValues(ctx context.Context, vals *values.Values) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if vals == nil {
		return ctx
	}
	return context.WithValue(ctx, commandValuesContextKey{}, vals)
}

func commandValuesFromContext(ctx context.Context) *values.Values {
	if ctx == nil {
		return nil
	}
	vals, _ := ctx.Value(commandValuesContextKey{}).(*values.Values)
	return vals
}

func (f *xgojaBotRuntimeFactory) NewRuntimeForVerb(ctx context.Context, registry *jsverbs.Registry, verb *jsverbs.VerbSpec) (*engine.Runtime, error) {
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
	return f.newRuntime(ctx, opts...)
}

func (f *xgojaBotRuntimeFactory) newRuntime(ctx context.Context, opts ...require.Option) (*engine.Runtime, error) {
	if f.factory == nil {
		return nil, fmt.Errorf("xgoja runtime factory is nil")
	}
	rt, err := f.factory.NewRuntime(ctx, opts...)
	if err != nil {
		return nil, err
	}
	if err := rt.AddCloser(func(context.Context) error {
		jsdiscord.ForgetRuntime(rt.VM)
		return nil
	}); err != nil {
		_ = rt.Close(context.Background())
		return nil, err
	}
	if err := initSelectedModules(ctx, commandValuesFromContext(ctx), rt, f.selectedModules); err != nil {
		_ = rt.Close(context.Background())
		return nil, err
	}
	return rt, nil
}

type xgojaHostRuntimeFactory struct {
	parent *xgojaBotRuntimeFactory
	values *values.Values
}

func (f xgojaHostRuntimeFactory) NewRuntime(ctx context.Context, opts ...require.Option) (*engine.Runtime, error) {
	if f.parent == nil {
		return nil, fmt.Errorf("xgoja runtime factory is nil")
	}
	return f.parent.newRuntime(contextWithCommandValues(ctx, f.values), opts...)
}

func initSelectedModules(ctx context.Context, vals *values.Values, rt *engine.Runtime, descriptors []providerapi.ModuleDescriptor) error {
	if rt == nil {
		return fmt.Errorf("runtime is nil")
	}
	return providerutil.InitRuntimeFromSections(ctx, vals, runtimeHandle{rt: rt}, descriptors)
}

type runtimeHandle struct {
	rt *engine.Runtime
}

func (h runtimeHandle) Runtime() *goja.Runtime {
	if h.rt == nil {
		return nil
	}
	return h.rt.VM
}

func (h runtimeHandle) EngineRuntime() *engine.Runtime {
	return h.rt
}

func (h runtimeHandle) Close(ctx context.Context) error {
	if h.rt == nil {
		return nil
	}
	return h.rt.Close(ctx)
}

func (h runtimeHandle) AddCloser(fn func(context.Context) error) error {
	if h.rt == nil {
		return fmt.Errorf("runtime is nil")
	}
	return h.rt.AddCloser(fn)
}

func wrapCommandsWithValues(commands []cmds.Command, factory *xgojaBotRuntimeFactory) []cmds.Command {
	wrapped := make([]cmds.Command, 0, len(commands))
	for _, command := range commands {
		wrapped = append(wrapped, wrapCommandWithValues(command, factory))
	}
	return wrapped
}

func wrapCommandWithValues(command cmds.Command, factory *xgojaBotRuntimeFactory) cmds.Command {
	if command == nil || factory == nil {
		return command
	}
	base := valueCommandBase{command: command, factory: factory}
	if _, ok := command.(cmds.GlazeCommand); ok {
		return valueGlazeCommand{valueCommandBase: base}
	}
	if _, ok := command.(cmds.WriterCommand); ok {
		return valueWriterCommand{valueCommandBase: base}
	}
	if _, ok := command.(cmds.BareCommand); ok {
		return valueBareCommand{valueCommandBase: base}
	}
	return command
}

type valueCommandBase struct {
	command cmds.Command
	factory *xgojaBotRuntimeFactory
}

func (c valueCommandBase) Description() *cmds.CommandDescription { return c.command.Description() }
func (c valueCommandBase) ToYAML(w io.Writer) error              { return c.command.ToYAML(w) }

type valueBareCommand struct{ valueCommandBase }

func (c valueBareCommand) Run(ctx context.Context, vals *values.Values) error {
	ctx = contextWithCommandValues(ctx, vals)
	ctx = botcli.ContextWithHostOptions(ctx, c.factory.HostOptionsWithValues(vals)...)
	return c.command.(cmds.BareCommand).Run(ctx, vals)
}

type valueWriterCommand struct{ valueCommandBase }

func (c valueWriterCommand) RunIntoWriter(ctx context.Context, vals *values.Values, w io.Writer) error {
	ctx = contextWithCommandValues(ctx, vals)
	ctx = botcli.ContextWithHostOptions(ctx, c.factory.HostOptionsWithValues(vals)...)
	return c.command.(cmds.WriterCommand).RunIntoWriter(ctx, vals, w)
}

type valueGlazeCommand struct{ valueCommandBase }

func (c valueGlazeCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	ctx = contextWithCommandValues(ctx, vals)
	ctx = botcli.ContextWithHostOptions(ctx, c.factory.HostOptionsWithValues(vals)...)
	return c.command.(cmds.GlazeCommand).RunIntoGlazeProcessor(ctx, vals, gp)
}
