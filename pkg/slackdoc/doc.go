// Package slackdoc embeds Slack bot development and UI guides.
package slackdoc

import (
	_ "embed"
	"github.com/go-go-golems/glazed/pkg/help"
)

//go:embed slack-bot-guide.md
var guide []byte

//go:embed slack-ui-dsl.md
var uiGuide []byte

func AddTo(h *help.HelpSystem) error {
	for _, data := range [][]byte{guide, uiGuide} {
		s, err := help.LoadSectionFromMarkdown(data)
		if err != nil {
			return err
		}
		h.AddSection(s)
	}
	return nil
}
