package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

const (
	ballSymbol       = "⚽"
	gloveSymbol      = "🧤"
	missSymbol       = "❌"
	leftSideSymbol   = "👈"
	rightSideSymbol  = "👉"
	middleSideSymbol = "👆"
)

const (
	newRoundDelay = time.Second * 2
)

func (n *Notifier) GameNewRound(
	ctx context.Context,
	gameID game.ID,
	state game.State,
) error {
	var (
		cRound      = state.CurrentRound()
		roundNumber = state.CurrentRoundNumber()
	)

	if err := n.sendShotImageForNewRound(
		ctx,
		state.PlayerByID(cRound.Attack.PlayerID),
		gameID,
		roundNumber,
		"Attack",
		shotimage.ImageTypePrepareAttack,
		cRound.Attack.CreatedAt.Add(newRoundDelay),
	); err != nil {
		return fmt.Errorf("send attacker shot image: %w", err)
	}

	if err := n.sendShotImageForNewRound(
		ctx,
		state.PlayerByID(cRound.Defend.PlayerID),
		gameID,
		roundNumber,
		"Defend",
		shotimage.ImageTypePrepareDefend,
		cRound.Defend.CreatedAt.Add(newRoundDelay),
	); err != nil {
		return fmt.Errorf("send defender shot image: %w", err)
	}

	return nil
}

func (n *Notifier) sendShotImageForNewRound(
	ctx context.Context,
	plr game.Player,
	gameID game.ID,
	roundNumber int,
	caption string,
	imageType shotimage.ImageType,
	visibilityAt time.Time,
) error {
	if plr.IsBot() {
		return nil
	}

	shotButtons, err := makeShotSideButtons(gameID, roundNumber)
	if err != nil {
		return fmt.Errorf("make buttons: %w", err)
	}

	payload, err := shotimage.ShotImage{
		Type: imageType,
	}.GOBEncode()

	if err != nil {
		return fmt.Errorf("make shot attacker payload: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:       plr.ChatID,
		Text:         caption,
		Type:         msginfo.MessageTypeShotImage,
		Payload:      payload,
		Buttons:      []button.ButtonRow{shotButtons},
		VisibilityAt: visibilityAt,
	}); err != nil {
		return fmt.Errorf("send shot image: %w", err)
	}

	return nil
}

func makeShotSideButtons(gameID game.ID, roundNumber int) (button.ButtonRow, error) {
	left, err := game.ShotSideButton(leftSideSymbol, gameID, roundNumber, game.ShotSideLeft)
	if err != nil {
		return nil, fmt.Errorf("left button: %w", err)
	}

	middle, err := game.ShotSideButton(middleSideSymbol, gameID, roundNumber, game.ShotSideMiddle)
	if err != nil {
		return nil, fmt.Errorf("right button: %w", err)
	}

	right, err := game.ShotSideButton(rightSideSymbol, gameID, roundNumber, game.ShotSideRight)
	if err != nil {
		return nil, fmt.Errorf("right button: %w", err)
	}

	return button.ButtonRow{left, middle, right}, nil
}
