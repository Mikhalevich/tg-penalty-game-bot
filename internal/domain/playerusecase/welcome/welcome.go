package welcome

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type PlayerProvider interface {
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
}

type Notifier interface {
	WelcomeNewPlayer(ctx context.Context, plr player.Player) error
}

type Welcome struct {
	playerProvider PlayerProvider
	notifier       Notifier
}

func New(
	playerProvider PlayerProvider,
	notifier Notifier,
) *Welcome {
	return &Welcome{
		playerProvider: playerProvider,
		notifier:       notifier,
	}
}

func (w *Welcome) Welcome(ctx context.Context, chatID msginfo.ChatID) error {
	plr, err := w.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	if err := w.notifier.WelcomeNewPlayer(ctx, plr); err != nil {
		return fmt.Errorf("notification: %w", err)
	}

	return nil
}
