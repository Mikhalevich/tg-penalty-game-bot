package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) ErrorMiddleware(next tgbot.Handler) tgbot.Handler {
	return tgbot.Handler(func(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
		if err := next(ctx, msg, sender); err != nil {
			if err := t.notifier.ParseError(
				ctx,
				msginfo.ChatIDFromInt64(msg.ChatID),
				err,
			); err != nil {
				return fmt.Errorf("parse error: %w", err)
			}
		}

		return nil
	})
}
