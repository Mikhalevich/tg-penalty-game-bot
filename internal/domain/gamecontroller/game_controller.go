package gamecontroller

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
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

type PlayerController interface {
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
	ChangeGameStatus(
		ctx context.Context,
		playerID player.ID,
		gameID game.ID,
		status player.GameStatus,
		changedAt time.Time,
	) error
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
	PlayerAlreadyInGame(ctx context.Context, plr player.Player) error
}

type MessageDeleteter interface {
	DeleteMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
	) error
}

type GameController struct {
	repo             Repository
	transactor       Transactor
	playerController PlayerController
	timeProvider     TimeProvider
	notifier         Notifier
	messageDeleter   MessageDeleteter
}

func New(
	repo Repository,
	transactor Transactor,
	playerController PlayerController,
	timeProvier TimeProvider,
	notifier Notifier,
	messageDeleter MessageDeleteter,
) *GameController {
	return &GameController{
		repo:             repo,
		transactor:       transactor,
		playerController: playerController,
		timeProvider:     timeProvier,
		notifier:         notifier,
		messageDeleter:   messageDeleter,
	}
}
