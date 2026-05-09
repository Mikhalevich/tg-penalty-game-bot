package model

import (
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/settings"
)

type SettingItem struct {
	ID               string      `db:"id"`
	IsEnabled        bool        `db:"is_enabled"`
	CreatedAt        time.Time   `db:"created_at"`
	Payload          jsonb.JSONB `db:"payload"`
	PayloadUpdatedAt time.Time   `db:"payload_updated_at"`
}

func (si SettingItem) ToDomSettingItem() settings.SettingItem {
	return settings.SettingItem{
		ID:               si.ID,
		IsEnabled:        si.IsEnabled,
		CreatedAt:        si.CreatedAt,
		Payload:          []byte(si.Payload),
		PayloadUpdatedAt: si.PayloadUpdatedAt,
	}
}

func ToDBSettingItem(item settings.SettingItem) SettingItem {
	return SettingItem{
		ID:               item.ID,
		IsEnabled:        item.IsEnabled,
		CreatedAt:        item.CreatedAt,
		Payload:          jsonb.NewString(string(item.Payload)),
		PayloadUpdatedAt: item.PayloadUpdatedAt,
	}
}
