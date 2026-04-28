package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/application"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

type Config struct {
	Token              string `yaml:"token" required:"true"`
	WebHookSecretToken string `yaml:"webhook_secret_token" required:"true"`
	URL                string `yaml:"url" required:"true"`
	CertificatePath    string `yaml:"certificate_path" required:"true"`
}

func main() {
	var cfg Config
	if err := application.LoadConfig(&cfg); err != nil {
		logger.StdLogger().WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	if err := setWebHook(context.Background(), cfg); err != nil {
		logger.StdLogger().WithError(err).Error("failed to set webhook")
		os.Exit(1)
	}
}

func setWebHook(ctx context.Context, cfg Config) error {
	botAPI, err := bot.New(cfg.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot api: %w", err)
	}

	file, err := os.Open(cfg.CertificatePath)
	if err != nil {
		return fmt.Errorf("open certificate file: %w", err)
	}

	if _, err := botAPI.SetWebhook(
		ctx,
		&bot.SetWebhookParams{
			URL:         cfg.URL,
			SecretToken: cfg.WebHookSecretToken,
			Certificate: &models.InputFileUpload{
				Filename: file.Name(),
				Data:     file,
			},
		},
	); err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}

	return nil
}
