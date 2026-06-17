package tournament

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Player struct {
	ID          player.ID
	IsBot       bool
	ChatID      msginfo.ChatID
	DisplayName string
}
