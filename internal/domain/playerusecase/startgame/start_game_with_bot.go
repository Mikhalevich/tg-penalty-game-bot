package startgame

import (
	"context"
	"fmt"
	"time"

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

	switch currentPlayer.GameStatus {
	case player.GameStatusInGame:
		if err := s.shotStater.ShotState(ctx, currentPlayer); err != nil {
			return fmt.Errorf("shot state: %w", err)
		}

		return nil

	case player.GameStatusReadyForGame:
		if err := s.notifier.PlayerAlreadyInGame(ctx, currentPlayer); err != nil {
			return fmt.Errorf("already in game notitication: %w", err)
		}

		return nil

	case player.GameStatusIdle:
	}

	if err := s.startGameWithBot(ctx, currentPlayer, s.timeProvider.Now()); err != nil {
		return fmt.Errorf("start game with bot: %w", err)
	}

	return nil
}

func (s *StartGame) startGameWithBot(
	ctx context.Context,
	currentPlayer player.Player,
	createdAt time.Time,
) error {
	currentGame := game.CreateGame(
		ctx,
		game.GameTypeFriendly,
		currentPlayer,
		player.Player{
			ID:          0,
			DisplayName: "bot",
		},
		createdAt,
	)

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.ChangePlayerGameStatus(
			ctx,
			currentPlayer.ID,
			currentGame.ID,
			player.GameStatusInGame,
			createdAt,
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
