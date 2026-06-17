package tournament_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/tournament"
)

func TestRoundRobinSchedule(t *testing.T) {
	t.Parallel()
	t.Run("two teams", func(t *testing.T) {
		t.Parallel()

		rounds := tournament.RoundRobinSchedule(2, func() game.ID {
			return game.IDFromString("")
		})
		require.Equal(t, []tournament.Round{
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
		}, rounds)
	})

	t.Run("three teams", func(t *testing.T) {
		t.Parallel()

		rounds := tournament.RoundRobinSchedule(3, func() game.ID {
			return game.IDFromString("")
		})
		require.Equal(t, []tournament.Round{
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
		}, rounds)
	})
}
