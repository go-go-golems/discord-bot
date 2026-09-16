package jsslack

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

// Input optionality belongs to its containing block, never to an element.
func inputBlock(vm *goja.Runtime, call goja.FunctionCall) map[string]any {
	id, label := argString(call, 0), argString(call, 1)
	element := valueMap(vm, call.Argument(2))
	if id == "" || label == "" || element == nil {
		panic(vm.NewTypeError("slack/ui.input: block id, label and element required"))
	}
	block := map[string]any{"type": "input", "block_id": id, "label": plainText(label), "element": element}
	if opts := valueMap(vm, call.Argument(3)); opts != nil {
		for key, value := range opts {
			switch key {
			case "optional", "dispatch_action":
				if _, ok := value.(bool); !ok {
					panic(vm.NewTypeError("slack/ui.input: %s must be boolean", key))
				}
			case "hint":
				text, ok := value.(string)
				if !ok {
					panic(vm.NewTypeError("slack/ui.input: hint must be text"))
				}
				value = plainText(text)
			default:
				panic(vm.NewTypeError("slack/ui.input: unknown option %s", key))
			}
			block[key] = value
		}
	}
	return block
}
func plainText(text string) map[string]any {
	return map[string]any{"type": "plain_text", "text": text, "emoji": true}
}
func actionElement(kind string) bool {
	switch kind {
	case "button", "static_select", "multi_static_select", "external_select", "multi_external_select", "users_select", "multi_users_select", "channels_select", "multi_channels_select", "conversations_select", "multi_conversations_select", "datepicker", "timepicker", "checkboxes", "radio_buttons", "overflow":
		return true
	}
	return false
}

// Options use Slack's wire keys. Each helper supplies type/action_id and keeps
// SDK types outside JS; the raw Block Kit escape hatch remains available.
func registerRichUI(vm *goja.Runtime, exports *goja.Object) {
	set := func(name string, fn any) { must(vm, exports.Set(name, fn)) }
	set("option", func(call goja.FunctionCall) goja.Value {
		label, value := argString(call, 0), argString(call, 1)
		if label == "" || value == "" {
			panic(vm.NewTypeError("slack/ui.option: label and value required"))
		}
		return vm.ToValue(map[string]any{"text": plainText(label), "value": value})
	})
	for name, kind := range map[string]string{
		"staticSelect": "static_select", "multiStaticSelect": "multi_static_select",
		"externalSelect": "external_select", "multiExternalSelect": "multi_external_select",
		"usersSelect": "users_select", "multiUsersSelect": "multi_users_select",
		"channelsSelect": "channels_select", "multiChannelsSelect": "multi_channels_select",
		"conversationsSelect": "conversations_select", "multiConversationsSelect": "multi_conversations_select",
		"datePicker": "datepicker", "timePicker": "timepicker", "checkboxes": "checkboxes", "radioButtons": "radio_buttons", "overflow": "overflow",
	} {
		name, kind := name, kind
		set(name, func(call goja.FunctionCall) goja.Value {
			id := argString(call, 0)
			if id == "" {
				panic(vm.NewTypeError("slack/ui.%s: action id required", name))
			}
			element := valueMap(vm, call.Argument(1))
			if element == nil {
				element = map[string]any{}
			}
			element["type"], element["action_id"] = kind, id
			if placeholder, ok := element["placeholder"].(string); ok {
				element["placeholder"] = plainText(placeholder)
			}
			if kind == "static_select" || kind == "multi_static_select" || kind == "checkboxes" || kind == "radio_buttons" || kind == "overflow" {
				options, ok := element["options"].([]any)
				maximum := 100
				if kind == "overflow" {
					maximum = 5
				}
				if kind == "checkboxes" || kind == "radio_buttons" {
					maximum = 10
				}
				if !ok || len(options) < 1 || len(options) > maximum {
					panic(vm.NewTypeError("slack/ui.%s: provide 1–%d options", name, maximum))
				}
			}
			return vm.ToValue(element)
		})
	}
	set("context", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 || len(call.Arguments) > 10 {
			panic(vm.NewTypeError("slack/ui.context: provide 1–10 text/image elements"))
		}
		elements := []map[string]any{}
		for _, v := range call.Arguments {
			m := valueMap(vm, v)
			kind := fmt.Sprint(m["type"])
			if kind != "plain_text" && kind != "mrkdwn" && kind != "image" {
				panic(vm.NewTypeError("slack/ui.context: expected text or image"))
			}
			elements = append(elements, m)
		}
		return vm.ToValue(map[string]any{"type": "context", "elements": elements})
	})
	set("image", func(call goja.FunctionCall) goja.Value {
		url, alt := argString(call, 0), argString(call, 1)
		if (!strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://")) || alt == "" {
			panic(vm.NewTypeError("slack/ui.image: HTTP(S) URL and alt text required"))
		}
		return vm.ToValue(map[string]any{"type": "image", "image_url": url, "alt_text": alt})
	})
	set("linkButton", func(call goja.FunctionCall) goja.Value {
		id, label, url := argString(call, 0), argString(call, 1), argString(call, 2)
		if id == "" || label == "" || (!strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://")) {
			panic(vm.NewTypeError("slack/ui.linkButton: action id, label and HTTP(S) URL required"))
		}
		return vm.ToValue(map[string]any{"type": "button", "action_id": id, "text": plainText(label), "url": url})
	})
	set("confirm", func(call goja.FunctionCall) goja.Value {
		element := valueMap(vm, call.Argument(0))
		if element == nil || !actionElement(fmt.Sprint(element["type"])) {
			panic(vm.NewTypeError("slack/ui.confirm: expected interactive element"))
		}
		title, text := argString(call, 1), argString(call, 2)
		if title == "" || text == "" {
			panic(vm.NewTypeError("slack/ui.confirm: title and text required"))
		}
		yes, no := argString(call, 3), argString(call, 4)
		if yes == "" {
			yes = "Confirm"
		}
		if no == "" {
			no = "Cancel"
		}
		element["confirm"] = map[string]any{"title": plainText(title), "text": plainText(text), "confirm": plainText(yes), "deny": plainText(no)}
		return vm.ToValue(element)
	})
}
