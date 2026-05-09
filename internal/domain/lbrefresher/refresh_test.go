package lbrefresher_test

import (
	"context"
	"errors"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/lbrefresher"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

var (
	//nolint:gochecknoglobals
	currentTime = time.Date(2026, 5, 9, 7, 7, 7, 7, time.UTC)
	//nolint:gochecknoglobals
	positions = []player.Position{
		{
			ID:          1,
			DisplayName: "1",
			Score:       100,
			UpdatedAt:   currentTime,
			Position:    1,
		},
		{
			ID:          2,
			DisplayName: "2",
			Score:       70,
			UpdatedAt:   currentTime,
			Position:    2,
		},
		{
			ID:          3,
			DisplayName: "3",
			Score:       50,
			UpdatedAt:   currentTime,
			Position:    3,
		},
	}
)

func (s *LeaderboardRefresherSuite) TestGetMonthFromSettingsError() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(0, errors.New("some error")),
	)

	err := refresher.Refresh(ctx)

	s.Require().EqualError(err, "check is month changed: get leaderboard month from settings: some error")
}

func (s *LeaderboardRefresherSuite) TestSetMonthSettingsError() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(0, perror.NotFound("month not found")),

		s.settingsProvider.EXPECT().
			SetLeaderboardMonth(ctx, 5).
			Return(errors.New("some set error")),
	)

	err := refresher.Refresh(ctx)

	s.Require().EqualError(err, "check is month changed: set leaderboard month setting: some set error")
}

func (s *LeaderboardRefresherSuite) TestMonthNotChangedFromSettings() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(5, nil),

		s.repo.EXPECT().RefreshScoreLeaderboard(ctx),
	)

	err := refresher.Refresh(ctx)

	s.Require().NoError(err)
}

func (s *LeaderboardRefresherSuite) TestMonthNotFoundInSettings() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(0, perror.NotFound("month not found")),

		s.settingsProvider.EXPECT().
			SetLeaderboardMonth(ctx, 5).
			Return(nil),

		s.repo.EXPECT().
			RefreshScoreLeaderboard(ctx).
			Return(nil),
	)

	err := refresher.Refresh(ctx)

	s.Require().NoError(err)
}

func (s *LeaderboardRefresherSuite) TestMonthIsChanged() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(4, nil),

		s.transactor.EXPECT().
			Transaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, callback func(ctx context.Context) error) error {
				return callback(ctx)
			}),

		s.repo.EXPECT().
			PlayerScorePositionsFrom(ctx, 0, 10).
			Return(positions, nil),

		s.repo.EXPECT().
			InsertHistoryScoreLeaderboard(
				ctx,
				lbrefresher.HistoryPositions{
					Year:      2026,
					Month:     5,
					Positions: positions,
				},
			).Return(nil),

		s.settingsProvider.EXPECT().
			SetLeaderboardMonth(ctx, 5).
			Return(nil),

		s.repo.EXPECT().
			RefreshScoreLeaderboard(ctx).
			Return(nil),
	)

	err := refresher.Refresh(ctx)

	s.Require().NoError(err)
}

func (s *LeaderboardRefresherSuite) TestMonthIsChangedTrxError() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(4, nil),

		s.transactor.EXPECT().
			Transaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, callback func(ctx context.Context) error) error {
				//nolint:errcheck
				callback(ctx)

				return errors.New("trx some error")
			}),

		s.repo.EXPECT().
			PlayerScorePositionsFrom(ctx, 0, 10).
			Return(positions, nil),

		s.repo.EXPECT().
			InsertHistoryScoreLeaderboard(
				ctx,
				lbrefresher.HistoryPositions{
					Year:      2026,
					Month:     5,
					Positions: positions,
				},
			).Return(nil),

		s.settingsProvider.EXPECT().
			SetLeaderboardMonth(ctx, 5).
			Return(nil),
	)

	err := refresher.Refresh(ctx)

	s.Require().EqualError(err, "update history leaderboard: transaction: trx some error")
}

func (s *LeaderboardRefresherSuite) TestMonthIsChangedTrxSetLeaderboardError() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(4, nil),

		s.transactor.EXPECT().
			Transaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, callback func(ctx context.Context) error) error {
				return callback(ctx)
			}),

		s.repo.EXPECT().
			PlayerScorePositionsFrom(ctx, 0, 10).
			Return(positions, nil),

		s.repo.EXPECT().
			InsertHistoryScoreLeaderboard(
				ctx,
				lbrefresher.HistoryPositions{
					Year:      2026,
					Month:     5,
					Positions: positions,
				},
			).Return(nil),

		s.settingsProvider.EXPECT().
			SetLeaderboardMonth(ctx, 5).
			Return(errors.New("some set error")),
	)

	err := refresher.Refresh(ctx)

	s.Require().EqualError(err, "update history leaderboard: transaction: set leaderboard month setting: some set error")
}

func (s *LeaderboardRefresherSuite) TestMonthIsChangedTrxInsertHistoryError() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(4, nil),

		s.transactor.EXPECT().
			Transaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, callback func(ctx context.Context) error) error {
				return callback(ctx)
			}),

		s.repo.EXPECT().
			PlayerScorePositionsFrom(ctx, 0, 10).
			Return(positions, nil),

		s.repo.EXPECT().
			InsertHistoryScoreLeaderboard(
				ctx,
				lbrefresher.HistoryPositions{
					Year:      2026,
					Month:     5,
					Positions: positions,
				},
			).Return(errors.New("insert history some error")),
	)

	err := refresher.Refresh(ctx)

	s.Require().EqualError(err,
		"update history leaderboard: transaction: insert history leaderboard: insert history some error")
}

func (s *LeaderboardRefresherSuite) TestMonthIsChangedTrxPlayerPositionsError() {
	var (
		ctx       = context.Background()
		refresher = lbrefresher.New(s.repo, s.transactor, s.settingsProvider, s.timeProvider)
	)

	gomock.InOrder(
		s.timeProvider.EXPECT().
			Now().
			Return(currentTime),

		s.settingsProvider.EXPECT().
			GetLeaderboardMonth(ctx).
			Return(4, nil),

		s.transactor.EXPECT().
			Transaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, callback func(ctx context.Context) error) error {
				return callback(ctx)
			}),

		s.repo.EXPECT().
			PlayerScorePositionsFrom(ctx, 0, 10).
			Return(nil, errors.New("some positions error")),
	)

	err := refresher.Refresh(ctx)

	s.Require().EqualError(err, "update history leaderboard: transaction: player positions by month: some positions error")
}
