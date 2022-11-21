package common

import (
	"testing"
	"time"
)

func TestInterpretTime(t *testing.T) {
	// endtime := time.Date(2014, 2, 4, 01, 30, 00, 0, time.UTC)            // 30 min later
	// item := NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	// auction1 := NewAuction(item, nil, nil)

	time := time.Date(2022, 12, 1, 15, 15, 00, 0, time.UTC)
	myTimeStr := "2022-12-01 15:15:00.000000"

	result, _ := InterpretTimeStr(myTimeStr)
	expected := time

	if *result != expected {
		t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", "bidRepo.GetBid()", expected, result)
	}

}
