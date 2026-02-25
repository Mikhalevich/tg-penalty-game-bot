package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (p *Postgres) InsertShots(ctx context.Context, shots []game.Shot) error {
	var (
		query = `
			INSERT INTO game(
				game_id,
				player_id,
				round,
				shot_type,
				created_at,
				side,
				completed_at
			) VALUES (
				:game_id,
				:player_id,
				:round,
				:shot_type,
				:created_at,
				:side,
				:completed_at
			)
		`
	)

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, model.ToDBShots(shots))
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNoRowsUpdated
	}

	return nil
}
