package slackcli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-go-golems/discord-bot/internal/slackconfig"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/pkg/errors"
)

// developerInstallURL is the method used by the open-source Slack CLI for a
// local developer installation. Slack does not publish it as a general Web API
// method, so keep this call isolated and provide a manual fallback in help.
const developerInstallURL = "https://slack.com/api/apps.developerInstall"

type developerInstallResponse struct {
	OK              bool   `json:"ok"`
	AppID           string `json:"app_id"`
	Error           string `json:"error"`
	APIAccessTokens struct {
		Bot      string `json:"bot"`
		AppLevel string `json:"app_level"`
		User     string `json:"user"`
	} `json:"api_access_tokens"`
}

func (c *command) installApp(ctx context.Context, d slackbot.Descriptor, s settings, w io.Writer) error {
	if strings.TrimSpace(s.TeamID) == "" {
		return errors.New("--team-id is required for app installation")
	}
	if s.ConfigDir == "" {
		s.ConfigDir = defaultConfigDir()
	}
	store := slackconfig.New(s.ConfigDir)
	cfg, creds, err := store.Load()
	if err != nil {
		return err
	}
	profileName, profile, err := slackconfig.ResolveProfile(cfg, s.Profile)
	if err != nil {
		return err
	}
	if profile.App == "" {
		return errors.Errorf("profile %q has no app; create an app first", profileName)
	}
	app, ok := cfg.Apps[profile.App]
	if !ok || app.AppID == "" {
		return errors.Errorf("profile %q references missing app %q; inspect or create the app first", profileName, profile.App)
	}
	management, ok := creds.Management[profile.Management]
	if !ok || management.AccessToken == "" {
		return errors.Errorf("profile %q has no management access token; import one first", profileName)
	}
	installationName := profileName + "-" + s.TeamID
	if existing, ok := cfg.Installations[installationName]; ok && (existing.App != profile.App || (existing.TeamID != "" && existing.TeamID != s.TeamID)) {
		return errors.Errorf("installation %q belongs to another app or workspace", installationName)
	}
	appCredential := creds.Apps[profile.App]
	if appCredential.AppID != "" && appCredential.AppID != app.AppID {
		return errors.Errorf("app %q has conflicting app IDs in local config", profile.App)
	}

	scopes, err := manifestBotScopes(Manifest(d))
	if err != nil {
		return err
	}
	body, err := json.Marshal(struct {
		AppID           string   `json:"app_id"`
		BotScopes       []string `json:"bot_scopes"`
		OutgoingDomains []string `json:"outgoing_domains"`
		TeamID          string   `json:"team_id"`
	}{app.AppID, scopes, []string{}, s.TeamID})
	if err != nil {
		return errors.Wrap(err, "encode developer installation request")
	}
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(s.TimeoutMS)*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, developerInstallURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+management.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	response, err := c.appClient.Do(req)
	if err != nil {
		return errors.New("app installation outcome unknown; inspect Slack before retrying")
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 1024*1024+1))
	if err != nil || len(data) > 1024*1024 {
		return errors.New("app installation outcome unknown; unreadable response; inspect Slack before retrying")
	}
	if response.StatusCode != http.StatusOK {
		return errors.Errorf("app installation returned HTTP %d; inspect Slack before retrying", response.StatusCode)
	}
	var result developerInstallResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return errors.New("app installation outcome unknown; invalid response; inspect Slack before retrying")
	}
	if !result.OK {
		code := result.Error
		if !slackErrorCode.MatchString(code) {
			code = "unknown_error"
		}
		return errors.Errorf("Slack app installation failed: %s", code)
	}
	if result.AppID != "" && result.AppID != app.AppID {
		return errors.Errorf("Slack app installation returned app %q, expected %q", result.AppID, app.AppID)
	}
	if result.AppID == "" {
		result.AppID = app.AppID
	}
	if result.APIAccessTokens.Bot == "" || result.APIAccessTokens.AppLevel == "" {
		return errors.New("Slack app installation returned no bot and app-level tokens")
	}

	cfg.Installations[installationName] = slackconfig.Installation{App: profile.App, TeamID: s.TeamID}
	profile.Installation = installationName
	cfg.Profiles[profileName] = profile
	appCredential.AppID = app.AppID
	appCredential.AppToken = result.APIAccessTokens.AppLevel
	creds.Apps[profile.App] = appCredential
	creds.Installations[installationName] = slackconfig.InstallationCredential{BotToken: result.APIAccessTokens.Bot, AppToken: result.APIAccessTokens.AppLevel}
	if err := store.Save(cfg, creds); err != nil {
		return errors.Wrap(err, "save installed app credentials")
	}
	return json.NewEncoder(w).Encode(map[string]any{
		"profile":      profileName,
		"installation": installationName,
		"app_id":       result.AppID,
		"team_id":      s.TeamID,
	})
}

func manifestBotScopes(manifest map[string]any) ([]string, error) {
	oauth, ok := manifest["oauth_config"].(map[string]any)
	if !ok {
		return nil, errors.New("manifest has no oauth_config.scopes.bot")
	}
	scopesContainer, ok := oauth["scopes"].(map[string]any)
	if !ok {
		return nil, errors.New("manifest has no oauth_config.scopes.bot")
	}
	values, ok := scopesContainer["bot"].([]string)
	if !ok || len(values) == 0 {
		return nil, errors.New("manifest has no bot scopes")
	}
	return values, nil
}
