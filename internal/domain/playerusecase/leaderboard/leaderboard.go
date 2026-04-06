package leaderboard

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	PlayerScorePositionsFrom(
		ctx context.Context,
		positionFrom int,
		limit int,
	) ([]player.Position, error)
	PlayerScoreMaxPosition(ctx context.Context) (int, error)
}

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
}

type Notifier interface {
	ShowLeaderboard(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		playerID player.ID,
		pageNumber int,
		pagesCount int,
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
