package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (p *Postgres) UpdateGame(ctx context.Context, game game.Game) error {
	var (
		query = `
			UPDATE game SET
				game_status = :game_status,
				payload = :payload,
				payload_version = payload_version + 1,
				payload_updated_at = :payload_updated_at
			WHERE
				game_id = :id AND
				payload_version = :payload_version
		`
	)

	dbGame, err := model.ToDBGame(game)
	if err != nil {
		return fmt.Errorf("convert to db game: %w", err)
	}

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, dbGame)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
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
