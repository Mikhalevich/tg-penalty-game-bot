package gamecontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	playersCountForNewGame = 2
)

// CreateGame create a new game and return created game and two pending players shots.
func (gc *GameController) CreateGame(
	ctx context.Context,
	players []player.Player,
) (game.Game, error) {
	if len(players) != playersCountForNewGame {
		return game.Game{}, perror.InvalidParam("invalid players count")
	}

	var (
		gameID      = game.IDFromString(uuid.NewString())
		now         = gc.timeProvider.Now()
		gamePlayers = make([]game.Player, 0, len(players))
	)

	for _, plr := range players {
		gamePlayers = append(gamePlayers, game.Player{
			ID:             plr.ID,
			DisplayName:    plr.DisplayName,
			ShotsAvailable: game.ShotsInitial,
		})
	}

	createdGame := game.Game{
		ID:              gameID,
		CreatedAt:       now,
		Status:          game.GameStatusInProgress,
		StatusChagnedAt: now,
		State: game.State{
			Players: gamePlayers,
		},
	}

	if err := createdGame.StartNextRound(now); err != nil {
		return game.Game{}, fmt.Errorf("start next round: %w", err)
	}

	return createdGame, nil
}
