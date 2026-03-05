package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) Start(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.welcome.Welcome(ctx, msginfo.ChatIDFromInt64(msg.ChatID)); err != nil {
		return fmt.Errorf("welcome: %w", err)
	}

	return nil
}
