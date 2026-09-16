package jsslack

import (
	"github.com/dop251/goja"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"strings"
)

func (h *Host) addOperations(vm *goja.Runtime, s *invocationState, slack, messages *goja.Object) {
	groups := map[string]*goja.Object{"messages": messages}
	for operation := range slackbot.OperationMethods() {
		parts := strings.SplitN(operation, ".", 2)
		group := groups[parts[0]]
		if group == nil {
			group = vm.NewObject()
			groups[parts[0]] = group
			must(vm, slack.Set(parts[0], group))
		}
		operation := operation
		must(vm, group.Set(parts[1], func(call goja.FunctionCall) goja.Value {
			var params map[string]any
			must(vm, decode(vm, call.Argument(0), &params))
			if _, ok := params["token"]; ok {
				panic(vm.NewTypeError("token is host-owned"))
			}
			return h.async(vm, s, operation, func() (any, error) {
				if h.operations == nil {
					return nil, slackbot.Fail("unavailable", operation, "no operation service")
				}
				return h.operations.Call(s.ctx, operation, params)
			})
		}))
	}
}
