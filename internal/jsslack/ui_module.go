package jsslack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
)

// registerUILoader mounts the small Slack-native Block Kit construction API.
// Builders return detached JSON objects; no Slack SDK value crosses into JS.
func registerUILoader(reg *require.Registry) {
	reg.RegisterNativeModule("slack/ui", uiLoader)
}

func uiLoader(vm *goja.Runtime, module *goja.Object) {
	exports := module.Get("exports").(*goja.Object)
	set := func(name string, fn any) { must(vm, exports.Set(name, fn)) }
	set("plain", func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(map[string]any{"type": "plain_text", "text": argString(call, 0), "emoji": true})
	})
	set("mrkdwn", func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(map[string]any{"type": "mrkdwn", "text": argString(call, 0)})
	})
	set("button", func(call goja.FunctionCall) goja.Value {
		return newSlackButtonBuilder(vm, argString(call, 0), argString(call, 1))
	})
	set("section", func(call goja.FunctionCall) goja.Value {
		text := valueMap(vm, call.Argument(0))
		if text == nil {
			panic(vm.NewTypeError("slack/ui.section: text object is required"))
		}
		return vm.ToValue(map[string]any{"type": "section", "text": text})
	})
	set("actions", func(call goja.FunctionCall) goja.Value {
		blockID := argString(call, 0)
		if blockID == "" {
			panic(vm.NewTypeError("slack/ui.actions: block id is required"))
		}
		if len(call.Arguments) < 2 || len(call.Arguments)-1 > 25 {
			panic(vm.NewTypeError("slack/ui.actions: provide 1–25 elements"))
		}
		elements := make([]map[string]any, 0, len(call.Arguments)-1)
		for _, arg := range call.Arguments[1:] {
			m := valueMap(vm, arg)
			if m == nil || fmt.Sprint(m["type"]) != "button" {
				panic(vm.NewTypeError("slack/ui.actions: expected button builders"))
			}
			elements = append(elements, m)
		}
		return vm.ToValue(map[string]any{"type": "actions", "block_id": blockID, "elements": elements})
	})
	set("divider", func(goja.FunctionCall) goja.Value {
		return vm.ToValue(map[string]any{"type": "divider"})
	})
	set("header", func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(map[string]any{"type": "header", "text": map[string]any{"type": "plain_text", "text": argString(call, 0), "emoji": true}})
	})
	set("message", func(call goja.FunctionCall) goja.Value {
		return newSlackMessageBuilder(vm, argString(call, 0))
	})
}

type slackMessageBuilder struct {
	text   string
	blocks []map[string]any
}

func newSlackMessageBuilder(vm *goja.Runtime, text string) goja.Value {
	b := &slackMessageBuilder{text: text}
	obj := vm.NewObject()
	must(vm, obj.Set("block", func(call goja.FunctionCall) goja.Value {
		m := valueMap(vm, call.Argument(0))
		if m == nil || strings.TrimSpace(fmt.Sprint(m["type"])) == "" {
			panic(vm.NewTypeError("slack/ui.message.block: expected a Block Kit object"))
		}
		if len(b.blocks) >= 50 {
			panic(vm.NewTypeError("slack/ui.message: maximum 50 blocks exceeded"))
		}
		b.blocks = append(b.blocks, m)
		return obj
	}))
	must(vm, obj.Set("build", func(goja.FunctionCall) goja.Value {
		return vm.ToValue(map[string]any{"text": b.text, "blocks": b.blocks})
	}))
	return obj
}

type slackButtonBuilder struct {
	actionID string
	label    string
	value    string
	style    string
}

func newSlackButtonBuilder(vm *goja.Runtime, actionID, label string) goja.Value {
	if actionID == "" || label == "" {
		panic(vm.NewTypeError("slack/ui.button: action id and label are required"))
	}
	b := &slackButtonBuilder{actionID: actionID, label: label, style: "primary"}
	obj := vm.NewObject()
	must(vm, obj.Set("value", func(call goja.FunctionCall) goja.Value {
		b.value = argString(call, 0)
		return obj
	}))
	must(vm, obj.Set("style", func(call goja.FunctionCall) goja.Value {
		style := strings.TrimSpace(strings.ToLower(argString(call, 0)))
		if style != "primary" && style != "danger" && style != "" {
			panic(vm.NewTypeError("slack/ui.button.style: expected primary or danger"))
		}
		if style != "" {
			b.style = style
		}
		return obj
	}))
	must(vm, obj.Set("build", func(goja.FunctionCall) goja.Value {
		m := map[string]any{"type": "button", "action_id": b.actionID, "text": map[string]any{"type": "plain_text", "text": b.label, "emoji": true}, "style": b.style}
		if b.value != "" {
			m["value"] = b.value
		}
		return vm.ToValue(m)
	}))
	return obj
}

func valueMap(vm *goja.Runtime, value goja.Value) map[string]any {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil
	}
	if obj, ok := value.(*goja.Object); ok {
		if build := obj.Get("build"); build != nil && !goja.IsUndefined(build) {
			if fn, ok := goja.AssertFunction(build); ok {
				built, err := fn(goja.Undefined())
				if err != nil {
					panic(err)
				}
				value = built
			}
		}
	}
	obj, ok := value.(*goja.Object)
	if !ok {
		return nil
	}
	b, err := obj.MarshalJSON()
	if err != nil {
		panic(vm.NewGoError(err))
	}
	var out map[string]any
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&out); err != nil {
		panic(vm.NewGoError(err))
	}
	return out
}

var _ engine.RuntimeModuleRegistrar = (*registrar)(nil)

func argString(call goja.FunctionCall, index int) string {
	if index >= len(call.Arguments) || call.Arguments[index] == nil || goja.IsUndefined(call.Arguments[index]) || goja.IsNull(call.Arguments[index]) {
		return ""
	}
	return call.Arguments[index].String()
}
