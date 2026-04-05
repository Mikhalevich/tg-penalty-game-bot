package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func (p *Postgres) OutboxIncrementRetryCount(
	ctx context.Context,
	ids []int,
	updatedAt time.Time,
) error {
	if len(ids) == 0 {
		return nil
	}

	var (
		query = `
			UPDATE outbox_messages
			SET
				retry_count = retry_count + 1,
				updated_at = ?
			WHERE
				id IN(?)
		`
	)

	query, args, err := sqlx.In(query, updatedAt, ids)
	if err != nil {
		return fmt.Errorf("sqlx in: %w", err)
	}

	trx := p.transactor.ExtContext(ctx)

	res, err := trx.ExecContext(ctx, trx.Rebind(query), args...)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return ErrNoRowsUpdated
	}

	return nil
}
