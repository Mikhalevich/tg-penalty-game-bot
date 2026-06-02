package changename

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (c *ChangeName) ChangeNameFromGroup(
	ctx context.Context,
	chatID msginfo.ChatID,
	fullName string,
	userName string,
) error {
	plr, err := c.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	var (
		lastNameChangedInterval = time.Since(plr.NameChangedAt)
	)

	if lastNameChangedInterval < c.changeNameTimeout {
		if err := c.notifier.ChangeNameDelay(ctx, plr, formatDelta(c.changeNameTimeout-lastNameChangedInterval)); err != nil {
			return fmt.Errorf("change name delay: %w", err)
		}

		return nil
	}

	if err := c.notifier.ChangeName(ctx, plr, true, fullName, userName); err != nil {
		return fmt.Errorf("change name notificatoin: %w", err)
	}

	return nil
}
