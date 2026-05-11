package app

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/redirectionbot/internal/config"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

func Start(
	ctx context.Context,
	botCfg config.Bot,
	redirectionLink string,
) error {
	opts := []bot.Option{
		bot.WithSkipGetMe(),
		bot.WithDefaultHandler(makeRedirectionHandler(redirectionLink)),
	}

	botAPI, err := bot.New(
		botCfg.Token,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	botAPI.Start(ctx)

	return nil
}

func makeRedirectionHandler(link string) bot.HandlerFunc {
	return func(ctx context.Context, botAPI *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		if _, err := botAPI.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("bot was moved to %s", link),
		}); err != nil {
			logger.FromContext(ctx).
				WithError(err).
				Error("send redirection message")
		}
	}
}
