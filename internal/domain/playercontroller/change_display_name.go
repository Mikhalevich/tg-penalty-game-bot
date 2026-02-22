package playercontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (p *PlayerController) ChangeDisplayName(
	ctx context.Context,
	chatID msginfo.ChatID,
	displayName string,
) error {
	if err := p.repo.ChangeDisplayName(ctx, chatID, displayName, p.timeProvider.Now()); err != nil {
		if p.repo.IsAlreadyExistsError(err) {
			return perror.AlreadyExists("name already reserved")
		}

		return fmt.Errorf("repo change display name: %w", err)
	}

	return nil
}
