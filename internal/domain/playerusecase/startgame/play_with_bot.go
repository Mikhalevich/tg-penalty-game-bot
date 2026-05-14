package startgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (s *StartGame) PlayWithBot(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
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
		return perror.InSearchGameState()

	case player.GameStatusIdle:
	}

	if err := s.notifier.SelectBotDifficulty(ctx, chatID); err != nil {
		return fmt.Errorf("select bot difficulty: %w", err)
	}

	return nil
}
