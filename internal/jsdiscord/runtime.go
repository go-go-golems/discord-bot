package jsdiscord

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/engine"
	"github.com/go-go-golems/go-go-goja/pkg/runtimebridge"
)

// RuntimeStateContextKey is the engine context key for the discord runtime state.
// It is set during module registration so that runtime inspectors or future
// extensions can retrieve the RuntimeState from the module context.
const RuntimeStateContextKey = "discord.runtime"

type Config struct {
	ModuleName string
}

var runtimeStates sync.Map // *goja.Runtime -> *RuntimeState

func StateForRuntime(vm *goja.Runtime) (*RuntimeState, bool) {
	if vm == nil {
		return nil, false
	}
	value, ok := runtimeStates.Load(vm)
	if !ok {
		return nil, false
	}
	state, ok := value.(*RuntimeState)
	return state, ok
}

type Registrar struct {
	config Config
}

func NewRegistrar(config Config) *Registrar {
	return &Registrar{config: config}
}

func (r *Registrar) ID() string {
	return "discord-js-registrar"
}

func (r *Registrar) RegisterRuntimeModule(ctx *engine.RuntimeModuleContext, reg *require.Registry) error {
	if reg == nil {
		return fmt.Errorf("require registry is nil")
	}
	moduleName := strings.TrimSpace(r.config.ModuleName)
	if moduleName == "" {
		moduleName = "discord"
	}
	state := NewRuntimeState(moduleName)
	if ctx != nil {
		ctx.SetValue(RuntimeStateContextKey, state)
	}
	reg.RegisterNativeModule(state.ModuleName(), state.Loader)
	return nil
}

type RuntimeState struct {
	moduleName     string
	store          *MemoryStore
	outboundMu     sync.RWMutex
	outbound       *DiscordOps
	defaultGuildID string
}

func NewRuntimeState(moduleName string) *RuntimeState {
	moduleName = strings.TrimSpace(moduleName)
	if moduleName == "" {
		moduleName = "discord"
	}
	return &RuntimeState{moduleName: moduleName, store: NewMemoryStore()}
}

func (s *RuntimeState) ModuleName() string {
	if s == nil || strings.TrimSpace(s.moduleName) == "" {
		return "discord"
	}
	return s.moduleName
}

func (s *RuntimeState) SetOutboundOps(ops *DiscordOps) {
	s.SetOutboundOpsForGuild(ops, "")
}

func (s *RuntimeState) SetOutboundOpsForGuild(ops *DiscordOps, guildID string) {
	if s == nil {
		return
	}
	s.outboundMu.Lock()
	defer s.outboundMu.Unlock()
	s.outbound = ops
	if strings.TrimSpace(guildID) != "" {
		s.defaultGuildID = strings.TrimSpace(guildID)
	}
}

func (s *RuntimeState) outboundOps() *DiscordOps {
	if s == nil {
		return nil
	}
	s.outboundMu.RLock()
	defer s.outboundMu.RUnlock()
	return s.outbound
}

func (s *RuntimeState) Store() *MemoryStore {
	if s == nil {
		return NewMemoryStore()
	}
	if s.store == nil {
		s.store = NewMemoryStore()
	}
	return s.store
}

func NewLoader(config Config) require.ModuleLoader {
	return NewRuntimeState(config.ModuleName).Loader
}

func (s *RuntimeState) Loader(vm *goja.Runtime, moduleObj *goja.Object) {
	runtimeStates.Store(vm, s)
	exports := moduleObj.Get("exports").(*goja.Object)
	_ = exports.Set("defineBot", func(call goja.FunctionCall) goja.Value {
		return s.defineBot(vm, call)
	})
	_ = exports.Set("channels", s.topLevelChannelsObject(vm))

	// Polyfill jsverbs metadata functions so bot scripts can coexist
	// with __verb__ / __section__ / __package__ declarations.
	for _, name := range []string{"__package__", "__section__", "__verb__", "doc"} {
		if v := vm.Get(name); v == nil || goja.IsUndefined(v) {
			_ = vm.Set(name, func(goja.FunctionCall) goja.Value { return goja.Undefined() })
		}
	}
}

func (s *RuntimeState) topLevelChannelsObject(vm *goja.Runtime) *goja.Object {
	channels := vm.NewObject()
	_ = channels.Set("send", func(channelID string, payload any) (any, error) {
		ops := s.outboundOps()
		if ops == nil || ops.ChannelSend == nil {
			return nil, fmt.Errorf("discord outbound channel API is not ready; the bot session must be running")
		}
		ctx := runtimebridge.CurrentContext(vm)
		if ctx == nil {
			ctx = context.Background()
		}
		return nil, ops.ChannelSend(ctx, channelID, payload)
	})
	_ = channels.Set("list", func(call goja.FunctionCall) goja.Value {
		ops := s.outboundOps()
		if ops == nil || ops.ChannelList == nil {
			panic(vm.NewGoError(fmt.Errorf("discord channel list API is not ready; the bot session must be running")))
		}
		guildID := s.defaultGuild()
		if len(call.Arguments) > 0 && !goja.IsUndefined(call.Argument(0)) && !goja.IsNull(call.Argument(0)) {
			guildID = strings.TrimSpace(call.Argument(0).String())
		}
		if guildID == "" {
			panic(vm.NewGoError(fmt.Errorf("discord.channels.list requires a guild ID when no default guild is configured")))
		}
		ctx := runtimebridge.CurrentContext(vm)
		if ctx == nil {
			ctx = context.Background()
		}
		channels, err := ops.ChannelList(ctx, guildID)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(channels)
	})
	return channels
}

func (s *RuntimeState) defaultGuild() string {
	if s == nil {
		return ""
	}
	s.outboundMu.RLock()
	defer s.outboundMu.RUnlock()
	return strings.TrimSpace(s.defaultGuildID)
}

func (s *RuntimeState) defineBot(vm *goja.Runtime, call goja.FunctionCall) goja.Value {
	if len(call.Arguments) != 1 {
		panic(vm.NewGoError(fmt.Errorf("discord.defineBot expects defineBot(builderFn)")))
	}
	builder, ok := goja.AssertFunction(call.Arguments[0])
	if !ok {
		panic(vm.NewGoError(fmt.Errorf("discord.defineBot builder is not a function")))
	}
	draft := newBotDraft(s)
	api := vm.NewObject()
	_ = api.Set("command", func(call goja.FunctionCall) goja.Value { return draft.command(vm, call) })
	_ = api.Set("userCommand", func(call goja.FunctionCall) goja.Value { return draft.userCommand(vm, call) })
	_ = api.Set("messageCommand", func(call goja.FunctionCall) goja.Value { return draft.messageCommand(vm, call) })
	_ = api.Set("subcommand", func(call goja.FunctionCall) goja.Value { return draft.subcommand(vm, call) })
	_ = api.Set("event", func(call goja.FunctionCall) goja.Value { return draft.event(vm, call) })
	_ = api.Set("component", func(call goja.FunctionCall) goja.Value { return draft.component(vm, call) })
	_ = api.Set("modal", func(call goja.FunctionCall) goja.Value { return draft.modal(vm, call) })
	_ = api.Set("autocomplete", func(call goja.FunctionCall) goja.Value { return draft.autocomplete(vm, call) })
	_ = api.Set("configure", func(call goja.FunctionCall) goja.Value { return draft.configure(vm, call) })
	if _, err := builder(goja.Undefined(), api); err != nil {
		panic(vm.NewGoError(err))
	}
	return draft.finalize(vm)
}
