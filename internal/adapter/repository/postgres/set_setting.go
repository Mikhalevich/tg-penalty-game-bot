package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/settings"
)

func (p *Postgres) SetSetting(ctx context.Context, item settings.SettingItem) error {
	var (
		query = `
			INSERT INTO settings(
				id,
				is_enabled,
				created_at,
				payload,
				payload_updated_at
			) VALUES (
				:id,
				:is_enabled,
				:created_at,
				:payload,
				:payload_updated_at
			) ON CONFLICT(id) DO
				UPDATE SET
					is_enabled = EXCLUDED.is_enabled,
					payload = EXCLUDED.payload,
					payload_updated_at = EXCLUDED.payload_updated_at
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.ToDBSettingItem(item),
	)

	if err != nil {
		return fmt.Errorf("insert or update setting: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
