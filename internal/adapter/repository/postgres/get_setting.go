package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/settings"
)

func (p *Postgres) GetSetting(ctx context.Context, settingID string) (settings.SettingItem, error) {
	var (
		query = `
			SELECT
				id,
				is_enabled,
				created_at,
				payload,
				payload_updated_at
			FROM
				settings
			WHERE
				id = $1
		`

		item model.SettingItem
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&item,
		query,
		settingID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return settings.SettingItem{}, perror.NotFound("setting not found")
		}

		return settings.SettingItem{}, fmt.Errorf("insert or update setting: %w", err)
	}

	return item.ToDomSettingItem(), nil
}
