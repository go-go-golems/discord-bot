package slackconfig

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const maxFileSize = 1 << 20

type Config struct {
	DefaultProfile string                  `yaml:"default_profile,omitempty"`
	Profiles       map[string]Profile      `yaml:"profiles,omitempty"`
	Management     map[string]Management   `yaml:"management,omitempty"`
	Apps           map[string]App          `yaml:"apps,omitempty"`
	Installations  map[string]Installation `yaml:"installations,omitempty"`
}
type Profile struct {
	Management   string `yaml:"management,omitempty"`
	App          string `yaml:"app,omitempty"`
	Installation string `yaml:"installation,omitempty"`
	Bot          string `yaml:"bot,omitempty"`
}
type Management struct {
	TeamID string `yaml:"team_id,omitempty"`
	UserID string `yaml:"user_id,omitempty"`
}
type App struct {
	AppID string `yaml:"app_id"`
}
type Installation struct {
	App    string `yaml:"app,omitempty"`
	TeamID string `yaml:"team_id,omitempty"`
}

type Credentials struct {
	Management    map[string]ManagementCredential   `json:"management,omitempty"`
	Apps          map[string]AppCredential          `json:"apps,omitempty"`
	Installations map[string]InstallationCredential `json:"installations,omitempty"`
}
type ManagementCredential struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	TeamID       string `json:"team_id,omitempty"`
	UserID       string `json:"user_id,omitempty"`
}
type AppCredential struct {
	AppID         string `json:"app_id,omitempty"`
	ClientID      string `json:"client_id,omitempty"`
	ClientSecret  string `json:"client_secret,omitempty"`
	SigningSecret string `json:"signing_secret,omitempty"`
	AppToken      string `json:"app_token,omitempty"`
}
type InstallationCredential struct {
	UserToken string `json:"user_token,omitempty"`
	BotToken  string `json:"bot_token,omitempty"`
	AppToken  string `json:"app_token,omitempty"`
}

type Store struct{ Root string }

func New(root string) Store             { return Store{Root: root} }
func (s Store) configPath() string      { return filepath.Join(s.Root, "config.yaml") }
func (s Store) credentialsPath() string { return filepath.Join(s.Root, "credentials.json") }

func (s Store) Load() (Config, Credentials, error) {
	c := Config{Profiles: map[string]Profile{}, Management: map[string]Management{}, Apps: map[string]App{}, Installations: map[string]Installation{}}
	cr := Credentials{Management: map[string]ManagementCredential{}, Apps: map[string]AppCredential{}, Installations: map[string]InstallationCredential{}}
	if b, err := readOptional(s.configPath()); err != nil {
		return c, cr, errors.Wrap(err, "read slack config")
	} else if len(b) > 0 {
		if err := yaml.Unmarshal(b, &c); err != nil {
			return c, cr, errors.Wrap(err, "decode slack config")
		}
	}
	if b, err := readOptional(s.credentialsPath()); err != nil {
		return c, cr, errors.Wrap(err, "read slack credentials")
	} else if len(b) > 0 {
		if err := json.Unmarshal(b, &cr); err != nil {
			return c, cr, errors.Wrap(err, "decode slack credentials")
		}
	}
	ensureMaps(&c, &cr)
	return c, cr, nil
}
func readOptional(path string) ([]byte, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, maxFileSize+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxFileSize {
		return nil, errors.Errorf("%s exceeds 1 MiB", path)
	}
	return b, nil
}
func ensureMaps(c *Config, cr *Credentials) {
	if c.Profiles == nil {
		c.Profiles = map[string]Profile{}
	}
	if c.Management == nil {
		c.Management = map[string]Management{}
	}
	if c.Apps == nil {
		c.Apps = map[string]App{}
	}
	if c.Installations == nil {
		c.Installations = map[string]Installation{}
	}
	if cr.Management == nil {
		cr.Management = map[string]ManagementCredential{}
	}
	if cr.Apps == nil {
		cr.Apps = map[string]AppCredential{}
	}
	if cr.Installations == nil {
		cr.Installations = map[string]InstallationCredential{}
	}
}
func (s Store) Save(c Config, cr Credentials) error {
	ensureMaps(&c, &cr)
	if err := os.MkdirAll(s.Root, 0700); err != nil {
		return errors.Wrap(err, "create slack config directory")
	}
	if err := writeAtomic(s.configPath(), func(w io.Writer) error { return yaml.NewEncoder(w).Encode(c) }); err != nil {
		return errors.Wrap(err, "write slack config")
	}
	if err := writeAtomic(s.credentialsPath(), func(w io.Writer) error { e := json.NewEncoder(w); e.SetIndent("", "  "); return e.Encode(cr) }); err != nil {
		return errors.Wrap(err, "write slack credentials")
	}
	return nil
}
func writeAtomic(path string, write func(io.Writer) error) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		return err
	}
	if err := write(tmp); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
func ResolveProfile(c Config, name string) (string, Profile, error) {
	if name == "" {
		name = c.DefaultProfile
	}
	if strings.TrimSpace(name) == "" {
		return "", Profile{}, errors.New("no profile selected; use --profile or configure default_profile")
	}
	p, ok := c.Profiles[name]
	if !ok {
		return "", Profile{}, errors.Errorf("unknown profile %q", name)
	}
	return name, p, nil
}
