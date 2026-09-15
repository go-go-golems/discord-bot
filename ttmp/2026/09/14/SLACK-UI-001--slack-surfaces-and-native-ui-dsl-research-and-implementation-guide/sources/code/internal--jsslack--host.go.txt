// Package jsslack owns the Slack JavaScript runtime. No transport credentials enter it.
package jsslack

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Options struct {
	Messages slackbot.MessageService
	Config   map[string]any
	Timeout  time.Duration
	Logger   zerolog.Logger
}
type Host struct {
	store      map[string]map[string]json.RawMessage // owner-only, indexed by workspace
	runtime    *engine.Runtime
	cancel     context.CancelFunc
	lifetime   context.Context
	gate       chan struct{}
	timeout    time.Duration
	messages   slackbot.MessageService
	logger     zerolog.Logger
	descriptor slackbot.Descriptor
	config     map[string]any
	handlers   map[string]goja.Callable // owner-only
	definition *goja.Object             // owner-only
	configured bool
	defined    bool
	loading    bool
	closeOnce  sync.Once
	closeErr   error
}

// Load constructs a runtime exposing only slack and local CommonJS modules.
func Load(ctx context.Context, path string, opts Options) (*Host, error) {
	return load(ctx, path, opts, false)
}

// Inspect executes declarations without resolving required runtime configuration.
func Inspect(ctx context.Context, path string, timeout time.Duration) (slackbot.Descriptor, error) {
	h, err := load(ctx, path, Options{Timeout: timeout}, true)
	if err != nil {
		return slackbot.Descriptor{}, err
	}
	defer func() { _ = h.Close(context.Background()) }()
	return h.Descriptor(), nil
}
func load(ctx context.Context, path string, opts Options, inspect bool) (*Host, error) {
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Second
	}
	if opts.Timeout < 0 {
		return nil, errors.New("timeout must be positive")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "script path")
	}
	lifetime, cancel := context.WithCancel(ctx)
	h := &Host{cancel: cancel, lifetime: lifetime, gate: make(chan struct{}, 1), timeout: opts.Timeout, messages: opts.Messages, logger: opts.Logger, handlers: map[string]goja.Callable{}, loading: true}
	h.store = map[string]map[string]json.RawMessage{}
	h.descriptor.ScriptPath = abs
	factory, err := engine.NewRuntimeFactoryBuilder(engine.WithImplicitDefaultRegistryModules(false), engine.WithDataOnlyDefaultRegistryModules(false)).WithModules(&registrar{h}).Build()
	if err != nil {
		cancel()
		return nil, err
	}
	rt, err := factory.NewRuntime(engine.WithStartupContext(ctx), engine.WithLifetimeContext(lifetime))
	if err != nil {
		cancel()
		return nil, err
	}
	h.runtime = rt
	loadCtx, stop := context.WithTimeout(lifetime, h.timeout)
	defer stop()
	_, err = h.call(loadCtx, "slack.load", func(vm *goja.Runtime) (any, error) {
		value, err := rt.Require.Require(abs)
		if err != nil {
			return nil, err
		}
		if !h.defined || value != h.definition {
			return nil, errors.New("script must export its defineBot result")
		}
		h.loading = false
		if err := h.descriptor.Validate(); err != nil {
			return nil, err
		}
		if !inspect {
			h.config, err = h.descriptor.Config(opts.Config)
		}
		return nil, err
	})
	if err != nil {
		_ = h.Close(context.Background())
		return nil, errors.Wrap(err, "load Slack bot "+abs)
	}
	h.logger.Debug().Str("bot", h.descriptor.Name).Msg("loaded Slack bot")
	return h, nil
}
func (h *Host) Descriptor() slackbot.Descriptor {
	// The descriptor is immutable after loading; deep copy for callers.
	b, _ := json.Marshal(h.descriptor)
	var d slackbot.Descriptor
	_ = json.Unmarshal(b, &d)
	return d
}
func (h *Host) Close(ctx context.Context) error {
	h.cancel()
	select {
	case h.gate <- struct{}{}:
		defer func() { <-h.gate }()
	case <-ctx.Done():
		return ctx.Err()
	}
	h.closeOnce.Do(func() { h.closeErr = h.runtime.Close(ctx) })
	return h.closeErr
}

// call interrupts CPU-bound JS on cancellation, and clears the interrupt before the next owner entry.
func (h *Host) call(ctx context.Context, op string, fn func(*goja.Runtime) (any, error)) (any, error) {
	return h.runtime.Owner.Call(ctx, op, func(_ context.Context, vm *goja.Runtime) (any, error) {
		done := make(chan struct{})
		stop := context.AfterFunc(ctx, func() { vm.Interrupt(ctx.Err()); close(done) })
		defer func() {
			if !stop() {
				<-done
			}
			vm.ClearInterrupt()
		}()
		return fn(vm)
	})
}

type registrar struct{ host *Host }

var _ engine.RuntimeModuleRegistrar = (*registrar)(nil)

func (*registrar) ID() string { return "slack" }
func (r *registrar) RegisterRuntimeModule(_ *engine.RuntimeModuleRegistrationContext, reg *require.Registry) error {
	reg.RegisterNativeModule("slack", r.host.loader)
	return nil
}
