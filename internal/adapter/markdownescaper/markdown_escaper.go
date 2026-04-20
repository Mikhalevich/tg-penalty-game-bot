package markdownescaper

import (
	"github.com/go-telegram/bot"
)

type MarkdownEscaper struct {
}

func New() *MarkdownEscaper {
	return &MarkdownEscaper{}
}

func (m *MarkdownEscaper) EscapeMarkdown(s string) string {
	return bot.EscapeMarkdown(s)
}
