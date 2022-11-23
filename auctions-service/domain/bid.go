package domain

import (
	"auctions-service/common"
	"encoding/json"
	"time"
)

type Bid struct {
	BidId         string
	ItemId        string
	BidderUserId  string
	TimeReceived  time.Time
	AmountInCents int64 // to avoid floating point errors, store money as cents (int)
	active        bool  // represents whether the bid is "activated" or "deactivated"
}

func NewBid(bidId, itemId, bidderUserId string, timeReceived time.Time, ammountInCents int64, active bool) *Bid {
	return &Bid{
		BidId:         bidId,
		ItemId:        itemId,
		BidderUserId:  bidderUserId,
		TimeReceived:  timeReceived.UTC(), // represent time in UTC
		AmountInCents: ammountInCents,
		active:        active,
	}
}

func (bid *Bid) Outbids(otherBid *Bid) bool {
	if (bid.TimeReceived).After(otherBid.TimeReceived) {
		if bid.AmountInCents > otherBid.AmountInCents {
			return true
		}
	}
	return false
}

func (bid *Bid) Activate() bool {
	if !bid.active {
		bid.active = true
		return true
	}
	return false
}

func (bid *Bid) Deactivate() bool {
	if bid.active {
		bid.active = false
		return true
	}
	return false
}

func (bid *Bid) IsActive() bool {
	return bid.active
}

func (bid *Bid) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		BidId         string `json:"bid_id"`
		ItemId        string `json:"item_id"`
		BidderUserId  string `json:"bidder_user_id"`
		TimeReceived  string `json:"time_received"`
		AmountInCents int64  `json:"amount_in_cents"`
		Active        bool   `json:"active"`
	}{
		BidId:         bid.BidId,
		ItemId:        bid.ItemId,
		BidderUserId:  bid.BidderUserId,
		TimeReceived:  common.TimeToSQLTimestamp6(bid.TimeReceived),
		AmountInCents: bid.AmountInCents,
		Active:        bid.active,
	})
}
