package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/tournament"
)

func (p *Postgres) GetTournament(
	ctx context.Context,
	tournamentID tournament.ID,
) (tournament.Tournament, error) {
	var (
		query = `
			SELECT
				id,
				created_at,
				tournament_status,
				payload,
				payload_version,
				payload_updated_at
			FROM
				tournament
			WHERE
				id = $1
			FOR UPDATE
		`

		dbTournament model.Tournament
	)

	if err := sqlx.GetContext(ctx, p.transactor.ExtContext(ctx), &dbTournament, query, tournamentID.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tournament.Tournament{}, perror.NotFound("tournament not found")
		}

		return tournament.Tournament{}, fmt.Errorf("select tournament by id: %w", err)
	}

	domTournament, err := dbTournament.ToDomTournament()
	if err != nil {
		return tournament.Tournament{}, fmt.Errorf("convert to domain tournament: %w", err)
	}

	return domTournament, nil
}
