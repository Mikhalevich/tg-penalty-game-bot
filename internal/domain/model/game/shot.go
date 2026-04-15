package game

import (
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type ShotType string

const (
	ShotTypeDefend ShotType = "defend"
	ShotTypeAttack ShotType = "attack"
)

func (st ShotType) String() string {
	return string(st)
}

type ShotSide string

const (
	ShotSideNoShot ShotSide = "no_shot"
	ShotSideLeft   ShotSide = "left"
	ShotSideRight  ShotSide = "right"
	ShotSideMiddle ShotSide = "middle"
	ShotSideMiss   ShotSide = "miss"
)

func (ss ShotSide) String() string {
	return string(ss)
}

type Shot struct {
	GameID       ID
	PlayerID     player.ID
	Round        int
	Type         ShotType
	CreatedAt    time.Time
	ExpectedSide ShotSide
	ActualSide   ShotSide
	CompletedAt  time.Time
}
