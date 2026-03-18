package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

func (n *Notifier) GameRoundFinish(ctx context.Context, currentGame game.Game) error {
	var (
		cRound = currentGame.CurrentRound()
		msg    = n.makeScoreMsg(goalMsg(currentGame), currentGame)
	)

	for _, plr := range currentGame.LivePlayers() {
		if err := n.sendShotImage(ctx, plr, msg, cRound.Attack.Side, cRound.Defend.Side); err != nil {
			return fmt.Errorf("send shot image: %w", err)
		}
	}

	return nil
}

func goalMsg(currentGame game.Game) string {
	if currentGame.CurrentRound().IsGoal {
		return "Goal"
	}

	return "Save"
}

func (n *Notifier) sendShotImage(
	ctx context.Context,
	plr game.Player,
	msg string,
	attackerSide game.ShotSide,
	defenderSide game.ShotSide,
) error {
	if plr.IsBot() {
		return nil
	}

	payload, err := shotimage.ShotImage{
		Type:         shotimage.ImageTypeShot,
		AttackerSide: attackerSide,
		DefenderSide: defenderSide,
	}.GOBEncode()

	if err != nil {
		return fmt.Errorf("make shot image payload: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:  plr.ChatID,
		Text:    msg,
		Type:    msginfo.MessageTypeShotImage,
		Payload: payload,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
