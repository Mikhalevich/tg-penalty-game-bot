package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (p *Postgres) GetButton(ctx context.Context, btnID button.ID) (*button.Button, error) {
	var (
		query = `
			SELECT
				id,
				caption,
				operation,
				is_delete_message,
				payload
			FROM
				button
			WHERE
				id = $1
		`

		btn model.Button
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&btn,
		query,
		btnID.String(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, perror.NotFound("button not found")
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return btn.ToDomButton(), nil
}
