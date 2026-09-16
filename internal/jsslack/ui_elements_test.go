package jsslack

import (
	"encoding/json"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
)

func buildUI(t *testing.T, expression string) map[string]any {
	t.Helper()
	vm := goja.New()
	module := vm.NewObject()
	exports := vm.NewObject()
	require.NoError(t, module.Set("exports", exports))
	uiLoader(vm, module)
	require.NoError(t, vm.Set("ui", exports))
	result, err := vm.RunString(expression)
	require.NoError(t, err)
	raw, err := result.ToObject(vm).MarshalJSON()
	require.NoError(t, err)
	var out map[string]any
	require.NoError(t, json.Unmarshal(raw, &out))
	return out
}
func TestOptionalInputIsBlockProperty(t *testing.T) {
	for _, expression := range []string{
		`ui.input("notes","Notes",ui.textInput("text").multiline().length(2,200),{optional:true,hint:"Details"})`,
		`ui.modal("edit","Edit").input("notes","Notes",ui.textInput("text").multiline().length(2,200),{optional:true,hint:"Details"}).submit("Save").build().blocks[0]`,
	} {
		got := buildUI(t, expression)
		require.Equal(t, true, got["optional"])
		element := got["element"].(map[string]any)
		require.NotContains(t, element, "optional")
		require.Equal(t, true, element["multiline"])
		require.Equal(t, float64(2), element["min_length"])
		require.Equal(t, float64(200), element["max_length"])
	}
}
func TestSelectsPreserveWireShape(t *testing.T) {
	got := buildUI(t, `ui.actions("controls",ui.staticSelect("choice",{placeholder:"Choose",options:[ui.option("One","1")]}),ui.usersSelect("owner"),ui.datePicker("date"))`)
	elements := got["elements"].([]any)
	require.Len(t, elements, 3)
	selectElement := elements[0].(map[string]any)
	require.Equal(t, "static_select", selectElement["type"])
	require.Equal(t, "plain_text", selectElement["placeholder"].(map[string]any)["type"])
	require.Equal(t, "1", selectElement["options"].([]any)[0].(map[string]any)["value"])
	require.Equal(t, "users_select", elements[1].(map[string]any)["type"])
}
func TestRichLayoutAndConfirmation(t *testing.T) {
	got := buildUI(t, `ui.message("Review").block(ui.section(ui.mrkdwn("Body"),ui.image("https://example.test/image.png","Preview"))).block(ui.context(ui.plain("Details"))).block(ui.actions("review",ui.confirm(ui.button("delete","Delete").style("danger"),"Delete?","Remove this entry?"),ui.linkButton("source","Source","https://example.test"))).build()`)
	blocks := got["blocks"].([]any)
	require.Len(t, blocks, 3)
	actions := blocks[2].(map[string]any)["elements"].([]any)
	button := actions[0].(map[string]any)
	require.Equal(t, "danger", button["style"])
	require.Equal(t, "Confirm", button["confirm"].(map[string]any)["confirm"].(map[string]any)["text"])
	require.Equal(t, "image", blocks[0].(map[string]any)["accessory"].(map[string]any)["type"])
	require.NotContains(t, buildUI(t, `ui.button("neutral","Neutral").build()`), "style")
}
