package domain

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestBid(t *testing.T) {

	time1 := time.Date(2014, 2, 4, 00, 00, 00, 0, time.UTC)
	time2 := time.Date(2014, 2, 4, 00, 00, 00, 0, time.UTC)    // same as time1
	time3 := time.Date(2014, 2, 4, 00, 00, 00, 1000, time.UTC) // 1 microsecond after
	time4 := time.Date(2014, 2, 4, 00, 00, 01, 0, time.UTC)    // 1 sec after
	// time5 := time.Date(2014, 2, 4, 00, 01, 00, 0, time.UTC)    // 1 min after
	// time6 := time.Date(2014, 2, 4, 01, 00, 00, 0, time.UTC)    // 1 hr after
	// time7 := time.Date(2014, 2, 5, 00, 00, 00, 0, time.UTC)    // 1 day after
	// time8 := time.Date(2014, 3, 4, 00, 00, 00, 0, time.UTC)    // 1 month after
	// time9 := time.Date(2015, 2, 4, 00, 00, 00, 0, time.UTC)    // 1 year after

	bid1 := *NewBid("101", "20", "asclark109", time1, int64(300), true)
	bid2 := *NewBid("102", "20", "mcostigan9", time2, int64(300), true)
	bid3 := *NewBid("103", "20", "katharine2", time3, int64(400), true)
	bid4 := *NewBid("104", "20", "katharine2", time4, int64(10), true)
	// bid5 := *NewBid("105", "katharine2", time5, int64(1000))
	// bid6 := *NewBid("106", "katharine2", time6, int64(1000))
	// bid7 := *NewBid("107", "katharine2", time7, int64(900))
	// bid8 := *NewBid("108", "mcostigan9", time8, int64(950))
	// bid9 := *NewBid("109", "mcostigan9", time9, int64(975))

	// fmt.Println("bid1: ", bid1)
	// fmt.Println("bid2: ", bid2)
	// fmt.Println("bid3: ", bid3)
	// fmt.Println("bid4: ", bid4)
	// fmt.Println("bid5: ", bid5)
	// fmt.Println("bid6: ", bid6)
	// fmt.Println("bid7: ", bid7)
	// fmt.Println("bid8: ", bid8)
	// fmt.Println("bid9: ", bid9)

	var tests = []struct {
		bidA     Bid
		bidB     Bid
		expected bool
	}{
		{bid2, bid1, false}, // bids w equal amount w equal time returns false arbitrarily (bid.outbids(otherbid) == false)
		{bid3, bid2, true},  // bids that have a later time and higher amount win (bid.outbids(otherbid) == true); 1 micrsecond later
		{bid4, bid3, false}, // bids that have a later time but lesser amount do not win (bid.outbids(otherbid) == false)
		{bid1, bid4, false}, // IMPORTANT bids that have an earlier time never outbid a bid that comes later...though this comparison would likely never happen
	}

	for num, test := range tests {
		bidA := test.bidA
		bidB := test.bidB
		testname := fmt.Sprintf("T=%v; %v.Outbids(%v)", num, bidA, bidB)
		t.Run(testname, func(t *testing.T) {
			result := bidA.Outbids(&bidB)
			if result != test.expected {
				t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "bidA.Outbids(&bidB)", strconv.FormatBool(test.expected), strconv.FormatBool(result))
			}
		})
	}
}
