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
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

const (
	gameFinishedDelay = time.Second * 2
)

func (n *Notifier) GameFinish(
	ctx context.Context,
	gameType game.GameType,
	state game.State,
	finishedAt time.Time,
) error {
	var (
		msg          = n.makeFinishGameMsg(state)
		visibilityAt = finishedAt.Add(gameFinishedDelay)
	)

	if err := n.sendFinishGameMsg(
		ctx,
		state.Player1,
		makeResultImage(state.Player1.GoalsScored, state.Player2.GoalsScored),
		msg,
		gameType,
		visibilityAt,
	); err != nil {
		return fmt.Errorf("finish notification for player1: %w", err)
	}

	if err := n.sendFinishGameMsg(
		ctx,
		state.Player2,
		makeResultImage(state.Player2.GoalsScored, state.Player1.GoalsScored),
		msg,
		gameType,
		visibilityAt,
	); err != nil {
		return fmt.Errorf("finish notification for player1: %w", err)
	}

	return nil
}

func makeResultImage(score1, score2 int) shotimage.ShotImage {
	switch {
	case score1 > score2:
		return shotimage.ShotImage{
			Type: shotimage.ImageTypeWin,
		}

	case score1 < score2:
		return shotimage.ShotImage{
			Type: shotimage.ImageTypeLose,
		}
	}

	return shotimage.ShotImage{
		Type: shotimage.ImageTypeDraw,
	}
}

func (n *Notifier) sendFinishGameMsg(
	ctx context.Context,
	plr game.Player,
	resImage shotimage.ShotImage,
	caption string,
	gameType game.GameType,
	visibilityAt time.Time,
) error {
	if plr.IsBot() {
		return nil
	}

	payload, err := resImage.GOBEncode()
	if err != nil {
		return fmt.Errorf("gob enbode: %w", err)
	}

	buttons, err := makeFinishGameButtons(gameType)
	if err != nil {
		return fmt.Errorf("make buttons: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:       plr.ChatID,
		Text:         caption,
		Type:         msginfo.MessageTypeShotImage,
		Payload:      payload,
		Buttons:      buttons,
		VisibilityAt: visibilityAt,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeFinishGameButtons(gameType game.GameType) ([]button.ButtonRow, error) {
	var caption string
	switch gameType {
	case game.GameTypeRating:
		caption = "Find game"

	case game.GameTypeBot:
		caption = "Play with bot"

	case game.GameTypeByLink:
		return nil, nil

	default:
		return nil, fmt.Errorf("invalid game type: %s", gameType)
	}

	btn, err := game.StartGameButton(caption, gameType)
	if err != nil {
		return nil, fmt.Errorf("create repeat button: %w", err)
	}

	return []button.ButtonRow{button.Row(btn)}, nil
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
