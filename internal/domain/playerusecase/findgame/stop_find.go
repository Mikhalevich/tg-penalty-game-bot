package findgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (fg *FindGame) StopFind(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	currentPlayer, err := fg.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	switch currentPlayer.GameStatus {
	case player.GameStatusIdle, player.GameStatusInGame:

	case player.GameStatusReadyForGame:
		if err := fg.stopSearchInReeadyToGameStatus(ctx, chatID, currentPlayer); err != nil {
			return fmt.Errorf("stop search in ready to game status: %w", err)
		}
	}

	if err := fg.messageDeleter.DeleteMessage(ctx, chatID, messageID); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}

func (fg *FindGame) stopSearchInReeadyToGameStatus(
	ctx context.Context,
	chatID msginfo.ChatID,
	currentPlayer player.Player,
) error {
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
