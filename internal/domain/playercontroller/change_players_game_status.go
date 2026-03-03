package playercontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *PlayerController) ChangePlayersGameStatus(
	ctx context.Context,
	playerIDs []player.ID,
	status player.GameStatus,
	changedAt time.Time,
) error {
	if err := p.repo.ChangePlayersGameStatus(
		ctx,
		playerIDs,
		status,
		changedAt,
	); err != nil {
		return fmt.Errorf("change playres game status: %w", err)
	}

	return nil
}
