// Package slackdoc embeds the offline Slack authoring guide.
package slackdoc

import (
	_ "embed"
	"github.com/go-go-golems/glazed/pkg/help"
)

//go:embed slack-offline.md
var guide []byte

func AddTo(h *help.HelpSystem) error {
	s, err := help.LoadSectionFromMarkdown(guide)
	if err != nil {
		return err
	}
	h.AddSection(s)
	return nil
}
