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
	ChangeStatusForFinishedGame(
		ctx context.Context,
		finishedGame game.Game,
		finishedAt time.Time,
	) error
}

type Notifier interface {
	GameFinish(ctx context.Context, gameType game.GameType, state game.State, finishedAt time.Time) error
	GameCanceled(ctx context.Context, players []game.Player) error
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

		switch currentGame.Status {
		case game.GameStatusInProgress:
			currentGame.MissForRestShotsAndCompleteGame(playerID, completedAt)

		case game.GameStatusPending:
			currentGame.Cancel(completedAt)

		case game.GameStatusCompleted, game.GameStatusCanceled:
			fallthrough

		default:
			return perror.InvalidGameState()
		}

		if err := s.playerStatusChanger.ChangeStatusForFinishedGame(
			ctx,
			currentGame,
			completedAt,
		); err != nil {
			return fmt.Errorf("change players game status: %w", err)
		}

		if err := s.sendNotifications(ctx, currentGame); err != nil {
			return fmt.Errorf("send notifications: %w", err)
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

func (s *LeaveGame) sendNotifications(ctx context.Context, currentGame game.Game) error {
	switch currentGame.Status {
	case game.GameStatusCompleted:
		if err := s.notifier.GameFinish(ctx, currentGame.Type, currentGame.State, currentGame.StateUpdatedAt); err != nil {
			return fmt.Errorf("game finish notification: %w", err)
		}

	case game.GameStatusCanceled:
		if err := s.notifier.GameCanceled(ctx, currentGame.State.LivePlayers()); err != nil {
			return fmt.Errorf("game canceled notification: %w", err)
		}

	case game.GameStatusInProgress, game.GameStatusPending:
		fallthrough

	default:
		return perror.InvalidGameState()
	}

	return nil
}
