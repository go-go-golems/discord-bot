package slackcli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/discord-bot/internal/slackconfig"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/sources"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func defaultConfigDir() string {
	d, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", ".config", "go-go-slack")
	}
	return filepath.Join(d, "go-go-slack")
}

func NewCredentialsCommand(logger zerolog.Logger, client *http.Client) *cobra.Command {
	root := &cobra.Command{Use: "credentials", Short: "Manage local Slack credentials"}
	root.AddCommand(newImportManagementCommand(), newProfilesCommand(), newStatusCommand(), newRefreshCommand(client))
	root.AddCommand(newImportRuntimeCommand())
	return root
}

func newImportRuntimeCommand() *cobra.Command {
	desc := cmds.NewCommandDescription("import-runtime", cmds.WithShort("Import bot, Socket Mode and optional administrative user tokens"), cmds.WithFlags(
		fields.New("config-dir", fields.TypeString, fields.WithDefault(defaultConfigDir()), fields.WithHelp("Local credentials directory")),
		fields.New("profile", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Profile name")),
		fields.New("installation", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Installation name")),
		fields.New("team-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Slack workspace ID")),
		fields.New("bot-token-file", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Bot token file")),
		fields.New("app-token-file", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Socket Mode app token file")),
		fields.New("user-token-file", fields.TypeString, fields.WithHelp("Optional separately authorized user token file")),
	))
	c := cli.NewCobraCommandFromCommandDescription(desc)
	parser, setupErr := cli.NewCobraParserFromSections(desc.Schema, &cli.CobraParserConfig{SkipCommandSettingsSection: true, MiddlewaresFunc: func(_ *values.Values, cmd *cobra.Command, args []string) ([]sources.Middleware, error) {
		return []sources.Middleware{sources.FromCobra(cmd), sources.FromArgs(args), sources.FromDefaults()}, nil
	}})
	if setupErr == nil {
		setupErr = parser.AddToCobraCommand(c)
	}
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if setupErr != nil {
			return setupErr
		}
		vals, err := parser.Parse(cmd, args)
		if err != nil {
			return err
		}
		var settings struct {
			Dir          string `glazed:"config-dir"`
			Profile      string `glazed:"profile"`
			Installation string `glazed:"installation"`
			TeamID       string `glazed:"team-id"`
			BotFile      string `glazed:"bot-token-file"`
			AppFile      string `glazed:"app-token-file"`
			UserFile     string `glazed:"user-token-file"`
		}
		if err := vals.DecodeSectionInto(schema.DefaultSlug, &settings); err != nil {
			return err
		}
		dir, profile, installation, teamID, botFile, appFile, userFile := settings.Dir, settings.Profile, settings.Installation, settings.TeamID, settings.BotFile, settings.AppFile, settings.UserFile
		if profile == "" || installation == "" || teamID == "" {
			return errors.New("--profile, --installation, and --team-id are required")
		}
		bot, err := readSecretFile(botFile)
		if err != nil {
			return err
		}
		app, err := readSecretFile(appFile)
		if err != nil {
			return err
		}
		store := slackconfig.New(dir)
		cfg, cr, err := store.Load()
		if err != nil {
			return err
		}
		name, p, err := slackconfig.ResolveProfile(cfg, profile)
		if err != nil {
			return err
		}
		if p.App == "" {
			return errors.Errorf("profile %q has no app; create or link one first", name)
		}
		if existing, ok := cfg.Installations[installation]; ok && (existing.App != p.App || (existing.TeamID != "" && existing.TeamID != teamID)) {
			return errors.Errorf("installation %q belongs to another app or workspace", installation)
		}
		cfg.Installations[installation] = slackconfig.Installation{App: p.App, TeamID: teamID}
		p.Installation = installation
		cfg.Profiles[name] = p
		ac := cr.Apps[p.App]
		ac.AppToken = app
		cr.Apps[p.App] = ac
		runtimeCredentials := cr.Installations[installation]
		runtimeCredentials.BotToken, runtimeCredentials.AppToken = bot, app
		if userFile != "" {
			runtimeCredentials.UserToken, err = readSecretFile(userFile)
			if err != nil {
				return err
			}
		}
		cr.Installations[installation] = runtimeCredentials
		if err := store.Save(cfg, cr); err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{"profile": name, "installation": installation, "team_id": teamID})
	}
	return c
}

// NewProfilesCommand exposes profile listing at the root, alongside the credentials group.
func NewProfilesCommand() *cobra.Command { return newProfilesCommand() }

type credentialSettings struct {
	Dir         string `glazed:"config-dir"`
	Profile     string `glazed:"profile"`
	Management  string `glazed:"management"`
	AccessFile  string `glazed:"access-token-file"`
	RefreshFile string `glazed:"refresh-token-file"`
}

// The pinned Glazed parser supplies schema flags while RunE preserves errors for the caller.
func credentialCommand(name, short string, flags []*fields.Definition, run func(*cobra.Command, credentialSettings) error) *cobra.Command {
	flags = append(flags, fields.New("config-dir", fields.TypeString, fields.WithDefault(defaultConfigDir()), fields.WithHelp("Local Slack configuration directory")))
	desc := cmds.NewCommandDescription(name, cmds.WithShort(short), cmds.WithFlags(flags...))
	c := cli.NewCobraCommandFromCommandDescription(desc)
	parser, setupErr := cli.NewCobraParserFromSections(desc.Schema, &cli.CobraParserConfig{
		SkipCommandSettingsSection: true,
		MiddlewaresFunc: func(_ *values.Values, cmd *cobra.Command, args []string) ([]sources.Middleware, error) {
			return []sources.Middleware{sources.FromCobra(cmd), sources.FromArgs(args), sources.FromDefaults()}, nil
		},
	})
	if setupErr == nil {
		setupErr = parser.AddToCobraCommand(c)
	}
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if setupErr != nil {
			return setupErr
		}
		vals, err := parser.Parse(cmd, args)
		if err != nil {
			return err
		}
		var settings credentialSettings
		if err := vals.DecodeSectionInto(schema.DefaultSlug, &settings); err != nil {
			return err
		}
		return run(cmd, settings)
	}
	return c
}

func newImportManagementCommand() *cobra.Command {
	return credentialCommand("import-management", "Import a configuration token pair", []*fields.Definition{
		fields.New("profile", fields.TypeString, fields.WithShortFlag("p"), fields.WithRequired(true), fields.WithHelp("Profile name")),
		fields.New("management", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Management identity name")),
		fields.New("access-token-file", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Configuration access-token file")),
		fields.New("refresh-token-file", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Configuration refresh-token file")),
	}, func(cmd *cobra.Command, settings credentialSettings) error {
		dir, profile, management, access, refresh := settings.Dir, settings.Profile, settings.Management, settings.AccessFile, settings.RefreshFile
		if profile == "" || management == "" {
			return errors.New("--profile and --management are required")
		}
		a, e := readSecretFile(access)
		if e != nil {
			return e
		}
		r, e := readSecretFile(refresh)
		if e != nil {
			return e
		}
		s := slackconfig.New(dir)
		cfg, cr, e := s.Load()
		if e != nil {
			return e
		}
		cfg.Profiles[profile] = mergeProfile(cfg.Profiles[profile], profile, management)
		cr.Management[management] = slackconfig.ManagementCredential{AccessToken: a, RefreshToken: r}
		if cfg.DefaultProfile == "" {
			cfg.DefaultProfile = profile
		}
		return s.Save(cfg, cr)
	})
}

func mergeProfile(p slackconfig.Profile, _ string, management string) slackconfig.Profile {
	p.Management = management
	return p
}
func newProfilesCommand() *cobra.Command {
	return credentialCommand("profiles", "List configured Slack profiles", nil, func(cmd *cobra.Command, settings credentialSettings) error {
		dir := settings.Dir
		cfg, _, e := slackconfig.New(dir).Load()
		if e != nil {
			return e
		}
		names := make([]string, 0, len(cfg.Profiles))
		for n := range cfg.Profiles {
			names = append(names, n)
		}
		sort.Strings(names)
		return json.NewEncoder(cmd.OutOrStdout()).Encode(names)
	})
}
func newStatusCommand() *cobra.Command {
	return credentialCommand("status", "Show safe credential status", []*fields.Definition{
		fields.New("profile", fields.TypeString, fields.WithHelp("Profile name")),
	}, func(cmd *cobra.Command, settings credentialSettings) error {
		dir, profile := settings.Dir, settings.Profile
		cfg, cr, e := slackconfig.New(dir).Load()
		if e != nil {
			return e
		}
		name, p, e := slackconfig.ResolveProfile(cfg, profile)
		if e != nil {
			return e
		}
		m := cr.Management[p.Management]
		out := map[string]any{"profile": name, "management": p.Management, "app": p.App, "installation": p.Installation, "has_user_token": cr.Installations[p.Installation].UserToken != "", "has_access_token": m.AccessToken != "", "has_refresh_token": m.RefreshToken != "", "expires_at": m.ExpiresAt}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
	})
}
func readSecretFile(path string) (string, error) {
	if path == "" {
		return "", errors.New("token file is required")
	}
	b, e := os.ReadFile(path)
	if e != nil {
		return "", errors.Wrap(e, "read token file")
	}
	v := strings.TrimSpace(string(b))
	if v == "" || strings.ContainsAny(v, " \t\r\n") {
		return "", errors.New("token file must contain one token")
	}
	return v, nil
}

func newRefreshCommand(client *http.Client) *cobra.Command {
	return credentialCommand("refresh", "Refresh a management token pair", []*fields.Definition{
		fields.New("profile", fields.TypeString, fields.WithHelp("Profile name")),
	}, func(cmd *cobra.Command, settings credentialSettings) error {
		dir, profile := settings.Dir, settings.Profile
		store := slackconfig.New(dir)
		cfg, cr, err := store.Load()
		if err != nil {
			return err
		}
		name, p, err := slackconfig.ResolveProfile(cfg, profile)
		if err != nil {
			return err
		}
		old, ok := cr.Management[p.Management]
		if !ok || old.RefreshToken == "" {
			return errors.Errorf("profile %q has no refresh token; import a new pair", name)
		}
		if client == nil {
			client = appHTTPClient()
		}
		form := url.Values{"refresh_token": []string{old.RefreshToken}}
		req, err := http.NewRequestWithContext(cmd.Context(), http.MethodPost, "https://slack.com/api/tooling.tokens.rotate", strings.NewReader(form.Encode()))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			return errors.New("credential refresh failed; import a new pair")
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.Errorf("credential refresh returned HTTP %d; import a new pair", resp.StatusCode)
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024+1))
		if err != nil || len(body) > 64*1024 {
			return errors.New("credential refresh returned an unreadable response; import a new pair")
		}
		var result struct {
			OK           bool   `json:"ok"`
			AccessToken  string `json:"token"`
			RefreshToken string `json:"refresh_token"`
			TeamID       string `json:"team_id"`
			UserID       string `json:"user_id"`
			IAT          int64  `json:"iat"`
			Exp          int64  `json:"exp"`
			Error        string `json:"error"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return errors.New("credential refresh returned invalid JSON; import a new pair")
		}
		if !result.OK || result.AccessToken == "" || result.RefreshToken == "" {
			code := result.Error
			if code == "" {
				code = "unknown_error"
			}
			return errors.Errorf("credential refresh failed: %s; import a new pair", code)
		}
		if old.TeamID != "" && result.TeamID != "" && old.TeamID != result.TeamID {
			return errors.New("credential refresh workspace mismatch; import a new pair")
		}
		if old.UserID != "" && result.UserID != "" && old.UserID != result.UserID {
			return errors.New("credential refresh user mismatch; import a new pair")
		}
		updated := old
		updated.AccessToken = result.AccessToken
		updated.RefreshToken = result.RefreshToken
		updated.TeamID = result.TeamID
		updated.UserID = result.UserID
		if result.Exp > 0 {
			updated.ExpiresAt = time.Unix(result.Exp, 0).UTC().Format(time.RFC3339)
		} else if result.IAT > 0 {
			updated.ExpiresAt = fmt.Sprintf("issued:%d", result.IAT)
		}
		cr.Management[p.Management] = updated
		if err := store.Save(cfg, cr); err != nil {
			return errors.Wrap(err, "save refreshed credentials")
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{"profile": name, "management": p.Management, "expires_at": updated.ExpiresAt})
	})
}
