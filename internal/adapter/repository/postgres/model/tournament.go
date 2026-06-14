package model

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/tournament"
)

type Tournament struct {
	ID               string      `db:"id"`
	CreatedAt        time.Time   `db:"created_at"`
	Status           string      `db:"tournament_status"`
	Payload          jsonb.JSONB `db:"payload"`
	PayloadVersion   int         `db:"payload_version"`
	PayloadUpdatedAt time.Time   `db:"payload_updated_at"`
}

func (t Tournament) ToDomTournament() (tournament.Tournament, error) {
	var state tournament.State
	if err := jsonb.ConvertTo(t.Payload, &state); err != nil {
		return tournament.Tournament{}, fmt.Errorf("convert payload to state: %w", err)
	}

	return tournament.Tournament{
		ID:             tournament.IDFromString(t.ID),
		CreatedAt:      t.CreatedAt,
		Status:         tournament.TournamentStatus(t.Status),
		State:          state,
		StateVersion:   t.PayloadVersion,
		StateUpdatedAt: t.PayloadUpdatedAt,
	}, nil
}

func ToDBTrounament(dom tournament.Tournament) (Tournament, error) {
	payload, err := jsonb.NewFromMarshaler(dom.State)
	if err != nil {
		return Tournament{}, fmt.Errorf("marshal payload: %w", err)
	}

	return Tournament{
		ID:               dom.ID.String(),
		CreatedAt:        dom.CreatedAt,
		Status:           dom.Status.String(),
		Payload:          payload,
		PayloadVersion:   dom.StateVersion,
		PayloadUpdatedAt: dom.StateUpdatedAt,
	}, nil
}
