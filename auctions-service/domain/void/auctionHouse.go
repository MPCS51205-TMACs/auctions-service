package domain

// type AuctionHouse struct {
// 	auctions          map[string]*Auction
// 	auctionRepository AuctionRepository
// }

// func NewAuctionHouse(auctionsToManage []*Auction, repo AuctionRepository) *AuctionHouse {
// 	var auctions map[string]*Auction = make(map[string]*Auction)
// 	for _, auction := range auctionsToManage {
// 		auctions[auction.Item.ItemId] = auction
// 	}
// 	return &AuctionHouse{
// 		auctions,
// 		repo,
// 	}
// }

// func (auctionhouse *AuctionHouse) directBid(newIncomingBid *Bid) {
// 	relevantItemId := newIncomingBid.ItemId
// 	relevantAuction := auctionhouse.auctions[relevantItemId]
// 	if relevantAuction != nil {
// 		relevantAuction.ProcessNewBid(newIncomingBid)
// 	}
// }

// func (auctionhouse *AuctionHouse) CreateAuction(item *Item) {
// 	newAuction := NewAuction(item, nil, nil)
// 	auctionhouse.auctionRepository.SaveAuction(newAuction)
// 	auctionhouse.addAuction(newAuction)
// }

// func (auctionhouse *AuctionHouse) CancelAuction(itemId string, timeWhenCancellationIssued time.Time) bool {
// 	var success bool
// 	relevantAuction := auctionhouse.auctions[itemId]
// 	if relevantAuction != nil {
// 		success := relevantAuction.Cancel(timeWhenCancellationIssued)
// 		if success {
// 			auctionhouse.auctionRepository.SaveAuction(relevantAuction)
// 		}
// 	}
// 	return success
// }

// func (auctionhouse *AuctionHouse) StopAuction(itemId string, timeWhenStopIssued time.Time) bool {
// 	var success bool
// 	relevantAuction := auctionhouse.auctions[itemId]
// 	if relevantAuction != nil {
// 		success := relevantAuction.Stop(timeWhenStopIssued)
// 		if success {
// 			auctionhouse.auctionRepository.SaveAuction(relevantAuction)
// 		}
// 	}
// 	return success
// }

// func (auctionhouse *AuctionHouse) TurnOnSessionManagement() {}

// func (auctionhouse *AuctionHouse) TurnOffSessionManagement() {}

// func (auctionhouse *AuctionHouse) addAuction(auction *Auction) {
// 	auctionhouse.auctions[auction.Item.ItemId] = auction
// }

// func (auctionhouse *AuctionHouse) removeAuction(itemId string) {
// 	if _, ok := auctionhouse.auctions[itemId]; ok {
// 		delete(auctionhouse.auctions, itemId)
// 	}
// }

// func (auctionhouse *AuctionHouse) finalizeAuction(itemId string) {
// 	relevantAuction := auctionhouse.auctions[itemId]
// 	if relevantAuction != nil {
// 		winningBid := relevantAuction.GetHighestBid()
// 		bidderUserId := winningBid.BidderUserId
// 		payment := winningBid.AmountInCents
// 		itemWon := winningBid.ItemId
// 		fmt.Println("STUBBED: packaging up data and sending to Closed Auction Metrics")
// 		fmt.Printf("STUBBED: sending out message to Shopping Cart of auction end; winnerUserId=%v,itemId=%v,amountInCents=%v\n", bidderUserId, itemWon, payment)
// 		auctionhouse.removeAuction(itemId)
// 	}
// 	panic("tried to finalize an Auction that I don't have in memory!")
// }
