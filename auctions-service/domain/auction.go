package domain

import (
	"log"
	"time"
)

// Enum that defines various states an auction can be in
type AuctionState string

const (
	PENDING   AuctionState = "PENDING" // has not yet started
	ACTIVE    AuctionState = "ACTIVE"  // is happening now
	CANCELED  AuctionState = "CANCELED"
	OVER      AuctionState = "OVER"      // is over (but winner has not been declared and auction has not been "archived away")
	FINALIZED AuctionState = "FINALIZED" // is over and archived away; can delete
	UNKNOWN   AuctionState = "UKNOWN"
)

type Auction struct {
	Item               *Item
	bids               []*Bid // slice of pointers to bids; new higher bids get appended on the end
	cancellation       *Cancellation
	sentStartSoonAlert bool
	sentEndSoonAlert   bool
	finalization       *Finalization
}

func NewAuction(item *Item, bids *[]*Bid, cancellation *Cancellation, sentStartSoonAlert, sentEndSoonAlert bool, finalization *Finalization) *Auction {
	if bids == nil {
		newBidsSlice := make([]*Bid, 0)
		bids = &newBidsSlice
	}
	return &Auction{
		Item:               item,
		bids:               *bids,            // nil if brand new
		cancellation:       cancellation,     // nil if brand new
		sentStartSoonAlert: sentEndSoonAlert, // false if brand new
		sentEndSoonAlert:   sentEndSoonAlert, // false if brand new
		finalization:       finalization,     // nil if brand new
	}
}

func (auction *Auction) ProcessNewBid(incomingBid *Bid) (AuctionState, bool) {
	timeBidReceived := incomingBid.TimeReceived
	stateWhenBidReceived := auction.getStateAtTime(timeBidReceived)

	// if the auction has been finalized, it is archived and we are no longer
	// considering new bids.
	if auction.HasFinalization() { // i.e. has state FINALIZED at some known point in time
		log.Printf("[Auction %s] ignoring bid. auction has been finalized.\n", auction.Item.ItemId)
		return FINALIZED, false
	}

	switch {
	case stateWhenBidReceived == PENDING:
		log.Printf("[Auction %s] ignoring bid. auction hadn't begun when bid was received.\n", auction.Item.ItemId)
		return PENDING, false
	case stateWhenBidReceived == CANCELED:
		log.Printf("[Auction %s] ignoring bid. auction was cancelled before bid was received.\n", auction.Item.ItemId)
		return CANCELED, false
	case stateWhenBidReceived == OVER:
		log.Printf("[Auction %s] ignoring bid. auction was over before bid was received.\n", auction.Item.ItemId)
		return OVER, false
	// case stateWhenBidReceived == FINALIZED: HANDLED ABOVE
	case stateWhenBidReceived == ACTIVE:
		highestActiveBid := auction.GetHighestActiveBid()
		if highestActiveBid == nil { // case: there are no active bids
			if incomingBid.AmountInCents >= auction.Item.StartPriceInCents { // bid amount must at least be start price
				log.Printf("[Auction %s] new top bid!\n", auction.Item.ItemId)
				auction.addBid(incomingBid)
				auction.alertSeller("you have a new top bid!")
				return ACTIVE, true
			} else {
				log.Printf("[Auction %s] ignoring bid. bid was under start price.\n", auction.Item.ItemId)
				return ACTIVE, false
			}
		} else { // case: auction already has at least one active bid
			if incomingBid.Outbids(highestActiveBid) {
				log.Printf("[Auction %s] new top bid!\n", auction.Item.ItemId)
				auction.addBid(incomingBid)
				auction.alertSeller("you have a new top bid!")
				auction.alertBidder("your top bid has been out-matched!", highestActiveBid)
				return ACTIVE, true
			} else {
				log.Printf("[Auction %s] ignoring bid. bid was under highest bid offer amount.\n", auction.Item.ItemId)
				return ACTIVE, false
			}
		}
	default:
		panic("see processNewBid()! couldn't process bid because didn't understand state of auction at time bid was received.")
	}
}

func (auction *Auction) addBid(bid *Bid) {
	auction.bids = append(auction.bids, bid)
}

func (auction *Auction) GetHighestActiveBid() *Bid {
	if len(auction.bids) == 0 {
		return nil
	}
	// by convention, top bids are appended at the end; so start at end and walk to the left.
	// find the first active bid
	idx := len(auction.bids) - 1
	for !auction.bids[idx].active {
		idx--
		if idx == -1 {
			return nil
		}
	}
	return auction.bids[idx]
}

func (auction *Auction) getStateAtTime(currTime time.Time) AuctionState {
	// if auction has been cancelled, then if the time was
	// after the time of the cancellation, then the state at that
	// time is cancelled
	if auction.HasFinalization() {
		if AfterOrOn(&currTime, &auction.finalization.TimeReceived) {
			return FINALIZED
		}
	}

	// if auction has been cancelled, then if the time was
	// after the time of the cancellation, then the state at that
	// time is cancelled
	if auction.HasCancellation() {
		if AfterOrOn(&currTime, &auction.cancellation.TimeReceived) {
			return CANCELED
		}
	}

	// if time is before the auction start time, the auction is pending
	if currTime.Before(auction.Item.StartTime) {
		return PENDING
	}

	// if time is between the start and end time, inclusive, the auction is active
	atOrAfterStart := AfterOrOn(&currTime, &auction.Item.StartTime)
	atOrBeforeEnd := BeforeOrOn(&currTime, &auction.Item.EndTime)
	if atOrAfterStart && atOrBeforeEnd {
		return ACTIVE
	}

	// if time is after auction end time (auction has not been cancelled nor finalized), then the auction is over
	// (already checked if auction has been cancelled)
	if currTime.After(auction.Item.EndTime) {
		return OVER
	}

	panic("Auction.GetStateAtTime() couldn't determine auction state at time!")
}

func (auction *Auction) alertSeller(msg string) {
	sellerUserId := auction.Item.SellerUserId
	log.Printf("[Auction %s] STUBBED: sending out request to notify seller (userId=%s,msg=%s)\n", auction.Item.ItemId, sellerUserId, msg)
}

func (auction *Auction) alertBidder(msg string, bid *Bid) {
	bidderUserId := bid.BidderUserId
	log.Printf("[Auction %s] STUBBED: sending out request to notify bidder (userId=%s,msg=%s)\n", auction.Item.ItemId, bidderUserId, msg)
}

func (auction *Auction) Cancel(timeWhenCancellationIssued time.Time) bool {

	// cant issue cancel if there is already a cancellation, or the auction is considered finalized
	if auction.HasCancellation() || auction.HasFinalization() {
		log.Printf("[Auction %s] can't cancel self because I am already canceled/finalized.\n", auction.Item.ItemId)
		return false // only allow 1 cancellation; don't allow any changes once finalized
	}

	// otherwise, the auction is pending, active, or over.
	// can only issue cancel if auction is pending or active and has no bids
	stateWhenCancellationIssued := auction.getStateAtTime(timeWhenCancellationIssued)
	switch {
	case stateWhenCancellationIssued == PENDING: //
		auction.cancellation = NewCancellation(timeWhenCancellationIssued)
		log.Printf("[Auction %s] canceling self (pending auction state).\n", auction.Item.ItemId)
		return true
	case stateWhenCancellationIssued == ACTIVE && !auction.HasActiveBid(): //
		auction.cancellation = NewCancellation(timeWhenCancellationIssued)
		log.Printf("[Auction %s] canceling self (active auction state but no active bids).\n", auction.Item.ItemId)
		return true
	default:
		log.Printf("[Auction %s] can't cancel self (auction is over or finalized).\n", auction.Item.ItemId)
		return false // state is OVER, or CANCELED already
	}
}

func (auction *Auction) HasActiveBid() bool {
	for _, bid := range auction.bids {
		if bid.active {
			return true
		}
	}
	return false
}

func (auction *Auction) Stop(timeWhenStopIssued time.Time) bool {

	// cant issue stop if there is already a cancellation, or the auction is considered finalized
	if auction.HasCancellation() || auction.HasFinalization() {
		log.Printf("[Auction %s] can't stop self because I am already canceled/finalized.\n", auction.Item.ItemId)
		return false // only allow 1 cancellation; don't allow any changes once finalized
	}

	// otherwise, the auction is pending, active, or over.
	// can only issue stop if auction is pending or active
	stateWhenStopIssued := auction.getStateAtTime(timeWhenStopIssued)
	switch {
	case stateWhenStopIssued == PENDING || stateWhenStopIssued == ACTIVE:
		auction.cancellation = NewCancellation(timeWhenStopIssued)
		log.Printf("[Auction %s] stopping self.\n", auction.Item.ItemId)
		return true
	default:
		log.Printf("[Auction %s] can't stop self because of my state.\n", auction.Item.ItemId)
		return false // state is OVER, CANCELED, FINALIZED; cant stop
	}
}

func (auction *Auction) HasCancellation() bool {
	return auction.cancellation != nil
}

func (auction *Auction) HasFinalization() bool {
	return auction.finalization != nil
}

func (auction *Auction) DeactivateUserBids(userId string, timeWhenUserDeactivated time.Time) (*[]*Bid, bool) {
	// note: this call will deactivate all of the user's bids in the auction even
	// bids that are placed after the timeWhenUserDeactivated. timeWhenUserDeactivated
	// is only used to determine whether the auction outcome was "set-in-stone" when
	// the request to deactivate user's bids comes in; this is the only situation where
	// we refused a deactivateUserBids request
	log.Printf("[Auction %s] de-activating user's bids (userId=%s)\n", auction.Item.ItemId, userId)
	stateWhenUserDeactivated := auction.getStateAtTime(timeWhenUserDeactivated)
	bidsToSave := []*Bid{}
	if stateWhenUserDeactivated == FINALIZED {
		return &bidsToSave, false
	} else {
		for _, bid := range auction.bids {
			if bid.BidderUserId == userId {
				gotDeactivated := bid.Deactivate()
				if gotDeactivated {
					bidsToSave = append(bidsToSave, bid)
				}
			}
		}
		return &bidsToSave, true
	}
}

func (auction *Auction) ActivateUserBids(userId string, timeWhenUserActivated time.Time) (*[]*Bid, bool) {

	log.Printf("[Auction %s] activating user's bids (userId=%s)\n", auction.Item.ItemId, userId)
	stateWhenUserActivated := auction.getStateAtTime(timeWhenUserActivated)
	bidsToSave := []*Bid{}
	if stateWhenUserActivated == FINALIZED {
		return &bidsToSave, false
	} else {
		for _, bid := range auction.bids {
			if bid.BidderUserId == userId {
				wasActivated := bid.Activate()
				if wasActivated {
					bidsToSave = append(bidsToSave, bid)
				}
			}
		}
		return &bidsToSave, true
	}
}

func (auction *Auction) IsOverOrCanceledAtTime(atTime time.Time) bool {
	stateAtTime := auction.getStateAtTime(atTime)
	if stateAtTime == OVER || stateAtTime == CANCELED {
		return true
	}
	return false
}

func (auction *Auction) IsPending(nowTime time.Time) bool {
	stateAtTime := auction.getStateAtTime(nowTime)
	return stateAtTime == PENDING
}

func (auction *Auction) IsActive(nowTime time.Time) bool {
	stateAtTime := auction.getStateAtTime(nowTime)
	return stateAtTime == ACTIVE
}

func (auction *Auction) IsCanceled(nowTime time.Time) bool {
	stateAtTime := auction.getStateAtTime(nowTime)
	return stateAtTime == CANCELED
}

func (auction *Auction) IsFinalized(nowTime time.Time) bool {
	stateAtTime := auction.getStateAtTime(nowTime)
	return stateAtTime == FINALIZED
}

func (auction *Auction) Finalize(timeWhenFinalizationIssued time.Time) bool {

	// cant issue finalization if this auction has already been finalized
	if auction.HasFinalization() {
		return false // only allow 1 finalization
	}

	// finalization only allowed when auction is canceled or over
	state := auction.getStateAtTime(timeWhenFinalizationIssued)
	switch {
	case state == CANCELED || state == OVER:
		log.Printf("[Auction %s] STUBBED finalizing self...\n", auction.Item.ItemId)
		auction.finalization = NewFinalization(timeWhenFinalizationIssued)
		return true
	default:
		return false // state is PENDING, ACTIVE, FINALIZED
	}

}

func (auction *Auction) hasBid(bidId string) bool {
	for _, bid := range auction.bids {
		if bid.BidId == bidId {
			return true
		}
	}
	return false
}

func (auction *Auction) OverlapsWith(leftBound *time.Time, rightBound *time.Time) bool {
	if rightBound.Before(auction.Item.StartTime) || leftBound.After(auction.Item.EndTime) {
		return false
	}
	return true
}

func (auction *Auction) SendStartSoonAlertIfApplicable() bool {
	nowTime := time.Now()
	stateNow := auction.getStateAtTime(nowTime) // send start soon alert if in pending state; send active now alert if in active state

	if !auction.sentStartSoonAlert {

		if stateNow == PENDING {

			timeUntilStart := auction.Item.StartTime.Sub(nowTime)
			hours := timeUntilStart.Hours()
			mins := timeUntilStart.Minutes()

			if hours < 1 { // send alert if auction is active and end is within 1 hour from now
				if mins < 60 {
					log.Printf("[Auction %s] STUBBED sending out starting soon alert; starts in (%f minutes)\n", auction.Item.ItemId, mins)
				} else {
					log.Printf("[Auction %s] STUBBED sending out starting soon alert; starts in (%f hours)\n", auction.Item.ItemId, hours)
				}

				auction.sentStartSoonAlert = true
				return true
			}

		} else if stateNow == ACTIVE {

			timeSinceStart := nowTime.Sub(auction.Item.StartTime)
			hours := timeSinceStart.Hours()
			mins := timeSinceStart.Minutes()

			if mins < 60 {
				log.Printf("[Auction %s] STUBBED sending out starting soon alert; auction started (%f minutes) ago!\n", auction.Item.ItemId, mins)
			} else {
				log.Printf("[Auction %s] STUBBED sending out starting soon alert; auction started (%f minutes) ago!\n", auction.Item.ItemId, hours)
			}

			auction.sentStartSoonAlert = true
			return true

		} else {
			// otherwise the auction state is over/finalized/canceled,
			// and we never managed to send out a start soon alert,
			// so just dont send any alert at all
			auction.sentStartSoonAlert = true
			return true
		}
	} // already sent start soon alert
	return false
}

// func (auction *Auction) timeUntilStart(currTime *time.Time) (*float64, *float64, bool) {
// 	if currTime.After(auction.Item.StartTime) {
// 		return nil, nil, false
// 	} else {
// 		durationTilStart := auction.Item.StartTime.Sub(*currTime) // e.g. 1 hrs 30 min

// 		hours := durationTilStart.Hours()
// 		minutes := durationTilStart.Hours()
// 		return &hours, &minutes, true
// 	}
// }

// func (auction *Auction) timeSinceStart(currTime *time.Time) (*float64, *float64, bool) {
// 	if currTime.Before(auction.Item.StartTime) {
// 		return nil, nil, false
// 	} else {
// 		durationTilStart := (*currTime).Sub(auction.Item.StartTime) // e.g. 1 hrs 30 min

// 		hours := durationTilStart.Hours()
// 		minutes := durationTilStart.Hours()
// 		return &hours, &minutes, true
// 	}
// }

// func (auction *Auction) timeUntilEnd(currTime *time.Time) (*float64, *float64, bool) {
// 	if currTime.After(auction.Item.EndTime) {
// 		return nil, nil, false
// 	} else {
// 		durationTilEnd := auction.Item.EndTime.Sub(*currTime) // e.g. 1 hrs 30 min

// 		hours := durationTilEnd.Hours()
// 		minutes := durationTilEnd.Hours()
// 		return &hours, &minutes, true
// 	}
// }

func (auction *Auction) SendEndSoonAlertIfApplicable() bool {

	nowTime := time.Now()
	stateNow := auction.getStateAtTime(nowTime) // send end soon alert if in active state; send ended earlier if in over, canceled, completed state

	if !auction.sentEndSoonAlert {

		if stateNow == ACTIVE {
			timeUntilEnd := auction.Item.EndTime.Sub(nowTime)
			hours := timeUntilEnd.Hours()
			mins := timeUntilEnd.Minutes()

			if hours < 1 { // send alert if auction is active and end is within 1 hour from now
				if mins < 60 {
					log.Printf("[Auction %s] STUBBED sending out ending soon alert; ends in (%f minutes)\n", auction.Item.ItemId, mins)
				} else {
					log.Printf("[Auction %s] STUBBED sending out ending soon alert; ends in (%f hours)\n", auction.Item.ItemId, hours)
				}

				auction.sentEndSoonAlert = true
				return true
			}
		} else if stateNow == OVER {

			timeSinceEnd := nowTime.Sub(auction.Item.EndTime)
			hours := timeSinceEnd.Hours()
			mins := timeSinceEnd.Minutes()

			if mins < 60 {
				log.Printf("[Auction %s] STUBBED sending out ending soon alert; auction ended (%f minutes) ago!\n", auction.Item.ItemId, mins)
			} else {
				log.Printf("[Auction %s] STUBBED sending out ending soon alert; auction ended (%f hours) ago!\n", auction.Item.ItemId, hours)
			}

			auction.sentEndSoonAlert = true
			return true

		} else if stateNow == CANCELED {

			timeSinceCancel := nowTime.Sub(auction.cancellation.TimeReceived)
			hours := timeSinceCancel.Hours()
			mins := timeSinceCancel.Minutes()

			if mins < 60 {
				log.Printf("[Auction %s] STUBBED sending out ending soon alert; auction was canceled (%f minutes) ago!\n", auction.Item.ItemId, mins)
			} else {
				log.Printf("[Auction %s] STUBBED sending out ending soon alert; auction was canceled (%f hours) ago!\n", auction.Item.ItemId, hours)
			}

			auction.sentEndSoonAlert = true
			return true

		} else if stateNow == FINALIZED {

			timeSinceFinalization := nowTime.Sub(auction.finalization.TimeReceived)
			hours := timeSinceFinalization.Hours()
			mins := timeSinceFinalization.Minutes()

			if mins < 60 {
				log.Printf("[Auction %s] STUBBED sending out ending soon alert; auction was finalized (%f minutes) ago!\n", auction.Item.ItemId, mins)
			} else {
				log.Printf("[Auction %s] STUBBED sending out ending soon alert; auction was finalized (%f hours) ago!\n", auction.Item.ItemId, hours)
			}

			auction.sentEndSoonAlert = true
			return true

		} else {
			// otherwise, the auction state is pending, and we haven't sent any "end soon" alerts.
			// later when auction becomes active, the alert will go out
			return false
		}
	} // already sent notification;
	return false

}

// misc helper functions

func AfterOrOn(someTime *time.Time, otherTime *time.Time) bool {
	return someTime.After(*otherTime) || someTime.Equal(*otherTime)
}

func BeforeOrOn(someTime *time.Time, otherTime *time.Time) bool {
	return someTime.Before(*otherTime) || someTime.Equal(*otherTime)
}
