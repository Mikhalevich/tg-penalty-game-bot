package gameshot

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (g *GameShot) Shot(
	ctx context.Context,
	shot game.Shot,
) error {
	if err := g.transactor.Transaction(ctx, func(ctx context.Context) error {
		currentGame, err := g.repo.GetGame(ctx, shot.GameID)
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}

		if currentGame.Status != game.GameStatusInProgress {
			return perror.InvalidGameState()
		}

		if err := g.processGameShot(
			ctx,
			&currentGame,
			shot,
		); err != nil {
			return fmt.Errorf("process game shot: %w", err)
		}

		if err := g.repo.UpdateGame(ctx, currentGame); err != nil {
			return fmt.Errorf("update game: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (g *GameShot) processGameShot(
	ctx context.Context,
	currentGame *game.Game,
	shot game.Shot,
) error {
	if err := currentGame.PlayerShot(shot); err != nil {
		return fmt.Errorf("player shot: %w", err)
	}

	if currentGame.IsGameWithBot() {
		if err := g.processBotShot(currentGame, shot.Round, shot.CompletedAt); err != nil {
			return fmt.Errorf("bot shot: %w", err)
		}
	}

	if !currentGame.TryToCompleteRound() {
		return nil
	}

	if err := g.notifier.GameStage(ctx, *currentGame); err != nil {
		return fmt.Errorf("complete round game stage: %w", err)
	}

	if err := currentGame.StartNextRound(shot.CompletedAt); err != nil {
		return fmt.Errorf("start next round: %w", err)
	}

	if err := g.notifier.GameStage(ctx, *currentGame); err != nil {
		return fmt.Errorf("start next round game stage: %w", err)
	}

	if currentGame.IsFinished() {
		if err := g.playerStatusChanger.ChangePlayersGameStatus(
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

func (g *GameShot) processBotShot(
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
