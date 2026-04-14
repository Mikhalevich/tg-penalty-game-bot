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
	if err := n.sendStartGameNotification(
		ctx,
		gameType,
		player1,
		player2,
	); err != nil {
		return fmt.Errorf("send start game to first player: %w", err)
	}

	if err := n.sendStartGameNotification(
		ctx,
		gameType,
		player2,
		player1,
	); err != nil {
		return fmt.Errorf("send start game to second player: %w", err)
	}

	return nil
}

func (n *Notifier) sendStartGameNotification(
	ctx context.Context,
	gameType game.GameType,
	plr game.Player,
	plrAgainst game.Player,
) error {
	if plr.IsBot() {
		return nil
	}

	msg := n.startGameAgainstMsg(plrAgainst.DisplayName)

	buttons, err := buttonsByGameType(gameType, plrAgainst)
	if err != nil {
		return fmt.Errorf("create buttons: %w", err)
	}

	if err := n.sendMsgToPlayer(
		ctx,
		plr,
		msg,
		buttons...,
	); err != nil {
		return fmt.Errorf("send msg to player: %w", err)
	}

	return nil
}

func buttonsByGameType(gameType game.GameType, plr game.Player) ([]button.ButtonRow, error) {
	if gameType != game.GameTypeRating {
		return []button.ButtonRow{
			{
				game.LeaveGameButton("Leave"),
			},
		}, nil
	}

	statsBtn, err := game.ShotStatsButton("View Statistics", plr.ID, plr.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("create view statistics button: %w", err)
	}

	return []button.ButtonRow{
		{
			game.LeaveGameButton("Leave"),
		},
		{
			statsBtn,
		},
	}, nil
}

func (n *Notifier) startGameAgainstMsg(displayName string) string {
	return fmt.Sprintf("Game started against *%s*", n.escaper.EscapeMarkdown(displayName))
}
