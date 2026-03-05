package timeprovider

import (
	"time"
)

type TimeProvider struct {
}

func New() *TimeProvider {
	return &TimeProvider{}
}

func (tp *TimeProvider) Now() time.Time {
	return time.Now()
}
