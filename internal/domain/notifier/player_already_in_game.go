package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) PlayerAlreadyInGame(ctx context.Context, plr player.Player) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: plr.ChatID,
		Text:   "Already in game",
		Type:   msginfo.MessageTypePlain,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
