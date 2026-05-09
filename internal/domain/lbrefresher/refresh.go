package lbrefresher

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

const (
	historyPositionsLimit = 10
)

func (lr *LeaderboardRefresher) Refresh(ctx context.Context) error {
	var (
		now = lr.timeProvider.Now()
	)

	isChanged, err := lr.isMonthChangedSinceLastHistoryUpdate(ctx, now.Month())
	if err != nil {
		return fmt.Errorf("check is month changed: %w", err)
	}

	if isChanged {
		if err := lr.updateHistoryLeaderboard(ctx, now); err != nil {
			return fmt.Errorf("update history leaderboard: %w", err)
		}
	}

	if err := lr.repo.RefreshScoreLeaderboard(ctx); err != nil {
		return fmt.Errorf("refresh score leaderboard: %w", err)
	}

	return nil
}

func (lr *LeaderboardRefresher) updateHistoryLeaderboard(
	ctx context.Context,
	now time.Time,
) error {
	if err := lr.transactor.Transaction(ctx, func(ctx context.Context) error {
		positions, err := lr.repo.PlayerScorePositionsFrom(ctx, 0, historyPositionsLimit)
		if err != nil {
			return fmt.Errorf("player positions by month: %w", err)
		}

		if err := lr.repo.InsertHistoryScoreLeaderboard(ctx, HistoryPositions{
			Year:      now.Year(),
			Month:     int(now.Month()),
			Positions: positions,
		}); err != nil {
			return fmt.Errorf("insert history leaderboard: %w", err)
		}

		if err := lr.settingsProvider.SetLeaderboardMonth(ctx, now.Month()); err != nil {
			return fmt.Errorf("set leaderboard month setting: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

// isMonthChangedSinceLastHistoryUpdate check is month changed since last time store history positions
// returns true is month was changed and false otherwise.
func (lr *LeaderboardRefresher) isMonthChangedSinceLastHistoryUpdate(
	ctx context.Context,
	currentMonth time.Month,
) (bool, error) {
	settingsMonth, err := lr.settingsProvider.GetLeaderboardMonth(ctx)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return false, fmt.Errorf("get leaderboard month from settings: %w", err)
		}

		if err := lr.settingsProvider.SetLeaderboardMonth(ctx, currentMonth); err != nil {
			return false, fmt.Errorf("set leaderboard month setting: %w", err)
		}

		return false, nil
	}

	return settingsMonth != currentMonth, nil
}
