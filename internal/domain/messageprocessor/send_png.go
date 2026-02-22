package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (m *MessageProcessor) SendPNG(
	ctx context.Context,
	chatID msginfo.ChatID,
	caption string,
	png []byte,
	rows ...button.ButtonRow,
) error {
	if err := m.SendMessage(ctx, Message{
		ChatID:  chatID,
		Text:    caption,
		Type:    MessageTypePNG,
		Payload: png,
		Buttons: rows,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
