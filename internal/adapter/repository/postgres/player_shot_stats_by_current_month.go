package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) PlayerShotStatsByCurrentMonth(
	ctx context.Context,
	playerID player.ID,
	shotType game.ShotType,
) ([]game.ShotSidePercent, error) {
	var (
		query = `
			WITH total_shots AS (
				SELECT
					shot_type,
					side,
					COUNT(*) AS total_count
				FROM
					shot
				WHERE
					player_id = $1 AND
					shot_type = $2 AND
					created_at >= DATE_TRUNC('month', CURRENT_DATE)
				GROUP BY
					shot_type, side
			)
			SELECT
				side,
				ROUND(total_count * 100 / SUM(total_count) OVER (), 2) AS percent
			FROM
				total_shots
		`

		shots []model.ShotSidePercent
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&shots,
		query,
		playerID,
		shotType,
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return model.ToDomShotSidePercents(shots), nil
}
