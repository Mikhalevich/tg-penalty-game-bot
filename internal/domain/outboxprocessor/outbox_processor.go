package outboxprocessor

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/outboxmsg"
)

type Repository interface {
	OutboxSelectForDispatchMessages(
		ctx context.Context,
		visibilityAt time.Time,
		limit int,
	) ([]outboxmsg.Message, error)

	OutboxUpdateStatus(
		ctx context.Context,
		ids []int,
		status outboxmsg.Status,
		updatedAt time.Time,
	) error

	OutboxIncrementRetryCount(
		ctx context.Context,
		ids []int,
		updatedAt time.Time,
	) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type Sender interface {
	SendMessage(ctx context.Context, msg msginfo.Message) error
}

type TimeProvider interface {
	Now() time.Time
}

type OutboxProcessor struct {
	repository   Repository
	transactor   Transactor
	sender       Sender
	timeProvider TimeProvider
}

func New(
	repository Repository,
	transactor Transactor,
	sender Sender,
	timeProvider TimeProvider,
) *OutboxProcessor {
	return &OutboxProcessor{
		repository:   repository,
		transactor:   transactor,
		sender:       sender,
		timeProvider: timeProvider,
	}
}
