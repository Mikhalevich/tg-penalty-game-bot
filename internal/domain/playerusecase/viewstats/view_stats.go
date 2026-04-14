package viewstats

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type PlayerShotStatser interface {
	PlayerStats(
		ctx context.Context,
		playerID player.ID,
	) (game.ShotStats, error)
}

type Notifier interface {
	StartGameWithStats(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		displayName string,
		shotStats game.ShotStats,
	) error
}

type ViewStats struct {
	playerStats PlayerShotStatser
	notifier    Notifier
}

func New(
	playerStats PlayerShotStatser,
	notifier Notifier,
) *ViewStats {
	return &ViewStats{
		playerStats: playerStats,
		notifier:    notifier,
	}
}
