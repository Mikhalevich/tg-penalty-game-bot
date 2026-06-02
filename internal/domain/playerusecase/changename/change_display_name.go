package changename

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (c *ChangeName) ChangeDisplayName(
	ctx context.Context,
	chatID msginfo.ChatID,
	displayName string,
	force bool,
) error {
	currentPlayer, err := c.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player info: %w", err)
	}

	isProcessed, err := c.validateParams(ctx, currentPlayer, displayName, force)
	if err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if isProcessed {
		return nil
	}

	if err := c.repo.ChangeDisplayName(ctx, chatID, displayName, c.timeProvider.Now()); err != nil {
		if c.repo.IsAlreadyExistsError(err) {
			if err := c.notifier.NameAlreadyRegistered(
				ctx,
				currentPlayer.ChatID,
				displayName,
			); err != nil {
				return fmt.Errorf("name already redistered notification: %w", err)
			}

			return nil
		}

		return fmt.Errorf("repo change display name: %w", err)
	}

	currentPlayer.DisplayName = displayName

	if err := c.notifier.NameChanged(ctx, currentPlayer); err != nil {
		return fmt.Errorf("name changed notification: %w", err)
	}

	return nil
}

// validateParams check preconditions for change name process.
// returns is_processed flag and error.
func (c *ChangeName) validateParams(
	ctx context.Context,
	plr player.Player,
	displayName string,
	force bool,
) (bool, error) {
	if utf8.RuneCountInString(displayName) > c.maxNameLen {
		if err := c.notifier.NameIsTooLong(ctx, plr.ChatID, c.maxNameLen); err != nil {
			return false, fmt.Errorf("name is too long: %w", err)
		}

		return true, nil
	}

	if force {
		ok, err := c.isNeedChangeNameDelay(ctx, plr)
		if err != nil {
			return false, fmt.Errorf("change name delay: %w", err)
		}

		return ok, nil
	}

	if !plr.IsChangeNameTriggered {
		return false, errors.New("change name not triggered")
	}

	return false, nil
}

func (c *ChangeName) isNeedChangeNameDelay(ctx context.Context, plr player.Player) (bool, error) {
	var (
		lastNameChangedInterval = time.Since(plr.NameChangedAt)
	)

	if lastNameChangedInterval < c.changeNameTimeout {
		if err := c.notifier.ChangeNameDelay(
			ctx,
			plr,
			formatDelta(c.changeNameTimeout-lastNameChangedInterval),
		); err != nil {
			return false, fmt.Errorf("change name delay: %w", err)
		}

		return true, nil
	}

	return false, nil
}
