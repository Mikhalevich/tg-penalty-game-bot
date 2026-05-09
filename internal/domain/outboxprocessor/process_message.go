package outboxprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/outboxmsg"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

func (o *OutboxProcessor) ProcessMessage(
	ctx context.Context,
	batchSize int,
	maxRetryCount int,
) error {
	now := o.timeProvider.Now()

	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		msgs, err := o.repository.OutboxSelectForDispatchMessages(ctx, now, batchSize)
		if err != nil {
			return fmt.Errorf("select outbox messages: %w", err)
		}

		results := o.sendMessages(ctx, msgs, maxRetryCount)

		if len(results.DispatchedIDs) > 0 {
			if err := o.repository.OutboxUpdateStatus(
				ctx,
				results.DispatchedIDs,
				outboxmsg.StatusDispatched,
				now,
			); err != nil {
				return fmt.Errorf("set dispatched: %w", err)
			}
		}

		if len(results.CanceledIDs) > 0 {
			if err := o.repository.OutboxUpdateStatus(
				ctx,
				results.CanceledIDs,
				outboxmsg.StatusCanceled,
				now,
			); err != nil {
				return fmt.Errorf("set canceled: %w", err)
			}
		}

		if len(results.IncrementRetryCountIDs) > 0 {
			if err := o.repository.OutboxIncrementRetryCount(
				ctx,
				results.IncrementRetryCountIDs,
				now,
			); err != nil {
				return fmt.Errorf("increment retry count: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

type sendResults struct {
	DispatchedIDs          []int
	CanceledIDs            []int
	IncrementRetryCountIDs []int
}

func (o *OutboxProcessor) sendMessages(
	ctx context.Context,
	msgs []outboxmsg.Message,
	maxRetryCount int,
) sendResults {
	results := sendResults{
		DispatchedIDs: make([]int, 0, len(msgs)),
	}

	for _, msg := range msgs {
		if err := o.sender.SendMessage(ctx, msg.Message); err != nil {
			processAndLogSendError(ctx, err, msg, maxRetryCount, &results)

			continue
		}

		results.DispatchedIDs = append(results.DispatchedIDs, msg.ID)
	}

	return results
}

func processAndLogSendError(
	ctx context.Context,
	err error,
	msg outboxmsg.Message,
	maxRetryCount int,
	results *sendResults,
) {
	loggerFromCtx := logger.FromContext(ctx).
		WithFields(
			logger.Fields{
				"chat_id":  msg.ChatID,
				"text":     msg.Text,
				"msg_type": msg.Type,
			},
		).
		WithError(err)

	loggerFromCtx.Error("send message")

	if msg.RetryCount < maxRetryCount {
		results.IncrementRetryCountIDs = append(results.IncrementRetryCountIDs, msg.ID)

		return
	}

	loggerFromCtx.Error("riched max retry count, skip message")

	results.CanceledIDs = append(results.CanceledIDs, msg.ID)
}
