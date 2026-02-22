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
			Error("send message error")
	}
}
