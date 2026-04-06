package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/outboxmsg"
)

func (p *Postgres) OutboxSelectForDispatchMessages(
	ctx context.Context,
	visibilityAt time.Time,
	limit int,
) ([]outboxmsg.Message, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				reply_msg_id,
				msg_text,
				msg_type,
				payload,
				buttons,
				retry_count
			FROM
				outbox_messages
			WHERE
				status = 'pending' AND
				visibility_at <= $1
			ORDER BY
				id
			LIMIT
				$2
			FOR UPDATE SKIP LOCKED
		`

		outboxMsgs []model.OutboxMessage
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&outboxMsgs,
		query,
		visibilityAt,
		limit,
	); err != nil {
		return nil, fmt.Errorf("select messages: %w", err)
	}

	msgs, err := model.ToOutboxMessages(outboxMsgs)
	if err != nil {
		return nil, fmt.Errorf("convert to outbox messages: %w", err)
	}

	return msgs, nil
}
