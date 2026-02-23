package button

type InlineKeyboardButton struct {
	ID      ID
	Caption string
}

type InlineKeyboardButtonRow []InlineKeyboardButton

func InlineRow(buttons ...InlineKeyboardButton) InlineKeyboardButtonRow {
	return buttons
}
