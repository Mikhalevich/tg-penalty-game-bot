package gamecontroller

import (
	"context"

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
) (game.Game, []game.Shot, error) {
	if len(players) != playersCountForNewGame {
		return game.Game{}, nil, perror.InvalidParam("invalid players count")
	}

	var (
		gameID       = game.IDFromString(uuid.NewString())
		now          = gc.timeProvider.Now()
		gamePlayers  = make([]game.Player, 0, len(players))
		pendingShots = make([]game.Shot, 0, playersCountForNewGame)
	)

	for plrIdx, plr := range players {
		gamePlayers = append(gamePlayers, game.Player{
			ID:             plr.ID,
			DisplayName:    plr.DisplayName,
			ShotsAvailable: game.ShotsInitial,
		})

		pendingShots = append(pendingShots, game.Shot{
			GameID:    gameID,
			PlayerID:  plr.ID,
			Round:     0,
			Type:      shotTypeByPlayerPos(plrIdx),
			CreatedAt: now,
			Side:      game.ShotSideNoShot,
		})
	}

	return game.Game{
		ID:              gameID,
		CreatedAt:       now,
		Status:          game.GameStatusInProgress,
		StatusChagnedAt: now,
		State: game.State{
			Players: gamePlayers,
		},
	}, pendingShots, nil
}

func shotTypeByPlayerPos(pos int) game.ShotType {
	if pos == 0 {
		return game.ShotTypeAttack
	}

	return game.ShotTypeDefend
}
