package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
	"github.com/jmoiron/sqlx"
)

func (p *Postgres) SetIdleForLongReadyToGamePlayerStatus(
	ctx context.Context,
	startSearchBefore time.Time,
	updatedAt time.Time,
) ([]player.Player, error) {
	var (
		query = `
			UPDATE player SET
				game_status = :game_status_idle,
				game_status_changed_at = :game_status_changed_at,
				current_game_id = :current_game_id,
			WHERE
				game_status = :game_status_ready_for_game AND
				game_status_changd_at <= :start_search_before
			RETURNING
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
		`

		players []model.Player
		trx     = p.transactor.ExtContext(ctx)
	)

	query, args, err := sqlx.Named(query, map[string]any{
		"game_status_idle":           player.GameStatusIdle,
		"game_status_changed_at":     updatedAt,
		"current_game_id":            "",
		"game_status_ready_for_game": player.GameStatusReadyForGame,
		"start_search_before":        startSearchBefore,
	})

	if err != nil {
		return nil, fmt.Errorf("prepare named query: %w", err)
	}

	if err := sqlx.SelectContext(
		ctx,
		trx,
		&players,
		trx.Rebind(query),
		args...,
	); err != nil {
		return nil, fmt.Errorf("select players: %w", err)
	}

	return model.ToDomainPlayers(players), nil
}
