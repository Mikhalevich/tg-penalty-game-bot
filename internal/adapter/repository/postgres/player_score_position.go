package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	halfLimit = 2
)

func (p *Postgres) PlayerScorePosition(
	ctx context.Context,
	playerID player.ID,
	limit int,
) ([]player.Position, error) {
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
				position >= (
					SELECT
						position - $1
					FROM
						score_leaderboard
					WHERE
						player_id = $2
				)
			LIMIT $3
		`

		positions []model.Position
	)

	if err := sqlx.SelectContext(ctx, p.transactor.ExtContext(ctx), &positions, query,
		limit/halfLimit,
		playerID,
		limit,
	); err != nil {
		return nil, fmt.Errorf("select player position: %w", err)
	}

	return model.ToDomPositions(positions), nil
}
