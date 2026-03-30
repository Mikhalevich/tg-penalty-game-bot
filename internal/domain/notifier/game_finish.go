package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) GameFinish(ctx context.Context, state game.State) error {
	msg := n.makeScoreMsg("Game Finished", state)

	for _, plr := range state.LivePlayers() {
		if err := n.sendMsgToPlayer(ctx, plr, msg); err != nil {
			return fmt.Errorf("send msg to player: %w", err)
		}
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

func (n *Notifier) makeScoreMsg(header string, state game.State) string {
	return fmt.Sprintf("*%s*\n %s \\: %s",
		header,
		playerGoalsScoredMsg(state.Player1.ID, state.Rounds),
		playerGoalsScoredMsg(state.Player2.ID, state.Rounds),
	)
}

func playerGoalsScoredMsg(playerID player.ID, rounds []game.Round) string {
	var builder strings.Builder
	for _, round := range rounds {
		if round.Attack.PlayerID != playerID {
			continue
		}

		if !round.IsCompleted() {
			break
		}

		switch round.Result {
		case game.RoundResultNotCompleted:
			// impossible

		case game.RoundResultGoal:
			builder.WriteString(ballSymbol)

		case game.RoundResultSave:
			builder.WriteString(gloveSymbol)

		case game.RoundResultMiss:
			builder.WriteString(missSymbol)
		}
	}

	return builder.String()
}
