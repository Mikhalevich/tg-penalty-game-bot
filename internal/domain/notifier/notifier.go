package notifier

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

type Sender interface {
	SendMessage(ctx context.Context, msg msginfo.Message) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type Notifier struct {
	sender  Sender
	escaper MarkdownEscaper
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
) *Notifier {
	return &Notifier{
		sender:  sender,
		escaper: escaper,
	}
}
