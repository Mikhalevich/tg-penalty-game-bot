package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/lbrefresher"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (p *Postgres) InsertHistoryScoreLeaderboard(
	ctx context.Context,
	positions lbrefresher.HistoryPositions,
) error {
	var (
		query = `
			INSERT INTO leaderboard_history(
				player_id,
				year,
				month,
				display_name,
				score,
				updated_at,
				position
			) VALUES (
				:player_id,
				:year,
				:month,
				:display_name,
				:score,
				:updated_at,
				:position
			)
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.ToDBHisotoryPositions(positions),
	)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
