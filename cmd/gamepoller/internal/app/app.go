package app

import (
	"context"
	"sync"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/gamepoller/internal/config"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

type Processor interface {
	ProcessReadyToGamePlayers(ctx context.Context, playersLimit int) error
}

type App struct {
	processor Processor
}

func New(processor Processor) *App {
	return &App{
		processor: processor,
	}
}

func (a *App) Run(
	ctx context.Context,
	cfg config.Worker,
) {
	var wgr sync.WaitGroup

	runWorkers(ctx, "players ready to game", cfg.Count, cfg.Interval, &wgr,
		func(ctx context.Context) error {
			return a.processor.ProcessReadyToGamePlayers(ctx, cfg.BatchSize)
		},
	)

	wgr.Wait()
}

func runWorkers(
	ctx context.Context,
	workerName string,
	workersCount int,
	pollerInterval time.Duration,
	wgr *sync.WaitGroup,
	processFn func(ctx context.Context) error,
) {
	for i := range workersCount {
		wgr.Go(func() {
			log := logger.FromContext(ctx).
				WithFields(logger.Fields{
					"worker_name":   workerName,
					"worker_number": i,
				})
			runPoller(
				logger.WithLogger(ctx, log),
				pollerInterval,
				processFn,
			)
		})
	}
}

func runPoller(
	ctx context.Context,
	interval time.Duration,
	processFn func(ctx context.Context) error,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := processFn(ctx); err != nil {
				logger.FromContext(ctx).
					WithError(err).
					Error("process error")
			}

		case <-ctx.Done():
			return
		}
	}
}
