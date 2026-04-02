package game

import (
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
)

type ShotSidePayload struct {
	GameID ID
	Round  int
	Side   ShotSide
}

func ShotSideButton(caption string, gameID ID, round int, side ShotSide) (button.Button, error) {
	btn, err := button.CreateButton(caption, button.OperationShotSide,
		ShotSidePayload{
			GameID: gameID,
			Round:  round,
			Side:   side,
		},
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}

func LeaveGameButton(caption string) button.Button {
	return button.CreateButtonWithoutPayload(caption, button.OperationLeaveGame)
}

func StopSearchGame(caption string) button.Button {
	return button.CreateButtonWithoutPayload(caption, button.OperationStopSearchGame)
}
