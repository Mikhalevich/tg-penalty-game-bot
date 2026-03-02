package playercontroller

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
	CreatePlayer(ctx context.Context, p player.Player) (player.ID, error)
	ChangeDisplayName(ctx context.Context, chatID msginfo.ChatID, displayName string, changedAt time.Time) error
	SetChangeDisplayNameTrigger(ctx context.Context, chatID msginfo.ChatID) error
	ChangePlayerGameStatus(
		ctx context.Context,
		playerID player.ID,
		gameID game.ID,
		status player.GameStatus,
		changedTime time.Time,
		prevoiusStatuses ...player.GameStatus,
	) error
	ChangePlayersGameStatus(
		ctx context.Context,
		playerIDs []player.ID,
		status player.GameStatus,
		changedAt time.Time,
	) error

	IsNotFoundError(err error) bool
	IsAlreadyExistsError(err error) bool
}

type Notifier interface {
	WelcomeNewPlayer(ctx context.Context, plr player.Player) error
	ChangeNameDelay(ctx context.Context, plr player.Player, waitPeriod time.Duration) error
	ChangeName(ctx context.Context, plr player.Player, fullName, userName string) error
	NameAlreadyRegistered(ctx context.Context, plr player.Player, msgID msginfo.MessageID) error
	NameChanged(ctx context.Context, plr player.Player, msgID msginfo.MessageID) error
	PlayerAlreadyInGame(ctx context.Context, plr player.Player) error
}

type NameGenerator interface {
	GenerateName() string
}

type TimeProvider interface {
	Now() time.Time
}

type PlayerController struct {
	repo              Repository
	notifier          Notifier
	nameGenerator     NameGenerator
	timeProvider      TimeProvider
	changeNameTimeout time.Duration
}

func New(
	repo Repository,
	notifier Notifier,
	nameGenerator NameGenerator,
	timeProvider TimeProvider,
	changeNameTimeout time.Duration,
) *PlayerController {
	return &PlayerController{
		repo:              repo,
		notifier:          notifier,
		nameGenerator:     nameGenerator,
		timeProvider:      timeProvider,
		changeNameTimeout: changeNameTimeout,
	}
}
