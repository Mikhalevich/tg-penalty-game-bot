package gamecontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (gc *GameController) Shot(
	ctx context.Context,
	shot game.Shot,
) error {
	if err := gc.transactor.Transaction(ctx, func(ctx context.Context) error {
		currentGame, err := gc.repo.GetGame(ctx, shot.GameID)
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}

		if currentGame.Status != game.GameStatusInProgress {
			return perror.InvalidGameState()
		}

		if err := gc.processGameShot(
			ctx,
			&currentGame,
			shot,
		); err != nil {
			return fmt.Errorf("process game shot: %w", err)
		}

		if err := gc.repo.UpdateGame(ctx, currentGame); err != nil {
			return fmt.Errorf("update game: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (gc *GameController) processGameShot(
	ctx context.Context,
	currentGame *game.Game,
	shot game.Shot,
) error {
	if err := currentGame.PlayerShot(shot); err != nil {
		return fmt.Errorf("player shot: %w", err)
	}

	if currentGame.IsGameWithBot() {
		if err := gc.processBotShot(currentGame, shot.Round, shot.CompletedAt); err != nil {
			return fmt.Errorf("bot shot: %w", err)
		}
	}

	if !currentGame.TryToCompleteRound() {
		return nil
	}

	if err := gc.notifier.GameStage(ctx, *currentGame); err != nil {
		return fmt.Errorf("complete round game stage: %w", err)
	}

	if err := currentGame.StartNextRound(shot.CompletedAt); err != nil {
		return fmt.Errorf("start next round: %w", err)
	}

	if err := gc.notifier.GameStage(ctx, *currentGame); err != nil {
		return fmt.Errorf("start next round game stage: %w", err)
	}

	if currentGame.IsFinished() {
		if err := gc.playerStatusChanger.ChangePlayersGameStatus(
			ctx,
			currentGame.PlayerIDs(),
			player.GameStatusIdle,
			shot.CompletedAt,
		); err != nil {
			return fmt.Errorf("change players game status: %w", err)
		}
	}

	return nil
}

func (gc *GameController) processBotShot(
	currentGame *game.Game,
	round int,
	completedAt time.Time,
) error {
	if err := currentGame.PlayerShot(game.Shot{
		GameID:      currentGame.ID,
		PlayerID:    0,
		Round:       round,
		Side:        game.ShotSideLeft,
		CompletedAt: completedAt,
	}); err != nil {
		return fmt.Errorf("player shot: %w", err)
	}

	return nil
}
