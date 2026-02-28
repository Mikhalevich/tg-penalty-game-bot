package postgres

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/jmoiron/sqlx"
)

func (p *Postgres) GetRoundShots(
	ctx context.Context,
	gameID game.ID,
	round int,
) ([]game.Shot, error) {
	var (
		query = `
			SELECT
				game_id,
				player_id,
				round,
				shot_type,
				created_at,
				side,
				completed_at
			FROM
				shot
			WHERE
				game_id = $1 AND
				round = $2
		`

		shots []model.Shot
	)

	if err := sqlx.SelectContext(ctx, p.transactor.ExtContext(ctx), &shots, query, gameID, round); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return model.ToShots(shots), nil
}
