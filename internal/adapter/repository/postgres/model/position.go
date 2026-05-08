package model

import (
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/lbrefresher"
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

func ToDomPositions(dbPos []Position) []player.Position {
	domPos := make([]player.Position, 0, len(dbPos))

	for _, pos := range dbPos {
		domPos = append(domPos, pos.ToDomPosition())
	}

	return domPos
}

type HistoryPosition struct {
	PlayerID    int       `db:"player_id"`
	Year        int       `db:"year"`
	Month       int       `db:"month"`
	DisplayName string    `db:"display_name"`
	Score       int       `db:"score"`
	UpdatedAt   time.Time `db:"updated_at"`
	Position    int       `db:"position"`
}

func ToDBHisotoryPositions(history lbrefresher.HistoryPositions) []HistoryPosition {
	dbPositions := make([]HistoryPosition, 0, len(history.Positions))

	for _, pos := range history.Positions {
		dbPositions = append(dbPositions, HistoryPosition{
			PlayerID:    pos.ID.Int(),
			Year:        history.Year,
			Month:       history.Month,
			DisplayName: pos.DisplayName,
			Score:       pos.Score,
			UpdatedAt:   pos.UpdatedAt,
			Position:    pos.Position,
		})
	}

	return dbPositions
}
