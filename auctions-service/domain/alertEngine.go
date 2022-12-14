package domain

import (
	"time"
)

type AlertEngine interface {
	SendAuctionStartSoonAlert(msg, itemId string, startTime, endTime time.Time)
	SendAuctionEndSoonAlert(msg, itemId string, startTime, endTime time.Time)
	SendAuctionEndAlert(finalizedAuction *AuctionData)
	SendNewTopBidAlert(itemId, sellerUserId, formerTopBidderUserId, newTopBidderUserId *string)
	// AlertSeller(msg, itemId, sellerUserId string)
	// AlertBidder(msg string, bid *Bid)
	TurnDown()
}
