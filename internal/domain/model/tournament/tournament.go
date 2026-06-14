package tournament

import (
	"time"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(id string) ID {
	return ID(id)
}

type TournamentStatus string

const (
	TournamentStatusPending    TournamentStatus = "pending"
	TournamentStatusInProgress TournamentStatus = "in_progress"
	TournamentStatusCompleted  TournamentStatus = "completed"
	TournamentStatusCanceled   TournamentStatus = "canceled"
)

func (ts TournamentStatus) String() string {
	return string(ts)
}

type Tournament struct {
	ID             ID
	CreatedAt      time.Time
	Status         TournamentStatus
	State          State
	StateVersion   int
	StateUpdatedAt time.Time
}
