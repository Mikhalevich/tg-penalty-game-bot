package postgres_test

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/tournament"
)

func (s *PostgresSuit) TestCRUTournament() {
	var (
		ctx = context.Background()
		id  = tournament.IDFromString(uuid.NewString())
	)

	err := s.pgDB.InsertTournament(ctx, tournament.Tournament{
		ID:             id,
		CreatedAt:      time.Date(2026, 06, 14, 14, 49, 7, 0, time.Local),
		Status:         tournament.TournamentStatusPending,
		State:          tournament.State{},
		StateVersion:   1,
		StateUpdatedAt: time.Date(2026, 06, 14, 14, 49, 7, 0, time.Local),
	})
	s.Require().NoError(err)

	actualCreatedTournament, err := s.pgDB.GetTournament(ctx, id)
	s.Require().NoError(err)
	s.Require().Equal(tournament.Tournament{
		ID:             id,
		CreatedAt:      time.Date(2026, 06, 14, 14, 49, 7, 0, time.Local),
		Status:         tournament.TournamentStatusPending,
		State:          tournament.State{},
		StateVersion:   1,
		StateUpdatedAt: time.Date(2026, 06, 14, 14, 49, 7, 0, time.Local),
	}, actualCreatedTournament)

	err = s.pgDB.UpdateTournament(ctx, tournament.Tournament{
		ID:             id,
		CreatedAt:      time.Date(2027, 07, 16, 16, 49, 7, 0, time.Local),
		Status:         tournament.TournamentStatusInProgress,
		State:          tournament.State{},
		StateVersion:   1,
		StateUpdatedAt: time.Date(2026, 06, 15, 15, 49, 7, 0, time.Local),
	})
	s.Require().NoError(err)

	actualUpdatedTournament, err := s.pgDB.GetTournament(ctx, id)
	s.Require().NoError(err)
	s.Require().Equal(tournament.Tournament{
		ID:             id,
		CreatedAt:      time.Date(2026, 06, 14, 14, 49, 7, 0, time.Local),
		Status:         tournament.TournamentStatusInProgress,
		State:          tournament.State{},
		StateVersion:   2,
		StateUpdatedAt: time.Date(2026, 06, 15, 15, 49, 7, 0, time.Local),
	}, actualUpdatedTournament)
}

func (s *PostgresSuit) TestGetTournamentNotFoundError() {
	domTournament, err := s.pgDB.GetTournament(
		context.Background(),
		tournament.IDFromString(uuid.NewString()),
	)

	s.Require().EqualError(err, "tournament not found")
	s.Require().Equal(tournament.Tournament{}, domTournament)
}

func (s *PostgresSuit) TestUpdateTournamentInvalidVersionPayloadError() {
	var (
		ctx = context.Background()
		id  = tournament.IDFromString(uuid.NewString())
	)

	err := s.pgDB.InsertTournament(ctx, tournament.Tournament{
		ID:             id,
		CreatedAt:      time.Date(2026, 06, 14, 14, 49, 7, 0, time.Local),
		Status:         tournament.TournamentStatusPending,
		State:          tournament.State{},
		StateVersion:   2,
		StateUpdatedAt: time.Date(2026, 06, 14, 14, 49, 7, 0, time.Local),
	})
	s.Require().NoError(err)

	err = s.pgDB.UpdateTournament(ctx, tournament.Tournament{
		ID:             id,
		CreatedAt:      time.Date(2026, 06, 14, 14, 49, 7, 0, time.Local),
		Status:         tournament.TournamentStatusInProgress,
		State:          tournament.State{},
		StateVersion:   3,
		StateUpdatedAt: time.Date(2026, 06, 15, 15, 49, 7, 0, time.Local),
	})
	s.Require().EqualError(err, "no rows updated")
}
