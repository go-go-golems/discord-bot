package jsslack

import (
	"context"
	"encoding/json"
	"github.com/dop251/goja"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/pkg/errors"
)

// InvokeVerb runs a local synchronous metadata verb. No transport services or
// credentials are supplied. Bot execution remains the CLI's run operation.
func (h *Host) InvokeVerb(ctx context.Context, name string) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()
	stop := context.AfterFunc(h.lifetime, cancel)
	defer stop()
	select {
	case h.gate <- struct{}{}:
		defer func() { <-h.gate }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if h.lifetime.Err() != nil {
		return nil, slackbot.Fail("context_closed", "verb", "host is closed")
	}
	return h.call(ctx, "slack.verb", func(vm *goja.Runtime) (any, error) {
		fn, ok := h.handlers["verb:"+name]
		if !ok {
			return nil, errors.Errorf("local verb %q not found", name)
		}
		raw, _ := json.Marshal(h.config)
		var config any
		_ = json.Unmarshal(raw, &config)
		value, err := fn(goja.Undefined(), vm.ToValue(map[string]any{"config": config}))
		if err != nil {
			return nil, err
		}
		if _, ok := value.Export().(*goja.Promise); ok {
			return nil, errors.New("local verbs must be synchronous")
		}
		if goja.IsUndefined(value) || goja.IsNull(value) {
			return nil, nil
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		var result any
		err = json.Unmarshal(encoded, &result)
		return result, err
	})
}
