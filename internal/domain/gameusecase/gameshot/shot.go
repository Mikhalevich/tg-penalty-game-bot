package gameshot

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

var (
	//nolint:gochecknoglobals
	possibleBotShotSides = []game.ShotSide{
		game.ShotSideLeft,
		game.ShotSideMiddle,
		game.ShotSideRight,
	}
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

	if err := g.notifier.GameRoundFinish(ctx, currentGame.State); err != nil {
		return fmt.Errorf("round finish: %w", err)
	}

	if err := g.startNewRound(ctx, currentGame, shot.CompletedAt); err != nil {
		return fmt.Errorf("start new round: %w", err)
	}

	return nil
}

func (g *GameShot) startNewRound(
	ctx context.Context,
	currentGame *game.Game,
	startedAt time.Time,
) error {
	if err := currentGame.StartNextRound(startedAt); err != nil {
		return fmt.Errorf("start next round: %w", err)
	}

	if currentGame.IsFinished() {
		if err := g.finishGame(ctx, currentGame, startedAt); err != nil {
			return fmt.Errorf("finish game: %w", err)
		}

		return nil
	}

	if err := g.notifier.GameNewRound(
		ctx,
		currentGame.ID,
		currentGame.State,
	); err != nil {
		return fmt.Errorf("start next round: %w", err)
	}

	return nil
}

func (g *GameShot) finishGame(
	ctx context.Context,
	currentGame *game.Game,
	finishedAt time.Time,
) error {
	if err := g.playerStatusChanger.ChangeStatusForFinishedGame(
		ctx,
		*currentGame,
		finishedAt,
	); err != nil {
		return fmt.Errorf("change players game status: %w", err)
	}

	if currentGame.IsRatingGame() {
		if err := g.repo.InsertShots(ctx, currentGame.State.CompletedShots()); err != nil {
			return fmt.Errorf("insert shots: %w", err)
		}
	}

	if err := g.notifier.GameFinish(ctx, currentGame.State, currentGame.StateUpdatedAt); err != nil {
		return fmt.Errorf("game finish: %w", err)
	}

	return nil
}

func (g *GameShot) processBotShot(
	currentGame *game.Game,
	round int,
	completedAt time.Time,
) error {
	if err := currentGame.PlayerShot(game.Shot{
		GameID:       currentGame.ID,
		PlayerID:     0,
		Round:        round,
		ExpectedSide: generateBotSide(),
		CompletedAt:  completedAt,
	}); err != nil {
		return fmt.Errorf("player shot: %w", err)
	}

	return nil
}

func generateBotSide() game.ShotSide {
	//nolint:gosec
	return possibleBotShotSides[rand.Int()%len(possibleBotShotSides)]
}
