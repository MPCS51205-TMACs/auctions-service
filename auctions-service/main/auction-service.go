package main

// this file is the entire interface of the auctions-service

import (
	"auctions-service/domain"
	"fmt"
	"log"
	"sync"
	"time"
)

type AuctionService struct {
	bidRepo          domain.BidRepository
	auctionRepo      domain.AuctionRepository
	inMemoryAuctions map[string]*domain.Auction
	mutex            *sync.Mutex // for now, using coarse-grained concurrency implementation
}

func NewAuctionService(bidRepo domain.BidRepository, auctionRepo domain.AuctionRepository) *AuctionService {
	inMemoryAuctions := map[string]*domain.Auction{}
	mutex := &sync.Mutex{}
	// comment out lines below to turn ON logging (logging currently turned OFF)
	// log.SetFlags(0)
	// log.SetOutput(ioutil.Discard)
	return &AuctionService{
		bidRepo:          bidRepo,
		auctionRepo:      auctionRepo,
		inMemoryAuctions: inMemoryAuctions,
		mutex:            mutex,
	}
}

type AuctionInteractionOutcome string

const (
	auctionAlreadyCreated                   AuctionInteractionOutcome = "ALREADY_CREATED"      // create
	auctionSuccessfullyCreated              AuctionInteractionOutcome = "CREATED_SUCCESSFULLY" // create
	auctionWouldStartTooSoon                AuctionInteractionOutcome = "STARTS_TOO_SOON"      // create
	auctionStartsInPast                     AuctionInteractionOutcome = "STARTS_IN_PAST"       // create
	badTimeSpecified                        AuctionInteractionOutcome = "BAD_TIME_SPECIFIED_TIME"
	auctionSuccessfullyCanceled             AuctionInteractionOutcome = "CANCELED_SUCCESSFULLY"   // cancel
	auctionSuccessfullyStopped              AuctionInteractionOutcome = "STOPPED_SUCCESSFULLY"    // stop
	auctionNotExist                         AuctionInteractionOutcome = "AUCTION_NOT_EXIST"       // cancel, stop
	auctionAlreadyCanceled                  AuctionInteractionOutcome = "ALREADY_CANCELED"        // cancel
	auctionAlreadyOver                      AuctionInteractionOutcome = "ALREADY_OVER"            // cancel, stop
	auctionAlreadyFinalized                 AuctionInteractionOutcome = "ALREADY_FINALIZED"       // cancel, stop
	auctionCancellationRequesterIsNotSeller AuctionInteractionOutcome = "REQUESTER_IS_NOT_SELLER" // cancel
	auctionProcessedBid                     AuctionInteractionOutcome = "BID_WAS_SEEN_BY_AUCTION" // cancel
)

func (auctionservice *AuctionService) CreateAuction(itemId, sellerUserId string, startTime, endTime *time.Time, startPriceInCents int64) AuctionInteractionOutcome {

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	log.Printf("[AuctionService] creating Auction (itemId=%s)...", itemId)

	// confirm well-specified time
	if !endTime.After(*startTime) {
		log.Printf("[AuctionService] fail. starttime is not < endtime")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return badTimeSpecified
	}

	creationTime := time.Now()

	// confirm auction does not start in the past
	if creationTime.After(*startTime) {
		log.Printf("[AuctionService] fail. Auction would start in the past.")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionStartsInPast
	}

	// if auction to be created will start in sooner than 5 minutes, do not proceed
	if creationTime.Add(time.Duration(15) * time.Second).After(*startTime) {
		log.Printf("[AuctionService] fail. Auction would start in <= 1 minute from now. push back start time to later time.")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionWouldStartTooSoon
	}

	// confirm an auction hasn't already been created for the item
	if auctionservice.auctionRepo.GetAuction(itemId) != nil {
		log.Printf("[AuctionService] fail. Auction already exists for item.")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionAlreadyCreated
	}

	newItem := domain.NewItem(itemId, sellerUserId, *startTime, *endTime, startPriceInCents)
	newAuction := auctionservice.auctionRepo.NewAuction(newItem, nil, nil, false, false, nil)

	auctionservice.auctionRepo.SaveAuction(newAuction)                   // save Auction
	auctionservice.inMemoryAuctions[newAuction.Item.ItemId] = newAuction // cache Auction
	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()
	// auctionservice.addAuction(newAuction)
	// auctionservice.auctionRepo.SaveAuction()
	log.Printf("[AuctionService] success. Auction created.")
	return auctionSuccessfullyCreated
}

func (auctionservice *AuctionService) CancelAuction(itemId string, requesterUserId string) AuctionInteractionOutcome {

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	log.Printf("[AuctionService] cancelling Auction (itemId=%s;requesterUserId=%s)...", itemId, requesterUserId)
	timeWhenCancelReceived := time.Now()

	relevantAuction, ok := auctionservice.inMemoryAuctions[itemId] // lookup in cache
	if !ok {
		relevantAuction = auctionservice.auctionRepo.GetAuction(itemId) // get from db if not cached
	} // dont bother caching though

	// confirm auction exists
	if relevantAuction == nil {
		log.Printf("[AuctionService] fail. Auction does not exist for itemId=%s ", itemId)
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionNotExist
	}

	// confirm the person requesting an auction be canceled is the seller of the item
	if relevantAuction.Item.SellerUserId != requesterUserId {
		log.Printf("[AuctionService] fail. Requester trying to cancel is not seller of itemId=%s", itemId)
		fmt.Println("AuctionService thinks seller is = '%s'; requester is '%s'", relevantAuction.Item.SellerUserId, requesterUserId)
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionCancellationRequesterIsNotSeller
	}

	// confirm auction isn't already finalized
	if relevantAuction.HasFinalization() {
		log.Printf("[AuctionService] fail. Auction already finalized")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionAlreadyFinalized
	}

	// confirm auction isn't already canceled
	if relevantAuction.HasCancellation() {
		log.Printf("[AuctionService] fail. Auction already canceled")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionAlreadyCanceled
	}

	// confirm auction isn't already over (at time)
	if relevantAuction.IsOverOrCanceledAtTime(timeWhenCancelReceived) {
		log.Printf("[AuctionService] fail. Auction already over")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionAlreadyOver
	}

	// confirm auction isn't active and has a bid (at time)
	if relevantAuction.HasActiveBid() && relevantAuction.IsActive(timeWhenCancelReceived) {
		log.Printf("[AuctionService] fail. Auction is ongoing and there is at least one bid")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionAlreadyOver
	}

	// otherwise, should be ok to cancel.

	wasCanceled := relevantAuction.Cancel(timeWhenCancelReceived) // should always return true...
	if wasCanceled {
		auctionservice.auctionRepo.SaveAuction(relevantAuction) // save Auction
		log.Printf("[AuctionService] success. Auction canceled")
		// dont cache Auction
	}

	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()

	if wasCanceled {
		return auctionSuccessfullyCanceled
	} else {
		panic("[AuctionService] see CancelAuction(). reached end of method without determining what happened (bug).")
	}
}

func (auctionservice *AuctionService) StopAuction(itemId string) AuctionInteractionOutcome {

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	log.Printf("[AuctionService] stopping Auction (itemId=%s)...", itemId)
	timeWhenStopReceived := time.Now()

	relevantAuction, ok := auctionservice.inMemoryAuctions[itemId] // lookup in cache
	toCache := false
	if !ok {
		relevantAuction = auctionservice.auctionRepo.GetAuction(itemId) // get from db if not cached
		toCache = true
	} // cache it if successful b/c it means we will need to finalized it.

	// confirm auction exists
	if relevantAuction == nil {
		log.Printf("[AuctionService] fail. Auction does not exist for itemId=%s ", itemId)
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionNotExist
	}

	// delegate security to main.go; assume main.go confirmed requester is an admin

	// confirm auction isn't already finalized
	if relevantAuction.HasFinalization() {
		log.Printf("[AuctionService] fail. Auction already finalized")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionAlreadyFinalized
	}

	// confirm auction isn't already canceled
	if relevantAuction.HasCancellation() {
		log.Printf("[AuctionService] fail. Auction already canceled")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionAlreadyCanceled
	}

	// confirm auction isn't already over
	if relevantAuction.IsOverOrCanceledAtTime(timeWhenStopReceived) {
		log.Printf("[AuctionService] fail. Auction already over")
		// log.Printf("[AuctionService] UNLOCK")
		auctionservice.mutex.Unlock()
		return auctionAlreadyOver
	}

	// otherwise, ok to stop.
	wasStopped := relevantAuction.Stop(time.Now()) // should always return true?
	if wasStopped {
		auctionservice.auctionRepo.SaveAuction(relevantAuction)
	}

	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()

	if wasStopped {
		if toCache {
			auctionservice.inMemoryAuctions[itemId] = relevantAuction // cache it
		}
		log.Printf("[AuctionService] success. Auction stopped")
		return auctionSuccessfullyStopped
	} else {
		panic("[AuctionService] see StopAuction(). reached end of method without determining what happened (bug).")
	}

}

func (auctionservice *AuctionService) ProcessNewBid(itemId string, bidderUserId string, timeReceived time.Time, amountInCents int64) (AuctionInteractionOutcome, domain.AuctionState, bool) {

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	log.Printf("[AuctionService] processing new bid (itemId=%s;bidderUserId=%s;amountInCents=%d;time=%v)...", itemId, bidderUserId, amountInCents, timeReceived)
	newId := auctionservice.bidRepo.NextBidId()
	newBid := domain.NewBid(newId, itemId, bidderUserId, timeReceived, amountInCents, true)

	relevantAuction, ok := auctionservice.inMemoryAuctions[itemId] // lookup in cache
	toCache := false
	if !ok {
		relevantAuction = auctionservice.auctionRepo.GetAuction(itemId) // get from db if not cached
		if relevantAuction == nil {
			// log.Printf("[AuctionService] UNLOCK")
			auctionservice.mutex.Unlock()
			return auctionNotExist, domain.UNKNOWN, false // unknown auction state == auction not exist
		}
		toCache = true
	} // cache the auction if the bid ends up successfully being placed.

	auctionState, wasNewTopBid := relevantAuction.ProcessNewBid(newBid)
	if wasNewTopBid {
		auctionservice.bidRepo.SaveBid(newBid) // only save bids that were determined to be new Top bids
		if toCache {
			auctionservice.inMemoryAuctions[itemId] = relevantAuction
		}
	}

	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()

	return auctionProcessedBid, auctionState, wasNewTopBid

}

func (auctionservice *AuctionService) ValidateBid(itemId string, bidderUserId string, timeReceived time.Time, amountInCents int64) (domain.AuctionState, bool) {

	// log.Printf("[AuctionService] LOCK")
	// NOT USING LOCK
	// auctionservice.mutex.Lock()

	log.Printf("[AuctionService] validating new bid request (itemId=%s;bidderUserId=%s;amountInCents=%d;time=%v)...", itemId, bidderUserId, amountInCents, timeReceived)

	relevantAuction, ok := auctionservice.inMemoryAuctions[itemId] // lookup in cache
	if !ok {
		auctionservice.mutex.Lock()
		relevantAuction = auctionservice.auctionRepo.GetAuction(itemId) // get from db if not cached
		if relevantAuction == nil {
			// log.Printf("[AuctionService] UNLOCK")
			auctionservice.mutex.Unlock()
			return domain.UNKNOWN, false // unknown auction state == auction not exist
		}
		auctionservice.inMemoryAuctions[itemId] = relevantAuction // cache it
		auctionservice.mutex.Unlock()
	}

	auctionState := relevantAuction.GetStateAtTime(timeReceived)
	// PENDING   AuctionState = "PENDING" // has not yet started
	// ACTIVE    AuctionState = "ACTIVE"  // is happening now
	// CANCELED  AuctionState = "CANCELED"
	// OVER      AuctionState = "OVER"      // is over (but winner has not been declared and auction has not been "archived away")
	// FINALIZED AuctionState = "FINALIZED" // is over and archived away; can delete
	// UNKNOWN   AuctionState = "UKNOWN"

	canAcceptBid := auctionState == domain.ACTIVE // it's only ok to place bid when auction is ACTIVE

	// log.Printf("[AuctionService] UNLOCK")
	// auctionservice.mutex.Unlock()

	return auctionState, canAcceptBid

}

func (auctionservice *AuctionService) GetItemsUserHasBidsOn(userId string) *[]string {
	log.Printf("[AuctionService] getting and returning items that userId=%s has bids on...", userId)
	bids := auctionservice.bidRepo.GetBidsByUserId(userId) // includes inactive bids
	itemIds := make([]string, 0)
	alreadySeenItemIds := map[string]interface{}{}
	for _, bid := range *bids {
		if _, ok := alreadySeenItemIds[bid.ItemId]; !ok {
			itemIds = append(itemIds, bid.ItemId)
			alreadySeenItemIds[bid.ItemId] = nil
		}
	}
	return &itemIds
}

func (auctionservice *AuctionService) GetAuction(itemId string) *domain.Auction {
	log.Println("[AuctionService] getting and returning auction for itemid=%s...", itemId)
	auction := auctionservice.auctionRepo.GetAuction(itemId)
	return auction
}

func (auctionservice *AuctionService) GetAuctions() *[]*domain.Auction {
	log.Println("[AuctionService] getting and returning all auctions (excluding bid data)...")
	nowTime := time.Now()
	old := nowTime.Add(-time.Duration(24*3000) * time.Hour)         // 3000 days in past (basically all previous auctions)
	future := nowTime.Add(time.Duration(24*3000) * time.Hour)       // 3000 days in future (basically all future auctions)
	auctions := auctionservice.auctionRepo.GetAuctions(old, future) // all auctions whose start->end time overlaps with this window
	// filter down to only active auctions (some auctions may be canceled / finalized even though their start->end time overlaps w now)
	return &auctions
}

func (auctionservice *AuctionService) GetActiveAuctions() *[]*domain.Auction {
	log.Println("[AuctionService] getting and returning active auctions...")
	nowTime := time.Now()
	auctions := auctionservice.auctionRepo.GetAuctions(nowTime, nowTime) // all auctions whose start->end time overlaps with nowTime
	// filter down to only active auctions (some auctions may be canceled / finalized even though their start->end time overlaps w now)
	activeAuctions := make([]*domain.Auction, 0)
	for _, auction := range auctions {
		if auction.IsActive(nowTime) {
			activeAuctions = append(activeAuctions, auction)
		}
	}
	return &auctions
}

func (auctionservice *AuctionService) ActivateUserBids(userId string) int {

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	log.Printf("[AuctionService] activating bids for userId=%s...", userId)

	timeWhenUserActivated := time.Now()

	userBids := auctionservice.bidRepo.GetBidsByUserId(userId)
	itemIds := make([]string, 0) // list of all items (auctions) the user has bids in
	alreadySeenItemIds := map[string]interface{}{}
	for _, bid := range *userBids {
		if _, ok := alreadySeenItemIds[bid.ItemId]; !ok {
			itemIds = append(itemIds, bid.ItemId)
			alreadySeenItemIds[bid.ItemId] = nil
		}
	}

	bidsToUpdateInRepo := []*domain.Bid{}
	numAuctionsWBidUpdates := 0

	for _, itemId := range itemIds {
		auction, ok := auctionservice.inMemoryAuctions[itemId]
		if !ok { // did not find auction in memory; bring into memory
			auction = auctionservice.auctionRepo.GetAuction(itemId)
		}
		bidsToSave, _ := auction.ActivateUserBids(userId, timeWhenUserActivated) // returns the bids whose state was changed
		bidsToUpdateInRepo = append(bidsToUpdateInRepo, *bidsToSave...)
		numAuctionsWBidUpdates++
	}

	for _, bid := range bidsToUpdateInRepo {
		auctionservice.bidRepo.SaveBid(bid)
	}

	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()

	return numAuctionsWBidUpdates

}

func (auctionservice *AuctionService) DeactivateUserBids(userId string) int {

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	log.Printf("[AuctionService] de-activating bids for userId=%s...", userId)

	timeWhenUserDeactivated := time.Now()

	userBids := auctionservice.bidRepo.GetBidsByUserId(userId)
	itemIds := make([]string, 0) // list of all items (auctions) the user has bids in
	alreadySeenItemIds := map[string]interface{}{}
	for _, bid := range *userBids {
		if _, ok := alreadySeenItemIds[bid.ItemId]; !ok {
			itemIds = append(itemIds, bid.ItemId)
			alreadySeenItemIds[bid.ItemId] = nil
		}
	}

	bidsToUpdateInRepo := []*domain.Bid{}
	numAuctionsWBidUpdates := 0

	for _, itemId := range itemIds {
		auction, ok := auctionservice.inMemoryAuctions[itemId]
		if !ok { // did not find auction in memory; bring into memory
			auction = auctionservice.auctionRepo.GetAuction(itemId)
		}
		bidsToSave, _ := auction.DeactivateUserBids(userId, timeWhenUserDeactivated) // returns the bids whose state was changed
		bidsToUpdateInRepo = append(bidsToUpdateInRepo, *bidsToSave...)
		numAuctionsWBidUpdates++
	}

	for _, bid := range bidsToUpdateInRepo {
		auctionservice.bidRepo.SaveBid(bid)
	}

	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()

	return numAuctionsWBidUpdates
}

func (auctionservice *AuctionService) LoadAuctionsIntoMemory(sinceTime time.Time, upToTime time.Time) {

	inMemAuctions := auctionservice.inMemoryAuctions
	var broughtIntoMemory int

	var InMemoryPending int
	var InMemoryActive int
	var InMemoryCanceled int
	var InMemoryFinalized int

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	nowTime := time.Now()
	for _, auction := range inMemAuctions {
		switch {
		case auction.IsPending(nowTime):
			InMemoryPending++
		case auction.IsActive(nowTime):
			InMemoryActive++
		case auction.IsCanceled(nowTime):
			InMemoryCanceled++
		case auction.IsFinalized(nowTime):
			InMemoryFinalized++
		}
	}

	auctions := auctionservice.auctionRepo.GetAuctions(sinceTime, upToTime)
	for _, auction := range auctions {
		if !auction.HasFinalization() { // dont bring into memory if it is a finalized auction
			if _, ok := (inMemAuctions)[auction.Item.ItemId]; !ok {
				(inMemAuctions)[auction.Item.ItemId] = auction
				broughtIntoMemory++
			}
		}
	}

	numInRepo := auctionservice.auctionRepo.NumAuctionsSaved()
	memoryState := fmt.Sprintf("%d/%d/%d/%d", InMemoryPending, InMemoryActive, InMemoryCanceled, InMemoryFinalized)
	log.Printf("[AuctionService] loaded %d new auctions into memory ([Pending/Active/Canceled/Finalized] %s in memory; %d total in repository;)", broughtIntoMemory, memoryState, numInRepo)

	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()
}

func (auctionservice *AuctionService) SendOutLifeCycleAlerts() {
	inMemAuctions := auctionservice.inMemoryAuctions
	var sentNotif1, sentNotif2 bool

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	log.Println("[AuctionService] sending out life cycle alerts...")

	for _, auction := range inMemAuctions {
		sentNotif1 = auction.SendStartSoonAlertIfApplicable()
		sentNotif2 = auction.SendEndSoonAlertIfApplicable()
		if sentNotif1 || sentNotif2 {
			auctionservice.auctionRepo.SaveAuction(auction) // save the knowledge that alert was sent out;
		}
	}

	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()
}

func (auctionservice *AuctionService) FinalizeAnyPastAuctions(finalizeDelay time.Duration) {
	inMemAuctions := auctionservice.inMemoryAuctions

	// log.Printf("[AuctionService] LOCK")
	auctionservice.mutex.Lock()

	log.Println("[AuctionService] finalizing (archiving) any past auctions...")

	for _, auction := range inMemAuctions {
		wasFinalized := auction.Finalize(time.Now())
		if wasFinalized {
			auctionservice.auctionRepo.SaveAuction(auction) // save the knowledge that we finalized the auction
		}
	}

	// log.Printf("[AuctionService] UNLOCK")
	auctionservice.mutex.Unlock()
}
