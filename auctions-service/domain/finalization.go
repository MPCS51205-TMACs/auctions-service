package domain

import (
	"auctions-service/common"
	"encoding/json"
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

func (finalization *Finalization) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TimeReceived string `json:"time_received"`
	}{
		TimeReceived: common.TimeToSQLTimestamp6(finalization.TimeReceived),
	})
}
