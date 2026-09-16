// Package slackcli provides offline Slack bot discovery, inspection and replay.
package slackcli

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/go-go-golems/discord-bot/pkg/slackhost"
	"github.com/pkg/errors"
)

// Discover uses explicit entry conventions: root *.js files and immediate child index.js files.
// It does not execute helper index.js files deeper in a bot's directory.
func Discover(ctx context.Context, root string, timeout time.Duration) ([]slackbot.Descriptor, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return nil, errors.Wrap(err, "read bot repository")
	}
	out := []slackbot.Descriptor{}
	seen := map[string]string{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") || entry.Name() == "node_modules" {
			continue
		}
		path := filepath.Join(abs, entry.Name())
		if entry.IsDir() {
			path = filepath.Join(path, "index.js")
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			} else if err != nil {
				return nil, err
			}
		} else if filepath.Ext(path) != ".js" {
			continue
		}
		d, err := slackhost.Inspect(ctx, path, timeout)
		if err != nil {
			return nil, errors.Wrap(err, "inspect "+path)
		}
		if previous, ok := seen[d.Name]; ok {
			return nil, errors.Errorf("duplicate bot %q: %s and %s", d.Name, previous, path)
		}
		seen[d.Name] = path
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
func Resolve(ctx context.Context, root, name string, timeout time.Duration) (slackbot.Descriptor, error) {
	all, err := Discover(ctx, root, timeout)
	if err != nil {
		return slackbot.Descriptor{}, err
	}
	for _, d := range all {
		if d.Name == name {
			return d, nil
		}
	}
	return slackbot.Descriptor{}, errors.Errorf("bot %q not found", name)
}
