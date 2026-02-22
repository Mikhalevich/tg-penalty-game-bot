package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) ChangePlayersGameStatus(
	ctx context.Context,
	playerIDs []player.ID,
	status player.GameStatus,
	changedAt time.Time,
) error {
	var (
		query = `
			UPDATE player SET
				game_status = ?,
				game_status_changed_at = ?
			WHERE
				id IN(?)
		`
	)

	query, args, err := sqlx.In(query, status, changedAt, playerIDs)
	if err != nil {
		return fmt.Errorf("sqlx in statement: %w", err)
	}

	trx := p.transactor.ExtContext(ctx)

	res, err := trx.ExecContext(ctx, trx.Rebind(query), args...)
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
