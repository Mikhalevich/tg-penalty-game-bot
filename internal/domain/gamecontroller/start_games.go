package gamecontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

func (gc *GameController) StartGames(ctx context.Context, games []game.Game) error {
	if err := gc.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := gc.repo.InsertGames(ctx, games); err != nil {
			return fmt.Errorf("insert games: %w", err)
		}

		for _, currentGame := range games {
			if err := gc.notifier.GameStage(ctx, currentGame); err != nil {
				return fmt.Errorf("game stage notification: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
