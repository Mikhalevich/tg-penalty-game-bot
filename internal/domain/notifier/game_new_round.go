package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

const (
	ballSymbol       = "⚽"
	gloveSymbol      = "🧤"
	leftSideSymbol   = "👈"
	rightSideSymbol  = "👉"
	middleSideSymbol = "🖐"
)

func (n *Notifier) GameNewRound(ctx context.Context, gameID game.ID, state game.State) error {
	var (
		cRound      = state.Rounds.Last()
		roundNumber = state.Rounds.Len()
	)

	attackerPlayer := state.PlayerByID(cRound.Attack.PlayerID)
	if !attackerPlayer.IsBot() {
		attackerButtons, err := makeShotSideButtons(gameID, roundNumber)
		if err != nil {
			return fmt.Errorf("make buttons: %w", err)
		}

		payload, err := shotimage.ShotImage{
			Type: shotimage.ImageTypeAttackerPrepare,
		}.GOBEncode()

		if err != nil {
			return fmt.Errorf("make shot attacker payload: %w", err)
		}

		if err := n.sender.SendMessage(ctx, msginfo.Message{
			ChatID:  attackerPlayer.ChatID,
			Text:    "Attack",
			Type:    msginfo.MessageTypeShotImage,
			Payload: payload,
			Buttons: []button.ButtonRow{attackerButtons},
		}); err != nil {
			return fmt.Errorf("send attacker msg: %w", err)
		}
	}

	defenderPlayer := state.PlayerByID(cRound.Defend.PlayerID)

	if !defenderPlayer.IsBot() {
		defenderButtons, err := makeShotSideButtons(gameID, roundNumber)
		if err != nil {
			return fmt.Errorf("make buttons: %w", err)
		}

		payload, err := shotimage.ShotImage{
			Type: shotimage.ImageTypeDefenderPrepare,
		}.GOBEncode()

		if err != nil {
			return fmt.Errorf("make shot defender payload: %w", err)
		}

		if err := n.sender.SendMessage(ctx, msginfo.Message{
			ChatID:  defenderPlayer.ChatID,
			Text:    "Defend",
			Type:    msginfo.MessageTypeShotImage,
			Payload: payload,
			Buttons: []button.ButtonRow{defenderButtons},
		}); err != nil {
			return fmt.Errorf("send defender msg: %w", err)
		}
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
