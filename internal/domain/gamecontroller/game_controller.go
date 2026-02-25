package gamecontroller

import (
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

type GameController struct {
	timeProvider TimeProvider
}

func New(
	timeProvier TimeProvider,
) *GameController {
	return &GameController{
		timeProvider: timeProvier,
	}
}
