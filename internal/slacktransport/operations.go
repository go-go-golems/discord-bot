package slacktransport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
)

var _ slackbot.OperationService = (*Client)(nil)

func (c *Client) Call(ctx context.Context, operation string, params map[string]any) (map[string]any, error) {
	method, ok := slackbot.OperationMethods()[operation]
	if !ok {
		return nil, slackbot.Fail("invalid_argument", operation, "unknown Slack operation")
	}
	if operation == "files.upload" {
		return c.uploadText(ctx, params)
	}
	return c.webCall(ctx, method, params, operation)
}

func (c *Client) webCall(ctx context.Context, method string, params map[string]any, operation string) (map[string]any, error) {
	values := url.Values{}
	for key, value := range params {
		if key == "token" {
			return nil, slackbot.Fail("invalid_argument", operation, "token is host-owned")
		}
		switch v := value.(type) {
		case string:
			values.Set(key, v)
		case bool, float64, int, int64:
			values.Set(key, fmt.Sprint(v))
		case nil:
			continue
		default:
			data, err := json.Marshal(value)
			if err != nil {
				return nil, slackbot.Fail("invalid_argument", operation, "parameters must be JSON")
			}
			values.Set(key, string(data))
		}
	}
	endpoint := "https://slack.com/api/" + method
	if c.origin != nil {
		endpoint = c.origin.String() + method
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, slackbot.Fail("invalid_argument", operation, "invalid request")
	}
	token := c.opts.BotToken
	if strings.HasPrefix(method, "admin.") || strings.HasSuffix(operation, "AsUser") {
		token = c.opts.UserToken
		if token == "" {
			return nil, slackbot.Fail("missing_user_token", operation, "separately authorized user token required")
		}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, safeError(ctx, err, operation)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == 429 {
		return nil, slackbot.Fail("rate_limited", operation, "Slack rate limit reached; retry later")
	}
	if resp.StatusCode != 200 {
		return nil, slackbot.Fail("service_error", operation, "Slack API HTTP failure")
	}
	var result map[string]any
	if err := json.NewDecoder(io.LimitReader(resp.Body, 8<<20)).Decode(&result); err != nil {
		return nil, slackbot.Fail("service_error", operation, "invalid Slack response")
	}
	if result["ok"] != true {
		code, _ := result["error"].(string)
		// Expose only stable known error codes; never arbitrary upstream text.
		switch code {
		case "permission_denied", "plan_upgrade_required", "feature_not_enabled", "team_access_not_granted", "missing_scope", "not_authed", "invalid_auth", "token_revoked", "not_in_channel", "channel_not_found", "user_not_found", "message_not_found", "cant_delete_message", "restricted_action", "not_allowed_token_type", "paid_teams_only", "invalid_arguments", "is_archived", "already_pinned", "no_pin":
		default:
			code = "service_error"
		}
		return nil, slackbot.Fail(code, operation, "Slack rejected the operation")
	}
	return result, nil
}

// uploadText implements Slack's external-upload sequence for generated UTF-8 files.
func (c *Client) uploadText(ctx context.Context, params map[string]any) (map[string]any, error) {
	name, _ := params["filename"].(string)
	content, _ := params["content"].(string)
	channel, _ := params["channel_id"].(string)
	if name == "" || channel == "" || len(content) == 0 || len(content) > 8<<20 {
		return nil, slackbot.Fail("invalid_argument", "files.upload", "filename, channel_id and 1 byte–8 MiB content required")
	}
	result, err := c.webCall(ctx, "files.getUploadURLExternal", map[string]any{"filename": name, "length": len(content)}, "files.upload")
	if err != nil {
		return nil, err
	}
	target, _ := result["upload_url"].(string)
	id, _ := result["file_id"].(string)
	u, err := url.Parse(target)
	if err != nil || id == "" || u.User != nil {
		return nil, slackbot.Fail("service_error", "files.upload", "invalid upload capability")
	}
	allowed := u.Scheme == "https" && u.Host == "files.slack.com"
	if c.origin != nil {
		allowed = u.Scheme == c.origin.Scheme && u.Host == c.origin.Host
	}
	if !allowed {
		return nil, slackbot.Fail("service_error", "files.upload", "unexpected upload destination")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	// No bot authorization header is sent to the returned upload capability.
	httpClient := *c.http
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error { return fmt.Errorf("upload redirects are disabled") }
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, safeError(ctx, err, "files.upload")
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	_ = resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, slackbot.Fail("service_error", "files.upload", "upload transfer failed")
	}
	completion := map[string]any{"channel_id": channel, "files": []map[string]string{{"id": id, "title": name}}}
	if thread, ok := params["thread_ts"].(string); ok && thread != "" {
		completion["thread_ts"] = thread
	}
	return c.webCall(ctx, "files.completeUploadExternal", completion, "files.upload")
}
