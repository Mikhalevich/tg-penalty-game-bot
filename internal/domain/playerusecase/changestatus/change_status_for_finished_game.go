package changestatus

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (c *ChangeStatus) ChangeStatusForFinishedGame(
	ctx context.Context,
	finishedGame game.Game,
	finishedAt time.Time,
) error {
	if finishedGame.Type == game.GameTypeFriendly {
		if err := c.repo.ChangePlayersGameStatus(
			ctx,
			finishedGame.PlayerIDs(),
			player.GameStatusIdle,
			finishedAt,
		); err != nil {
			return fmt.Errorf("change players game status: %w", err)
		}

		return nil
	}

	goalsDiff := finishedGame.State.Player1.GoalsScored - finishedGame.State.Player2.GoalsScored

	if err := c.repo.SetPlayerIdleStatusWithScore(
		ctx,
		finishedGame.State.Player1.ID,
		goalsDiff,
		finishedAt,
	); err != nil {
		return fmt.Errorf("set player 1 idle status: %w", err)
	}

	if err := c.repo.SetPlayerIdleStatusWithScore(
		ctx,
		finishedGame.State.Player2.ID,
		-goalsDiff,
		finishedAt,
	); err != nil {
		return fmt.Errorf("set player 2 idle status: %w", err)
	}

	return nil
}
