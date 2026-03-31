package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) GetPlayerByChatID(
	ctx context.Context,
	chatID msginfo.ChatID,
) (player.Player, error) {
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
				game_status_changed_at,
				current_game_id,
				score
			FROM
				player
			WHERE
				chat_id = $1
		`

		plr model.Player
	)

	if err := sqlx.GetContext(ctx, p.transactor.ExtContext(ctx), &plr, query, chatID.Int64()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return player.Player{}, ErrNotFound
		}

		return player.Player{}, fmt.Errorf("get context: %w", err)
	}

	return plr.ToDomainPlayer(), nil
}
