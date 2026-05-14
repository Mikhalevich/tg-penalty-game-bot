package game

import (
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type ShotSidePayload struct {
	GameID ID
	Round  int
	Side   ShotSide
}

func ShotSideButton(caption string, gameID ID, round int, side ShotSide) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationShotSide,
		true,
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
	return button.CreateButtonWithoutPayload(caption, button.OperationLeaveGame, true)
}

type ShotStatsPayload struct {
	PlayerID    player.ID
	DisplayName string
}

func ShotStatsButton(caption string, playerID player.ID, playerDisplayName string) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationShotStats,
		false,
		ShotStatsPayload{
			PlayerID:    playerID,
			DisplayName: playerDisplayName,
		},
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}

func StopSearchGame(caption string) button.Button {
	return button.CreateButtonWithoutPayload(caption, button.OperationStopSearchGame, true)
}

type RepeatGamePayload struct {
	GameType GameType
}

func StartGameButton(caption string, gameType GameType) (button.Button, error) {
	btn, err := button.CreateButton(caption, button.OperationStartGame, false, RepeatGamePayload{
		GameType: gameType,
	})

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}

func StartGameWithDeleteAfterPressButton(caption string, gameType GameType) (button.Button, error) {
	btn, err := button.CreateButton(caption, button.OperationStartGame, true, RepeatGamePayload{
		GameType: gameType,
	})

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}

type SelectBotDifficultyPayload struct {
	Difficulty BotDifficulty
}

func SelectBotDifficultyButton(caption string, difficulty BotDifficulty) (button.Button, error) {
	btn, err := button.CreateButton(caption, button.OperatoinSelectBotDifficulty, true, SelectBotDifficultyPayload{
		Difficulty: difficulty,
	})

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}
