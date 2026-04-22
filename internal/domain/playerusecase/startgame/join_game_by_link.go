package startgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (s *StartGame) JoinGameByLink(
	ctx context.Context,
	chatID msginfo.ChatID,
	gameID game.ID,
) error {
	currentPlayer, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	if currentPlayer.GameStatus != player.GameStatusIdle {
		return perror.AlreadyInGame()
	}

	ok, err := s.validateGameForJoin(ctx, chatID, gameID)
	if err != nil {
		return fmt.Errorf("validate game for join: %w", err)
	}

	if !ok {
		return nil
	}

	now := s.timeProvider.Now()

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.ChangePlayerGameStatus(
			ctx,
			currentPlayer.ID,
			gameID,
			player.GameStatusInGame,
			now,
			player.GameStatusIdle,
		); err != nil {
			return fmt.Errorf("change player game status: %w", err)
		}

		if err := s.gameJoiner.JoinGame(
			ctx,
			gameID,
			game.CreateGamePlayerFromPlayer(currentPlayer),
			now,
		); err != nil {
			return fmt.Errorf("join game: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

// validateGameForJoin validates game status for avalability to join
// returns true if player can join to game or false otherwise.
func (s *StartGame) validateGameForJoin(
	ctx context.Context,
	chatID msginfo.ChatID,
	gameID game.ID,
) (bool, error) {
	currentGame, err := s.gameGetter.GetGame(ctx, gameID)
	if err != nil {
		return false, fmt.Errorf("get game: %w", err)
	}

	switch currentGame.Status {
	case game.GameStatusInProgress, game.GameStatusCompleted:
		if err := s.notifier.LinkActivated(ctx, chatID); err != nil {
			return false, fmt.Errorf("link activated: %w", err)
		}

		return false, nil

	case game.GameStatusCanceled:
		if err := s.notifier.LinkCanceled(ctx, chatID); err != nil {
			return false, fmt.Errorf("link canceled: %w", err)
		}

		return false, nil

	case game.GameStatusPending:
	}

	return true, nil
}
