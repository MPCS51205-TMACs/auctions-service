package domain

import (
	"fmt"
	"log"
	"time"
	// acquired by doing 'go get github.com/rabbitmq/amqp091-go'
)

// sends messages to console
type ConsoleAlertEngine struct {
}

func NewConsoleAlertEngine() *ConsoleAlertEngine {
	return &ConsoleAlertEngine{}
}

func (alertEngine *ConsoleAlertEngine) SendAuctionStartSoonAlert(msg, itemId string, startTime, endTime time.Time) {
	log.Printf("[AlertEngine] Sending AuctionStartingSoon Alert to Console (item_id=%s)\n", itemId)
	alertEngine.sendToConsole(msg)
	log.Printf("[AlertEngine] [x] Sent AuctionStartingSoon Alert (item_id=%s)\n", itemId)
}

func (alertEngine *ConsoleAlertEngine) SendAuctionEndSoonAlert(msg, itemId string, startTime, endTime time.Time) {
	log.Printf("[AlertEngine] Sending AuctionEndingSoon Alert to Console (item_id=%s)\n", itemId)
	alertEngine.sendToConsole(msg)
	log.Printf("[AlertEngine] [x] Sent AuctionEndingSoon Alert (item_id=%s)\n", itemId)
}

func (alertEngine *ConsoleAlertEngine) SendAuctionEndAlert(auctionData *AuctionData) {
	log.Printf("[AlertEngine] Sending AuctionData to Console (item_id=%s)\n", auctionData.Item.ItemId)
	msg := fmt.Sprint(auctionData)
	alertEngine.sendToConsole(msg)
	log.Printf("[AlertEngine] [x] Sent AuctionData (item_id=%s)\n", auctionData.Item.ItemId)
}

func (alertEngine *ConsoleAlertEngine) SendNewTopBidAlert(itemId, sellerUserId, formerTopBidderUserId, newTopBidderUserId *string) {
	log.Printf("[AlertEngine] Sending NewTopBid Alert to Console (item_id=%s)\n", *itemId)
	msg := fmt.Sprintf("itemId=%v, sellerUserId=%v, formerTopBidderUserId=%v, newTopBidderUserId=%v", itemId, sellerUserId, formerTopBidderUserId, newTopBidderUserId)
	alertEngine.sendToConsole(msg)
	log.Printf("[AlertEngine] [x] Sent NewTopBid Alert (item_id=%s)\n", *itemId)
}

// func (alertEngine *ConsoleAlertEngine) AlertSeller(msg, itemId, sellerUserId string) {
// 	log.Printf("[AlertEngine] Sending AlertSeller Alert to Console (user_id=%s)\n", sellerUserId)
// 	alertEngine.sendToConsole(msg)
// 	log.Printf("[AlertEngine] [x] Sent AlertSeller Alert (user_id=%s)\n", sellerUserId)
// }

// func (alertEngine *ConsoleAlertEngine) AlertBidder(msg string, bid *Bid) {
// 	bidderUserId := bid.BidderUserId
// 	log.Printf("[AlertEngine] Sending AlertBidder Alert to Console (user_id=%s)\n", bidderUserId)
// 	alertEngine.sendToConsole(msg)
// 	log.Printf("[AlertEngine] [x] Sent AlertBidder Alert (user_id=%s)\n", bidderUserId)
// }

func (alertEngine *ConsoleAlertEngine) sendToConsole(msg string) {
	log.Printf(msg)
}

func (alertEngine *ConsoleAlertEngine) TurnDown() {
	log.Printf("[AlertEngine] shutting down...")
}
