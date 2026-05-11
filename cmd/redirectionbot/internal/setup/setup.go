package setup

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/redirectionbot/internal/app"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/redirectionbot/internal/config"
)

func StartBot(ctx context.Context, cfg config.Config) error {
	if err := app.Start(
		ctx,
		cfg.Bot,
		cfg.RedirectionLink,
	); err != nil {
		return fmt.Errorf("app start: %w", err)
	}

	return nil
}
