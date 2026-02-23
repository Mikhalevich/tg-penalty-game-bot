package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (t *TGHandler) DefaultHandler(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	var (
		chatID = msginfo.ChatIDFromInt64(msg.ChatID)
		msgID  = msginfo.MessageIDFromInt(msg.MessageID)
	)

	currentPlayer, err := t.playerController.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player info: %w", err)
	}

	if currentPlayer.IsChangeNameTriggered {
		if err := t.changeDisplayName(ctx, chatID, msgID, msg.Text); err != nil {
			return fmt.Errorf("change display name: %w", err)
		}

		return nil
	}

	return nil
}

func (t *TGHandler) changeDisplayName(
	ctx context.Context,
	chatID msginfo.ChatID,
	msgID msginfo.MessageID,
	displayName string,
) error {
	if err := t.playerController.ChangeDisplayName(ctx, chatID, msgID, displayName); err != nil {
		if !perror.IsType(err, perror.TypeAlreadyExists) {
			return fmt.Errorf("change display name: %w", err)
		}

		return nil
	}

	return nil
}
