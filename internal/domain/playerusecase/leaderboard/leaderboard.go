package leaderboard

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	PlayerScorePosition(ctx context.Context, playerID player.ID) (player.Position, error)
}

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
}

type Notifier interface {
	ShowLeaderbord(ctx context.Context, chatID msginfo.ChatID, positions []player.Position) error
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

func (l *Leaderboard) MyPosition(ctx context.Context, chatID msginfo.ChatID) error {
	currentPlayer, err := l.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by id: %w", err)
	}

	position, err := l.repo.PlayerScorePosition(ctx, currentPlayer.ID)
	if err != nil {
		return fmt.Errorf("player score position: %w", err)
	}

	if err := l.notifier.ShowLeaderbord(ctx, chatID, []player.Position{position}); err != nil {
		return fmt.Errorf("show leaderboard: %w", err)
	}

	return nil
}
