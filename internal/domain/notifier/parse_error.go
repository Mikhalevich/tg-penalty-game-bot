package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (n *Notifier) ParseError(
	ctx context.Context,
	chatID msginfo.ChatID,
	err error,
) error {
	pErr, ok := perror.ParseError(err)
	if !ok {
		return err
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text:   pErr.Message,
		Type:   msginfo.MessageTypePlain,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
