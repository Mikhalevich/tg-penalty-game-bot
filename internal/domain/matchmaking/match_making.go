package matchmaking

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	SelectReadyToGamePlayers(ctx context.Context, limit int) ([]player.Player, error)
	InsertGames(ctx context.Context, games []game.Game) error
	ChangeOrInsertPlayersGameStatus(
		ctx context.Context,
		players []player.Player,
	) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type TimeProvider interface {
	Now() time.Time
}

type Notifier interface {
	GameStage(ctx context.Context, gm game.Game) error
}

type MatchMaking struct {
	repo         Repository
	transactor   Transactor
	timeProvider TimeProvider
	notifier     Notifier
}

func New(
	repo Repository,
	transactor Transactor,
	timeProvider TimeProvider,
	notifier Notifier,
) *MatchMaking {
	return &MatchMaking{
		repo:         repo,
		transactor:   transactor,
		timeProvider: timeProvider,
		notifier:     notifier,
	}
}
