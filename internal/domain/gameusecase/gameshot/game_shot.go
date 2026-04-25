package gameshot

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type Repository interface {
	GetGame(ctx context.Context, gameID game.ID) (game.Game, error)
	UpdateGame(ctx context.Context, game game.Game) error
	InsertShots(ctx context.Context, shots []game.Shot) error
	IsNoRowsUpdated(err error) bool
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerStatusChanger interface {
	ChangeStatusForFinishedGame(
		ctx context.Context,
		finishedGame game.Game,
		finishedAt time.Time,
	) error
}

type Notifier interface {
	GameNewRound(ctx context.Context, gameID game.ID, state game.State) error
	GameRoundFinish(ctx context.Context, state game.State) error
	GameFinish(ctx context.Context, gameType game.GameType, state game.State, finishedAt time.Time) error
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
