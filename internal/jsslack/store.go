package jsslack

import (
	"encoding/json"
	"sort"

	"github.com/dop251/goja"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
)

func (h *Host) storeObject(vm *goja.Runtime, s *invocationState) *goja.Object {
	obj := vm.NewObject()
	access := func(call goja.FunctionCall) (map[string]json.RawMessage, string) {
		if s.ctx.Err() != nil {
			panic(jsError(vm, slackbot.Fail("context_closed", "store", "invocation is closed")))
		}
		key, ok := call.Argument(0).Export().(string)
		if !ok || key == "" {
			panic(vm.NewTypeError("store key must be a nonempty string"))
		}
		m := h.store[s.input.TeamID]
		if m == nil {
			m = map[string]json.RawMessage{}
			h.store[s.input.TeamID] = m
		}
		return m, key
	}
	must(vm, obj.Set("get", func(c goja.FunctionCall) goja.Value {
		m, k := access(c)
		b, ok := m[k]
		if !ok {
			return goja.Undefined()
		}
		var v any
		must(vm, json.Unmarshal(b, &v))
		return vm.ToValue(v)
	}))
	must(vm, obj.Set("set", func(c goja.FunctionCall) goja.Value {
		m, k := access(c)
		b, err := c.Argument(1).ToObject(vm).MarshalJSON()
		must(vm, err)
		if len(b) == 0 {
			panic(vm.NewTypeError("store value must be JSON serializable"))
		}
		m[k] = append(json.RawMessage(nil), b...)
		return goja.Undefined()
	}))
	must(vm, obj.Set("delete", func(c goja.FunctionCall) goja.Value {
		m, k := access(c)
		_, ok := m[k]
		delete(m, k)
		return vm.ToValue(ok)
	}))
	must(vm, obj.Set("keys", func(goja.FunctionCall) goja.Value {
		if s.ctx.Err() != nil {
			panic(jsError(vm, slackbot.Fail("context_closed", "store", "invocation is closed")))
		}
		keys := []string{}
		for k := range h.store[s.input.TeamID] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return vm.ToValue(keys)
	}))
	return obj
}
