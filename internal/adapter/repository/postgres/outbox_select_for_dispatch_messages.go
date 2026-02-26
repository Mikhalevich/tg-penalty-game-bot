package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/outboxmsg"
)

func (p *Postgres) OutboxSelectForDispatchMessages(
	ctx context.Context,
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
				buttons
			FROM
				outbox_messages
			WHERE
				is_dispatched = FALSE
			ORDER BY
				id
			LIMIT
				$1
			FOR UPDATE SKIP LOCKED
		`

		outboxMsgs []model.OutboxMessage
	)

	if err := sqlx.SelectContext(ctx, p.transactor.ExtContext(ctx), &outboxMsgs, query, limit); err != nil {
		return nil, fmt.Errorf("select messages: %w", err)
	}

	msgs, err := model.ToOutboxMessages(outboxMsgs)
	if err != nil {
		return nil, fmt.Errorf("convert to outbox messages: %w", err)
	}

	return msgs, nil
}
