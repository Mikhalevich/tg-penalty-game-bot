package player

import (
	"time"
)

type Position struct {
	ID          ID
	DisplayName string
	Score       int
	UpdatedAt   time.Time
	Position    int
}
