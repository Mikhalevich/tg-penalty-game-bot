package startgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
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
	GameLink(ctx context.Context, chatID msginfo.ChatID, gameID game.ID) error
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
			if err := s.sendNotificationForCreatedGame(ctx, currentGame); err != nil {
				return fmt.Errorf("send notification: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (s *StartGame) sendNotificationForCreatedGame(ctx context.Context, currentGame game.Game) error {
	switch currentGame.Status {
	case game.GameStatusPending:
		if err := s.notifier.GameLink(ctx, currentGame.State.Player1.ChatID, currentGame.ID); err != nil {
			return fmt.Errorf("game link: %w", err)
		}

	case game.GameStatusInProgress:
		if err := s.notifier.GameStart(ctx, currentGame.State.Player1, currentGame.State.Player2); err != nil {
			return fmt.Errorf("start game: %w", err)
		}

		if err := s.notifier.GameNewRound(ctx, currentGame.ID, currentGame.State); err != nil {
			return fmt.Errorf("new round: %w", err)
		}

	case game.GameStatusCanceled, game.GameStatusCompleted:
		fallthrough

	default:
		return fmt.Errorf("invalid status: %s", currentGame.Status.String())
	}

	return nil
}
