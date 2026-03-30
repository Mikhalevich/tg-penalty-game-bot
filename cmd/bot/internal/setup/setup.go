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
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/gameshot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/getgame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/joingame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/leavegame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/startgame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/notifier"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/changename"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/changestatus"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/findgame"
	playergameshot "github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/gameshot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/gameshotstate"
	playerleavegame "github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/leavegame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/playerprovider"
	playerstartgame "github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/startgame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/welcome"
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
		nameGenerator       = randomgenerator.New(cfg.RandomNameGenerator.Prefix, cfg.RandomNameGenerator.Length)
		timeProvider        = timeprovider.New()
		msgSender           = messagesender.New(botAPI)
		msgProcessor        = messageprocessor.New(msgSender, msgSender, btnRepo, nil)
		notificationService = notifier.New(pgDB, msgSender)
		playerProvider      = playerprovider.New(pgDB, nameGenerator, timeProvider)
		welcomeService      = welcome.New(playerProvider, notificationService)
		changeNameService   = changename.New(pgDB, playerProvider, timeProvider,
			notificationService, cfg.ChangeNameInterval)
		getGameService         = getgame.New(pgDB)
		gameShotStateService   = gameshotstate.New(getGameService, notificationService)
		findGameSerivce        = findgame.New(pgDB, playerProvider, timeProvider, gameShotStateService)
		changeStatusService    = changestatus.New(pgDB)
		startGameService       = startgame.New(pgDB, pgDB.Transactor(), changeStatusService, notificationService)
		joinGameService        = joingame.New(pgDB, pgDB.Transactor(), changeStatusService, notificationService)
		playerStartGameService = playerstartgame.New(pgDB, pgDB.Transactor(),
			playerProvider, startGameService, joinGameService, timeProvider, gameShotStateService, notificationService)
		gameShotService        = gameshot.New(pgDB, pgDB.Transactor(), changeStatusService, notificationService)
		playerGameShotService  = playergameshot.New(playerProvider, gameShotService, timeProvider, msgProcessor)
		leaveGameService       = leavegame.New(pgDB, pgDB.Transactor(), changeStatusService, notificationService)
		leavePlayerGameService = playerleavegame.New(playerProvider, leaveGameService, timeProvider)
	)

	if err := app.Start(
		ctx,
		cfg.Bot.Token,
		msgProcessor,
		playerProvider,
		welcomeService,
		changeNameService,
		findGameSerivce,
		playerStartGameService,
		playerGameShotService,
		leavePlayerGameService,
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
