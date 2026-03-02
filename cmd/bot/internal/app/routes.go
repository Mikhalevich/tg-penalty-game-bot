package app

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tghandler"
)

func makeRoutes(tbot *tgbot.TGBot, handler *tghandler.TGHandler) {
	tbot.AddTextCommand("/start", handler.Start)

	tbot.AddMenuCommand("/change_name", "chanage your name", handler.ChangeName)
	tbot.AddMenuCommand("/find_online_game", "find online player for game", handler.FindGame)

	tbot.AddDefaultHandler(handler.DefaultHandler)
	tbot.AddDefaultCallbackQueryHander(handler.DefaultCallbackQuery)
}
