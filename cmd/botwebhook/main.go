package main

import (
	"context"
	"encoding/json"
	"flag"
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
	var (
		isGet    = flag.Bool("get", false, "receive info about webhook")
		isSet    = flag.Bool("set", false, "set webhook")
		isRemove = flag.Bool("remove", false, "remove webhook by bot token")
	)

	var cfg Config
	if err := application.LoadConfig(&cfg); err != nil {
		logger.StdLogger().WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	botAPI, err := bot.New(cfg.Token, bot.WithSkipGetMe())
	if err != nil {
		logger.StdLogger().WithError(err).Error("initialization bot api")
		os.Exit(1)
	}

	if *isGet {
		if err := getWebHookInfo(context.Background(), botAPI); err != nil {
			logger.StdLogger().WithError(err).Error("failed to get webhook info")
			os.Exit(1)
		}

		return
	}

	if *isSet {
		if err := setWebHook(context.Background(), botAPI, cfg); err != nil {
			logger.StdLogger().WithError(err).Error("failed to set webhook")
			os.Exit(1)
		}

		return
	}

	if *isRemove {
		if err := removeWebHook(context.Background(), botAPI); err != nil {
			logger.StdLogger().WithError(err).Error("failed to remove webhook")
			os.Exit(1)
		}

		return
	}

	logger.StdLogger().Info("you need to specify --get or --set or --remove flag")

	os.Exit(1)
}

func setWebHook(ctx context.Context, botAPI *bot.Bot, cfg Config) error {
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

func removeWebHook(ctx context.Context, botAPI *bot.Bot) error {
	if _, err := botAPI.SetWebhook(
		ctx,
		&bot.SetWebhookParams{
			URL: "",
		},
	); err != nil {
		return fmt.Errorf("remove webhook: %w", err)
	}

	return nil
}

func getWebHookInfo(ctx context.Context, botAPI *bot.Bot) error {
	info, err := botAPI.GetWebhookInfo(ctx)
	if err != nil {
		return fmt.Errorf("get webhook info: %w", err)
	}

	buf, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	//nolint:forbidigo
	fmt.Println(string(buf))

	return nil
}
