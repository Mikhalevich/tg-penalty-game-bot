package setup

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/gamepoller/internal/app"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/gamepoller/internal/config"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gamecontroller"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/matchmaking"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/notifier"
)

func StartWorker(ctx context.Context, cfg config.Config) error {
	botAPI, err := bot.New(
		cfg.Bot.Token,
		bot.WithSkipGetMe(),
	)
	if err != nil {
		return fmt.Errorf("creating bot api: %w", err)
	}

	pgDB, dbCleanup, err := MakePostgres(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}

	defer dbCleanup()

	var (
		msgSender            = messagesender.New(botAPI)
		gameNofifier         = notifier.New(nil, msgSender)
		timeProvider         = timeprovider.New()
		gameController       = gamecontroller.New(nil, nil, timeProvider)
		matchmakingProcessor = matchmaking.New(
			pgDB,
			pgDB.Transactor(),
			gameController,
			gameNofifier,
		)
	)

	app.New(matchmakingProcessor).Run(ctx, cfg.Worker)

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
