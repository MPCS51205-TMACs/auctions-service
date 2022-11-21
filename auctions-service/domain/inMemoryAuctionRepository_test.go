package domain

import (
	"testing"
	"time"
)

func TestGetAuction(t *testing.T) {
	// endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	// item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	// auction1 := NewAuction(item, nil, nil)

	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid("100", "201", "asclark", timeReceived, 4000, true) // $40
	bid2 := NewBid("101", "201", "mark11", timeReceived, 4000, true)  // $40
	bid3 := NewBid("102", "201", "sammy2", timeReceived, 4000, true)  // $40
	bid4 := NewBid("103", "201", "asclark", timeReceived, 4000, true) // $40
	bidRepo := NewInMemoryBidRepository(true)
	bidRepo.SaveBid(bid1)
	bidRepo.SaveBid(bid2)
	bidRepo.SaveBid(bid3)
	bidRepo.SaveBid(bid4)
	bids201 := []*Bid{bid1, bid2, bid3, bid4}
	starttime := timeReceived.Add(-time.Duration(10) * time.Minute)
	endtime := starttime.Add(time.Duration(10) * time.Hour)
	item201 := NewItem("201", "sellerMike", starttime, endtime, int64(2000))
	auction := NewAuction(item201, &bids201, nil, false, false, nil)
	auctionRepo := NewInMemoryAuctionRepository()
	auctionRepo.SaveAuction(auction)

	result := auctionRepo.GetAuction(item201.ItemId)
	expected := auction

	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "auctionRepo.GetAuction()", expected, result)
	}
}

func TestSaveAuction(t *testing.T) {

	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid("100", "201", "asclark", timeReceived, 4000, true) // $40
	bid2 := NewBid("101", "201", "mark11", timeReceived, 4000, true)  // $40
	bid3 := NewBid("102", "201", "sammy2", timeReceived, 4000, true)  // $40
	bid4 := NewBid("103", "201", "asclark", timeReceived, 4000, true) // $40
	bidRepo := NewInMemoryBidRepository(true)
	bidRepo.SaveBid(bid1)
	bidRepo.SaveBid(bid2)
	bidRepo.SaveBid(bid3)
	bidRepo.SaveBid(bid4)
	bids201 := []*Bid{bid1, bid2, bid3, bid4}
	starttime := timeReceived.Add(-time.Duration(10) * time.Minute)
	endtime := starttime.Add(time.Duration(10) * time.Hour)
	item201 := NewItem("201", "sellerMike", starttime, endtime, int64(2000))
	auction := NewAuction(item201, &bids201, nil, false, false, nil)
	auctionRepo := NewInMemoryAuctionRepository()

	if auctionRepo.NumAuctionsSaved() != 0 {
		t.Errorf("\nRan:%s\nExpected:%d\nGot:%d", "auctionRepo.NumAuctionsSaved()", 1, auctionRepo.NumAuctionsSaved())
	}
	auctionRepo.SaveAuction(auction)
	if auctionRepo.NumAuctionsSaved() != 1 {
		t.Errorf("\nRan:%s\nExpected:%d\nGot:%d", "auctionRepo.NumAuctionsSaved()", 1, auctionRepo.NumAuctionsSaved())
	}

	// now re-save the same auction, and confirm there is only one auction in the repo
	auctionRepo.SaveAuction(auction) // should be idempotent
	if auctionRepo.NumAuctionsSaved() != 1 {
		t.Errorf("\nRan:%s\nExpected:%d\nGot:%d", "auctionRepo.NumAuctionsSaved()", 1, auctionRepo.NumAuctionsSaved())
	}

	timeReceived2 := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid5 := NewBid("104", "202", "asclark", timeReceived2, 4000, true) // $40
	bid6 := NewBid("105", "202", "mark11", timeReceived2, 4000, true)  // $40
	bid7 := NewBid("106", "202", "sammy2", timeReceived2, 4000, true)  // $40
	bid8 := NewBid("107", "202", "asclark", timeReceived2, 4000, true) // $40
	bidRepo.SaveBid(bid5)
	bidRepo.SaveBid(bid6)
	bidRepo.SaveBid(bid7)
	bidRepo.SaveBid(bid8)
	bids202 := []*Bid{bid5, bid6, bid7, bid8}
	starttime2 := timeReceived.Add(-time.Duration(10) * time.Minute)
	endtime2 := starttime.Add(time.Duration(10) * time.Hour)
	item202 := NewItem("202", "sellerMike", starttime2, endtime2, int64(2000))
	auction2 := NewAuction(item202, &bids202, nil, false, false, nil)

	// now save the same auction, and confirm there are two auctions saved in the repo
	auctionRepo.SaveAuction(auction2)
	if auctionRepo.NumAuctionsSaved() != 2 {
		t.Errorf("\nRan:%s\nExpected:%d\nGot:%d", "auctionRepo.NumAuctionsSaved()", 1, auctionRepo.NumAuctionsSaved())
	}

	// re-save both auctions and confirm still 2 auctions in repo
	auctionRepo.SaveAuction(auction)
	auctionRepo.SaveAuction(auction2)
	if auctionRepo.NumAuctionsSaved() != 2 {
		t.Errorf("\nRan:%s\nExpected:%d\nGot:%d", "auctionRepo.NumAuctionsSaved()", 1, auctionRepo.NumAuctionsSaved())
	}
}

func TestGetAuctions(t *testing.T) {
	// endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	// item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	// auction1 := NewAuction(item, nil, nil)

	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid("100", "201", "asclark", timeReceived, 4000, true) // $40
	bid2 := NewBid("101", "201", "mark11", timeReceived, 4000, true)  // $40
	bid3 := NewBid("102", "201", "sammy2", timeReceived, 4000, true)  // $40
	bid4 := NewBid("103", "201", "asclark", timeReceived, 4000, true) // $40

	bidRepo := NewInMemoryBidRepository(true)

	bidRepo.SaveBid(bid1)
	bidRepo.SaveBid(bid2)
	bidRepo.SaveBid(bid3)
	bidRepo.SaveBid(bid4)

	result := len(*bidRepo.GetBidsByUserId("asclark"))
	expected := 2

	// confirm bid 1 saved
	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "bidRepo.GetBid()", expected, result)
	}
}
