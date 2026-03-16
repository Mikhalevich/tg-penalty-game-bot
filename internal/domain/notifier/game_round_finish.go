package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (n *Notifier) GameRoundFinish(ctx context.Context, currentGame game.Game) error {
	if err := n.sendMsgToPlayers(
		ctx,
		currentGame,
		n.makeScoreMsg(goalMsg(currentGame), currentGame),
	); err != nil {
		return fmt.Errorf("send msg to players: %w", err)
	}

	return nil
}

func goalMsg(currentGame game.Game) string {
	if currentGame.CurrentRound().IsGoal {
		return "Goal"
	}

	return "Save"
}
