package jsslack

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"

	"github.com/dop251/goja"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/pkg/errors"
)

// decode rejects unknown fields and non-JSON values rather than reflecting arbitrary Go objects.
func decode(vm *goja.Runtime, value goja.Value, target any) error {
	obj, ok := value.(*goja.Object)
	if !ok || obj.ClassName() != "Object" {
		return errors.New("expected an object")
	}
	b, err := value.ToObject(vm).MarshalJSON()
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	return dec.Decode(target)
}
func must(vm *goja.Runtime, err error) {
	if err != nil {
		panic(vm.NewGoError(err))
	}
}
func (h *Host) loader(vm *goja.Runtime, module *goja.Object) {
	exports := module.Get("exports").(*goja.Object)
	must(vm, exports.Set("defineBot", func(call goja.FunctionCall) goja.Value {
		if h.defined || !h.loading {
			panic(vm.NewTypeError("defineBot may be called once during loading"))
		}
		h.defined = true
		fn, ok := goja.AssertFunction(call.Argument(0))
		if !ok {
			panic(vm.NewTypeError("defineBot requires a function"))
		}
		api := vm.NewObject()
		must(vm, api.Set("configure", func(c goja.FunctionCall) goja.Value {
			if !h.loading || h.configured {
				panic(vm.NewTypeError("configure may be called once during loading"))
			}
			var cfg struct {
				Name        string             `json:"name"`
				Description string             `json:"description"`
				Run         slackbot.RunSchema `json:"run"`
			}
			must(vm, decode(vm, c.Argument(0), &cfg))
			h.descriptor.Name = cfg.Name
			h.descriptor.Description = cfg.Description
			h.descriptor.Run = cfg.Run
			h.configured = true
			return goja.Undefined()
		}))
		must(vm, api.Set("command", func(c goja.FunctionCall) goja.Value {
			name, ok := c.Argument(0).Export().(string)
			if !ok || !strings.HasPrefix(name, "/") || len(name) > 32 || strings.ContainsAny(name, " \t\n") || len(name) < 2 {
				panic(vm.NewTypeError("command requires a /name of 2–32 characters"))
			}
			var spec struct {
				Description string `json:"description"`
			}
			must(vm, decode(vm, c.Argument(1), &spec))
			if strings.TrimSpace(spec.Description) == "" {
				panic(vm.NewTypeError("command description is required"))
			}
			h.register(vm, "command:"+name, c.Argument(2))
			h.descriptor.Commands = append(h.descriptor.Commands, slackbot.Command{Name: name, Description: spec.Description})
			return goja.Undefined()
		}))
		must(vm, api.Set("event", func(c goja.FunctionCall) goja.Value {
			name, ok := c.Argument(0).Export().(string)
			if !ok || name != "app_mention" {
				panic(vm.NewTypeError("only app_mention is supported"))
			}
			h.register(vm, "event:"+name, c.Argument(1))
			h.descriptor.Events = append(h.descriptor.Events, name)
			return goja.Undefined()
		}))
		result, err := fn(goja.Undefined(), api)
		must(vm, err)
		if _, ok := result.Export().(*goja.Promise); ok {
			panic(vm.NewTypeError("bot registration must be synchronous"))
		}
		sort.Slice(h.descriptor.Commands, func(i, j int) bool { return h.descriptor.Commands[i].Name < h.descriptor.Commands[j].Name })
		sort.Strings(h.descriptor.Events)
		h.definition = vm.NewObject()
		return h.definition
	}))
}
func (h *Host) register(vm *goja.Runtime, key string, value goja.Value) {
	if !h.loading {
		panic(vm.NewTypeError("registration is closed"))
	}
	if _, ok := h.handlers[key]; ok {
		panic(vm.NewTypeError("duplicate registration: " + key))
	}
	fn, ok := goja.AssertFunction(value)
	if !ok {
		panic(vm.NewTypeError(key + ": handler must be a function"))
	}
	h.handlers[key] = fn
}
