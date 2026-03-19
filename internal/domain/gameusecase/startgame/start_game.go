package startgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type Repository interface {
	InsertGames(ctx context.Context, games []game.Game) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type Notifier interface {
	GameNewRound(ctx context.Context, gameID game.ID, state game.State) error
	GameStart(ctx context.Context, player1, player2 game.Player) error
}

type StartGame struct {
	repo       Repository
	transactor Transactor
	notifier   Notifier
}

func New(
	repo Repository,
	transactor Transactor,
	notifier Notifier,
) *StartGame {
	return &StartGame{
		repo:       repo,
		transactor: transactor,
		notifier:   notifier,
	}
}

func (s *StartGame) StartGames(ctx context.Context, games []game.Game) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.InsertGames(ctx, games); err != nil {
			return fmt.Errorf("insert games: %w", err)
		}

		for _, currentGame := range games {
			if err := s.notifier.GameStart(ctx, currentGame.State.Player1, currentGame.State.Player2); err != nil {
				return fmt.Errorf("start game notification: %w", err)
			}

			if err := s.notifier.GameNewRound(ctx, currentGame.ID, currentGame.State); err != nil {
				return fmt.Errorf("new round notification: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
