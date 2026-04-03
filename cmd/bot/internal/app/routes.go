package app

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tghandler"
)

func makeRoutes(tbot *tgbot.TGBot, handler *tghandler.TGHandler) {
	tbot.AddTextCommand("start", handler.Start)

	tbot.AddMenuCommand("change_name", "chanage your name", handler.ChangeName)
	tbot.AddMenuCommand("find_online_game", "find online player for game", handler.FindGame)
	tbot.AddMenuCommand("play_with_bot", "start game with bot", handler.PlayWithBot)
	tbot.AddMenuCommand("create_game_link", "start game by share link", handler.CreateGameLink)
	tbot.AddMenuCommand("leave_game", "leave current game", handler.LeaveGame)
	tbot.AddMenuCommand("leaderboard", "show leaderboard", handler.Leaderboard)

	tbot.AddDefaultHandler(handler.DefaultHandler)
	tbot.AddDefaultCallbackQueryHander(handler.DefaultCallbackQuery)
}
