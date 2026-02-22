package app

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tghandler"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

func Start(
	ctx context.Context,
	token string,
	playerController tghandler.PlayerController,
	messageProcessor tghandler.ButtonProvider,
) error {
	var (
		botHandler = tghandler.New(
			playerController,
			messageProcessor,
		)
	)

	tbot, err := tgbot.New(token, logger.FromContext(ctx))
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	makeRoutes(tbot, botHandler)

	if err := tbot.Start(ctx); err != nil {
		return fmt.Errorf("bot start: %w", err)
	}

	return nil
}
