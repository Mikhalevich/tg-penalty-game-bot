package gameshot

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	GetGame(ctx context.Context, gameID game.ID) (game.Game, error)
	UpdateGame(ctx context.Context, game game.Game) error
	IsNoRowsUpdated(err error) bool
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerStatusChanger interface {
	ChangePlayersGameStatus(
		ctx context.Context,
		playerIDs []player.ID,
		status player.GameStatus,
		changedAt time.Time,
	) error
}

type Notifier interface {
	GameNewRound(ctx context.Context, currentGame game.Game) error
	GameRoundFinish(ctx context.Context, currentGame game.Game) error
	GameFinish(ctx context.Context, currentGame game.Game) error
}

type GameShot struct {
	repo                Repository
	transactor          Transactor
	playerStatusChanger PlayerStatusChanger
	notifier            Notifier
}

func New(
	repo Repository,
	transactor Transactor,
	playerStatusChanger PlayerStatusChanger,
	notifier Notifier,
) *GameShot {
	return &GameShot{
		repo:                repo,
		transactor:          transactor,
		playerStatusChanger: playerStatusChanger,
		notifier:            notifier,
	}
}
