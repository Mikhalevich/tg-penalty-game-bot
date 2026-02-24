package game

import (
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(id string) ID {
	return ID(id)
}

type Game struct {
	ID              ID
	CreatedAt       time.Time
	Status          GameStatus
	StatusChagnedAt time.Time
	State           State
	StateVersion    int
}

type Player struct {
	ID             player.ID
	DisplayName    string
	ShotsAvailable int
	GoalsScored    int
}

type RoundShot struct {
	PlayerID player.ID
	Side     ShotSide
}

type Round struct {
	Defend RoundShot
	Attack RoundShot
	IsGoal bool
}

type State struct {
	Players []Player
	Rounds  []Round
}
