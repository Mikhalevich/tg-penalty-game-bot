package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (n *Notifier) LinkCanceled(ctx context.Context, chatID msginfo.ChatID) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text:   "Game canceled by this link",
		Type:   msginfo.MessageTypePlain,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
