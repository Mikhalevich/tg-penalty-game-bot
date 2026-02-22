package playercontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (p *PlayerController) SetChangeDisplayNameTrigger(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	plr, err := p.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	var (
		lastNameChangedInterval = time.Since(plr.NameChangedAt)
	)

	if lastNameChangedInterval < p.changeNameTimeout {
		return perror.TooManyRequests(
			"to many requests to change display name",
			formatDelta(p.changeNameTimeout-lastNameChangedInterval),
		)
	}

	if err := p.repo.SetChangeDisplayNameTrigger(ctx, plr.ChatID); err != nil {
		return fmt.Errorf("repo set change display name trigger: %w", err)
	}

	return nil
}

func formatDelta(delta time.Duration) time.Duration {
	return delta.Truncate(time.Second)
}
