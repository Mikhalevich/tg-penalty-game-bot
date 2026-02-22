package timeprovider

import (
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playercontroller"
)

var (
	_ playercontroller.TimeProvider = (*TimeProvider)(nil)
)

type TimeProvider struct {
}

func New() *TimeProvider {
	return &TimeProvider{}
}

func (tp *TimeProvider) Now() time.Time {
	return time.Now()
}
