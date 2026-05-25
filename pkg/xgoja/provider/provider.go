package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/discord-bot/internal/jsdiscord"
	"github.com/go-go-golems/discord-bot/pkg/botcli"
	glazedcli "github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
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
	sections, err := collectModuleSections(ctx.SelectedModules, profile, ctx.Name)
	if err != nil {
		return nil, err
	}

	opts := []botcli.CommandOption{}
	var runtimeFactory *xgojaBotRuntimeFactory
	if factory, ok := ctx.RuntimeFactory.(xgojaRuntimeFactory); ok {
		runtimeFactory = &xgojaBotRuntimeFactory{factory: factory, profile: profile, selectedModules: ctx.SelectedModules}
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

func collectModuleSections(descriptors []providerapi.ModuleDescriptor, profile string, commandProvider string) ([]schema.Section, error) {
	sections := []schema.Section{}
	seen := map[string]string{schema.DefaultSlug: "bot command schema"}
	for _, descriptor := range descriptors {
		for _, capability := range descriptor.Capabilities {
			sectionCapability, ok := capability.(providerapi.ConfigSectionCapability)
			if !ok {
				continue
			}
			moduleSections, err := sectionCapability.ConfigSections(providerapi.SectionContext{
				CommandProviderID: commandProvider,
				RuntimeProfile:    profile,
				PackageID:         descriptor.PackageID,
				ModuleID:          descriptor.ModuleID,
			})
			if err != nil {
				return nil, fmt.Errorf("collect config sections for %s.%s capability %s: %w", descriptor.PackageID, descriptor.ModuleID, capability.CapabilityID(), err)
			}
			for _, section := range moduleSections {
				if section == nil {
					return nil, fmt.Errorf("%s.%s capability %s returned nil section", descriptor.PackageID, descriptor.ModuleID, capability.CapabilityID())
				}
				slug := strings.TrimSpace(section.GetSlug())
				if slug == "" {
					return nil, fmt.Errorf("%s.%s capability %s returned empty section slug", descriptor.PackageID, descriptor.ModuleID, capability.CapabilityID())
				}
				if previous, ok := seen[slug]; ok {
					return nil, fmt.Errorf("duplicate config section slug %q from %s.%s capability %s; already provided by %s", slug, descriptor.PackageID, descriptor.ModuleID, capability.CapabilityID(), previous)
				}
				seen[slug] = fmt.Sprintf("%s.%s capability %s", descriptor.PackageID, descriptor.ModuleID, capability.CapabilityID())
				sections = append(sections, section)
			}
		}
	}
	return sections, nil
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
	factory         xgojaRuntimeFactory
	profile         string
	selectedModules []providerapi.ModuleDescriptor
	mu              sync.Mutex
	values          *values.Values
}

func (f *xgojaBotRuntimeFactory) HostOptions() []jsdiscord.HostOption {
	return []jsdiscord.HostOption{jsdiscord.WithRuntimeFactory(xgojaHostRuntimeFactory{parent: f})}
}

func (f *xgojaBotRuntimeFactory) setValues(vals *values.Values) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values = vals
}

func (f *xgojaBotRuntimeFactory) currentValues() *values.Values {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.values
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
	rt, err := f.factory.NewRuntime(ctx, f.profile, opts...)
	if err != nil {
		return nil, err
	}
	if vals := f.currentValues(); vals != nil {
		if err := initSelectedModules(ctx, vals, rt, f.selectedModules); err != nil {
			_ = rt.Close(context.Background())
			return nil, err
		}
	}
	return rt, nil
}

type xgojaHostRuntimeFactory struct {
	parent *xgojaBotRuntimeFactory
}

func (f xgojaHostRuntimeFactory) NewRuntime(ctx context.Context, opts ...require.Option) (*engine.Runtime, error) {
	if f.parent == nil {
		return nil, fmt.Errorf("xgoja runtime factory is nil")
	}
	return f.parent.newRuntime(ctx, opts...)
}

func initSelectedModules(ctx context.Context, vals *values.Values, rt *engine.Runtime, descriptors []providerapi.ModuleDescriptor) error {
	if rt == nil {
		return fmt.Errorf("runtime is nil")
	}
	handle := runtimeHandle{rt: rt}
	for _, descriptor := range descriptors {
		for _, capability := range descriptor.Capabilities {
			initializer, ok := capability.(providerapi.RuntimeInitializerCapability)
			if !ok {
				continue
			}
			if err := initializer.InitRuntimeFromSections(ctx, vals, handle); err != nil {
				return fmt.Errorf("initialize runtime from sections for %s.%s capability %s: %w", descriptor.PackageID, descriptor.ModuleID, capability.CapabilityID(), err)
			}
		}
	}
	return nil
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
	c.factory.setValues(vals)
	defer c.factory.setValues(nil)
	return c.command.(cmds.BareCommand).Run(ctx, vals)
}

type valueWriterCommand struct{ valueCommandBase }

func (c valueWriterCommand) RunIntoWriter(ctx context.Context, vals *values.Values, w io.Writer) error {
	c.factory.setValues(vals)
	defer c.factory.setValues(nil)
	return c.command.(cmds.WriterCommand).RunIntoWriter(ctx, vals, w)
}

type valueGlazeCommand struct{ valueCommandBase }

func (c valueGlazeCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	c.factory.setValues(vals)
	defer c.factory.setValues(nil)
	return c.command.(cmds.GlazeCommand).RunIntoGlazeProcessor(ctx, vals, gp)
}
