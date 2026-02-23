package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) NameAlreadyRegistered(
	ctx context.Context,
	plr player.Player,
	msgID msginfo.MessageID,
) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:     plr.ChatID,
		ReplyMsgID: msgID,
		Text:       "Name already registered, please try again",
		Type:       msginfo.MessageTypePlain,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
