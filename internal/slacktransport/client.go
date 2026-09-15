// Package slacktransport connects the Slack host to Slack Socket Mode and Web API services.
package slacktransport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// LocalOptions never enables external networking. Tokens are supplied explicitly by Go callers.
type LocalOptions struct {
	APIURL, BotToken, AppToken, TeamID, AppID string
	AllowedChannels                           []string
	Logger                                    zerolog.Logger
}

// RemoteOptions configures a real Slack workspace connection. Tokens are
// supplied explicitly by the caller; this package does not read environment
// variables or credential files.
type RemoteOptions struct {
	BotToken, AppToken, TeamID, AppID string
	AllowedChannels                   []string
	Logger                            zerolog.Logger
}
type Client struct {
	api           *slack.Client
	socket        *socketmode.Client
	http          *http.Client
	transport     *http.Transport
	origin        *url.URL
	responseHosts map[string]struct{}
	ackMu         sync.Mutex
	ackReplay     map[string]ackReplay
	opts          LocalOptions
}

type ackReplay struct {
	payload any
	expires time.Time
}

// jsonBlock preserves a validated framework block while satisfying the Slack
// SDK's Block interface. The SDK's UnknownBlock intentionally drops fields it
// does not know, so the transport uses this small lossless wrapper instead.
type jsonBlock struct{ slackbot.Block }

func (b jsonBlock) BlockType() slack.MessageBlockType { return slack.MessageBlockType(b.Type()) }
func (b jsonBlock) ID() string {
	if id, ok := b.Block["block_id"].(string); ok {
		return id
	}
	return ""
}
func (b jsonBlock) MarshalJSON() ([]byte, error) { return json.Marshal(map[string]any(b.Block)) }

var _ slackbot.MessageService = (*Client)(nil)
var _ slackbot.ViewService = (*Client)(nil)

func NewLocal(opts LocalOptions) (*Client, error) {
	u, err := url.Parse(opts.APIURL)
	if err != nil || u.Scheme != "http" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Port() == "" || !strings.HasSuffix(u.Path, "/") {
		return nil, errors.New("api URL must be a loopback HTTP URL with explicit port and trailing slash")
	}
	ip := net.ParseIP(u.Hostname())
	if ip == nil || !ip.IsLoopback() {
		return nil, errors.New("api URL requires a literal loopback address")
	}
	if opts.BotToken == "" || opts.AppToken == "" || opts.TeamID == "" || opts.AppID == "" {
		return nil, errors.New("bot token, app token, team ID and app ID are required")
	}
	opts.AllowedChannels = append([]string(nil), opts.AllowedChannels...)
	c := &Client{origin: u, responseHosts: map[string]struct{}{u.Host: {}}, ackReplay: map[string]ackReplay{}, opts: opts}
	dial := func(ctx context.Context, network, address string) (net.Conn, error) {
		if address != u.Host {
			return nil, errors.New("local transport refused unexpected destination")
		}
		return (&net.Dialer{Timeout: 3 * time.Second}).DialContext(ctx, network, address)
	}
	c.transport = &http.Transport{DialContext: dial}
	c.http = &http.Client{Transport: c.transport, Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("local transport refuses redirects") }}
	c.api = slack.New(opts.BotToken, slack.OptionAppLevelToken(opts.AppToken), slack.OptionAPIURL(opts.APIURL), slack.OptionHTTPClient(c.http))
	c.socket = socketmode.New(c.api, socketmode.OptionDialer(&websocket.Dialer{NetDialContext: dial, HandshakeTimeout: 5 * time.Second}))
	return c, nil
}

// NewRemote creates a client for Slack's public API and Socket Mode endpoints.
func NewRemote(opts RemoteOptions) (*Client, error) {
	if opts.BotToken == "" || opts.AppToken == "" || opts.TeamID == "" || opts.AppID == "" {
		return nil, errors.New("bot token, app token, team ID and app ID are required")
	}
	httpClient := &http.Client{Timeout: 30 * time.Second}
	api := slack.New(opts.BotToken, slack.OptionAppLevelToken(opts.AppToken), slack.OptionHTTPClient(httpClient))
	return &Client{
		api: api, socket: socketmode.New(api), http: httpClient,
		responseHosts: map[string]struct{}{"hooks.slack.com": {}, "hooks.slack-gov.com": {}}, ackReplay: map[string]ackReplay{},
		opts: LocalOptions{BotToken: opts.BotToken, AppToken: opts.AppToken, TeamID: opts.TeamID, AppID: opts.AppID, AllowedChannels: append([]string(nil), opts.AllowedChannels...), Logger: opts.Logger},
	}, nil
}

func (c *Client) rememberAck(id string, payload any) {
	c.ackMu.Lock()
	defer c.ackMu.Unlock()
	now := time.Now()
	for key, entry := range c.ackReplay {
		if !entry.expires.After(now) {
			delete(c.ackReplay, key)
		}
	}
	c.ackReplay[id] = ackReplay{payload: payload, expires: now.Add(5 * time.Minute)}
}

func (c *Client) replayAck(id string) (any, bool) {
	c.ackMu.Lock()
	defer c.ackMu.Unlock()
	entry, ok := c.ackReplay[id]
	if !ok || !entry.expires.After(time.Now()) {
		if ok {
			delete(c.ackReplay, id)
		}
		return nil, false
	}
	return entry.payload, true
}
func (c *Client) Close() { c.transport.CloseIdleConnections() }
func (c *Client) Post(ctx context.Context, m slackbot.PostMessage) (slackbot.MessageRef, error) {
	if err := m.Validate(); err != nil {
		return slackbot.MessageRef{}, err
	}
	options := []slack.MsgOption{slack.MsgOptionText(m.Text, false)}
	if m.Blocks != nil {
		blocks := make([]slack.Block, 0, len(m.Blocks))
		for _, block := range m.Blocks {
			blocks = append(blocks, jsonBlock{Block: block})
		}
		options = append(options, slack.MsgOptionBlocks(blocks...))
	}
	if m.ThreadTS != "" {
		options = append(options, slack.MsgOptionTS(m.ThreadTS))
	}
	channel, ts, err := c.api.PostMessageContext(ctx, m.ChannelID, options...)
	if err != nil {
		return slackbot.MessageRef{}, safeError(ctx, err, "messages.post")
	}
	return slackbot.MessageRef{ChannelID: channel, TS: ts}, nil
}

func (c *Client) Open(ctx context.Context, triggerID string, view slackbot.ModalView) (slackbot.ViewRef, error) {
	if err := view.Validate(); err != nil {
		return slackbot.ViewRef{}, err
	}
	blocks := make([]slack.Block, 0, len(view.Blocks))
	for _, block := range view.Blocks {
		blocks = append(blocks, jsonBlock{Block: block})
	}
	request := slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           textBlock(view.Title),
		Blocks:          slack.Blocks{BlockSet: blocks},
		PrivateMetadata: view.PrivateMetadata,
		CallbackID:      view.CallbackID,
	}
	if view.Close != nil {
		request.Close = textBlock(view.Close)
	}
	if view.Submit != nil {
		request.Submit = textBlock(view.Submit)
	}
	response, err := c.api.OpenViewContext(ctx, triggerID, request)
	if err != nil {
		return slackbot.ViewRef{}, safeError(ctx, err, "views.open")
	}
	return slackbot.ViewRef{ID: response.View.ID, Hash: response.View.Hash}, nil
}

func textBlock(block slackbot.Block) *slack.TextBlockObject {
	t := "plain_text"
	if value, ok := block["type"].(string); ok && value != "" {
		t = value
	}
	text, _ := block["text"].(string)
	emoji := true
	return &slack.TextBlockObject{Type: t, Text: text, Emoji: &emoji}
}

// Only reviewed error codes cross the service boundary; SDK error strings can contain remote content.
func safeError(ctx context.Context, err error, operation string) error {
	if ctx.Err() != nil {
		return slackbot.Fail("context_closed", operation, "operation canceled")
	}
	var limited *slack.RateLimitedError
	if errors.As(err, &limited) {
		return slackbot.Fail("rate_limited", operation, "service rate limit reached")
	}
	switch err.Error() {
	case "channel_not_found", "not_in_channel", "invalid_auth", "missing_scope", "token_revoked", "is_archived":
		return slackbot.Fail(err.Error(), operation, "Slack rejected the operation")
	}
	return slackbot.Fail("delivery_unknown", operation, "operation outcome could not be confirmed")
}

type responder struct {
	client *Client
	target string
}

var _ slackbot.Responder = (*responder)(nil)

func (r *responder) Reply(ctx context.Context, text slackbot.Text) error {
	return r.ReplyMessage(ctx, slackbot.MessagePayload{Text: text.Text})
}

func (r *responder) ReplyMessage(ctx context.Context, message slackbot.MessagePayload) error {
	if err := message.Validate("reply"); err != nil {
		return err
	}
	bodyValue := map[string]any{"response_type": "ephemeral", "text": message.Text}
	if message.Blocks != nil {
		blocks := make([]map[string]any, len(message.Blocks))
		for i, block := range message.Blocks {
			blocks[i] = map[string]any(block)
		}
		bodyValue["blocks"] = blocks
	}
	body, err := json.Marshal(bodyValue)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.target, bytes.NewReader(body))
	if err != nil {
		return slackbot.Fail("invalid_argument", "reply", "invalid response capability")
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := r.client.http.Do(req)
	if err != nil {
		return safeError(ctx, err, "reply")
	}
	defer func() { _ = response.Body.Close() }()
	_, err = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
	if err != nil {
		return safeError(ctx, err, "reply")
	}
	if response.StatusCode == http.StatusTooManyRequests {
		return slackbot.Fail("rate_limited", "reply", "service rate limit reached")
	}
	if response.StatusCode != http.StatusOK {
		return slackbot.Fail("delivery_unknown", "reply", "response could not be confirmed")
	}
	return nil
}
func (c *Client) responseCapability(target string) (slackbot.Responder, error) {
	u, err := url.Parse(target)
	if err != nil || u.User != nil || u.Fragment != "" || u.Host == "" || u.Path == "" {
		return nil, errors.New("invalid response capability")
	}
	if c.origin != nil {
		if u.Scheme != "http" || u.Host != c.origin.Host {
			return nil, errors.New("invalid local response capability")
		}
	} else if u.Scheme != "https" {
		return nil, errors.New("invalid Slack response capability")
	}
	if _, ok := c.responseHosts[u.Host]; !ok {
		return nil, errors.New("invalid local response capability")
	}
	return &responder{client: c, target: target}, nil
}
