package findgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (fg *FindGame) StopFind(ctx context.Context, chatID msginfo.ChatID) error {
	currentPlayer, err := fg.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	switch currentPlayer.GameStatus {
	case player.GameStatusIdle, player.GameStatusInGame:
		return nil

	case player.GameStatusReadyForGame:
	}

	if err := fg.repo.ChangePlayerGameStatus(
		ctx,
		currentPlayer.ID,
		"",
		player.GameStatusIdle,
		fg.timeProvider.Now(),
		player.GameStatusReadyForGame,
	); err != nil {
		return fmt.Errorf("change player game status: %w", err)
	}

	if err := fg.notifier.StopSearchGame(ctx, chatID); err != nil {
		return fmt.Errorf("stop search notification: %w", err)
	}

	return nil
}
