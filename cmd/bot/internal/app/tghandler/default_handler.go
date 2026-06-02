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
	)

	currentPlayer, err := t.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player info: %w", err)
	}

	if currentPlayer.IsChangeNameTriggered {
		if err := t.changeDisplayName(ctx, chatID, msg.Text); err != nil {
			return fmt.Errorf("change display name: %w", err)
		}

		return nil
	}

	return nil
}

func (t *TGHandler) changeDisplayName(
	ctx context.Context,
	chatID msginfo.ChatID,
	displayName string,
) error {
	if err := t.changeName.ChangeDisplayName(ctx, chatID, displayName, false); err != nil {
		if !perror.IsType(err, perror.TypeAlreadyExists) {
			return fmt.Errorf("change display name: %w", err)
		}

		return nil
	}

	return nil
}
