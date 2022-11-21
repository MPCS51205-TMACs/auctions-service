package main

import "auctions-service/domain"

type ResponseGetItemsByUserId struct {
	// could optionally include the userId here as well e.g.
	// "UserId string `json:"userid"`"
	ItemIds []string `json:"itemids"`
}

// type ResponseGetItemsByUserId struct {
// 	// could optionally include the userId here as well e.g.
// 	// "UserId string `json:"userid"`"
// 	ItemIds []string `json:"itemids"`
// }

type ResponseStopAuction struct {
	Msg string `json:"message"`
}

type RequestStopAuction struct {
	// ItemId            string `json:"itemid"`
	RequesterUserId string `json:"requesteruserid"`
}

type RequestCreateAuction struct {
	ItemId            string `json:"itemid"`
	SellerUserId      string `json:"selleruserid"`
	StartTime         string `json:"starttime"`
	EndTime           string `json:"endtime"`
	StartPriceInCents int64  `json:"startpriceincents"`
}

type RequestProcessNewBid struct {
	ItemId        string `json:"itemid"`
	BidderUserId  string `json:"selleruserid"`
	AmountInCents int64  `json:"amountincents"`
}

type ResponseProcessNewBid struct {
	Msg          string `json:"message"`
	WasNewTopBid bool   `json:"was_new_top_bid"`
}

type ResponseCreateAuction struct {
	Msg string `json:"message"`
}

type ResponseGetActiveAuctions struct {
	ActiveAuctions []JsonAuction `json:"activeauctions"`
}

type JsonAuction struct {
	ItemId            string `json:"itemid"`
	SellerUserId      string `json:"selleruserid"`
	StartTime         string `json:"starttime"`
	EndTime           string `json:"endtime"`
	StartPriceInCents int64  `json:"startpriceincents"`
}

func ExportAuction(auction *domain.Auction) *JsonAuction {
	layout := "2006-01-02 15:04:05.000000"
	return &JsonAuction{
		ItemId:            auction.Item.ItemId,
		SellerUserId:      auction.Item.SellerUserId,
		StartPriceInCents: auction.Item.StartPriceInCents,
		StartTime:         auction.Item.StartTime.Format(layout),
		EndTime:           auction.Item.EndTime.Format(layout),
	}
}
