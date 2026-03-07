package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

func (m *MessageProcessor) SendMessage(
	ctx context.Context,
	msg msginfo.Message,
) error {
	inlineButtons, err := m.SetButtonRows(ctx, msg.Buttons...)
	if err != nil {
		return fmt.Errorf("set button rows: %w", err)
	}

	if msg.Type == msginfo.MessageTypeShotImage {
		if err := m.processShotImage(ctx, msg, inlineButtons); err != nil {
			return fmt.Errorf("process shot image: %w", err)
		}

		return nil
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

func (m *MessageProcessor) processShotImage(
	ctx context.Context,
	msg msginfo.Message,
	inlineButtons []button.InlineKeyboardButtonRow,
) error {
	shot, err := shotimage.GOBDecode(msg.Payload)
	if err != nil {
		return fmt.Errorf("get shot image: %w", err)
	}

	image, err := m.shotImageProvider.Image(ctx, shot)
	if err != nil {
		logger.FromContext(ctx).
			WithError(err).
			Error("receiving shot image, fallback to text message")

		if err := m.sender.SendMessage(ctx, msginfo.SenderMessage{
			ChatID:     msg.ChatID,
			ReplyMsgID: msg.ReplyMsgID,
			Text:       msg.Text,
			Type:       msginfo.MessageTypeMarkdown,
			Buttons:    inlineButtons,
		}); err != nil {
			return fmt.Errorf("send fallback text message: %w", err)
		}

		return nil
	}

	if err := m.sender.SendMessage(ctx, msginfo.SenderMessage{
		ChatID:     msg.ChatID,
		ReplyMsgID: msg.ReplyMsgID,
		Text:       msg.Text,
		Type:       msginfo.MessageTypePNG,
		Payload:    image,
		Buttons:    inlineButtons,
	}); err != nil {
		return fmt.Errorf("send image: %w", err)
	}

	return nil
}
