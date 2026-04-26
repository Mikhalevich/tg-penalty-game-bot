package leaderboard

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (l *Leaderboard) MyPage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	currentPlayer, err := l.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by id: %w", err)
	}

	myPosition, err := l.repo.PlayerScoreMyPosition(ctx, currentPlayer.ID)
	if err != nil {
		return fmt.Errorf("my position: %w", err)
	}

	maxPosition, err := l.repo.PlayerScoreMaxPosition(ctx)
	if err != nil {
		return fmt.Errorf("max position: %w", err)
	}

	pagesCount := calculatePageCount(maxPosition, playersByPage)
	if pagesCount == 0 {
		return perror.NotExists("leaderboard is empty")
	}

	pageNumber := calculatePageByPosition(myPosition.Position, playersByPage)

	positions, err := l.repo.PlayerScorePositionsFrom(
		ctx,
		calculatePositionFrom(pageNumber, playersByPage),
		playersByPage,
	)

	if err != nil {
		return fmt.Errorf("load player positions: %w", err)
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
