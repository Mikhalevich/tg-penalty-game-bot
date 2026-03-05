package changename

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	ChangeDisplayName(ctx context.Context, chatID msginfo.ChatID, displayName string, changedAt time.Time) error
	SetChangeDisplayNameTrigger(ctx context.Context, chatID msginfo.ChatID) error

	IsAlreadyExistsError(err error) bool
}

type PlayerProvider interface {
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
}

type TimeProvider interface {
	Now() time.Time
}

type Notifier interface {
	ChangeNameDelay(ctx context.Context, plr player.Player, waitPeriod time.Duration) error
	ChangeName(ctx context.Context, plr player.Player, fullName, userName string) error
	NameChanged(ctx context.Context, plr player.Player, msgID msginfo.MessageID) error
	NameAlreadyRegistered(ctx context.Context, plr player.Player, msgID msginfo.MessageID) error
}

type ChangeName struct {
	repo              Repository
	playerProvider    PlayerProvider
	timeProvider      TimeProvider
	notifier          Notifier
	changeNameTimeout time.Duration
}

func New(
	repo Repository,
	playerProvider PlayerProvider,
	timeProvider TimeProvider,
	notifier Notifier,
	changeNameTimeout time.Duration,
) *ChangeName {
	return &ChangeName{
		repo:              repo,
		playerProvider:    playerProvider,
		timeProvider:      timeProvider,
		notifier:          notifier,
		changeNameTimeout: changeNameTimeout,
	}
}
