package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if msg.Data == "" {
		return nil
	}

	var (
		chatID = msginfo.ChatIDFromInt64(msg.ChatID)
		msgID  = msginfo.MessageIDFromInt(msg.MessageID)
	)

	btn, err := t.messageProcessor.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		return fmt.Errorf("get button: %w", err)
	}

	if btn.Operation == button.OperationChangeName {
		return t.processChangeNameButton(ctx, chatID, msgID, btn)
	}

	return nil
}

func (t *TGHandler) processChangeNameButton(
	ctx context.Context,
	chatID msginfo.ChatID,
	msgID msginfo.MessageID,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[button.ChangeNamePayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload: %w", err)
	}

	if err := t.changeDisplayName(ctx, chatID, msgID, payload.DisplayName); err != nil {
		return fmt.Errorf("change display name: %w", err)
	}

	return nil
}
