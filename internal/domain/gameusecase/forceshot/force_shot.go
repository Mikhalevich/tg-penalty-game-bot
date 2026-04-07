package forceshot

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

type Repository interface {
	GetShotExpiredRatingGames(ctx context.Context, expiredTime time.Time, limit int) ([]game.Game, error)
	UpdateGame(ctx context.Context, game game.Game) error
}

type TimeProvider interface {
	Now() time.Time
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
	GameNewRound(ctx context.Context, gameID game.ID, state game.State) error
	GameRoundFinish(ctx context.Context, state game.State) error
	GameFinish(ctx context.Context, state game.State, finishedAt time.Time) error
}

type ForceShot struct {
	repo                Repository
	timeProvider        TimeProvider
	playerStatusChanger PlayerStatusChanger
	notifier            Notifier
}

func New(
	repo Repository,
	timeProvider TimeProvider,
	playerStatusChanger PlayerStatusChanger,
	notifier Notifier,
) *ForceShot {
	return &ForceShot{
		repo:                repo,
		timeProvider:        timeProvider,
		playerStatusChanger: playerStatusChanger,
		notifier:            notifier,
	}
}

func (s *ForceShot) ForceShotForGames(
	ctx context.Context,
	expirationDuration time.Duration,
	limit int,
) error {
	var (
		now                = s.timeProvider.Now()
		shotExpirationTime = now.Add(-expirationDuration)
	)

	games, err := s.repo.GetShotExpiredRatingGames(
		ctx,
		shotExpirationTime,
		limit,
	)

	if err != nil {
		return fmt.Errorf("get games: %w", err)
	}

	if len(games) == 0 {
		return nil
	}

	for _, currentGame := range games {
		if err := s.updateGameShot(ctx, currentGame, now); err != nil {
			logger.FromContext(ctx).WithError(err).Error("update game shot")
		}
	}

	return nil
}

func (s *ForceShot) updateGameShot(
	ctx context.Context,
	currentGame game.Game,
	currentTime time.Time,
) error {
	currentGame.MakeMissShots(currentTime)

	if !currentGame.TryToCompleteRound() {
		return nil
	}

	if err := s.notifier.GameRoundFinish(ctx, currentGame.State); err != nil {
		return fmt.Errorf("round finish: %w", err)
	}

	if err := s.startNextRound(ctx, &currentGame, currentTime); err != nil {
		return fmt.Errorf("start next round: %w", err)
	}

	if err := s.repo.UpdateGame(ctx, currentGame); err != nil {
		return fmt.Errorf("update game: %w", err)
	}

	return nil
}

func (s *ForceShot) startNextRound(
	ctx context.Context,
	currentGame *game.Game,
	startedAt time.Time,
) error {
	if err := currentGame.StartNextRound(startedAt); err != nil {
		return fmt.Errorf("start next round: %w", err)
	}

	if !currentGame.IsFinished() {
		if err := s.notifier.GameNewRound(ctx, currentGame.ID, currentGame.State); err != nil {
			return fmt.Errorf("start next round: %w", err)
		}

		return nil
	}

	if err := s.playerStatusChanger.ChangePlayersGameStatus(
		ctx,
		currentGame.PlayerIDs(),
		player.GameStatusIdle,
		startedAt,
	); err != nil {
		return fmt.Errorf("change players game status: %w", err)
	}

	if err := s.notifier.GameFinish(ctx, currentGame.State, currentGame.StateUpdatedAt); err != nil {
		return fmt.Errorf("game finish: %w", err)
	}

	return nil
}
