package domain

import (
	"fmt"
	"testing"
	"time"
)

func TestGetBid(t *testing.T) {
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

	result := bidRepo.GetBid(bid1.BidId).BidId
	expected := bid1.BidId

	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "bidRepo.GetBid()", expected, result)
	}

	result = bidRepo.GetBid(bid3.BidId).BidId
	expected = bid3.BidId

	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "bidRepo.GetBid()", expected, result)
	}

}

func TestSaveBid(t *testing.T) {
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

	result := bidRepo.GetBid(bid1.BidId).BidId
	expected := bid1.BidId

	// confirm bid 1 saved
	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "bidRepo.GetBid()", expected, result)
	}

	// edit bid 1 and then re-save it (in practice bids are not edited really)
	// this is just to confirm changes have been saved. then re-load the bid
	// and confirm changes exist. Also confirm total number of bids have not increased
	bid1.AmountInCents = 2000
	bidRepo.SaveBid(bid1)

	result2 := bidRepo.GetBid(bid1.BidId).AmountInCents
	expected2 := bid1.AmountInCents

	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "bidRepo.GetBid().AmountInCents", expected2, result2)
	}

	result3 := len(*(bidRepo.GetBidsByItemId("201")))
	expected3 := 4 // 4 bids

	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "bidRepo.GetBid().AmountInCents", expected3, result3)
	}
}

func TestGetBidsByUserId(t *testing.T) {
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

func TestGetBidsByItemId(t *testing.T) {
	// endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	// item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	// auction1 := NewAuction(item, nil, nil)

	timeReceived := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	bid1 := NewBid("100", "201", "asclark", timeReceived, 4000, true) // $40
	bid2 := NewBid("101", "201", "mark11", timeReceived, 4000, true)  // $40
	bid3 := NewBid("102", "201", "sammy2", timeReceived, 4000, true)  // $40
	bid4 := NewBid("103", "201", "asclark", timeReceived, 4000, true) // $40
	bid5 := NewBid("104", "202", "asclark", timeReceived, 4000, true) // $40

	bidRepo := NewInMemoryBidRepository(true)

	bidRepo.SaveBid(bid1)
	bidRepo.SaveBid(bid2)
	bidRepo.SaveBid(bid3)
	bidRepo.SaveBid(bid4)
	bidRepo.SaveBid(bid5)

	foundBids := *bidRepo.GetBidsByItemId("201")
	for _, bid := range foundBids {
		fmt.Printf("%v", bid)
	}
	// fmt.Println(foundBids)
	result := len(foundBids)
	expected := 4

	// confirm bid 1 saved
	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "bidRepo.GetBid()", expected, result)
	}
}
