package tgbot

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

type TGBot struct {
	bot              *bot.Bot
	logger           logger.Logger
	middlewares      []Middleware
	commands         []models.BotCommand
	defaultHandlerFn Handler
}

func New(token string, logger logger.Logger) (*TGBot, error) {
	tgbot := &TGBot{
		logger: logger,
	}

	botAPI, err := bot.New(
		token,
		bot.WithSkipGetMe(),
		bot.WithDefaultHandler(tgbot.makeDefaultHandler()),
	)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	tgbot.bot = botAPI

	return tgbot, nil
}
