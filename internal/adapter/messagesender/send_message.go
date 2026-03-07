package messagesender

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (m *messageSender) SendMessage(
	ctx context.Context,
	msg msginfo.SenderMessage,
) error {
	switch msg.Type {
	case msginfo.MessageTypePlain, msginfo.MessageTypeMarkdown:
		if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:          msg.ChatID.Int64(),
			Text:            msg.Text,
			ParseMode:       parseMode(msg.Type),
			ReplyParameters: replyParameters(msg.ReplyMsgID),
			ReplyMarkup:     makeButtonsMarkup(msg.Buttons...),
		}); err != nil {
			return fmt.Errorf("send text message: %w", err)
		}

	case msginfo.MessageTypePNG:
		if err := m.SendPNGMarkdown(ctx, msg.ChatID, msg.Text, msg.Payload, msg.Buttons...); err != nil {
			return fmt.Errorf("send png: %w", err)
		}

	case msginfo.MessageTypeShotImage:
		return fmt.Errorf("invalid message type: %v", msg.Type)

	default:
		return fmt.Errorf("invalid message type: %v", msg.Type)
	}

	return nil
}

func parseMode(mt msginfo.MessageType) models.ParseMode {
	if mt == msginfo.MessageTypeMarkdown {
		return models.ParseModeMarkdown
	}

	return ""
}

func replyParameters(replyMsgID msginfo.MessageID) *models.ReplyParameters {
	if replyMsgID.Int() == 0 {
		return nil
	}

	return &models.ReplyParameters{
		MessageID: replyMsgID.Int(),
	}
}
