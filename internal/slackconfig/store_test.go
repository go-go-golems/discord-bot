package slackconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreRoundTripAndAtomicFiles(t *testing.T) {
	root := t.TempDir()
	s := New(filepath.Join(root, "config"))
	c := Config{DefaultProfile: "dev", Profiles: map[string]Profile{"dev": {Management: "owner"}}, Management: map[string]Management{"owner": {TeamID: "T1", UserID: "U1"}}}
	cr := Credentials{Management: map[string]ManagementCredential{"owner": {AccessToken: "access", RefreshToken: "refresh"}}}
	require.NoError(t, s.Save(c, cr))
	got, gotCreds, err := s.Load()
	require.NoError(t, err)
	require.Equal(t, "owner", got.Profiles["dev"].Management)
	require.Equal(t, "access", gotCreds.Management["owner"].AccessToken)
	for _, name := range []string{"config.yaml", "credentials.json"} {
		info, err := os.Stat(filepath.Join(s.Root, name))
		require.NoError(t, err)
		require.Equal(t, os.FileMode(0600), info.Mode().Perm())
	}
}

func TestResolveProfileRequiresExplicitOrDefault(t *testing.T) {
	_, _, err := ResolveProfile(Config{Profiles: map[string]Profile{}}, "")
	require.EqualError(t, err, "no profile selected; use --profile or configure default_profile")
}
