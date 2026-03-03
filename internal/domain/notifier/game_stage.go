package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	ballSymbol  = "⚽"
	gloveSymbol = "🧤"
)

func (n *Notifier) GameStage(ctx context.Context, currentGame game.Game) error {
	if currentGame.IsFinished() {
		if err := n.finishedGameNotifications(ctx, currentGame); err != nil {
			return fmt.Errorf("finished game notifications: %w", err)
		}

		return nil
	}

	if currentGame.CurrentRound().IsCompleted {
		if err := n.finishedRoundNotifications(ctx, currentGame); err != nil {
			return fmt.Errorf("finished round notifications: %w", err)
		}

		return nil
	}

	if err := n.inGameNotifications(ctx, currentGame); err != nil {
		return fmt.Errorf("in game notifications: %w", err)
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

func (n *Notifier) finishedGameNotifications(ctx context.Context, currentGame game.Game) error {
	if err := n.sendMsgToPlayers(
		ctx,
		currentGame,
		n.makeScoreMsg("Game Finished", currentGame),
	); err != nil {
		return fmt.Errorf("send msg to players: %w", err)
	}

	return nil
}

func (n *Notifier) finishedRoundNotifications(ctx context.Context, currentGame game.Game) error {
	if err := n.sendMsgToPlayers(
		ctx,
		currentGame,
		n.makeScoreMsg(goalMsg(currentGame), currentGame),
	); err != nil {
		return fmt.Errorf("send msg to players: %w", err)
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

func goalMsg(currentGame game.Game) string {
	if currentGame.CurrentRound().IsGoal {
		return "Goal"
	}

	return "Save"
}

func (n *Notifier) inGameNotifications(ctx context.Context, currentGame game.Game) error {
	cRound := currentGame.CurrentRound()

	attackerPlayer := playerByPlayerID(cRound.Attack.PlayerID, currentGame)
	if !attackerPlayer.IsBot() {
		attackerButtons, err := makeShotSideButtons(currentGame.ID, currentGame.CurrentRoundNumber())
		if err != nil {
			return fmt.Errorf("make buttons: %w", err)
		}

		if err := n.sender.SendMessage(ctx, msginfo.Message{
			ChatID:  attackerPlayer.ChatID,
			Text:    "Attack",
			Type:    msginfo.MessageTypeMarkdown,
			Buttons: []button.ButtonRow{attackerButtons},
		}); err != nil {
			return fmt.Errorf("send attacker msg: %w", err)
		}
	}

	defenderPlayer := playerByPlayerID(cRound.Defend.PlayerID, currentGame)

	if !defenderPlayer.IsBot() {
		defenderButtons, err := makeShotSideButtons(currentGame.ID, currentGame.CurrentRoundNumber())
		if err != nil {
			return fmt.Errorf("make buttons: %w", err)
		}

		if err := n.sender.SendMessage(ctx, msginfo.Message{
			ChatID:  defenderPlayer.ChatID,
			Text:    "Defend",
			Type:    msginfo.MessageTypeMarkdown,
			Buttons: []button.ButtonRow{defenderButtons},
		}); err != nil {
			return fmt.Errorf("send defender msg: %w", err)
		}
	}

	return nil
}

func makeShotSideButtons(gameID game.ID, roundNumber int) (button.ButtonRow, error) {
	left, err := game.ShotSideButton("left", gameID, roundNumber, game.ShotSideLeft)
	if err != nil {
		return nil, fmt.Errorf("left button: %w", err)
	}

	middle, err := game.ShotSideButton("middle", gameID, roundNumber, game.ShotSideMiddle)
	if err != nil {
		return nil, fmt.Errorf("right button: %w", err)
	}

	right, err := game.ShotSideButton("right", gameID, roundNumber, game.ShotSideRight)
	if err != nil {
		return nil, fmt.Errorf("right button: %w", err)
	}

	return button.ButtonRow{left, middle, right}, nil
}

func playerByPlayerID(playerID player.ID, currentGame game.Game) game.Player {
	if currentGame.State.Player1.ID == playerID {
		return currentGame.State.Player1
	}

	return currentGame.State.Player2
}
