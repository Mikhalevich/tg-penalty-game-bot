package playerprovider

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
	CreatePlayer(ctx context.Context, p player.Player) (player.ID, error)

	IsNotFoundError(err error) bool
	IsAlreadyExistsError(err error) bool
}

type NameGenerator interface {
	GenerateName() string
}

type TimeProvider interface {
	Now() time.Time
}

type PlayerProvider struct {
	repo          Repository
	nameGenerator NameGenerator
	timeProvider  TimeProvider
}

func New(
	repo Repository,
	nameGenerator NameGenerator,
	timeProvider TimeProvider,
) *PlayerProvider {
	return &PlayerProvider{
		repo: repo,
	}
}
