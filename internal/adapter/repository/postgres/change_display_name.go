package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

const (
	uniqPlayerDisplayNameConstraint = "player_display_name_u_idx"
)

func (p *Postgres) ChangeDisplayName(
	ctx context.Context,
	chatID msginfo.ChatID,
	displayName string,
	changedAt time.Time,
) error {
	var (
		query = `
			UPDATE player
			SET
				display_name = :display_name,
				is_change_name_triggered = FALSE,
				name_changed_at = :name_changed_at
			WHERE
				chat_id = :chat_id
		`
	)

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, map[string]any{
		"display_name":    displayName,
		"chat_id":         chatID.Int64(),
		"name_changed_at": changedAt,
	})

	if err != nil {
		if p.dbDriver.IsConstraintError(err, uniqPlayerDisplayNameConstraint) {
			return ErrAlreadyExists
		}

		return fmt.Errorf("named exec context: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
