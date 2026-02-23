package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (m *MessageProcessor) SendMessage(
	ctx context.Context,
	msg msginfo.Message,
) error {
	inlineButtons, err := m.SetButtonRows(ctx, msg.Buttons...)
	if err != nil {
		return fmt.Errorf("set button rows: %w", err)
	}

	if err := m.sender.SendMessage(ctx, msginfo.SenderMessage{
		ChatID:     msg.ChatID,
		ReplyMsgID: msg.ReplyMsgID,
		Text:       msg.Text,
		Type:       msg.Type,
		Payload:    msg.Payload,
		Buttons:    inlineButtons,
	}); err != nil {
		return fmt.Errorf("sender send message: %w", err)
	}

	return nil
}
