package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
)

func (p *Postgres) SetButtonRows(
	ctx context.Context,
	rows ...button.ButtonRow,
) error {
	var (
		query = `
			INSERT INTO button(
				id,
				caption,
				operation,
				is_delete_message,
				payload
			) VALUES (
				:id,
				:caption,
				:operation,
				:is_delete_message,
				:payload
			)
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.ToDBButtons(rows),
	)
	if err != nil {
		return fmt.Errorf("insert buttons: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
