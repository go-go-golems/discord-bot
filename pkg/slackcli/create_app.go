package slackcli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-go-golems/discord-bot/internal/slackconfig"
	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/pkg/errors"
)

const createAppURL = "https://slack.com/api/apps.manifest.create"

var slackErrorCode = regexp.MustCompile(`^[a-z][a-z0-9_]{0,79}$`)

func appHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil // Use explicit token-file auth; do not inspect environment proxy settings.
	return &http.Client{Transport: transport, Timeout: 30 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return errors.New("app creation redirects are disabled")
		},
	}
}

type appCreationResponse struct {
	OK                bool            `json:"ok"`
	AppID             string          `json:"app_id"`
	OAuthAuthorizeURL string          `json:"oauth_authorize_url"`
	Credentials       json.RawMessage `json:"credentials"`
	Error             string          `json:"error"`
	Errors            []struct {
		Message string `json:"message"`
		Pointer string `json:"pointer"`
	} `json:"errors"`
}

func configAccessToken(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "open config-token-file")
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return "", errors.New("config-token-file must be a regular file")
	}
	b, err := io.ReadAll(io.LimitReader(f, 16385))
	if err != nil {
		return "", errors.New("could not read config-token-file")
	}
	token := strings.TrimSpace(string(b))
	if len(b) > 16384 || token == "" || strings.ContainsAny(token, " \t\r\n") {
		return "", errors.New("config-token-file must contain one access token (maximum 16 KiB)")
	}
	return token, nil
}

func (c *command) createApp(ctx context.Context, d slackbot.Descriptor, s settings, w io.Writer) error {
	var err error
	var token string
	managed := false
	var store slackconfig.Store
	var cfg slackconfig.Config
	var creds slackconfig.Credentials
	var profileName string
	if s.Profile != "" {
		if s.ConfigTokenFile != "" {
			return errors.New("--profile and --config-token-file cannot be combined")
		}
		if s.ConfigDir == "" {
			s.ConfigDir = defaultConfigDir()
		}
		store = slackconfig.New(s.ConfigDir)
		var err error
		cfg, creds, err = store.Load()
		if err != nil {
			return err
		}
		var profile slackconfig.Profile
		profileName, profile, err = slackconfig.ResolveProfile(cfg, s.Profile)
		if err != nil {
			return err
		}
		if profile.App != "" {
			return errors.Errorf("profile %q already references app %q; use a new profile to create another app", profileName, profile.App)
		}
		management, ok := creds.Management[profile.Management]
		if !ok || management.AccessToken == "" {
			return errors.Errorf("profile %q has no management access token; import one first", profileName)
		}
		token, managed = management.AccessToken, true
	} else {
		var err error
		token, err = configAccessToken(s.ConfigTokenFile)
		if err != nil {
			return err
		}
	}
	var credentials *os.File
	credentialsSaved := false
	if s.CredentialsFile != "" {
		// Reserve before calling Slack: an existing/unwritable destination must not create an app.
		credentials, err = os.OpenFile(s.CredentialsFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return errors.Wrap(err, "create credentials-file")
		}
		defer func() {
			_ = credentials.Close()
			if !credentialsSaved {
				_ = os.Remove(s.CredentialsFile)
			}
		}()
	}
	manifest, err := json.Marshal(Manifest(d))
	if err != nil {
		return err
	}
	// Slack expects an API argument named manifest whose value is the serialized manifest.
	body, err := json.Marshal(map[string]string{"manifest": string(manifest)})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.TimeoutMS)*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createAppURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	response, err := c.appClient.Do(req)
	if err != nil {
		return errors.New("app creation outcome unknown; check api.slack.com/apps before retrying (request failed or canceled)")
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode == http.StatusTooManyRequests {
		return errors.New("Slack rate limited app creation; no automatic retry was made")
	}
	if response.StatusCode != http.StatusOK {
		return errors.Errorf("app creation returned HTTP %d; check api.slack.com/apps before retrying", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 1024*1024+1))
	if err != nil || len(data) > 1024*1024 {
		return errors.New("app creation outcome unknown; unreadable response; check api.slack.com/apps before retrying")
	}
	var result appCreationResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return errors.New("app creation outcome unknown; invalid response; check api.slack.com/apps before retrying")
	}
	if !result.OK {
		code := result.Error
		if !slackErrorCode.MatchString(code) || strings.Contains(code, token) {
			code = "unknown_error"
		}
		message := "Slack app creation failed: " + code
		if code == "invalid_manifest" {
			for i, detail := range result.Errors {
				if i >= 8 {
					break
				}
				text := strings.ReplaceAll(detail.Pointer+": "+detail.Message, token, "[redacted]")
				if len(text) > 512 {
					text = text[:512]
				}
				message += "\n" + text
			}
		}
		return errors.New(message)
	}
	if result.AppID == "" {
		return errors.New("app creation outcome unknown; response missing app_id; check api.slack.com/apps before retrying")
	}
	if credentials != nil {
		// Retain the app ID alongside credentials, without copying the configuration access token.
		if err := json.NewEncoder(credentials).Encode(map[string]any{"app_id": result.AppID, "credentials": result.Credentials}); err != nil {
			return errors.Errorf("app %s created, but credentials-file write failed; recover credentials in app settings", result.AppID)
		}
		if err := credentials.Close(); err != nil {
			return errors.Errorf("app %s created, but credentials-file close failed; inspect the file and app settings", result.AppID)
		}
		credentialsSaved = true
	}
	if managed {
		var returned struct {
			ClientID      string `json:"client_id"`
			ClientSecret  string `json:"client_secret"`
			SigningSecret string `json:"signing_secret"`
			AppToken      string `json:"app_token"`
		}
		if len(result.Credentials) > 0 && string(result.Credentials) != "null" {
			_ = json.Unmarshal(result.Credentials, &returned)
		}
		cfg.Apps[profileName] = slackconfig.App{AppID: result.AppID}
		cfg.Profiles[profileName] = mergeProfile(cfg.Profiles[profileName], profileName, cfg.Profiles[profileName].Management)
		p := cfg.Profiles[profileName]
		p.App = profileName
		cfg.Profiles[profileName] = p
		creds.Apps[profileName] = slackconfig.AppCredential{AppID: result.AppID, ClientID: returned.ClientID, ClientSecret: returned.ClientSecret, SigningSecret: returned.SigningSecret, AppToken: returned.AppToken}
		if err := store.Save(cfg, creds); err != nil {
			return errors.Errorf("app %s created, but local profile save failed: %v", result.AppID, err)
		}
	}
	output := struct {
		AppID             string `json:"app_id"`
		SettingsURL       string `json:"settings_url"`
		OAuthAuthorizeURL string `json:"oauth_authorize_url"`
		CredentialsFile   string `json:"credentials_file,omitempty"`
	}{result.AppID, "https://api.slack.com/apps/" + result.AppID, result.OAuthAuthorizeURL, s.CredentialsFile}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		return errors.Errorf("app %s created, but output failed; do not recreate the app", result.AppID)
	}
	return nil
}
