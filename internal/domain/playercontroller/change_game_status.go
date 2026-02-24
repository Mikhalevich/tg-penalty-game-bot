package playercontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *PlayerController) ChangeGameStatus(
	ctx context.Context,
	playerID player.ID,
	gameID game.ID,
	status player.GameStatus,
) error {
	if err := p.repo.ChangePlayerGameStatus(
		ctx,
		playerID,
		gameID,
		status,
		p.timeProvider.Now(),
		calculateLegalPreviousStatutses(status)...,
	); err != nil {
		return fmt.Errorf("change game status: %w", err)
	}

	return nil
}

func calculateLegalPreviousStatutses(status player.GameStatus) []player.GameStatus {
	switch status {
	case player.GameStatusIdle:
		return []player.GameStatus{player.GameStatusReadyForGame, player.GameStatusInGame}

	case player.GameStatusReadyForGame:
		return []player.GameStatus{player.GameStatusIdle}

	case player.GameStatusInGame:
		return []player.GameStatus{player.GameStatusIdle, player.GameStatusReadyForGame}
	}

	return nil
}
