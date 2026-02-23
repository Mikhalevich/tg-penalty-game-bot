package setup

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/buttonrespository"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/namegenerator/randomgenerator"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/notifier"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playercontroller"
)

func StartBot(ctx context.Context, cfg config.Config) error {
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

	btnRepo, err := MakeRedisButtonRepository(ctx, cfg.ButtonRedis)
	if err != nil {
		return fmt.Errorf("make button redis: %w", err)
	}

	var (
		nameGenerator    = randomgenerator.New(cfg.RandomNameGenerator.Prefix, cfg.RandomNameGenerator.Length)
		tp               = timeprovider.New()
		msgSender        = messagesender.New(botAPI)
		msgProcessor     = messageprocessor.New(msgSender, msgSender, btnRepo)
		notification     = notifier.New(msgProcessor, msgProcessor)
		playerController = playercontroller.New(pgDB, notification, nameGenerator, tp, cfg.ChangeNameInterval)
	)

	if err := app.Start(
		ctx,
		cfg.Bot.Token,
		playerController,
		msgProcessor,
	); err != nil {
		return fmt.Errorf("app start: %w", err)
	}

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

func MakeRedisButtonRepository(
	ctx context.Context,
	cfg config.ButtonRedis,
) (*buttonrespository.ButtonRepository, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.DB,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("redis instrument tracing: %w", err)
	}

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return buttonrespository.New(rdb, cfg.TTL), nil
}
