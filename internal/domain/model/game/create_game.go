package game

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

// CreateGame create a new game(in progress status), start first round and return created game.
func CreateGame(
	ctx context.Context,
	player1 player.Player,
	player2 player.Player,
	createdAt time.Time,
) (Game, error) {
	createdGame := Game{
		ID:              IDFromString(uuid.NewString()),
		CreatedAt:       createdAt,
		Status:          GameStatusInProgress,
		StatusChagnedAt: createdAt,
		State: State{
			Player1: createGamePlayerFromPlayer(player1),
			Player2: createGamePlayerFromPlayer(player2),
		},
	}

	if err := createdGame.StartNextRound(createdAt); err != nil {
		return Game{}, fmt.Errorf("start next round: %w", err)
	}

	return createdGame, nil
}

// CreatePendingGame create a new game in pending status, no round is starting and return created game.
func CreatePendingGame(
	ctx context.Context,
	plr player.Player,
	createdAt time.Time,
) Game {
	return Game{
		ID:              IDFromString(uuid.NewString()),
		CreatedAt:       createdAt,
		Status:          GameStatusPending,
		StatusChagnedAt: createdAt,
		State: State{
			Player1: createGamePlayerFromPlayer(plr),
		},
	}
}

func createGamePlayerFromPlayer(plr player.Player) Player {
	return Player{
		ID:             plr.ID,
		ChatID:         plr.ChatID,
		DisplayName:    plr.DisplayName,
		ShotsAvailable: ShotsInitial,
	}
}
