package gamecontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (gc *GameController) Shot(
	ctx context.Context,
	chatID msginfo.ChatID,
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

	if err := gc.repo.UpdateShot(ctx, game.Shot{
		GameID:      gameID,
		PlayerID:    currentPlayer.ID,
		Round:       round,
		Side:        side,
		CompletedAt: gc.timeProvider.Now(),
	}); err != nil {
		if gc.repo.IsNoRowsUpdated(err) {
			return fmt.Errorf("already shot: %w", err)
		}

		return fmt.Errorf("update shot: %w", err)
	}

	shots, err := gc.repo.GetRoundShotsByGame(ctx, gameID, round)
	if err != nil {
		return fmt.Errorf("get round shots by game: %w", err)
	}

	if isShotsCompleted(shots) {
		if err := gc.completeRound(ctx, gameID, round, shots); err != nil {
			return fmt.Errorf("complete round: %w", err)
		}
	}

	return nil
}

func isShotsCompleted(shots []game.Shot) bool {
	for _, s := range shots {
		if s.Side == game.ShotSideNoShot {
			return false
		}
	}

	return true
}

func (gc *GameController) completeRound(
	ctx context.Context,
	gameID game.ID,
	round int,
	shots []game.Shot,
) error {
	currentGame, err := gc.repo.GetGame(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get game: %w", err)
	}

	if currentGame.Status != game.GameStatusInProgress {
		return perror.InvalidGameState()
	}

	if currentGame.CurrentRound() != round {
		return perror.InvalidRound()
	}

	if err := currentGame.CompleteRound(shots); err != nil {
		return fmt.Errorf("complete round: %w", err)
	}

	if err := gc.repo.UpdateGame(ctx, currentGame); err != nil {
		return fmt.Errorf("update game: %w", err)
	}

	return nil
}
