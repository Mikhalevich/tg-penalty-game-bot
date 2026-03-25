package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/jmoiron/sqlx"
)

func (p *Postgres) GetShotExpiredRatingGames(
	ctx context.Context,
	expiredTime time.Time,
	limit int,
) ([]game.Game, error) {
	var (
		query = `
			SELECT
				id,
				created_at,
				game_type,
				game_status,
				payload,
				payload_version,
				payload_updated_at
			FROM
				game
			WHERE
				game_type = $1 AND
				game_status = $2 AND
				payload_updated_at < $3
			ORDER BY
				payload_updated_at
			LIMIT
				$4
		`

		dbGames []model.Game
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&dbGames,
		query,
		game.GameTypeRating,
		game.GameStatusInProgress,
		expiredTime,
		limit,
	); err != nil {
		return nil, fmt.Errorf("select games: %w", err)
	}

	games, err := model.ToDomainGames(dbGames)
	if err != nil {
		return nil, fmt.Errorf("convert to domain games: %w", err)
	}

	return games, nil
}
