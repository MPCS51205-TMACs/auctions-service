package domain

import (
	"time"
)

type Item struct {
	ItemId            string
	SellerUserId      string
	StartTime         time.Time
	EndTime           time.Time
	StartPriceInCents int64 // to avoid floating point errors, store money as cents (int); e.g. 7200 = $72.00
}

func NewItem(itemId, sellerUserId string, startTime, endTime time.Time, startPriceInCents int64) *Item {
	return &Item{
		ItemId:            itemId,
		SellerUserId:      sellerUserId,
		StartTime:         startTime.UTC(), // represent time in UTC
		EndTime:           endTime.UTC(),
		StartPriceInCents: startPriceInCents,
	}
}
