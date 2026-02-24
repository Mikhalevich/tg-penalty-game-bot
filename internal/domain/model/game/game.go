package game

import (
	"time"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(id string) ID {
	return ID(id)
}

type Game struct {
	GameID          ID
	CreatedAt       time.Time
	Status          GameStatus
	StatusChagnedAt time.Time
}
