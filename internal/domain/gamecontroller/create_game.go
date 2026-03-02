package gamecontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

// CreateGame create a new game and return created game and two pending players shots.
func (gc *GameController) CreateGame(
	ctx context.Context,
	player1 player.Player,
	player2 player.Player,
) (game.Game, error) {
	var (
		gameID = game.IDFromString(uuid.NewString())
		now    = gc.timeProvider.Now()
	)

	createdGame := game.Game{
		ID:              gameID,
		CreatedAt:       now,
		Status:          game.GameStatusInProgress,
		StatusChagnedAt: now,
		State: game.State{
			Player1: createGamePlayerFromPlayer(player1),
			Player2: createGamePlayerFromPlayer(player2),
		},
	}

	if err := createdGame.StartNextRound(now); err != nil {
		return game.Game{}, fmt.Errorf("start next round: %w", err)
	}

	return createdGame, nil
}

func createGamePlayerFromPlayer(plr player.Player) game.Player {
	return game.Player{
		ID:             plr.ID,
		ChatID:         plr.ChatID,
		DisplayName:    plr.DisplayName,
		ShotsAvailable: game.ShotsInitial,
	}
}
