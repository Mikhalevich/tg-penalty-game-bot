package viewstats

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (vs *ViewStats) ViewStatsOnStartGameMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	playerID player.ID,
	playerDisplayName string,
) error {
	shotStats, err := vs.playerStats.PlayerStats(ctx, playerID)
	if err != nil {
		return fmt.Errorf("player stats: %w", err)
	}

	if err := vs.notifier.StartGameWithStats(
		ctx,
		chatID,
		messageID,
		playerDisplayName,
		shotStats,
	); err != nil {
		return fmt.Errorf("start game with stats notification: %w", err)
	}

	return nil
}
