package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/tournament"
)

func (p *Postgres) InsertTournament(
	ctx context.Context,
	domTournnament tournament.Tournament,
) error {
	var (
		query = `
			INSERT INTO tournament(
				id,
				created_at,
				tournament_status,
				payload,
				payload_version,
				payload_updated_at
			) VALUES (
				:id,
				:created_at,
				:tournament_status,
				:payload,
				:payload_version,
				:payload_updated_at
			)
		`
	)

	dbTournament, err := model.ToDBTrounament(domTournnament)
	if err != nil {
		return fmt.Errorf("convert to db tournament: %w", err)
	}

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		dbTournament,
	)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
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
