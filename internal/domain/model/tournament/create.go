package tournament

import (
	"time"

	"github.com/google/uuid"
)

func Create(
	createdAt time.Time,
	plr Player,
) Tournament {
	return Tournament{
		ID:        IDFromString(uuid.NewString()),
		CreatedAt: createdAt,
		Status:    TournamentStatusPending,
		State: State{
			Teams: []Team{
				{
					Player: plr,
				},
			},
		},
		StateVersion:   1,
		StateUpdatedAt: createdAt,
	}
}
