package messagesender

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
)

var (
	_ messageprocessor.Sender          = (*messageSender)(nil)
	_ messageprocessor.MarkdownEscaper = (*messageSender)(nil)
)

type messageSender struct {
	bot *bot.Bot
}

func New(bot *bot.Bot) *messageSender {
	return &messageSender{
		bot: bot,
	}
}

func makeButtonsMarkup(rows ...button.InlineKeyboardButtonRow) models.ReplyMarkup {
	if len(rows) == 0 {
		return nil
	}

	keyboard := make([][]models.InlineKeyboardButton, 0, len(rows))

	for _, row := range rows {
		buttonRow := make([]models.InlineKeyboardButton, 0, len(row))

		for _, b := range row {
			buttonRow = append(buttonRow, models.InlineKeyboardButton{
				Text:         b.Caption,
				CallbackData: b.ID.String(),
			})
		}

		keyboard = append(keyboard, buttonRow)
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}
