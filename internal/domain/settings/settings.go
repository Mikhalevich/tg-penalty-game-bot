package settings

import (
	"context"
	"time"
)

type SettingItem struct {
	ID               string
	IsEnabled        bool
	CreatedAt        time.Time
	Payload          []byte
	PayloadUpdatedAt time.Time
}

type Repository interface {
	SetSetting(ctx context.Context, item SettingItem) error
	GetSetting(ctx context.Context, id string) (SettingItem, error)
}

type TimeProvider interface {
	Now() time.Time
}

type Settings struct {
	repo         Repository
	timeProvider TimeProvider
}

func New(
	repo Repository,
	timeProvider TimeProvider,
) *Settings {
	return &Settings{
		repo:         repo,
		timeProvider: timeProvider,
	}
}
