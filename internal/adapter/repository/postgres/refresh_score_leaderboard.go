package postgres

import (
	"context"
	"fmt"
)

func (p *Postgres) RefreshScoreLeaderboard(ctx context.Context) error {
	var (
		query = `
			REFRESH MATERIALIZED VIEW score_leaderboard
		`
	)

	if _, err := p.transactor.ExtContext(ctx).ExecContext(ctx, query); err != nil {
		return fmt.Errorf("refresh view: %w", err)
	}

	return nil
}
