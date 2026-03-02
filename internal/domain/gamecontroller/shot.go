package gamecontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

//nolint:cyclop
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

	currentGame, err := gc.repo.GetGame(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get game: %w", err)
	}

	if currentGame.Status != game.GameStatusInProgress {
		return perror.InvalidGameState()
	}

	now := gc.timeProvider.Now()

	if err := currentGame.PlayerShot(game.Shot{
		GameID:      currentGame.ID,
		PlayerID:    currentPlayer.ID,
		Round:       round,
		Side:        side,
		CompletedAt: now,
	}); err != nil {
		return fmt.Errorf("player shot: %w", err)
	}

	if currentGame.TryToCompleteRound() {
		if err := gc.notifier.GameStage(ctx, currentGame); err != nil {
			return fmt.Errorf("complete round game stage: %w", err)
		}

		if err := currentGame.StartNextRound(now); err != nil {
			return fmt.Errorf("start next round: %w", err)
		}

		if err := gc.notifier.GameStage(ctx, currentGame); err != nil {
			return fmt.Errorf("start next round game stage: %w", err)
		}
	}

	if err := gc.repo.UpdateGame(ctx, currentGame); err != nil {
		return fmt.Errorf("update game: %w", err)
	}

	if err := gc.messageDeleter.DeleteMessage(ctx, chatID, msgID); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
