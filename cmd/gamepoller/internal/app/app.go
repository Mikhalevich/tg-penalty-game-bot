package app

import (
	"context"
	"sync"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/gamepoller/internal/config"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

type MatchmakingProcessor interface {
	ProcessReadyToGamePlayers(ctx context.Context, playersLimit int) error
}

type ExpiredShotsProcessor interface {
	ForceShotForGames(ctx context.Context, expirationDuration time.Duration, limit int) error
}

type LeaderboradRefresher interface {
	RefreshScoreLeaderboard(ctx context.Context) error
}

type ExpiredSearchGamer interface {
	StopSearch(ctx context.Context, searchPeriod time.Duration) error
}

type App struct {
	matchmakingProcessor  MatchmakingProcessor
	expiredShotsProcessor ExpiredShotsProcessor
	leaderboradRefresher  LeaderboradRefresher
	expiredSearchGamer    ExpiredSearchGamer
}

func New(
	matchmakingProcessor MatchmakingProcessor,
	expiredShotsProcessor ExpiredShotsProcessor,
	leaderboradRefresher LeaderboradRefresher,
	expiredSearchGamer ExpiredSearchGamer,
) *App {
	return &App{
		matchmakingProcessor:  matchmakingProcessor,
		expiredShotsProcessor: expiredShotsProcessor,
		leaderboradRefresher:  leaderboradRefresher,
		expiredSearchGamer:    expiredSearchGamer,
	}
}

func (a *App) Run(
	ctx context.Context,
	matchmakingCfg config.Worker,
	expiredShotsCfg config.ExpiredShotsWorker,
	leaderboardCfg config.Worker,
	expiredSearchGameCfg config.ExpiredSearchGameWorker,
) {
	var wgr sync.WaitGroup

	runWorkers(ctx, "players ready to game", matchmakingCfg.Count, matchmakingCfg.Interval, &wgr,
		func(ctx context.Context) error {
			return a.matchmakingProcessor.ProcessReadyToGamePlayers(
				ctx,
				matchmakingCfg.BatchSize,
			)
		},
	)

	runWorkers(ctx, "expires shots", expiredShotsCfg.Count, expiredShotsCfg.Interval, &wgr,
		func(ctx context.Context) error {
			return a.expiredShotsProcessor.ForceShotForGames(
				ctx,
				expiredShotsCfg.ShotExpireDuration,
				expiredShotsCfg.BatchSize,
			)
		},
	)

	runWorkers(ctx, "leaderboard refresher", leaderboardCfg.Count, leaderboardCfg.Interval, &wgr,
		func(ctx context.Context) error {
			return a.leaderboradRefresher.RefreshScoreLeaderboard(ctx)
		},
	)

	runWorkers(ctx, "expired search game", expiredSearchGameCfg.Count, expiredSearchGameCfg.Interval, &wgr,
		func(ctx context.Context) error {
			return a.expiredSearchGamer.StopSearch(ctx, expiredSearchGameCfg.SearchExpireDuration)
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
