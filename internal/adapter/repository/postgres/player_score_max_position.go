package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (p *Postgres) PlayerScoreMaxPosition(ctx context.Context) (int, error) {
	var (
		query = `
			SELECT
				COALESCE(MAX(position), 0)
			FROM
				score_leaderboard
		`

		maxPosition int
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&maxPosition,
		query,
	); err != nil {
		return 0, fmt.Errorf("select max position: %w", err)
	}

	return maxPosition, nil
}
