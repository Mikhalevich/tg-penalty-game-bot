package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) SelectReadyToGamePlayers(
	ctx context.Context,
	limit int,
) ([]player.Player, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				display_name,
				created_at,
				is_change_name_triggered,
				name_changed_at,
				game_status,
				game_status_changed_at
			FROM
				player
			WHERE
				game_status = $1
			ORDER BY
				game_status_changed_at
			FOR UPDATE SKIP LOCKED
		`

		players []model.Player
		trx     = p.transactor.ExtContext(ctx)
	)

	if err := sqlx.SelectContext(ctx, trx, &players, query, player.GameStatusReadyForGame.String()); err != nil {
		return nil, fmt.Errorf("select players: %w", err)
	}

	return model.ToDomainPlayers(players), nil
}
