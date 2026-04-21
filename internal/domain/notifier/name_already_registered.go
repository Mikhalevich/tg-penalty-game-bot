package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (n *Notifier) NameAlreadyRegistered(
	ctx context.Context,
	chatID msginfo.ChatID,
	displayName string,
) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text: fmt.Sprintf("Name *%s* already registered, please try again",
			n.escaper.EscapeMarkdown(displayName)),
		Type: msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
