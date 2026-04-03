package model

import (
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Position struct {
	ID          int       `db:"player_id"`
	DisplayName string    `db:"display_name"`
	Score       int       `db:"score"`
	UpdatedAt   time.Time `db:"updated_at"`
	Position    int       `db:"position"`
}

func (p Position) ToDomPosition() player.Position {
	return player.Position{
		ID:          player.IDFromInt(p.ID),
		DisplayName: p.DisplayName,
		Score:       p.Score,
		UpdatedAt:   p.UpdatedAt,
		Position:    p.Position,
	}
}
