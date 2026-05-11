package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (n *Notifier) GameLink(ctx context.Context, chatID msginfo.ChatID, gameID game.ID) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text:   makeStartGameLink(gameID),
		Type:   msginfo.MessageTypePlain,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeStartGameLink(gameID game.ID) string {
	return fmt.Sprintf("https://t.me/SoccerPenaltyBot?start=join_%s", gameID.String())
}
