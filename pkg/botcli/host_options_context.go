package botcli

import (
	"context"

	"github.com/go-go-golems/discord-bot/internal/jsdiscord"
)

type hostOptionsContextKey struct{}

// ContextWithHostOptions appends host options for a single command invocation.
// It is intended for command-provider adapters that need parsed command values
// when constructing host-managed bot runtimes.
func ContextWithHostOptions(ctx context.Context, opts ...jsdiscord.HostOption) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(opts) == 0 {
		return ctx
	}
	current := hostOptionsFromContext(ctx)
	combined := append(append([]jsdiscord.HostOption(nil), current...), opts...)
	return context.WithValue(ctx, hostOptionsContextKey{}, combined)
}

func hostOptionsFromContext(ctx context.Context) []jsdiscord.HostOption {
	if ctx == nil {
		return nil
	}
	opts, _ := ctx.Value(hostOptionsContextKey{}).([]jsdiscord.HostOption)
	return opts
}
