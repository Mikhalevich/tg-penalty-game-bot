package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/messageprocessor/button"
)

func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if msg.Data == "" {
		return nil
	}

	btn, err := t.messageProcessor.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		return fmt.Errorf("get button: %w", err)
	}

	if btn.Operation == button.OperationChangeName {
		return nil
	}

	return nil
}
