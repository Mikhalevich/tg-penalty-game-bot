package changename

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	ChangeDisplayName(ctx context.Context, chatID msginfo.ChatID, displayName string, changedAt time.Time) error
	SetChangeDisplayNameTrigger(ctx context.Context, chatID msginfo.ChatID, trigger bool) error

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
	ChangeName(ctx context.Context, plr player.Player, isGroup bool, fullName, userName string) error
	NameChanged(ctx context.Context, plr player.Player) error
	NameAlreadyRegistered(ctx context.Context, chatID msginfo.ChatID, displayName string) error
	NameIsTooLong(ctx context.Context, chatID msginfo.ChatID, maxNameLen int) error
}

type ChangeName struct {
	repo              Repository
	playerProvider    PlayerProvider
	timeProvider      TimeProvider
	notifier          Notifier
	changeNameTimeout time.Duration
	maxNameLen        int
}

func New(
	repo Repository,
	playerProvider PlayerProvider,
	timeProvider TimeProvider,
	notifier Notifier,
	changeNameTimeout time.Duration,
	maxNameLen int,
) *ChangeName {
	return &ChangeName{
		repo:              repo,
		playerProvider:    playerProvider,
		timeProvider:      timeProvider,
		notifier:          notifier,
		changeNameTimeout: changeNameTimeout,
		maxNameLen:        maxNameLen,
	}
}
