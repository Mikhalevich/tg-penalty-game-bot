package tournament

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Player struct {
	ID            player.ID
	IsBot         bool
	BotDifficulty game.BotDifficulty
	ChatID        msginfo.ChatID
	DisplayName   string
}

func (p Player) toDomPlayer() player.Player {
	return player.Player{
		ID:          p.ID,
		ChatID:      p.ChatID,
		DisplayName: p.DisplayName,
	}
}
