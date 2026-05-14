package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

const (
	botDifficultyButtonsCount = 3
)

var (
	//nolint:gochecknoglobals
	botDifficultyButtonsTemplate = [...]botDifficultyButtonsHeader{
		{
			Caption:    "Easy",
			Difficulty: game.BotDifficultyEasy,
		},
		{
			Caption:    "Normal",
			Difficulty: game.BotDifficultyNormal,
		},
		{
			Caption:    "Hard",
			Difficulty: game.BotDifficultyHard,
		},
		{
			Caption:    "Insane",
			Difficulty: game.BotDifficultyInsane,
		},
	}
)

func (n *Notifier) SelectBotDifficulty(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	buttons, err := makeBotDifficultyButtons()
	if err != nil {
		return fmt.Errorf("create button by full name: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:  chatID,
		Text:    "Select bot difficulty",
		Type:    msginfo.MessageTypePlain,
		Buttons: buttons,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

type botDifficultyButtonsHeader struct {
	Caption    string
	Difficulty game.BotDifficulty
}

func makeBotDifficultyButtons() ([]button.ButtonRow, error) {
	buttons := make([]button.ButtonRow, 0, botDifficultyButtonsCount)

	for _, tmpl := range botDifficultyButtonsTemplate {
		btn, err := game.SelectBotDifficultyButton(tmpl.Caption, tmpl.Difficulty)
		if err != nil {
			return nil, fmt.Errorf("create bot difficulty button: %w", err)
		}

		buttons = append(buttons, button.Row(btn))
	}

	return buttons, nil
}
