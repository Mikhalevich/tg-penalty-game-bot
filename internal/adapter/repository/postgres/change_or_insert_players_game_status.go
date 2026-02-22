package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) ChangeOrInsertPlayersGameStatus(
	ctx context.Context,
	players []player.Player,
) error {
	var (
		query = `
			INSERT INTO player(
				chat_id,
				display_name,
				created_at,
				game_status,
				game_status_changed_at,
				current_game_id
			) VALUES (
				:chat_id,
				:display_name,
				:created_at,
				:game_status,
				:game_status_changed_at,
				:current_game_id
			) ON CONFLICT(chat_id)
				DO UPDATE SET
					game_status = EXCLUDED.game_status,
					game_status_changed_at = EXCLUDED.game_status_changed_at,
					current_game_id = EXCLUDED.current_game_id
		`
	)

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, model.ToDBPlayers(players))
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
