package startgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (s *StartGame) StopFind(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	currentPlayer, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	switch currentPlayer.GameStatus {
	case player.GameStatusIdle, player.GameStatusInGame:

	case player.GameStatusReadyForGame:
		if err := s.stopSearchInReeadyToGameStatus(ctx, chatID, currentPlayer); err != nil {
			return fmt.Errorf("stop search in ready to game status: %w", err)
		}
	}

	return nil
}

func (s *StartGame) stopSearchInReeadyToGameStatus(
	ctx context.Context,
	chatID msginfo.ChatID,
	currentPlayer player.Player,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.ChangePlayerGameStatus(
			ctx,
			currentPlayer.ID,
			"",
			player.GameStatusIdle,
			s.timeProvider.Now(),
			player.GameStatusReadyForGame,
		); err != nil {
			return fmt.Errorf("change player game status: %w", err)
		}

		if err := s.notifier.StopSearchGame(ctx, chatID); err != nil {
			return fmt.Errorf("stop search notification: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
