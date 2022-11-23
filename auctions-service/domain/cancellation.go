package domain

import (
	"auctions-service/common"
	"encoding/json"
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

func (cancellation *Cancellation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TimeReceived string `json:"time_received"`
	}{
		TimeReceived: common.TimeToSQLTimestamp6(cancellation.TimeReceived),
	})
}
