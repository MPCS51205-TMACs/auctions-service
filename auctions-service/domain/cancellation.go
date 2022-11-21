package domain

import (
	"time"
)

type Cancellation struct {
	TimeReceived time.Time
}

func NewCancellation(timeReceived time.Time) *Cancellation {
	return &Cancellation{
		timeReceived.UTC(),
	}
}
