package outboxprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

func (o *OutboxProcessor) ProcessMessage(ctx context.Context, batchSize int) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		msgs, err := o.repository.OutboxSelectForDispatchMessages(ctx, batchSize)
		if err != nil {
			return fmt.Errorf("select outbox messages: %w", err)
		}

		ids := make([]int, 0, len(msgs))

		for _, msg := range msgs {
			if err := o.sender.SendMessage(ctx, msg.Message); err != nil {
				logger.FromContext(ctx).
					WithFields(
						logger.Fields{
							"chat_id":  msg.ChatID,
							"text":     msg.Text,
							"msg_type": msg.Type,
						},
					).
					WithError(err).
					Error("send message")

				continue
			}

			ids = append(ids, msg.ID)
		}

		if len(ids) > 0 {
			if err := o.repository.OutboxSetDispatched(ctx, ids, o.timeProvider.Now()); err != nil {
				return fmt.Errorf("set dispatched: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
