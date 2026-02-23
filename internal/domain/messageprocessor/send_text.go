package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (m *MessageProcessor) SendTextPlain(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	rows ...button.ButtonRow,
) error {
	if err := m.SendMessage(ctx, msginfo.Message{
		ChatID:  chatID,
		Text:    text,
		Type:    msginfo.MessageTypePlain,
		Buttons: rows,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (m *MessageProcessor) SendTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	rows ...button.ButtonRow,
) error {
	if err := m.SendMessage(ctx, msginfo.Message{
		ChatID:  chatID,
		Text:    text,
		Type:    msginfo.MessageTypeMarkdown,
		Buttons: rows,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
