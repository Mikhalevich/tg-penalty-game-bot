package leaderboard

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (l *Leaderboard) PlayerPosition(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	currentPlayer, err := l.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by id: %w", err)
	}

	positions, err := l.repo.PlayerScorePosition(ctx, currentPlayer.ID, positionsLimit)
	if err != nil {
		return fmt.Errorf("player score position: %w", err)
	}

	if len(positions) == 0 {
		if err := l.notifier.LeaderboardRestrict(ctx, chatID); err != nil {
			return fmt.Errorf("show leaderboard restriction: %w", err)
		}

		return nil
	}

	if err := l.notifier.ShowPlayerLeaderbord(
		ctx,
		chatID,
		messageID,
		currentPlayer.ID,
		positions,
	); err != nil {
		return fmt.Errorf("show leaderboard: %w", err)
	}

	return nil
}
