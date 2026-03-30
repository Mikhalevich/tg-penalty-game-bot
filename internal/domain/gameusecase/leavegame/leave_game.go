package leavegame

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	GetGame(ctx context.Context, gameID game.ID) (game.Game, error)
	UpdateGame(ctx context.Context, game game.Game) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerStatusChanger interface {
	ChangePlayersGameStatus(
		ctx context.Context,
		playerIDs []player.ID,
		status player.GameStatus,
		changedAt time.Time,
	) error
}

type Notifier interface {
	GameFinish(ctx context.Context, state game.State) error
}

type LeaveGame struct {
	repo                Repository
	transactor          Transactor
	playerStatusChanger PlayerStatusChanger
	notifier            Notifier
}

func New(
	repo Repository,
	transactor Transactor,
	playerStatusChanger PlayerStatusChanger,
	notifier Notifier,
) *LeaveGame {
	return &LeaveGame{
		repo:                repo,
		transactor:          transactor,
		playerStatusChanger: playerStatusChanger,
		notifier:            notifier,
	}
}

func (s *LeaveGame) LeaveGame(
	ctx context.Context,
	playerID player.ID,
	gameID game.ID,
	completedAt time.Time,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		currentGame, err := s.repo.GetGame(ctx, gameID)
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}

		if currentGame.Status != game.GameStatusInProgress {
			return perror.InvalidGameState()
		}

		currentGame.MissForRestShotsAndCompleteGame(playerID, completedAt)

		if err := s.playerStatusChanger.ChangePlayersGameStatus(
			ctx,
			currentGame.PlayerIDs(),
			player.GameStatusIdle,
			completedAt,
		); err != nil {
			return fmt.Errorf("change players game status: %w", err)
		}

		if err := s.notifier.GameFinish(ctx, currentGame.State); err != nil {
			return fmt.Errorf("game finish notification: %w", err)
		}

		if err := s.repo.UpdateGame(ctx, currentGame); err != nil {
			return fmt.Errorf("update game: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
