package notifier

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (n *Notifier) GameStage(ctx context.Context, currentGame game.Game) error {
	return nil
}
