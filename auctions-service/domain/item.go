package domain

import (
	"auctions-service/common"
	"encoding/json"
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

func (item *Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ItemId            string `json:"item_id"`
		SellerUserId      string `json:"seller_user_id"`
		StartTime         string `json:"start_time"`
		EndTime           string `json:"end_time"`
		StartPriceInCents int64  `json:"start_price_in_cents"`
	}{
		ItemId:            item.ItemId,
		SellerUserId:      item.SellerUserId,
		StartTime:         common.TimeToSQLTimestamp6(item.StartTime),
		EndTime:           common.TimeToSQLTimestamp6(item.EndTime),
		StartPriceInCents: item.StartPriceInCents,
	})
}
