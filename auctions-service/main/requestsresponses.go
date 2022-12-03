package main

import (
	"auctions-service/domain"
	"time"
)

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

type RequestCreateNewBid struct {
	ItemId        string `json:"itemid"`
	BidderUserId  string `json:"bidderuserid"`
	AmountInCents int64  `json:"amountincents"`
}

type ResponseCreateNewBid struct {
	Msg string `json:"message"`
	// WasNewTopBid bool   `json:"was_new_top_bid"`
}

type RequestProcessNewBid struct {
	ItemId        string `json:"itemid"`
	BidderUserId  string `json:"bidderuserid"`
	AmountInCents int64  `json:"amountincents"`
	TimeReceived  string `json:"timereceived"`
}

// type ResponseProcessNewBid struct {
// 	Msg          string `json:"message"`
// 	WasNewTopBid bool   `json:"was_new_top_bid"`
// }

type ResponseCreateAuction struct {
	Msg string `json:"message"`
}

type ResponseGetAuctions struct {
	ActiveAuctions []JsonAuction `json:"auctions"`
}

type JsonAuction struct {
	ItemId            string `json:"itemid"`
	SellerUserId      string `json:"selleruserid"`
	StartTime         string `json:"starttime"`
	EndTime           string `json:"endtime"`
	StartPriceInCents int64  `json:"startpriceincents"`
	TopBid            *domain.Bid
	State             domain.AuctionState `json:"state"`
}

func ExportAuction(auction *domain.Auction) *JsonAuction {
	layout := "2006-01-02 15:04:05.000000"
	return &JsonAuction{
		ItemId:            auction.Item.ItemId,
		SellerUserId:      auction.Item.SellerUserId,
		StartPriceInCents: auction.Item.StartPriceInCents,
		StartTime:         auction.Item.StartTime.Format(layout),
		EndTime:           auction.Item.EndTime.Format(layout),
		TopBid:            auction.GetHighestActiveBid(),
		State:             auction.GetStateAtTime(time.Now()),
	}
}

// rabbit events published by other contexts
type ItemCounterfeitEvent struct {
	ItemId string `json:"itemid"`
}

type ItemInapropriateEvent struct {
	ItemId string `json:"itemid"`
}

type UserUpdateEvent struct {
	UserId     string `json:"userid"`
	UserUpdate string `json:"userupdate"` // this might require its own class to parse
}

type UserDeleteEvent struct {
	UserId string `json:"userid"`
}
