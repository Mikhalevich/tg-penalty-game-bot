package postgres_test

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/lbrefresher"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (s *PostgresSuit) TestInsertAndGetHistoryScoreLeaderboard() {
	var (
		ctx = context.Background()
		//nolint:gosmopolitan
		aprilTime = time.Date(2026, 4, 9, 7, 7, 7, 0, time.Local)
		//nolint:gosmopolitan
		mayTime        = time.Date(2026, 5, 9, 7, 7, 7, 0, time.Local)
		positionsApril = lbrefresher.HistoryPositions{
			Year:  2026,
			Month: 4,
			Positions: []player.Position{
				{
					ID:          1,
					DisplayName: "1",
					Score:       100,
					UpdatedAt:   aprilTime,
					Position:    1,
				},
				{
					ID:          2,
					DisplayName: "2",
					Score:       70,
					UpdatedAt:   aprilTime,
					Position:    2,
				},
			},
		}
		positionsMay = lbrefresher.HistoryPositions{
			Year:  2026,
			Month: 5,
			Positions: []player.Position{
				{
					ID:          1,
					DisplayName: "1",
					Score:       100,
					UpdatedAt:   mayTime,
					Position:    1,
				},
				{
					ID:          2,
					DisplayName: "2",
					Score:       70,
					UpdatedAt:   mayTime,
					Position:    2,
				},
				{
					ID:          3,
					DisplayName: "3",
					Score:       50,
					UpdatedAt:   mayTime,
					Position:    3,
				},
			},
		}
	)

	err := s.pgDB.InsertHistoryScoreLeaderboard(ctx, positionsApril)
	s.Require().NoError(err)

	actualPositionsApril, err := s.pgDB.HistoryScoreLeaderboard(ctx, 2026, 4, 10)
	s.Require().NoError(err)
	s.Require().Equal(positionsApril.Positions, actualPositionsApril)

	err = s.pgDB.InsertHistoryScoreLeaderboard(ctx, positionsMay)
	s.Require().NoError(err)

	actualPositionsMay, err := s.pgDB.HistoryScoreLeaderboard(ctx, 2026, 5, 10)
	s.Require().NoError(err)
	s.Require().Equal(positionsMay.Positions, actualPositionsMay)

	actualPositionsSep, err := s.pgDB.HistoryScoreLeaderboard(ctx, 2026, 9, 10)
	s.Require().NoError(err)
	s.Require().Empty(actualPositionsSep)
}
