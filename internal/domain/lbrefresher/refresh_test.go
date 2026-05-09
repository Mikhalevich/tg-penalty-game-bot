package lbrefresher_test

import (
	"context"
	"errors"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (s *LeaderboardRefresherSuite) TestGetMonthFromSettingsError() {
	var (
		ctx = context.Background()
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().Now().Return(time.Date(2026, 5, 9, 7, 7, 7, 7, time.UTC)),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(time.January, errors.New("some error")),
	)

	err := s.refresher.Refresh(ctx)

	s.Require().EqualError(err, "check is month changed: get leaderboard month from settings: some error")
}

func (s *LeaderboardRefresherSuite) TestSetMonthSettingsError() {
	var (
		ctx = context.Background()
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().Now().Return(time.Date(2026, 5, 9, 7, 7, 7, 7, time.UTC)),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(time.January, perror.NotFound("month not found")),

		s.settingsProvider.EXPECT().
			SetLeaderboardMonth(ctx, time.May).
			Return(errors.New("some set error")),
	)

	err := s.refresher.Refresh(ctx)

	s.Require().EqualError(err, "check is month changed: set leaderboard month setting: some set error")
}
