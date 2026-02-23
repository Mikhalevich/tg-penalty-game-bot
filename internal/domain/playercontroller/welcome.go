package playercontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (p *PlayerController) Welcome(ctx context.Context, chatID msginfo.ChatID) error {
	plr, err := p.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	if err := p.notifier.WelcomeNewPlayer(ctx, plr); err != nil {
		return fmt.Errorf("notification: %w", err)
	}

	return nil
}
