package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (n *Notifier) GameCanceled(ctx context.Context, players []game.Player) error {
	for _, plr := range players {
		if err := n.sendMsgToPlayer(ctx, plr, "Game canceled"); err != nil {
			return fmt.Errorf("send msg to player: %w", err)
		}
	}

	return nil
}
