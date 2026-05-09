package setup

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/gamepoller/internal/app"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/gamepoller/internal/config"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/markdownescaper"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/forceshot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/startgame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/lbrefresher"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/matchmaking"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/notifier"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/changestatus"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/settings"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/stopsearchgame"
)

func StartWorker(ctx context.Context, cfg config.Config) error {
	pgDB, dbCleanup, err := MakePostgres(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}

	defer dbCleanup()

	var (
		notificationService  = notifier.New(pgDB, markdownescaper.New())
		changeStatusService  = changestatus.New(pgDB)
		timeProvider         = timeprovider.New()
		startGameService     = startgame.New(pgDB, pgDB.Transactor(), changeStatusService, notificationService)
		matchmakingProcessor = matchmaking.New(
			pgDB,
			pgDB.Transactor(),
			startGameService,
			timeProvider,
		)
		expiredShotsProcessor      = forceshot.New(pgDB, timeProvider, changeStatusService, notificationService)
		settingsProvider           = settings.New(pgDB, timeProvider)
		lbRefresher                = lbrefresher.New(pgDB, pgDB.Transactor(), settingsProvider, timeProvider)
		expiredSearchGameProcessor = stopsearchgame.New(pgDB, timeProvider, notificationService)
	)

	app.New(
		matchmakingProcessor,
		expiredShotsProcessor,
		lbRefresher,
		expiredSearchGameProcessor,
	).Run(
		ctx,
		cfg.MatchmakingWorker,
		cfg.ExpiredShotsWorker,
		cfg.LeaderboardRecalculatorWorker,
		cfg.ExpiredSearchGameWorker,
	)

	return nil
}

func MakePostgres(cfg config.Postgres) (*postgres.Postgres, func(), error) {
	driver := driver.NewPgx()

	dbConn, err := otelsql.Open(driver.Name(), cfg.Connection)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	var (
		sqlxDBConn = sqlx.NewDb(dbConn, driver.Name())
		transactor = transaction.New(transaction.NewSqlxDB(sqlxDBConn))
		p          = postgres.New(driver, transactor)
	)

	return p, func() {
		dbConn.Close()
	}, nil
}
