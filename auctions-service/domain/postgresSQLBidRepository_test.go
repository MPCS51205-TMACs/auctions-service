package domain

import (
	"fmt"
	"testing"
	"time"
)

// note: these tests currently use the production database.
// they all pass if no errors occur. tests / code should be refactored
// to allow genuine tests to run without affecting production database.

func TestGetBidSQL(t *testing.T) {
	bidRepo := NewPostgresSQLBidRepository(true)
	nextid := bidRepo.NextBidId()
	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid(nextid, "201", "asclark", timeReceived, 4000, true) // $40

	bidRepo.SaveBid(bid1)

	resultbid := bidRepo.GetBid(nextid)

	if resultbid == nil {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%v", "bidRepo.GetBid()", "non-nil bid", resultbid)
	}

	result := resultbid.BidId
	expected := bid1.BidId

	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "bidRepo.GetBid()", expected, result)
	}

	bidRepo.DeleteBid(bid1.BidId) // to create idempotence
}

func TestSaveBidSQL(t *testing.T) {

	bidRepo := NewPostgresSQLBidRepository(true)

	nextid := bidRepo.NextBidId()
	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid(nextid, "201", "asclark", timeReceived, 4000, true) // $40

	bidRepo.SaveBid(bid1)

	resultbid := bidRepo.GetBid(nextid)

	if resultbid == nil {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%v", "bidRepo.GetBid()", "non-nil bid", resultbid)
	}

	result := resultbid.BidId
	expected := bid1.BidId

	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "bidRepo.GetBid()", expected, result)
	}

	bid1.active = false // make edit to bid1

	bidRepo.SaveBid(bid1) // save edited version

	resultbid = bidRepo.GetBid(nextid)

	if resultbid == nil {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%v", "bidRepo.GetBid()", "non-nil bid", resultbid)
	}

	if resultbid.active {
		t.Error("expected database to save changes to object (active -> false), but object got saved with active == true")
	}

	bidRepo.DeleteBid(bid1.BidId) // to create idempotence

}

func TestGetBidsByUserIdSQL(t *testing.T) {
	bidRepo := NewPostgresSQLBidRepository(true)
	nextid := bidRepo.NextBidId()
	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid(nextid, "201", "bmarcus1010101", timeReceived, 4000, true) // $40
	bidRepo.SaveBid(bid1)

	bids := bidRepo.GetBidsByUserId("bmarcus1010101") // assumes userid bmarcus1010101 not actually userid in production db
	if len(*bids) != 1 {
		t.Error("expected database to retreive 1 bid from database for user; instead got: ", len(*bids))
	}

	bidRepo.DeleteBid(bid1.BidId) // to create idempotence
}

func TestGetBidsByItemIdSQL(t *testing.T) {
	bidRepo := NewPostgresSQLBidRepository(true)
	nextid := bidRepo.NextBidId()
	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid(nextid, "CRAZYLONGITEMID", "bmarcus1010101", timeReceived, 4000, true) // $40
	bidRepo.SaveBid(bid1)

	bids := bidRepo.GetBidsByItemId("CRAZYLONGITEMID") // assumes itemid CRAZYLONGITEMID not actually itemid in production db
	if len(*bids) != 1 {
		t.Error("expected database to retreive 1 bid from database for item; instead got: ", len(*bids))
		for _, bid := range *bids {
			fmt.Println(bid)
		}
	}

	bidRepo.DeleteBid(bid1.BidId) // to make test idempotent

	bids = bidRepo.GetBidsByItemId("CRAZYLONGITEMID") // assumes itemid CRAZYLONGITEMID not actually itemid in production db

	if len(*bids) != 0 {
		t.Error("expected database to retreive 0 bids from database for item; instead got: ", len(*bids))
		for _, bid := range *bids {
			fmt.Println(bid)
		}
	}
}

func TestNextBidIdSQL(t *testing.T) {
	bidRepo := NewPostgresSQLBidRepository(true)
	nextID := bidRepo.NextBidId()
	fmt.Println(nextID) // not a real test just auto-passes
}

func TestSaveBidsSQL(t *testing.T) {
	bidRepo := NewPostgresSQLBidRepository(true)
	nextid1 := bidRepo.NextBidId()
	nextid2 := bidRepo.NextBidId()
	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid(nextid1, "CRAZYLONGITEMID", "bmarcus1010101", timeReceived, 4000, true) // $40
	bid2 := NewBid(nextid2, "CRAZYLONGITEMID", "bmarcus1010101", timeReceived, 4000, true) // $40

	bids := []*Bid{bid1, bid2}
	bidRepo.SaveBids(&bids)

	resultBids := bidRepo.GetBidsByItemId("CRAZYLONGITEMID")

	if len(*resultBids) != 2 {
		t.Error("expected database to retreive 2 bids from database for itemid; instead got: ", len(*resultBids))
		for _, bid := range *resultBids {
			fmt.Println(bid)
		}
	}

	for _, bid := range bids {
		bidRepo.DeleteBid(bid.BidId)
	}

	resultBids = bidRepo.GetBidsByItemId("CRAZYLONGITEMID")

	if len(*resultBids) != 0 {
		t.Error("expected database to retreive 0 bids from database for itemid; instead got: ", len(*resultBids))
		for _, bid := range *resultBids {
			fmt.Println(bid)
		}
	}

}
