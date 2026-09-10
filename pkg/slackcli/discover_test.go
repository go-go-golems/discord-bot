package slackcli

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDiscoverExcludesHelpersAndFindsDuplicates(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "ping", "lib"), 0700))
	source := []byte(`module.exports=require("slack").defineBot(({configure})=>configure({name:"ping"}));`)
	require.NoError(t, os.WriteFile(filepath.Join(root, "ping", "index.js"), source, 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "ping", "lib", "index.js"), []byte(`throw new Error("helper must not run")`), 0600))
	all, err := Discover(context.Background(), root, time.Second)
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.NoError(t, os.WriteFile(filepath.Join(root, "other.js"), source, 0600))
	_, err = Discover(context.Background(), root, time.Second)
	require.ErrorContains(t, err, "duplicate bot")
}
