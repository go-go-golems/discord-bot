package slackcli

import (
	"context"
	"os"
	"time"

	"github.com/go-go-golems/discord-bot/internal/slackconfig"
	"github.com/go-go-golems/discord-bot/internal/slacktransport"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/go-go-golems/discord-bot/pkg/slackhost"
	"github.com/pkg/errors"
)

func (c *command) runRemote(ctx context.Context, d slackbot.Descriptor, s settings) error {
	configDir := s.ConfigDir
	if configDir == "" {
		configDir = defaultConfigDir()
	}
	cfg, creds, err := slackconfig.New(configDir).Load()
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
	if profile.Installation == "" {
		return errors.Errorf("profile %q has no installation; install the app first", profileName)
	}
	app, ok := cfg.Apps[profile.App]
	if !ok || app.AppID == "" {
		return errors.Errorf("profile %q references missing app %q", profileName, profile.App)
	}
	installation, ok := cfg.Installations[profile.Installation]
	if !ok || installation.App != profile.App || installation.TeamID == "" {
		return errors.Errorf("profile %q references an invalid installation %q", profileName, profile.Installation)
	}
	if s.TeamID != "" && s.TeamID != installation.TeamID {
		return errors.Errorf("profile %q installation belongs to team %q, not %q", profileName, installation.TeamID, s.TeamID)
	}
	installationCred, ok := creds.Installations[profile.Installation]
	if !ok || installationCred.BotToken == "" {
		return errors.Errorf("installation %q has no bot token; import or install runtime credentials first", profile.Installation)
	}
	appCred := creds.Apps[profile.App]
	if appCred.AppID != "" && appCred.AppID != app.AppID {
		return errors.Errorf("app %q has conflicting app IDs in local config", profile.App)
	}
	appToken := appCred.AppToken
	if appToken == "" {
		appToken = installationCred.AppToken
	}
	if appToken == "" {
		return errors.Errorf("app %q has no Socket Mode app token; install or import runtime credentials first", profile.App)
	}

	c.logger.Info().Int("pid", os.Getpid()).Str("profile", profileName).Str("bot", d.Name).Str("script", d.ScriptPath).Str("app_id", app.AppID).Str("team_id", installation.TeamID).Msg("Starting Slack bot with stored installation")
	client, err := slacktransport.NewRemote(slacktransport.RemoteOptions{
		BotToken: installationCred.BotToken,
		AppToken: appToken,
		TeamID:   installation.TeamID,
		AppID:    app.AppID,
		Logger:   c.logger,
	})
	if err != nil {
		return err
	}
	defer client.Close()
	config := map[string]any{}
	if s.ConfigFile != "" {
		if err := readJSON(s.ConfigFile, &config); err != nil {
			return err
		}
	}
	host, err := slackhost.Load(ctx, d.ScriptPath, slackhost.Options{
		Messages: client,
		Views:    client,
		Config:   config,
		Timeout:  time.Duration(s.TimeoutMS) * time.Millisecond,
		Logger:   c.logger,
	})
	if err != nil {
		return err
	}
	defer func() { _ = host.Close(context.Background()) }()

	if !s.SkipManifestUpdate {
		management := creds.Management[profile.Management]
		if management.AccessToken == "" {
			return errors.New("manifest update requires a management access token; import one or use --skip-manifest-update")
		}
		c.logger.Info().Str("app_id", app.AppID).Str("bot", d.Name).Msg("Updating Slack app manifest")
		permissionsUpdated, err := c.updateManifest(ctx, d, s, app.AppID, management.AccessToken)
		if err != nil {
			return err
		}
		c.logger.Info().Str("app_id", app.AppID).Bool("permissions_updated", permissionsUpdated).Msg("Slack app manifest updated")
		if permissionsUpdated {
			return errors.Errorf("manifest updated; reinstall with bots install %s --profile %s --team-id %s before running again (include the same --config-dir and --bot-repository if customized)", s.Name, profileName, installation.TeamID)
		}
	} else {
		c.logger.Info().Msg("Skipping Slack app manifest update")
	}
	return client.Run(ctx, host)
}
