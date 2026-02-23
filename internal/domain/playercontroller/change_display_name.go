package playercontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (p *PlayerController) ChangeDisplayName(
	ctx context.Context,
	chatID msginfo.ChatID,
	displayName string,
	msgID msginfo.MessageID,
) error {
	currentPlayer, err := p.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player info: %w", err)
	}

	if !currentPlayer.IsChangeNameTriggered {
		return fmt.Errorf("change name not triggered: %w", err)
	}

	if err := p.repo.ChangeDisplayName(ctx, chatID, displayName, p.timeProvider.Now()); err != nil {
		if p.repo.IsAlreadyExistsError(err) {
			if err := p.notifier.NameAlreadyRegistered(ctx, currentPlayer, msgID); err != nil {
				return fmt.Errorf("name already redistered notification: %w", err)
			}
		}

		return fmt.Errorf("repo change display name: %w", err)
	}

	currentPlayer.DisplayName = displayName

	if err := p.notifier.NameChanged(ctx, currentPlayer, msgID); err != nil {
		return fmt.Errorf("name changed notification: %w", err)
	}

	return nil
}
