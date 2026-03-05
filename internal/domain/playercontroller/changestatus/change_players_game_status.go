package changestatus

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (c *ChangeStatus) ChangePlayersGameStatus(
	ctx context.Context,
	playerIDs []player.ID,
	status player.GameStatus,
	changedAt time.Time,
) error {
	if err := c.repo.ChangePlayersGameStatus(
		ctx,
		playerIDs,
		status,
		changedAt,
	); err != nil {
		return fmt.Errorf("change playres game status: %w", err)
	}

	return nil
}
