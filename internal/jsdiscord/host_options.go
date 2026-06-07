package jsdiscord

import (
	"context"
	"fmt"

	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
)

// HostOption customizes how a JavaScript bot host is created.
type HostOption func(*hostOptions) error

type RuntimeFactory interface {
	NewRuntime(ctx context.Context, opts ...require.Option) (*engine.Runtime, error)
}

type hostOptions struct {
	runtimeModuleRegistrars []engine.RuntimeModuleRegistrar
	runtimeFactory          RuntimeFactory
}

func applyHostOptions(opts ...HostOption) (hostOptions, error) {
	cfg := hostOptions{}
	for i, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&cfg); err != nil {
			return hostOptions{}, fmt.Errorf("apply host option %d: %w", i, err)
		}
	}
	return cfg, nil
}

// WithRuntimeModuleRegistrars appends per-runtime native module registrars.
func WithRuntimeModuleRegistrars(registrars ...engine.RuntimeModuleRegistrar) HostOption {
	return func(cfg *hostOptions) error {
		for i, registrar := range registrars {
			if registrar == nil {
				return fmt.Errorf("runtime module registrar at index %d is nil", i)
			}
		}
		cfg.runtimeModuleRegistrars = append(cfg.runtimeModuleRegistrars, registrars...)
		return nil
	}
}

func WithRuntimeFactory(factory RuntimeFactory) HostOption {
	return func(cfg *hostOptions) error {
		if factory == nil {
			return fmt.Errorf("runtime factory is nil")
		}
		cfg.runtimeFactory = factory
		return nil
	}
}
