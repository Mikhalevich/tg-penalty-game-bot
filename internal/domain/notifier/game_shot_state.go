package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

func (n *Notifier) GameShotState(
	ctx context.Context,
	chatID msginfo.ChatID,
	gameID game.ID,
	shotType game.ShotType,
	roundNumber int,
) error {
	buttons, err := makeShotSideButtons(gameID, roundNumber)
	if err != nil {
		return fmt.Errorf("make buttons: %w", err)
	}

	imageType, msgText := imageTypeAndTextByShotType(shotType)

	payload, err := shotimage.ShotImage{
		Type: imageType,
	}.GOBEncode()

	if err != nil {
		return fmt.Errorf("make shot image payload: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:  chatID,
		Text:    msgText,
		Type:    msginfo.MessageTypeShotImage,
		Payload: payload,
		Buttons: []button.ButtonRow{buttons},
	}); err != nil {
		return fmt.Errorf("send msg: %w", err)
	}

	return nil
}

func imageTypeAndTextByShotType(shotType game.ShotType) (shotimage.ImageType, string) {
	switch shotType {
	case game.ShotTypeAttack:
		return shotimage.ImageTypePrepareAttack, "Attack"

	case game.ShotTypeDefend:
		return shotimage.ImageTypePrepareDefend, "Defend"
	}

	return shotimage.ImageType(""), ""
}
