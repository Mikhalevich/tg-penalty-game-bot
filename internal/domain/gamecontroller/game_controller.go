package gamecontroller

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	GetGame(ctx context.Context, gameID game.ID) (game.Game, error)
	UpdateGame(ctx context.Context, game game.Game) error
	InsertGames(ctx context.Context, games []game.Game) error
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

type TimeProvider interface {
	Now() time.Time
}

type Notifier interface {
	GameStage(ctx context.Context, currentGame game.Game) error
}

type GameController struct {
	repo                Repository
	transactor          Transactor
	playerStatusChanger PlayerStatusChanger
	timeProvider        TimeProvider
	notifier            Notifier
}

func New(
	repo Repository,
	transactor Transactor,
	playerStatusChanger PlayerStatusChanger,
	timeProvier TimeProvider,
	notifier Notifier,
) *GameController {
	return &GameController{
		repo:                repo,
		transactor:          transactor,
		playerStatusChanger: playerStatusChanger,
		timeProvider:        timeProvier,
		notifier:            notifier,
	}
}
