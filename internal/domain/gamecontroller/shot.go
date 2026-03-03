package gamecontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (gc *GameController) Shot(
	ctx context.Context,
	chatID msginfo.ChatID,
	msgID msginfo.MessageID,
	gameID game.ID,
	round int,
	side game.ShotSide,
) error {
	currentPlayer, err := gc.playerController.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player: %w", err)
	}

	if !currentPlayer.IsInGame(gameID.String()) {
		return fmt.Errorf("not in game %s", gameID)
	}

	if err := gc.transactor.Transaction(ctx, func(ctx context.Context) error {
		currentGame, err := gc.repo.GetGame(ctx, gameID)
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}

		if currentGame.Status != game.GameStatusInProgress {
			return perror.InvalidGameState()
		}

		if err := gc.processGameShot(
			ctx,
			&currentGame,
			currentPlayer,
			round,
			side,
			gc.timeProvider.Now(),
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

	if err := gc.messageDeleter.DeleteMessage(ctx, chatID, msgID); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}

func (gc *GameController) processGameShot(
	ctx context.Context,
	currentGame *game.Game,
	currentPlayer player.Player,
	round int,
	side game.ShotSide,
	now time.Time,
) error {
	if err := currentGame.PlayerShot(game.Shot{
		GameID:      currentGame.ID,
		PlayerID:    currentPlayer.ID,
		Round:       round,
		Side:        side,
		CompletedAt: now,
	}); err != nil {
		return fmt.Errorf("player shot: %w", err)
	}

	if currentGame.IsGameWithBot() {
		if err := gc.processBotShot(currentGame, round, now); err != nil {
			return fmt.Errorf("bot shot: %w", err)
		}
	}

	if !currentGame.TryToCompleteRound() {
		return nil
	}

	if err := gc.notifier.GameStage(ctx, *currentGame); err != nil {
		return fmt.Errorf("complete round game stage: %w", err)
	}

	if err := currentGame.StartNextRound(now); err != nil {
		return fmt.Errorf("start next round: %w", err)
	}

	if err := gc.notifier.GameStage(ctx, *currentGame); err != nil {
		return fmt.Errorf("start next round game stage: %w", err)
	}

	if currentGame.IsFinished() {
		if err := gc.playerController.ChangePlayersGameStatus(
			ctx,
			currentGame.PlayerIDs(),
			player.GameStatusIdle,
			now,
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
