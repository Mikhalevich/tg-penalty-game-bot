package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (p *Postgres) SendMessage(
	ctx context.Context,
	msg msginfo.Message,
) error {
	dbOutbox, err := model.ToDBOutboxMessage(msg)
	if err != nil {
		return fmt.Errorf("make db outbox message: %w", err)
	}

	if err := p.insertOutboxMessage(ctx, dbOutbox); err != nil {
		return fmt.Errorf("insert outbox message: %w", err)
	}

	return nil
}

func (p *Postgres) insertOutboxMessage(ctx context.Context, msg model.OutboxMessage) error {
	var (
		query = `
			INSERT INTO outbox_messages(
				chat_id,
				reply_msg_id,
				msg_text,
				msg_type,
				payload,
				buttons,
				visibility_at
			) VALUES (
				:chat_id,
				:reply_msg_id,
				:msg_text,
				:msg_type,
				:payload,
				:buttons,
				:visibility_at
			)
		`
	)

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, msg)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNoRowsUpdated
	}

	return nil
}
