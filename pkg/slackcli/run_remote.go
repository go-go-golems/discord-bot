package slackcli

import (
	"context"
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
		Config:   config,
		Timeout:  time.Duration(s.TimeoutMS) * time.Millisecond,
		Logger:   c.logger,
	})
	if err != nil {
		return err
	}
	defer func() { _ = host.Close(context.Background()) }()
	return client.Run(ctx, host)
}
