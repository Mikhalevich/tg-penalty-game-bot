package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

func (n *Notifier) GameRoundFinish(ctx context.Context, state game.State) error {
	var (
		cRound = state.CurrentRound()
		msg    = n.makeScoreMsg(goalMsg(cRound), state)
	)

	for _, plr := range state.LivePlayers() {
		if err := n.sendShotImageForFinishRound(
			ctx,
			plr,
			msg,
			shotImageTypeByPlayerID(plr.ID, cRound.Attack.PlayerID),
			cRound.Attack.Side,
			cRound.Defend.Side,
		); err != nil {
			return fmt.Errorf("send shot image: %w", err)
		}
	}

	return nil
}

func goalMsg(round game.Round) string {
	switch round.Result {
	case game.RoundResultGoal:
		return "Goal"

	case game.RoundResultSave:
		return "Save"

	case game.RoundResultMiss:
		return "Miss"

	case game.RoundResultNotCompleted:
	}

	return ""
}

func shotImageTypeByPlayerID(playerID, attackerID player.ID) shotimage.ImageType {
	if playerID == attackerID {
		return shotimage.ImageTypeAttack
	}

	return shotimage.ImageTypeDefend
}

func (n *Notifier) sendShotImageForFinishRound(
	ctx context.Context,
	plr game.Player,
	msg string,
	imageType shotimage.ImageType,
	attackerSide game.ShotSide,
	defenderSide game.ShotSide,
) error {
	if plr.IsBot() {
		return nil
	}

	payload, err := shotimage.ShotImage{
		Type:         imageType,
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
