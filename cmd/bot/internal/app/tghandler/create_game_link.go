package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) CreateGameLink(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.startGame.StartGameByLink(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("start game by link: %w", err)
	}

	return nil
}
