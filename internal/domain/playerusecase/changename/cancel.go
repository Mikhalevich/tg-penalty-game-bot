package changename

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (c *ChangeName) Cancel(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	plr, err := c.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	if !plr.IsChangeNameTriggered {
		return nil
	}

	if err := c.repo.SetChangeDisplayNameTrigger(ctx, chatID, false); err != nil {
		return fmt.Errorf("repo down change display name trigger: %w", err)
	}

	return nil
}
