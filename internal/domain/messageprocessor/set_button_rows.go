package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
)

func (m *MessageProcessor) SetButtonRows(
	ctx context.Context,
	rows ...button.ButtonRow,
) ([]button.InlineKeyboardButtonRow, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	inlineButtonRows := make([]button.InlineKeyboardButtonRow, 0, len(rows))

	if err := m.buttonRepository.SetButtonRows(ctx, rows...); err != nil {
		return nil, fmt.Errorf("set button rows: %w", err)
	}

	for _, row := range rows {
		buttonRow := make([]button.InlineKeyboardButton, 0, len(row))

		for _, btn := range row {
			buttonRow = append(buttonRow, button.InlineKeyboardButton{
				ID:      btn.ID,
				Caption: btn.Caption,
			})
		}

		inlineButtonRows = append(inlineButtonRows, buttonRow)
	}

	return inlineButtonRows, nil
}
