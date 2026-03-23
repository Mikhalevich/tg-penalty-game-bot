package startgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (s *StartGame) StartGameWithBot(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	currentPlayer, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	if currentPlayer.GameStatus != player.GameStatusIdle {
		if err := s.notifier.PlayerAlreadyInGame(ctx, currentPlayer); err != nil {
			return fmt.Errorf("already in game notitication: %w", err)
		}

		return nil
	}

	now := s.timeProvider.Now()

	currentGame, err := game.CreateGame(
		ctx,
		game.GameTypeFriendly,
		currentPlayer,
		player.Player{
			ID:          0,
			DisplayName: "bot",
		},
		now,
	)

	if err != nil {
		return fmt.Errorf("create game: %w", err)
	}

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.ChangePlayerGameStatus(
			ctx,
			currentPlayer.ID,
			currentGame.ID,
			player.GameStatusInGame,
			now,
			player.GameStatusIdle,
		); err != nil {
			return fmt.Errorf("change player game status: %w", err)
		}

		if err := s.gameRunner.StartGames(ctx, []game.Game{currentGame}); err != nil {
			return fmt.Errorf("start games: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
