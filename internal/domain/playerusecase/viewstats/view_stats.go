package viewstats

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
}

type GameProvider interface {
	GetGameByID(ctx context.Context, gameID game.ID) (game.Game, error)
}

type PlayerShotStatser interface {
	PlayerStats(
		ctx context.Context,
		playerID player.ID,
	) (game.ShotStats, error)
}

type MessageDeleter interface {
	DeleteMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
	) error
}

type Notifier interface {
	StartGameWithStats(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		playerAgainst game.Player,
		shotStats game.ShotStats,
	) error
}

type ViewStats struct {
	playerProvider PlayerProvider
	gameProvider   GameProvider
	playerStats    PlayerShotStatser
	messageDeleter MessageDeleter
	notifier       Notifier
}

func New(
	playerProvider PlayerProvider,
	gameProvider GameProvider,
	playerStats PlayerShotStatser,
	messageDeleter MessageDeleter,
	notifier Notifier,
) *ViewStats {
	return &ViewStats{
		playerProvider: playerProvider,
		gameProvider:   gameProvider,
		playerStats:    playerStats,
		messageDeleter: messageDeleter,
		notifier:       notifier,
	}
}
