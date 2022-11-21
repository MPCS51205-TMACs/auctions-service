package domain

import (
	"time"
)

type Finalization struct {
	TimeReceived time.Time
}

func NewFinalization(timeReceived time.Time) *Finalization {
	return &Finalization{
		timeReceived.UTC(),
	}
}
