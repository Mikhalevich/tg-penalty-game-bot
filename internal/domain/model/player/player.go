package player

import (
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

type ID int

func (id ID) Int() int {
	return int(id)
}

func IDFromInt(id int) ID {
	return ID(id)
}

type Player struct {
	ID                    ID
	ChatID                msginfo.ChatID
	DisplayName           string
	CreatedAt             time.Time
	IsChangeNameTriggered bool
	NameChangedAt         time.Time
	GameStatus            GameStatus
	GameStatusChangedAt   time.Time
	CurrentGameID         string
	Score                 int
}

func (p Player) IsInGame(gameID string) bool {
	return p.GameStatus == GameStatusInGame &&
		p.CurrentGameID == gameID
}
