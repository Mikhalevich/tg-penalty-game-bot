package app

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tghandler"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

func Start(
	ctx context.Context,
	botCfg config.Bot,
	buttonProvider tghandler.ButtonProvider,
	playerProvider tghandler.PlayerProvider,
	welcome tghandler.Welcome,
	changeName tghandler.ChangeName,
	startGame tghandler.StartGame,
	gameShot tghandler.GameShot,
	leaveGame tghandler.LeaveGame,
	leaderboard tghandler.Leaderboard,
	shotStats tghandler.ShotStats,
	notifier tghandler.Notifier,
) error {
	var (
		botHandler = tghandler.New(
			buttonProvider,
			playerProvider,
			welcome,
			changeName,
			startGame,
			gameShot,
			leaveGame,
			leaderboard,
			shotStats,
			notifier,
		)
	)

	tbot, err := tgbot.New(botCfg.Token, botCfg.WebHookToken, logger.FromContext(ctx))
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	makeRoutes(tbot, botHandler)

	if err := tbot.Start(ctx); err != nil {
		return fmt.Errorf("bot start: %w", err)
	}

	return nil
}
