package jsslack

import (
	"context"
	"testing"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/stretchr/testify/require"
)

func TestCommandShortcutOptionsAndModalResults(t *testing.T) {
	rec := &slackbot.Recorder{}
	h := loadTestHost(t, `const {defineBot}=require("slack"),ui=require("slack/ui");module.exports=defineBot(({configure,command,shortcut,options,view})=>{
 configure({name:"test",scopes:["users:read"]});
 command("/test",{description:"Open"},async ctx=>{await ctx.openModal(ui.modal("edit","Edit").input("a","A",ui.textInput("text")).submit("Save").build());});
 shortcut({callbackId:"quote",type:"message",name:"Quote",description:"Quote message"},async ctx=>{await ctx.openModal(ui.modal("quote-result","Quote").block(ui.section(ui.plain(ctx.shortcut.message.text))).build());});
 options("search",async ctx=>{await ctx.ack.options([ui.option(ctx.query,"result")]);});
 view("edit",async ctx=>{await ctx.ack.update(ui.modal("done","Saved").block(ui.section(ui.plain(ctx.values.text("a","text")))).build());});
 });`, Options{Views: rec})
	cmd := command()
	cmd.TriggerID = "command-trigger"
	require.NoError(t, h.Dispatch(context.Background(), cmd, nil))
	require.Equal(t, "command-trigger", rec.Operations()[0].Trigger)
	require.NoError(t, h.Dispatch(context.Background(), slackbot.Invocation{TeamID: "T", UserID: "U", TriggerID: "shortcut-trigger", Shortcut: &slackbot.Shortcut{Type: "message_action", CallbackID: "quote", Message: map[string]any{"text": "Original"}}}, nil))
	require.Equal(t, "shortcut-trigger", rec.Operations()[1].Trigger)
	require.NoError(t, h.Dispatch(context.Background(), slackbot.Invocation{TeamID: "T", UserID: "U", Interaction: &slackbot.Interaction{Type: "block_suggestion", CallbackID: "search", Query: "Query", Ack: rec}}, nil))
	require.Equal(t, "Query", rec.Operations()[2].Ack.Options[0]["text"].(map[string]any)["text"])
	require.NoError(t, h.Dispatch(context.Background(), slackbot.Invocation{TeamID: "T", UserID: "U", Interaction: &slackbot.Interaction{Type: "view_submission", CallbackID: "edit", Ack: rec, Values: map[string]map[string]any{"a": {"text": map[string]any{"value": "Saved value"}}}}}, nil))
	require.Equal(t, "update", rec.Operations()[3].Ack.Kind)
	require.Equal(t, "done", rec.Operations()[3].Ack.View.CallbackID)
	require.Len(t, h.Descriptor().Shortcuts, 1)
	require.Equal(t, []string{"users:read"}, h.Descriptor().Scopes)
}
func TestReplaceOriginalRequiresAction(t *testing.T) {
	rec := &slackbot.Recorder{}
	h := loadTestHost(t, `module.exports=require("slack").defineBot(({configure,action})=>{configure({name:"test"});action("next",async ctx=>{await ctx.replaceOriginal({text:"Page 2"});});});`, Options{})
	require.NoError(t, h.Dispatch(context.Background(), slackbot.Invocation{TeamID: "T", UserID: "U", Action: &slackbot.Action{ActionID: "next"}}, rec))
	require.Equal(t, "replace_original", rec.Operations()[0].Kind)
}
