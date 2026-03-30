package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (n *Notifier) GameStart(ctx context.Context, player1, player2 game.Player) error {
	leaveBtn := button.ButtonRow{game.LeaveGameButton("Leave")}

	if err := n.sendMsgToPlayer(ctx, player1, n.startGameAgainstMsg(player2), leaveBtn); err != nil {
		return fmt.Errorf("send msg to first player: %w", err)
	}

	if err := n.sendMsgToPlayer(ctx, player2, n.startGameAgainstMsg(player1), leaveBtn); err != nil {
		return fmt.Errorf("send msg to second player: %w", err)
	}

	return nil
}

func (n *Notifier) startGameAgainstMsg(plr game.Player) string {
	return fmt.Sprintf("Game started against *%s*", n.escaper.EscapeMarkdown(plr.DisplayName))
}
