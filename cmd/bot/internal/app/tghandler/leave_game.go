package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) LeaveGame(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.leaveGame.LeaveGame(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(0),
	); err != nil {
		return fmt.Errorf("leave game: %w", err)
	}

	return nil
}
