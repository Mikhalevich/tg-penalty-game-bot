package playercontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (pc *PlayerController) SetPlayerReadyToGame(ctx context.Context, chatID msginfo.ChatID) error {
	currentPlayer, err := pc.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	switch currentPlayer.GameStatus {
	case player.GameStatusInGame:
		if err := pc.notifier.PlayerAlreadyInGame(ctx, currentPlayer); err != nil {
			return fmt.Errorf("send player already in game notification: %w", err)
		}

		return nil

	case player.GameStatusReadyForGame:
		return nil

	case player.GameStatusIdle:
	}

	if err := pc.ChangeGameStatus(
		ctx,
		currentPlayer.ID,
		"",
		player.GameStatusReadyForGame,
		pc.timeProvider.Now(),
	); err != nil {
		return fmt.Errorf("change player game status: %w", err)
	}

	return nil
}
