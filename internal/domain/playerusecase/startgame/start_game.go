package startgame

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	ChangePlayerGameStatus(
		ctx context.Context,
		playerID player.ID,
		gameID game.ID,
		status player.GameStatus,
		changedTime time.Time,
		prevoiusStatuses ...player.GameStatus,
	) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
}

type GameRunner interface {
	StartGames(ctx context.Context, games []game.Game) error
}

type TimeProvider interface {
	Now() time.Time
}

type Notifier interface {
	PlayerAlreadyInGame(ctx context.Context, plr player.Player) error
}

type StartGame struct {
	repo           Repository
	transactor     Transactor
	playerProvider PlayerProvider
	gameRunner     GameRunner
	timeProvider   TimeProvider
	notifier       Notifier
}

func New(
	repo Repository,
	transactor Transactor,
	playerProvider PlayerProvider,
	gameRunner GameRunner,
	timeProvider TimeProvider,
	notifier Notifier,
) *StartGame {
	return &StartGame{
		repo:           repo,
		transactor:     transactor,
		playerProvider: playerProvider,
		gameRunner:     gameRunner,
		timeProvider:   timeProvider,
		notifier:       notifier,
	}
}
