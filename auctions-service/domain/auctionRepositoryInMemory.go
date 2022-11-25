package domain

import (
	"time"
)

type inMemoryAuctionRepository struct {
	auctions    []*Auction
	alertEngine AlertEngine
}

func NewInMemoryAuctionRepository(alertEngine AlertEngine) AuctionRepository {
	auctions := []*Auction{}
	return &inMemoryAuctionRepository{auctions, alertEngine}
}

func (repo *inMemoryAuctionRepository) GetAuction(itemId string) *Auction {
	for _, auction := range repo.auctions {
		if auction.Item.ItemId == itemId {
			return auction
		}
	}
	return nil
}

func (repo *inMemoryAuctionRepository) GetAuctions(leftBound time.Time, rightBound time.Time) []*Auction {
	relevantAuctions := []*Auction{}
	for _, auction := range repo.auctions {
		if auction.OverlapsWith(&leftBound, &rightBound) {
			relevantAuctions = append(relevantAuctions, auction)
		}
	}
	return relevantAuctions
}

func (repo *inMemoryAuctionRepository) SaveAuction(auctionToSave *Auction) {
	for idx, auction := range repo.auctions {
		if auction.Item.ItemId == auctionToSave.Item.ItemId {
			repo.auctions[idx] = auctionToSave // overwrite
			return
		}
	}
	// else its new
	repo.auctions = append(repo.auctions, auctionToSave)
}

func (repo *inMemoryAuctionRepository) NewAuction(Item *Item, Bids *[]*Bid, Cancellation *Cancellation, SentStartSoonAlert, SentEndSoonAlert bool, Finalization *Finalization) *Auction {
	return NewAuction(Item, Bids, Cancellation, SentStartSoonAlert, SentEndSoonAlert, Finalization, repo.alertEngine)
}

func (repo *inMemoryAuctionRepository) NumAuctionsSaved() int {
	return len(repo.auctions)
}
