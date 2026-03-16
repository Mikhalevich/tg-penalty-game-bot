package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) GameFinish(ctx context.Context, currentGame game.Game) error {
	if err := n.sendMsgToPlayers(
		ctx,
		currentGame,
		n.makeScoreMsg("Game Finished", currentGame),
	); err != nil {
		return fmt.Errorf("send msg to players: %w", err)
	}

	return nil
}

func (n *Notifier) sendMsgToPlayers(ctx context.Context, currentGame game.Game, msg string) error {
	if err := n.sendMsgToPlayer(ctx, currentGame.State.Player1, msg); err != nil {
		return fmt.Errorf("send msg to first player: %w", err)
	}

	if err := n.sendMsgToPlayer(ctx, currentGame.State.Player2, msg); err != nil {
		return fmt.Errorf("send msg to second player: %w", err)
	}

	return nil
}

func (n *Notifier) sendMsgToPlayer(ctx context.Context, plr game.Player, msg string) error {
	if plr.IsBot() {
		return nil
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: plr.ChatID,
		Text:   msg,
		Type:   msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (n *Notifier) makeScoreMsg(header string, currentGame game.Game) string {
	return fmt.Sprintf("*%s*\n %s \\: %s",
		header,
		playerGoalsScoredMsg(currentGame.State.Player1.ID, currentGame),
		playerGoalsScoredMsg(currentGame.State.Player2.ID, currentGame),
	)
}

func playerGoalsScoredMsg(playerID player.ID, currentGame game.Game) string {
	var builder strings.Builder
	for _, round := range currentGame.State.Rounds {
		if round.Attack.PlayerID != playerID {
			continue
		}

		if !round.IsCompleted {
			continue
		}

		if round.IsGoal {
			builder.WriteString(ballSymbol)
		} else {
			builder.WriteString(gloveSymbol)
		}
	}

	return builder.String()
}
