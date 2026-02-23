package playercontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (p *PlayerController) SetChangeDisplayNameTrigger(
	ctx context.Context,
	chatID msginfo.ChatID,
	fullName string,
	userName string,
) error {
	plr, err := p.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	var (
		lastNameChangedInterval = time.Since(plr.NameChangedAt)
	)

	if lastNameChangedInterval < p.changeNameTimeout {
		if err := p.notifier.ChangeNameDelay(ctx, plr, formatDelta(p.changeNameTimeout-lastNameChangedInterval)); err != nil {
			return fmt.Errorf("change name delay: %w", err)
		}

		return nil
	}

	if err := p.repo.SetChangeDisplayNameTrigger(ctx, plr.ChatID); err != nil {
		return fmt.Errorf("repo set change display name trigger: %w", err)
	}

	if err := p.notifier.ChangeName(ctx, plr, fullName, userName); err != nil {
		return fmt.Errorf("change name notificatoin: %w", err)
	}

	return nil
}

func formatDelta(delta time.Duration) time.Duration {
	return delta.Truncate(time.Second)
}
