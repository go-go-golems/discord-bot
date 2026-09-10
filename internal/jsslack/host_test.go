package jsslack

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/stretchr/testify/require"
)

type fakeServices struct {
	mu      sync.Mutex
	posts   []slackbot.PostMessage
	replies []slackbot.Text
	post    func(context.Context, slackbot.PostMessage) (slackbot.MessageRef, error)
}

var _ slackbot.MessageService = (*fakeServices)(nil)
var _ slackbot.Responder = (*fakeServices)(nil)

func (f *fakeServices) Post(ctx context.Context, p slackbot.PostMessage) (slackbot.MessageRef, error) {
	if f.post != nil {
		return f.post(ctx, p)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.posts = append(f.posts, p)
	return slackbot.MessageRef{ChannelID: p.ChannelID, TS: "1741234567.000999"}, nil
}
func (f *fakeServices) Reply(_ context.Context, p slackbot.Text) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.replies = append(f.replies, p)
	return nil
}
func script(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "index.js")
	require.NoError(t, os.WriteFile(p, []byte(body), 0600))
	return p
}
func loadTestHost(t *testing.T, body string, opts Options) *Host {
	t.Helper()
	h, err := Load(context.Background(), script(t, body), opts)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, h.Close(context.Background())) })
	return h
}
func bot(handler string) string {
	return `const {defineBot}=require("slack");module.exports=defineBot(({configure,command})=>{configure({name:"test"});command("/test",{description:"test"},` + handler + `);});`
}
func command() slackbot.Invocation {
	return slackbot.Invocation{ID: "i", TeamID: "T", ChannelID: "C", UserID: "U", Command: "/test"}
}
func TestInspectOfflineExample(t *testing.T) {
	p, err := filepath.Abs("../../examples/slack-bots/ping/index.js")
	require.NoError(t, err)
	h, err := Load(context.Background(), p, Options{})
	require.NoError(t, err)
	defer func() { _ = h.Close(context.Background()) }()
	d := h.Descriptor()
	require.Equal(t, "ping", d.Name)
	require.Len(t, d.Commands, 1)
	require.Equal(t, []string{"app_mention"}, d.Events)
	d.Run.Fields["greeting"] = slackbot.Field{Type: "bool"}
	require.Equal(t, "string", h.Descriptor().Run.Fields["greeting"].Type)
}
func TestRepliesAndThreadPrecision(t *testing.T) {
	f := &fakeServices{}
	h := loadTestHost(t, `const {defineBot}=require("slack");module.exports=defineBot(({configure,command,event})=>{
 configure({name:"test",run:{fields:{greeting:{type:"string",default:"pong"}}}});
 command("/test",{description:"test"},async ctx=>({text:ctx.config.greeting+":"+Object.keys(ctx.config).join(",")}));
 event("app_mention",async ctx=>{await ctx.reply({text:"mention"});});});`, Options{Messages: f, Config: map[string]any{"botToken": "secret", "appToken": "secret2"}})
	require.NoError(t, h.Dispatch(context.Background(), command(), f))
	require.Equal(t, "pong:greeting", f.replies[0].Text)
	i := command()
	i.Command = ""
	i.Event = "app_mention"
	i.TS = "1741234567.000123"
	require.NoError(t, h.Dispatch(context.Background(), i, nil))
	require.Equal(t, i.TS, f.posts[0].ThreadTS)
	i.ThreadTS = "1741234567.000001"
	require.NoError(t, h.Dispatch(context.Background(), i, nil))
	require.Equal(t, i.ThreadTS, f.posts[1].ThreadTS)
}
func TestReplyTwiceFailsWithoutSecondSend(t *testing.T) {
	f := &fakeServices{}
	h := loadTestHost(t, bot(`async ctx=>{await ctx.reply({text:"one"});return {text:"two"};}`), Options{})
	require.ErrorContains(t, h.Dispatch(context.Background(), command(), f), "already_replied")
	require.Len(t, f.replies, 1)
}
func TestAsyncPostAndRejection(t *testing.T) {
	f := &fakeServices{}
	h := loadTestHost(t, bot(`async ctx=>{const m=await ctx.slack.messages.post({channelId:"C",text:"hello",threadTs:"123.000001"});return {text:m.ts};}`), Options{Messages: f})
	require.NoError(t, h.Dispatch(context.Background(), command(), f))
	require.Equal(t, "1741234567.000999", f.replies[0].Text)
	bad := loadTestHost(t, bot(`async ()=>{throw new Error("meaningful rejection");}`), Options{})
	require.ErrorContains(t, bad.Dispatch(context.Background(), command(), f), "meaningful rejection")
}
func TestRegistrationFailures(t *testing.T) {
	for _, tc := range []struct{ name, source, contains string }{
		{"duplicate", `module.exports=require("slack").defineBot(({configure,event})=>{configure({name:"x"});event("app_mention",()=>{});event("app_mention",()=>{});});`, "duplicate"},
		{"unknown event", `module.exports=require("slack").defineBot(({configure,event})=>{configure({name:"x"});event("message",()=>{});});`, "only app_mention"},
		{"async registration", `module.exports=require("slack").defineBot(async()=>{});`, "synchronous"},
		{"invalid schema", `module.exports=require("slack").defineBot(({configure})=>configure({name:"x",run:{fields:{a:{type:"bogus"}}}}));`, "unsupported type"},
		{"no fs capability", `require("fs");`, "module"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load(context.Background(), script(t, tc.source), Options{})
			require.ErrorContains(t, err, tc.contains)
		})
	}
}
func TestInspectionAndHandlerDeadlines(t *testing.T) {
	_, err := Load(context.Background(), script(t, `while(true){}`), Options{Timeout: 30 * time.Millisecond})
	require.Error(t, err)
	f := &fakeServices{}
	h := loadTestHost(t, bot(`ctx=>{if(ctx.text==="loop"){while(true){}}return {text:"ok"};}`), Options{Timeout: 40 * time.Millisecond})
	i := command()
	i.Text = "loop"
	require.Error(t, h.Dispatch(context.Background(), i, f))
	i.Text = "ok"
	require.NoError(t, h.Dispatch(context.Background(), i, f))
	require.Len(t, f.replies, 1)
}
func TestCancelPendingNetworkAndReuse(t *testing.T) {
	f := &fakeServices{post: func(ctx context.Context, _ slackbot.PostMessage) (slackbot.MessageRef, error) {
		<-ctx.Done()
		return slackbot.MessageRef{}, ctx.Err()
	}}
	h := loadTestHost(t, bot(`async ctx=>{await ctx.slack.messages.post({channelId:"C",text:"hello"});}`), Options{Messages: f, Timeout: 40 * time.Millisecond})
	require.Error(t, h.Dispatch(context.Background(), command(), nil))
	require.NoError(t, h.Close(context.Background()))
	require.Error(t, h.Dispatch(context.Background(), command(), nil))
}
func TestUnknownPayloadFieldAndStaleContext(t *testing.T) {
	h := loadTestHost(t, bot(`ctx=>ctx.reply({text:"hello",threadTS:"bad"})`), Options{})
	require.ErrorContains(t, h.Dispatch(context.Background(), command(), &fakeServices{}), "unknown field")
	s := loadTestHost(t, `let old;module.exports=require("slack").defineBot(({configure,command})=>{configure({name:"x"});command("/test",{description:"test"},ctx=>{if(old)return old.reply({text:"stale"});old=ctx;});});`, Options{})
	require.NoError(t, s.Dispatch(context.Background(), command(), nil))
	require.ErrorContains(t, s.Dispatch(context.Background(), command(), nil), "context_closed")
}

func TestStoreIsDetachedAndWorkspaceScoped(t *testing.T) {
	f := &fakeServices{}
	h := loadTestHost(t, bot(`ctx=>{const old=ctx.store.get("x");const value={n:1};ctx.store.set("x",value);value.n=99;const read=ctx.store.get("x");read.n=88;return {text:String(old===undefined)+":"+ctx.store.get("x").n};}`), Options{})
	require.NoError(t, h.Dispatch(context.Background(), command(), f))
	require.NoError(t, h.Dispatch(context.Background(), command(), f))
	i := command()
	i.TeamID = "other"
	require.NoError(t, h.Dispatch(context.Background(), i, f))
	require.Equal(t, []slackbot.Text{{Text: "true:1"}, {Text: "false:1"}, {Text: "true:1"}}, f.replies)
}

func TestInspectionAllowsRequiredConfig(t *testing.T) {
	p := script(t, `module.exports=require("slack").defineBot(({configure})=>configure({name:"x",run:{fields:{required:{type:"string",required:true}}}}));`)
	_, err := Inspect(context.Background(), p, time.Second)
	require.NoError(t, err)
	_, err = Load(context.Background(), p, Options{})
	require.ErrorContains(t, err, "required")
}

func TestStoreNullAndScalar(t *testing.T) {
	f := &fakeServices{}
	h := loadTestHost(t, bot(`ctx=>{ctx.store.set("n",null);ctx.store.set("s","hello");return {text:String(ctx.store.get("n"))+":"+ctx.store.get("s")};}`), Options{})
	require.NoError(t, h.Dispatch(context.Background(), command(), f))
	require.Equal(t, "null:hello", f.replies[0].Text)
}
