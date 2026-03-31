package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) SetPlayerIdleStatusWithScore(
	ctx context.Context,
	playerID player.ID,
	scoreDelta int,
	changedAt time.Time,
) error {
	var (
		query = `
			UPDATE player SET
				game_status = :game_status_idle,
				game_status_changed_at = :game_status_changed_at,
				current_game_id = :current_game_id,
				score = CASE
					WHEN (score + :score_delta) < 0 THEN 0
					ELSE score + :score_delta
				END
			WHERE
				id = :id AND
				game_status = :previous_game_status
		`
	)

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query,
		map[string]any{
			"game_status_idle":       player.GameStatusIdle,
			"game_status_changed_at": changedAt,
			"current_game_id":        "",
			"score_delta":            scoreDelta,
			"id":                     playerID,
			"previous_game_status":   player.GameStatusInGame,
		})
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
