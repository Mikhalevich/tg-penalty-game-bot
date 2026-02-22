package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/messageprocessor/button"
)

func (m *MessageProcessor) SetButton(
	ctx context.Context,
	btn button.Button,
) (button.InlineKeyboardButton, error) {
	if err := m.buttonRepository.SetButton(ctx, btn); err != nil {
		return button.InlineKeyboardButton{}, fmt.Errorf("button repository set button: %w", err)
	}

	return button.InlineKeyboardButton{
		ID:      btn.ID,
		Caption: btn.Caption,
	}, nil
}
