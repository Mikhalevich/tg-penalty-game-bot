package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (n *Notifier) GameStart(
	ctx context.Context,
	gameType game.GameType,
	player1, player2 game.Player,
) error {
	if err := n.sendStartGameNotification(ctx, gameType, player1, n.startGameAgainstMsg(player2)); err != nil {
		return fmt.Errorf("send start game to first player: %w", err)
	}

	if err := n.sendStartGameNotification(ctx, gameType, player2, n.startGameAgainstMsg(player1)); err != nil {
		return fmt.Errorf("send start game to second player: %w", err)
	}

	return nil
}

func (n *Notifier) sendStartGameNotification(
	ctx context.Context,
	gameType game.GameType,
	plr game.Player,
	msg string,
) error {
	if plr.IsBot() {
		return nil
	}

	if err := n.sendMsgToPlayer(
		ctx,
		plr,
		msg,
		buttonsByGameType(gameType)...,
	); err != nil {
		return fmt.Errorf("send msg to player: %w", err)
	}

	return nil
}

func buttonsByGameType(gameType game.GameType) []button.ButtonRow {
	if gameType == game.GameTypeRating {
		return []button.ButtonRow{
			{
				game.LeaveGameButton("Leave"),
			},
			{
				game.ShotStatsButton("View Statistics"),
			},
		}
	}

	return []button.ButtonRow{
		{
			game.LeaveGameButton("Leave"),
		},
	}
}

func (n *Notifier) startGameAgainstMsg(plr game.Player) string {
	return fmt.Sprintf("Game started against *%s*", n.escaper.EscapeMarkdown(plr.DisplayName))
}
