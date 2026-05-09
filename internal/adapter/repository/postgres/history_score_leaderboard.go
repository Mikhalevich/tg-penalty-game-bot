package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) HistoryScoreLeaderboard(
	ctx context.Context,
	year int,
	month int,
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
				leaderboard_score_history
			WHERE
				year = $1 AND
				month = $2
			ORDER BY
				position
			LIMIT
				$3
		`

		positions []model.Position
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&positions,
		query,
		year,
		month,
		limit,
	); err != nil {
		return nil, fmt.Errorf("select positions: %w", err)
	}

	return model.ToDomPositions(positions), nil
}
