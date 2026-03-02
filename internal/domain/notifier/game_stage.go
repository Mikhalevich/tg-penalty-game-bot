package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) GameStage(ctx context.Context, currentGame game.Game) error {
	if currentGame.IsFinished() {
		if err := n.finishedGameNotifications(ctx, currentGame); err != nil {
			return fmt.Errorf("finished game notifications: %w", err)
		}
	}

	if currentGame.CurrentRound().IsCompleted {
		if err := n.finishedRoundNotifications(ctx, currentGame); err != nil {
			return fmt.Errorf("finished round notifications: %w", err)
		}
	}

	if err := n.inGameNotifications(ctx, currentGame); err != nil {
		return fmt.Errorf("in game notifications: %w", err)
	}

	return nil
}

func (n *Notifier) sendMsgToPlayers(ctx context.Context, currentGame game.Game, msg string) error {
	if err := n.sendMsgToPlayer(ctx, currentGame.State.Player1.ChatID, msg); err != nil {
		return fmt.Errorf("send msg to first player: %w", err)
	}

	if err := n.sendMsgToPlayer(ctx, currentGame.State.Player2.ChatID, msg); err != nil {
		return fmt.Errorf("send msg to second player: %w", err)
	}

	return nil
}

func (n *Notifier) sendMsgToPlayer(ctx context.Context, chatID msginfo.ChatID, msg string) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text:   msg,
		Type:   msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (n *Notifier) finishedGameNotifications(ctx context.Context, currentGame game.Game) error {
	finishedGameMsg := fmt.Sprintf("Game Finished\n%s %d \\: %s %d",
		n.escaper.EscapeMarkdown(currentGame.State.Player1.DisplayName), currentGame.State.Player1.GoalsScored,
		n.escaper.EscapeMarkdown(currentGame.State.Player2.DisplayName), currentGame.State.Player2.GoalsScored,
	)

	if err := n.sendMsgToPlayers(ctx, currentGame, finishedGameMsg); err != nil {
		return fmt.Errorf("send msg to players: %w", err)
	}

	return nil
}

func (n *Notifier) finishedRoundNotifications(ctx context.Context, currentGame game.Game) error {
	if err := n.sendMsgToPlayers(ctx, currentGame, makeFinishedRoundMsg(currentGame)); err != nil {
		return fmt.Errorf("send msg to players: %w", err)
	}

	return nil
}

func makeFinishedRoundMsg(currentGame game.Game) string {
	if currentGame.CurrentRound().IsGoal {
		return "Goal"
	}

	return "Defend"
}

func (n *Notifier) inGameNotifications(ctx context.Context, currentGame game.Game) error {
	cRound := currentGame.CurrentRound()

	buttons, err := makeShotSideButtons(currentGame.ID, currentGame.CurrentRoundNumber())
	if err != nil {
		return fmt.Errorf("make buttons: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:  chatIDByPlayerID(cRound.Attack.PlayerID, currentGame),
		Text:    "Attack",
		Type:    msginfo.MessageTypeMarkdown,
		Buttons: []button.ButtonRow{buttons},
	}); err != nil {
		return fmt.Errorf("send attacker msg: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:  chatIDByPlayerID(cRound.Defend.PlayerID, currentGame),
		Text:    "Defend",
		Type:    msginfo.MessageTypeMarkdown,
		Buttons: []button.ButtonRow{buttons},
	}); err != nil {
		return fmt.Errorf("send defender msg: %w", err)
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

func chatIDByPlayerID(playerID player.ID, currentGame game.Game) msginfo.ChatID {
	if currentGame.State.Player1.ID == playerID {
		return currentGame.State.Player1.ChatID
	}

	return currentGame.State.Player2.ChatID
}
