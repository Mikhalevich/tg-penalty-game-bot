package tournament_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/tournament"
)

func TestRoundRobinSchedule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name             string
		TeamCount        int
		ExpectedSchedule []tournament.Round
	}{
		{
			Name:      "two teams",
			TeamCount: 2,
			ExpectedSchedule: []tournament.Round{
				{
					{
						Home: tournament.MatchTeamInfo{
							Idx: 0,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 1,
						},
					},
				},
			},
		},
		{
			Name:      "three teams",
			TeamCount: 3,
			ExpectedSchedule: []tournament.Round{
				{
					{
						Home: tournament.MatchTeamInfo{
							Idx: 1,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 2,
						},
					},
				},
				{
					{
						Home: tournament.MatchTeamInfo{
							Idx: 2,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 0,
						},
					},
				},
				{
					{
						Home: tournament.MatchTeamInfo{
							Idx: 0,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 1,
						},
					},
				},
			},
		},
		{
			Name:      "four teams",
			TeamCount: 4,
			ExpectedSchedule: []tournament.Round{
				{
					{
						Home: tournament.MatchTeamInfo{
							Idx: 0,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 3,
						},
					},
					{
						Home: tournament.MatchTeamInfo{
							Idx: 1,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 2,
						},
					},
				},
				{
					{
						Home: tournament.MatchTeamInfo{
							Idx: 2,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 0,
						},
					},
					{
						Home: tournament.MatchTeamInfo{
							Idx: 3,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 1,
						},
					},
				},
				{
					{
						Home: tournament.MatchTeamInfo{
							Idx: 0,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 1,
						},
					},
					{
						Home: tournament.MatchTeamInfo{
							Idx: 2,
						},
						Away: tournament.MatchTeamInfo{
							Idx: 3,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			rounds := tournament.RoundRobinSchedule(test.TeamCount, func() game.ID {
				return game.IDFromString("")
			})

			require.Equal(t, test.ExpectedSchedule, rounds)
		})
	}
}
