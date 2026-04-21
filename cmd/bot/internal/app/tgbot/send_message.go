package tgbot

import (
	"context"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

func (t *TGBot) SendMessage(ctx context.Context, chatID int64, msg string) {
	if _, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			Error("send message")
	}
}

func (t *TGBot) DeleteMessage(
	ctx context.Context,
	chatID int64,
	messageID int,
) {
	if _, err := t.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: messageID,
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			Error("delete message")
	}
}
