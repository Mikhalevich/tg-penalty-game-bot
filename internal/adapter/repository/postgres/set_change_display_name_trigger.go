package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (p *Postgres) SetChangeDisplayNameTrigger(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	var (
		query = `
			UPDATE player
			SET
				is_change_name_triggered = TRUE
			WHERE
				chat_id = :chat_id
		`
	)

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, map[string]any{
		"chat_id": chatID.Int64(),
	})

	if err != nil {
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
