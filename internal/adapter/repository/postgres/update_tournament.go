package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/tournament"
)

func (p *Postgres) UpdateTournament(
	ctx context.Context,
	domTournament tournament.Tournament,
) error {
	var (
		query = `
			UPDATE tournament SET
				tournament_status = :tournament_status,
				payload = :payload,
				payload_version = payload_version + 1,
				payload_updated_at = :payload_updated_at
			WHERE
				id = :id AND
				payload_version = :payload_version
		`
	)

	dbTournament, err := model.ToDBTrounament(domTournament)
	if err != nil {
		return fmt.Errorf("convert to db tournament: %w", err)
	}

	res, err := sqlx.NamedExecContext(ctx, p.transactor.ExtContext(ctx), query, dbTournament)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
