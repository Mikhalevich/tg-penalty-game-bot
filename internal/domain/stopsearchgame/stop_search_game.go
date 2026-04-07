package stopsearchgame

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

type Repository interface {
	SetIdleForLongReadyToGamePlayerStatus(
		ctx context.Context,
		startSearchBefore time.Time,
		updatedAt time.Time,
	) ([]player.Player, error)
}

type TimeProvider interface {
	Now() time.Time
}

type Notifier interface {
	StopSearchGame(ctx context.Context, chatID msginfo.ChatID) error
}

type StopSearchGame struct {
	repo         Repository
	timeProvider TimeProvider
	notifier     Notifier
}

func New(
	repo Repository,
	timeProvider TimeProvider,
	notifier Notifier,
) *StopSearchGame {
	return &StopSearchGame{
		repo:         repo,
		timeProvider: timeProvider,
		notifier:     notifier,
	}
}

func (s *StopSearchGame) StopSearch(ctx context.Context, searchPeriod time.Duration) error {
	var (
		now               = s.timeProvider.Now()
		startSearchBefore = now.Add(-searchPeriod)
	)

	players, err := s.repo.SetIdleForLongReadyToGamePlayerStatus(
		ctx,
		startSearchBefore,
		now,
	)

	if err != nil {
		return fmt.Errorf("set idle status: %w", err)
	}

	for _, plr := range players {
		if err := s.notifier.StopSearchGame(ctx, plr.ChatID); err != nil {
			logger.FromContext(ctx).
				WithError(err).
				Error("sto search game notification")
		}
	}

	return nil
}
