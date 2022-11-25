package domain

import "time"

type AuctionRepository interface {
	GetAuction(itemId string) *Auction
	GetAuctions(leftBound time.Time, rightBound time.Time) []*Auction
	SaveAuction(auctionToSave *Auction)
	NewAuction(Item *Item, Bids *[]*Bid, Cancellation *Cancellation, SentStartSoonAlert, SentEndSoonAlert bool, Finalization *Finalization) *Auction
	NumAuctionsSaved() int
}
