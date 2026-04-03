package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) Leaderboard(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.leaderboard.MyPosition(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("my position: %w", err)
	}

	return nil
}
