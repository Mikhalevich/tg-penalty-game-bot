package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) ChangeName(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if msg.Args != "" {
		if err := t.changeName.ChangeDisplayName(ctx, msginfo.ChatID(msg.ChatID), msg.Args); err != nil {
			return fmt.Errorf("change name to %q: %w", msg.Args, err)
		}

		return nil
	}

	if err := t.changeName.SetChangeDisplayNameTrigger(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msg.User.FullName(),
		msg.User.Username,
	); err != nil {
		return fmt.Errorf("set change name trigger: %w", err)
	}

	return nil
}
