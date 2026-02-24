package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (p *Postgres) GetGame(ctx context.Context, gameID game.ID) (game.Game, error) {
	var (
		query = `
			SELECT
				id,
				created_at,
				game_status,
				game_status_changed_at,
				payload,
				payload_version
			FROM
				game
			WHERE
				game_id = $1
		`

		dbGame model.Game
	)

	if err := sqlx.GetContext(ctx, p.transactor.ExtContext(ctx), &dbGame, query, gameID.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return game.Game{}, ErrNotFound
		}

		return game.Game{}, fmt.Errorf("select game by id: %w", err)
	}

	domGame, err := dbGame.ToGame()
	if err != nil {
		return game.Game{}, fmt.Errorf("convert to domain game: %w", err)
	}

	return domGame, nil
}
