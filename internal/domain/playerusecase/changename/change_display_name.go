package changename

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (c *ChangeName) ChangeDisplayName(
	ctx context.Context,
	chatID msginfo.ChatID,
	msgID msginfo.MessageID,
	displayName string,
) error {
	if utf8.RuneCountInString(displayName) > c.maxNameLen {
		if err := c.notifier.NameIsTooLong(ctx, chatID, c.maxNameLen); err != nil {
			return fmt.Errorf("name is too long: %w", err)
		}

		return nil
	}

	currentPlayer, err := c.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player info: %w", err)
	}

	if !currentPlayer.IsChangeNameTriggered {
		return fmt.Errorf("change name not triggered: %w", err)
	}

	if err := c.repo.ChangeDisplayName(ctx, chatID, displayName, c.timeProvider.Now()); err != nil {
		if c.repo.IsAlreadyExistsError(err) {
			if err := c.notifier.NameAlreadyRegistered(ctx, currentPlayer, msgID); err != nil {
				return fmt.Errorf("name already redistered notification: %w", err)
			}
		}

		return fmt.Errorf("repo change display name: %w", err)
	}

	currentPlayer.DisplayName = displayName

	if err := c.notifier.NameChanged(ctx, currentPlayer, msgID); err != nil {
		return fmt.Errorf("name changed notification: %w", err)
	}

	return nil
}
