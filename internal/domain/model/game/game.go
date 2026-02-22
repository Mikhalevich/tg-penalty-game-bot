package game

import (
	"time"
)

type GameID string

func (id GameID) String() string {
	return string(id)
}

func GameIDFromString(id string) GameID {
	return GameID(id)
}

type Game struct {
	GameID          GameID
	CreatedAt       time.Time
	Status          GameStatus
	StatusChagnedAt time.Time
}
