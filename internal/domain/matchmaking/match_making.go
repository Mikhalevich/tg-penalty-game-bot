package matchmaking

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	SelectReadyToGamePlayers(ctx context.Context, limit int) ([]player.Player, error)
	ChangeOrInsertPlayersGameStatus(
		ctx context.Context,
		players []player.Player,
	) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type GameRunner interface {
	StartGames(ctx context.Context, games []game.Game) error
}

type TimeProvider interface {
	Now() time.Time
}

type MatchMaking struct {
	repo         Repository
	transactor   Transactor
	gameRunner   GameRunner
	timeProvider TimeProvider
}

func New(
	repo Repository,
	transactor Transactor,
	gameRunner GameRunner,
	timeProvider TimeProvider,
) *MatchMaking {
	return &MatchMaking{
		repo:         repo,
		transactor:   transactor,
		gameRunner:   gameRunner,
		timeProvider: timeProvider,
	}
}
