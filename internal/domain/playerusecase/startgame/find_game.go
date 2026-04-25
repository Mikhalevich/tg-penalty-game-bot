package startgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (s *StartGame) FindGame(ctx context.Context, chatID msginfo.ChatID) error {
	currentPlayer, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	switch currentPlayer.GameStatus {
	case player.GameStatusInGame:
		if err := s.shotStater.ShotState(ctx, currentPlayer); err != nil {
			return fmt.Errorf("shot state: %w", err)
		}

		return nil

	case player.GameStatusReadyForGame:
		return nil

	case player.GameStatusIdle:
	}

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.ChangePlayerGameStatus(
			ctx,
			currentPlayer.ID,
			"",
			player.GameStatusReadyForGame,
			s.timeProvider.Now(),
			player.GameStatusIdle,
		); err != nil {
			return fmt.Errorf("change player game status: %w", err)
		}

		if err := s.notifier.SearchGame(ctx, chatID); err != nil {
			return fmt.Errorf("search notification: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
