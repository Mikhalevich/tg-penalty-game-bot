package leaderboard

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

const (
	playersByPage = 10
)

func (l *Leaderboard) Page(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	pageNumber int,
) error {
	if pageNumber < 1 {
		return fmt.Errorf("invalid page number: %d", pageNumber)
	}

	maxPosition, err := l.repo.PlayerScoreMaxPosition(ctx)
	if err != nil {
		return fmt.Errorf("max position: %w", err)
	}

	pagesCount := calculatePageCount(maxPosition, playersByPage)

	if pageNumber > pagesCount {
		return fmt.Errorf("invalid page number: %d, all pages: %d", pageNumber, pagesCount)
	}

	positions, err := l.repo.PlayerScorePositionsFrom(
		ctx,
		calculatePositionFrom(pageNumber, playersByPage),
		playersByPage,
	)

	if err != nil {
		return fmt.Errorf("load player positions: %w", err)
	}

	if len(positions) == 0 {
		if err := l.notifier.LeaderboardRestrict(ctx, chatID); err != nil {
			return fmt.Errorf("show leaderboard restriction: %w", err)
		}

		return nil
	}

	currentPlayer, err := l.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by id: %w", err)
	}

	if err := l.notifier.ShowLeaderboard(
		ctx,
		chatID,
		messageID,
		currentPlayer.ID,
		pageNumber,
		pagesCount,
		positions,
	); err != nil {
		return fmt.Errorf("show leaderboard: %w", err)
	}

	return nil
}

func calculatePageCount(maxPosition, pageSize int) int {
	var (
		fullPages    = maxPosition / pageSize
		lastPageSize = maxPosition % pageSize
	)

	if lastPageSize > 0 {
		return fullPages + 1
	}

	return fullPages
}

func calculatePositionFrom(pageNumber, pageSize int) int {
	return (pageNumber-1)*pageSize + 1
}
