package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) CreatePlayer(ctx context.Context, plr player.Player) (player.ID, error) {
	var (
		query = `
			INSERT INTO player(
				chat_id,
				display_name,
				created_at
			) VALUES (
				:chat_id,
				:display_name,
				:created_at
			)
			RETURNING id
		`

		trx      = p.transactor.ExtContext(ctx)
		playerID player.ID
	)

	query, args, err := sqlx.Named(query, model.ToDBPlayer(plr))
	if err != nil {
		return 0, fmt.Errorf("sqlx named: %w", err)
	}

	if err := sqlx.GetContext(ctx, trx, &playerID, trx.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRowsUpdated
		}

		return 0, fmt.Errorf("get context: %w", err)
	}

	return playerID, nil
}
