package outboxmsg

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

type Message struct {
	msginfo.Message

	ID int
}
