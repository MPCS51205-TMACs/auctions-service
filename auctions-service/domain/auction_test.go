package domain

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestProcessBid(t *testing.T) {
	startime := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)             // 30 min later
	item1 := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	item2 := NewItem("102", "asclark109", startime, endtime, int64(2000)) // $20 start price
	auction1 := NewAuction(item1, nil, nil, false, false, nil)            // will go to completion
	auction2 := NewAuction(item2, nil, nil, false, false, nil)            // will get cancelled halfway through

	time1 := startime.Add(-time.Duration(30) * time.Minute)     // 30 min before auction start;   $0.12; ignored
	time2 := startime.Add(-time.Duration(1) * time.Microsecond) // 1 microsecond before start;  $400.00; ignored
	time3 := startime                                           // on auction start;              $0.12; ignored because start price too low
	time4 := startime.Add(time.Duration(1) * time.Microsecond)  // 1 microsecond later;          $20.00; new top bid
	time5 := time4.Add(time.Duration(1) * time.Minute)          // 1 minute later;               $25.00; new top bid
	time6 := time5.Add(time.Duration(1) * time.Minute)          // 1 minute later;               $22.25; ignored because amount too low
	time7 := endtime                                            // exactly at the end;           $30.00; new top bid
	time8 := endtime.Add(time.Duration(1) * time.Microsecond)   // 1 microsecond after end;     $300.00; ignored because after auction end

	amount1 := int64(12)
	amount2 := int64(40000)
	amount3 := int64(12)
	amount4 := int64(2000)
	amount5 := int64(2500)
	amount6 := int64(2225)
	amount7 := int64(3000)
	amount8 := int64(30000)

	times := []time.Time{
		time1,
		time2,
		time3,
		time4,
		time5,
		time6,
		time7,
		time8,
	}

	amounts := []int64{
		amount1,
		amount2,
		amount3,
		amount4,
		amount5,
		amount6,
		amount7,
		amount8,
	}

	var bidsForAuction1 []*Bid = make([]*Bid, len(amounts))
	var bidsForAuction2 []*Bid = make([]*Bid, len(amounts))

	auction2.Cancel(time6) // cancel auction 2 after a few bids have come in;
	// this should cause the final bid comming in when the auction ends to be ignored;

	fmt.Println("Auction 1 processing bids")
	idx := 1
	for i := 0; i < len(amounts); i++ {
		bidsForAuction1[i] = NewBid(fmt.Sprint(idx), "101", "asclark", times[i], amounts[i], true)
		auction1.ProcessNewBid(bidsForAuction1[i])
		idx++
	}

	fmt.Println("Auction 2 processing bids")
	idx = 1
	for i := 0; i < len(amounts); i++ {
		bidsForAuction2[i] = NewBid(fmt.Sprint(idx), "102", "asclark", times[i], amounts[i], true)
		auction2.ProcessNewBid(bidsForAuction2[i])
		idx++
	}

	// expect the very last bid placed in auction1 (bid7) to be highest (legal) Bid
	result := auction1.GetHighestActiveBid().BidId
	expected := "7"
	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.GetStateAtTime(currTime)", expected, result)
	}

	// in Auction2 we don't really care about who the highest bid is because it was cancelled
	// but to confirm this Auction ignored incoming bids after it was cancelled, confirm that
	// bid5 was the last bid placed (i.e. the highest bid)
	// expect the very last bid placed in auction1 (bid7) to be highest (legal) Bid
	result = auction2.GetHighestActiveBid().BidId
	expected = "5"
	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.GetStateAtTime(currTime)", expected, result)
	}

}

func TestGetHighestActiveBid(t *testing.T) {
	startime := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	auction1 := NewAuction(item, nil, nil, false, false, nil)

	timeBidsReceived := startime.Add(time.Duration(15) * time.Second) // 15 seconds into auctions tart

	amounts := []int64{400, 12, 32545, 453, 351, 534, 684, 651, 6, 516}
	var bids []*Bid = make([]*Bid, len(amounts))

	var maxBid Bid
	for i := 0; i < len(amounts); i++ {
		bids[i] = NewBid("40", "101", "104", timeBidsReceived, amounts[i], true)
		if i == 0 {
			maxBid = *bids[0]
		} else {
			if bids[i].AmountInCents > maxBid.AmountInCents {
				maxBid = *bids[i]
			}
		}
	}

	// confirm no highest active bid

	var expected *Bid = nil // should be nil
	var result *Bid
	if result = auction1.GetHighestActiveBid(); result != nil {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "auction.GetHighestActiveBid()", expected, result)
	}

	// add a bid, confirm bid is highest active bid

	firstbid := NewBid("40", "101", "mary", timeBidsReceived, int64(2000), true)
	auction1.addBid(firstbid)
	expected = firstbid // should be firstbid (not nil)
	if result = auction1.GetHighestActiveBid(); result == nil {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "auction.GetHighestActiveBid()", expected, result)
	}

	// now deactivate that user's bid(s), and confirm there is no highest active bid

	auction1.DeactivateUserBids("mary", timeBidsReceived)
	expected = nil // should be back to nil
	if result = auction1.GetHighestActiveBid(); result != nil {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "auction.GetHighestActiveBid()", expected, result)
	}

	// now re-eactivate that user's bid(s), and confirm the original bid is the highest active bid

	auction1.ActivateUserBids("mary", timeBidsReceived)
	expected = firstbid // should be back to firstbid (not nil)
	if result = auction1.GetHighestActiveBid(); result == nil {
		t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "auction.GetHighestActiveBid()", expected, result)
	}

	// now, add a bunch of bids active bids; confirm highest bidder is correct

	for _, bid := range bids {
		auction1.addBid(bid)
	}

	result = auction1.GetHighestActiveBid()
	expected = &maxBid

	if result.BidId != expected.BidId {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.GetHighestActiveBid()", expected.BidId, result.BidId)
	}

}

func TestAddBidHasBid(t *testing.T) {
	startime := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	auction1 := NewAuction(item, nil, nil, false, false, nil)

	timeBid1Received := startime.Add(time.Duration(15) * time.Minute) // 10 min after auction start
	timeBid2Received := startime.Add(time.Duration(16) * time.Minute) // 11 min after auction start
	bid1 := NewBid("40", "101", "104", timeBid1Received, 4000, true)  // $40
	bid2 := NewBid("40", "101", "104", timeBid2Received, 5000, true)

	// confirm auction doesn't have the bid
	result := auction1.hasBid(bid1.BidId)
	expected := false
	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.hasBid(bidId string)", strconv.FormatBool(expected), strconv.FormatBool(result))
	}

	// add bid1
	auction1.addBid(bid1)

	// confirm auction has bid1
	result = auction1.hasBid(bid1.BidId)
	expected = true
	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.hasBid(bidId string)", strconv.FormatBool(expected), strconv.FormatBool(result))
	}

	// add bid2
	auction1.addBid(bid1)

	// confirm auction has bid1 and bid2
	result = auction1.hasBid(bid1.BidId) && auction1.hasBid(bid2.BidId)
	expected = true
	if result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.hasBid(bidId string)", strconv.FormatBool(expected), strconv.FormatBool(result))
	}

}

func TestStopAuction(t *testing.T) {
	startime := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	auction1 := NewAuction(item, nil, nil, false, false, nil)
	auction2 := NewAuction(item, nil, nil, false, false, nil)
	bid := NewBid("100", "101", "clark", startime, int64(3000), true)
	bids := make([]*Bid, 1)
	bids[0] = bid
	auction2wbid := NewAuction(item, &bids, nil, false, false, nil) // have this auction have an active bid
	auction3 := NewAuction(item, nil, nil, false, false, nil)
	auction4 := NewAuction(item, nil, nil, false, false, nil) // have this auction already be canceled
	auction4.Cancel(startime)
	auction5 := NewAuction(item, nil, nil, false, false, nil) // have this auction be finalized
	auction5.Finalize(endtime.Add(time.Duration(1) * time.Minute))

	stopTime1 := time.Date(2014, 2, 4, 00, 00, 00, 0, time.UTC) // before auction starts
	stopTime2 := time.Date(2014, 2, 4, 01, 10, 00, 0, time.UTC) // while auction going
	stopTime3 := time.Date(2014, 2, 4, 01, 40, 00, 0, time.UTC) // after auction over

	var tests = []struct {
		auction  *Auction
		stopTime time.Time
		expected bool
	}{
		{auction1, stopTime1, true},      // can stop if auction hasn't started (pending)
		{auction2, stopTime2, true},      // can stop if auction is active and no bids! (active, no bids)
		{auction2wbid, stopTime2, false}, // can't stop if auction is active and there is an active bid! (active, has bid)
		{auction3, stopTime3, false},     // can't cancel after auction over (over)
		{auction4, stopTime2, false},     // can't cancel auction if canceled (canceled)
		{auction4, stopTime2, false},     // can't cancel auction if already canceld (finalized)
	}

	for num, test := range tests {
		auction := test.auction
		stopTime := test.stopTime
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			result := auction.Cancel(stopTime)
			if result != test.expected {
				t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.Stop()", strconv.FormatBool(test.expected), strconv.FormatBool(result))
			}
		})
	}

}

func TestCancelAuction(t *testing.T) {
	startime := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	auction1 := NewAuction(item, nil, nil, false, false, nil)
	auction2 := NewAuction(item, nil, nil, false, false, nil)
	auction3 := NewAuction(item, nil, nil, false, false, nil)
	auction4 := NewAuction(item, nil, nil, false, false, nil) // have this auction already be canceled
	auction4.Cancel(startime)
	auction5 := NewAuction(item, nil, nil, false, false, nil) // have this auction be finalized
	auction5.Finalize(endtime.Add(time.Duration(1) * time.Minute))

	cancelTime1 := time.Date(2014, 2, 4, 00, 00, 00, 0, time.UTC) // before auction starts
	cancelTime2 := time.Date(2014, 2, 4, 01, 10, 00, 0, time.UTC) // while auction going
	cancelTime3 := time.Date(2014, 2, 4, 01, 40, 00, 0, time.UTC) // after auction over

	var tests = []struct {
		auction    *Auction
		cancelTime time.Time
		expected   bool
	}{
		{auction1, cancelTime1, true},  // can cancel if auction hasn't started (pending)
		{auction2, cancelTime2, true},  // can cancel if auction is active (active)
		{auction3, cancelTime3, false}, // can't cancel after auction over (over)
		{auction4, cancelTime2, false}, // can't cancel auction if canceled (canceled)
		{auction4, cancelTime2, false}, // can't cancel auction if already canceld (finalized)
	}

	for num, test := range tests {
		auction := test.auction
		cancelTime := test.cancelTime
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			result := auction.Cancel(cancelTime)
			if result != test.expected {
				t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.Cancel()", strconv.FormatBool(test.expected), strconv.FormatBool(result))
			}
		})
	}

}

func TestFinalizeAuction(t *testing.T) {
	startime := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	auction1 := NewAuction(item, nil, nil, false, false, nil)
	auction2 := NewAuction(item, nil, nil, false, false, nil)
	auction3 := NewAuction(item, nil, nil, false, false, nil)
	auction4 := NewAuction(item, nil, nil, false, false, nil) // have this auction already be canceled
	auction4.Cancel(startime)
	auction5 := NewAuction(item, nil, nil, false, false, nil) // have this auction be finalized
	auction5.Finalize(endtime.Add(time.Duration(1) * time.Minute))

	finalizetime1 := time.Date(2014, 2, 4, 00, 00, 00, 0, time.UTC) // before auction starts
	finalizetime2 := time.Date(2014, 2, 4, 01, 10, 00, 0, time.UTC) // while auction going
	finalizetime3 := time.Date(2014, 2, 4, 01, 40, 00, 0, time.UTC) // after auction over

	var tests = []struct {
		auction      *Auction
		finalizeTime time.Time
		expected     bool
	}{
		{auction1, finalizetime1, false}, // can't finalize if auction hasn't started (pending)
		{auction2, finalizetime2, false}, // can't finalize if auction is active (active)
		{auction3, finalizetime3, true},  // can finalize after auction over (over)
		{auction4, finalizetime2, true},  // can finalize auction if canceled (canceled)
		{auction4, finalizetime2, false}, // can't finalize auction if already finalized (finalized)
	}

	for num, test := range tests {
		auction := test.auction
		finalizeTime := test.finalizeTime
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			result := auction.Finalize(finalizeTime)
			if result != test.expected {
				t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.Finalize()", strconv.FormatBool(test.expected), strconv.FormatBool(result))
			}
		})
	}

}

func TestGetStateAtTime(t *testing.T) {

	startime := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // end: 30 min later
	finalizetime := endtime.Add(time.Duration(20) * time.Minute)         // finalized 20 min after end
	item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price

	auction := NewAuction(item, nil, nil, false, false, nil) // uncanceled auction
	auction.Finalize(finalizetime)

	auction_cancelled := NewAuction(item, nil, nil, false, false, nil) // auction cancelled 15 min into auction start (30 min long auction)
	canceltime := time.Date(2014, 2, 4, 01, 15, 00, 0, time.UTC)
	finalizetime2 := canceltime.Add(time.Duration(5) * time.Minute) // finalized 5 min after cancellation
	auction_cancelled.Cancel(canceltime)
	auction_cancelled.Finalize(finalizetime2)

	time1 := time.Date(2014, 2, 4, 00, 30, 00, 0, time.UTC)     // 30 min before auction start
	time2 := startime.Add(-time.Duration(1) * time.Microsecond) // 1 microsecond before start
	time3 := startime                                           // on auction start
	time4 := startime.Add(time.Duration(15) * time.Minute)      // halfway thru auction
	time5 := endtime                                            // on auction end
	time6 := endtime.Add(time.Duration(1) * time.Microsecond)   // 1 microsecond after end
	time7 := endtime.Add(time.Duration(30) * time.Minute)       // 30 minutes after end

	var tests = []struct {
		auction  *Auction
		currTime time.Time
		expected AuctionState
	}{
		{auction, time1, PENDING},
		{auction, time2, PENDING},
		{auction, time3, ACTIVE},
		{auction, time4, ACTIVE},
		{auction, time5, ACTIVE},
		{auction, time6, OVER},
		{auction, time7, FINALIZED},
		{auction_cancelled, time1, PENDING},
		{auction_cancelled, time2, PENDING},
		{auction_cancelled, time3, ACTIVE},
		{auction_cancelled, time4, CANCELED},
		{auction_cancelled, time5, FINALIZED},
		{auction_cancelled, time6, FINALIZED},
	}

	for num, test := range tests {
		auction := test.auction
		currTime := test.currTime
		expectedState := test.expected
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			calcdState := auction.getStateAtTime(currTime)
			if calcdState != expectedState {
				t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "auction.GetStateAtTime(currTime)", expectedState, calcdState)
			}
		})
	}
}

func TestOverlapsWith(t *testing.T) {

	startime := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)
	endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price

	auction := NewAuction(item, nil, nil, false, false, nil) // uncanceled auction

	time1 := time.Date(2014, 2, 4, 00, 30, 00, 0, time.UTC)     // 30 min before auction start
	time2 := startime.Add(-time.Duration(1) * time.Microsecond) // 1 microsecond before start
	time3 := startime                                           // on auction start
	time4 := startime.Add(time.Duration(15) * time.Minute)      // halfway thru auction
	time5 := endtime                                            // on auction end
	time6 := endtime.Add(time.Duration(time.Microsecond))       // 1 microsecond after end

	var tests = []struct {
		auction  *Auction
		currTime time.Time
		expected bool
	}{
		{auction, time1, false},
		{auction, time2, false},
		{auction, time3, true},
		{auction, time4, true},
		{auction, time5, true},
		{auction, time6, false},
	}

	for num, test := range tests {
		auction := test.auction
		currTime := test.currTime
		expected := test.expected
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			result := auction.OverlapsWith(&currTime, &currTime)
			if result != expected {
				t.Errorf("\nRan:%s\nExpected:%v\nGot:%v", "auction.GetStateAtTime(currTime)", expected, result)
			}
		})
	}

}
