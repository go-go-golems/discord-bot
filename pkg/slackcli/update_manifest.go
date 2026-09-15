package slackcli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/pkg/errors"
)

// updateManifest replaces the app configuration with the selected bot's manifest.
// Runtime tokens remain in the credential store; permission changes need reinstall.
func (c *command) updateManifest(ctx context.Context, d slackbot.Descriptor, s settings, appID, token string) (bool, error) {
	manifest, err := json.Marshal(Manifest(d))
	if err != nil {
		return false, err
	}
	body, err := json.Marshal(map[string]string{"app_id": appID, "manifest": string(manifest)})
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.TimeoutMS)*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://slack.com/api/apps.manifest.update", bytes.NewReader(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := c.appClient.Do(req)
	if err != nil {
		return false, errors.New("Slack manifest update request failed; inspect app settings before retrying")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, errors.Errorf("Slack manifest update returned HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024+1))
	if err != nil || len(data) > 1024*1024 {
		return false, errors.New("Slack manifest update returned an unreadable response")
	}
	var result struct {
		OK                 bool   `json:"ok"`
		AppID              string `json:"app_id"`
		PermissionsUpdated bool   `json:"permissions_updated"`
		Error              string `json:"error"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return false, errors.New("Slack manifest update returned invalid JSON")
	}
	if !result.OK {
		code := result.Error
		if !slackErrorCode.MatchString(code) || strings.Contains(code, token) {
			code = "unknown_error"
		}
		if code == "token_expired" || code == "invalid_auth" {
			return false, errors.Errorf("Slack manifest update failed: %s; run credentials refresh with the same profile and config directory, then retry", code)
		}
		return false, errors.Errorf("Slack manifest update failed: %s", code)
	}
	if result.AppID != appID {
		return false, errors.New("Slack manifest update returned an unexpected app ID")
	}
	return result.PermissionsUpdated, nil
}
