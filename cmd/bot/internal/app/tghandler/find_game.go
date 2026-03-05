package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) FindGame(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.findGame.FindGame(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("set player ready to game: %w", err)
	}

	return nil
}
