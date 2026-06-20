package startgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (s *StartGame) StartGameByLink(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	currentPlayer, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	if currentPlayer.GameStatus != player.GameStatusIdle {
		return perror.AlreadyInGame()
	}

	now := s.timeProvider.Now()

	pendingGame := game.CreatePendingGame(
		currentPlayer,
		now,
	)

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.ChangePlayerGameStatus(
			ctx,
			currentPlayer.ID,
			pendingGame.ID,
			player.GameStatusInGame,
			now,
			player.GameStatusIdle,
		); err != nil {
			return fmt.Errorf("change player game status: %w", err)
		}

		if err := s.gameRunner.StartGames(ctx, []game.Game{pendingGame}); err != nil {
			return fmt.Errorf("start games: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
