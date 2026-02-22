package button

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

type ChangeNamePayload struct {
	DisplayName string
}

func ChangeName(chatID msginfo.ChatID, caption, displayName string) (Button, error) {
	return createButton(chatID, caption, OperationChangeName,
		ChangeNamePayload{
			DisplayName: displayName,
		},
	)
}
