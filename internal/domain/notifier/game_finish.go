package notifier

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	gameFinishedDelay = time.Second * 2
)

func (n *Notifier) GameFinish(ctx context.Context, state game.State, finishedAt time.Time) error {
	var (
		msg          = n.makeFinishGameMsg(state)
		visibilityAt = finishedAt.Add(gameFinishedDelay)
	)

	for _, plr := range state.LivePlayers() {
		if err := n.sendMsgToPlayerWithVisiblity(ctx, plr, msg, visibilityAt); err != nil {
			return fmt.Errorf("send msg to player: %w", err)
		}
	}

	return nil
}

func (n *Notifier) sendMsgToPlayer(
	ctx context.Context,
	plr game.Player,
	msg string,
	buttons ...button.ButtonRow,
) error {
	if plr.IsBot() {
		return nil
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:  plr.ChatID,
		Text:    msg,
		Type:    msginfo.MessageTypeMarkdown,
		Buttons: buttons,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (n *Notifier) sendMsgToPlayerWithVisiblity(
	ctx context.Context,
	plr game.Player,
	msg string,
	visibilityAt time.Time,
	buttons ...button.ButtonRow,
) error {
	if plr.IsBot() {
		return nil
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:       plr.ChatID,
		Text:         msg,
		Type:         msginfo.MessageTypeMarkdown,
		Buttons:      buttons,
		VisibilityAt: visibilityAt,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (n *Notifier) makeFinishGameMsg(state game.State) string {
	return fmt.Sprintf("*Game Finished*\n*%s* \\(*%d*\\) \\- *%s* \\(*%d*\\)\n%s",
		n.escaper.EscapeMarkdown(state.Player1.DisplayName),
		state.Player1.GoalsScored,
		n.escaper.EscapeMarkdown(state.Player2.DisplayName),
		state.Player2.GoalsScored,
		makeScoreMsg(state),
	)
}

func makeScoreMsg(state game.State) string {
	return fmt.Sprintf("%s \\: %s",
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
