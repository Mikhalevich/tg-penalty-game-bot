package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) PlayerScorePosition(ctx context.Context, playerID player.ID) (player.Position, error) {
	var (
		query = `
			SELECT
				player_id,
				display_name,
				score,
				updated_at,
				position
			FROM
				score_leaderboard
			WHERE
				player_id = $1
		`

		position model.Position
	)

	if err := sqlx.GetContext(ctx, p.transactor.ExtContext(ctx), &position, query, playerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return player.Position{}, ErrNotFound
		}

		return player.Position{}, fmt.Errorf("select player position: %w", err)
	}

	return position.ToDomPosition(), nil
}
