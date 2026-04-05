package leaderboard

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	positionsLimit = 10
)

type Repository interface {
	PlayerScorePosition(ctx context.Context, playerID player.ID, limit int) ([]player.Position, error)
	PlayerScoreTop(ctx context.Context, limit int) ([]player.Position, error)
}

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
}

type Notifier interface {
	ShowPlayerLeaderbord(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		playerID player.ID,
		positions []player.Position,
	) error
	ShowTopLeaderbord(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		playerID player.ID,
		positions []player.Position,
	) error
	LeaderboardRestrict(ctx context.Context, chatID msginfo.ChatID) error
}

type Leaderboard struct {
	repo           Repository
	playerProvider PlayerProvider
	notifier       Notifier
}

func New(
	repo Repository,
	playerProvider PlayerProvider,
	notifier Notifier,
) *Leaderboard {
	return &Leaderboard{
		repo:           repo,
		playerProvider: playerProvider,
		notifier:       notifier,
	}
}
