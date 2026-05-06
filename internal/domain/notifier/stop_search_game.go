package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (n *Notifier) StopSearchGame(ctx context.Context, chatID msginfo.ChatID) error {
	playRatingBtn, err := game.StartGameWithDeleteAfterPressButton("Play rating game", game.GameTypeRating)
	if err != nil {
		return fmt.Errorf("create rating button: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text:   "Search stopped",
		Type:   msginfo.MessageTypePlain,
		Buttons: []button.ButtonRow{
			button.Row(playRatingBtn),
		},
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
