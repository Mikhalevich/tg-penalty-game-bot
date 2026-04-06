package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) PlayerScorePositionsFrom(
	ctx context.Context,
	positionFrom int,
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
				position >= $1
			LIMIT
				$2
		`

		positions []model.Position
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&positions,
		query,
		positionFrom,
		limit,
	); err != nil {
		return nil, fmt.Errorf("select player positions: %w", err)
	}

	return model.ToDomPositions(positions), nil
}
