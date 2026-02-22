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
		gameID game.GameID,
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

type NameGenerator interface {
	GenerateName() string
}

type TimeProvider interface {
	Now() time.Time
}

type PlayerController struct {
	repo              Repository
	nameGenerator     NameGenerator
	timeProvider      TimeProvider
	changeNameTimeout time.Duration
}

func New(
	repo Repository,
	nameGenerator NameGenerator,
	timeProvider TimeProvider,
	changeNameTimeout time.Duration,
) *PlayerController {
	return &PlayerController{
		repo:              repo,
		nameGenerator:     nameGenerator,
		timeProvider:      timeProvider,
		changeNameTimeout: changeNameTimeout,
	}
}
