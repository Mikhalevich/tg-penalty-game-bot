package gamecontroller

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (gc *GameController) StartGameWithBot(ctx context.Context, chatID msginfo.ChatID) error {
	currentPlayer, err := gc.playerController.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	if currentPlayer.GameStatus != player.GameStatusIdle {
		if err := gc.notifier.PlayerAlreadyInGame(ctx, currentPlayer); err != nil {
			return fmt.Errorf("already in game notitication: %w", err)
		}

		return nil
	}

	currentGame, err := gc.CreateGame(ctx, currentPlayer, player.Player{
		ID:          0,
		DisplayName: "bot",
	})

	if err != nil {
		return fmt.Errorf("create game: %w", err)
	}

	if err := gc.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := gc.repo.InsertGames(ctx, []game.Game{currentGame}); err != nil {
			return fmt.Errorf("insert games: %w", err)
		}

		if err := gc.playerController.ChangeGameStatus(
			ctx,
			currentPlayer.ID,
			currentGame.ID,
			player.GameStatusInGame,
			gc.timeProvider.Now(),
		); err != nil {
			return fmt.Errorf("change player game status: %w", err)
		}

		if err := gc.notifier.GameStage(ctx, currentGame); err != nil {
			return fmt.Errorf("game stage notification: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
