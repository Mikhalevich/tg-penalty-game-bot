package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (p *Postgres) InsertGames(ctx context.Context, games []game.Game) error {
	var (
		query = `
			INSERT INTO game(
				id,
				created_at,
				game_status,
				game_status_changed_at,
				payload,
				payload_version
			) VALUES (
				:id,
				:created_at,
				:game_status,
				:game_status_changed_at,
				:payload,
				:payload_version
			)
		`
	)

	dbGames, err := model.ToDBGames(games)
	if err != nil {
		return fmt.Errorf("convert to db games: %w", err)
	}

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, dbGames)
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
