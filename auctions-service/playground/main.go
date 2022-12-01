package main

import (
	"auctions-service/domain"
	"fmt"
	"log"
	"time"

	// acquired by doing 'go get github.com/gorilla/mux.git'
	_ "github.com/lib/pq"                 // postgres
	amqp "github.com/rabbitmq/amqp091-go" // acquired by doing 'go get github.com/rabbitmq/amqp091-go'
	// go get -u github.com/jkeys089/jserial
	// go get -u github.com/golang-jwt/jwt/v4
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func fillReposWDummyData(bidRepo domain.BidRepository, auctionRepo domain.AuctionRepository) []*domain.Auction {
	// fill bid repo with some bids
	time1 := time.Now()
	time2 := time1                                          // same as time1
	time3 := time1.Add(time.Duration(1) * time.Microsecond) // 1 microsecond after
	time4 := time1.Add(time.Duration(1) * time.Second)      // 1 sec after
	bid1 := *domain.NewBid("101", "20", "asclark109", time1, int64(300), true)
	bid2 := *domain.NewBid("102", "20", "mcostigan9", time2, int64(300), true)
	bid3 := *domain.NewBid("103", "20", "katharine2", time3, int64(400), true)
	bid4 := *domain.NewBid("104", "20", "katharine2", time4, int64(10), true)
	bidRepo.SaveBid(&bid1)
	bidRepo.SaveBid(&bid2)
	bidRepo.SaveBid(&bid3)
	bidRepo.SaveBid(&bid4)

	startime := time1.Add(-time.Duration(1) * time.Hour)
	endtime := time1.Add(time.Duration(10) * time.Hour)
	item1 := domain.NewItem("101", "asclark109", startime, endtime, int64(2000)) // $20 start price
	item2 := domain.NewItem("102", "asclark109", startime, endtime, int64(2000)) // $20 start price
	auction1 := auctionRepo.NewAuction(item1, nil, nil, false, false, nil)       // will go to completion
	auction2 := auctionRepo.NewAuction(item2, nil, nil, false, false, nil)       // will get cancelled halfway through

	nowtime := time.Now()
	latertime := nowtime.Add(time.Duration(4) * time.Hour)                        // 4 hrs from now
	item3 := domain.NewItem("103", "asclark109", nowtime, latertime, int64(2000)) // $20 start price
	auctionactive := auctionRepo.NewAuction(item3, nil, nil, false, false, nil)

	latertime2 := nowtime.Add(time.Duration(2) * time.Hour)                        // 2 hrs from now
	item4 := domain.NewItem("104", "asclark109", nowtime, latertime2, int64(2000)) // $20 start price
	auctionactive2 := auctionRepo.NewAuction(item4, nil, nil, false, false, nil)

	canceled := auction2.Cancel(startime)
	finalized := auction2.Finalize(startime.Add(time.Duration(1) * time.Second))
	fmt.Println(canceled)
	fmt.Println(finalized)

	auctionRepo.SaveAuction(auction1)
	auctionRepo.SaveAuction(auction2)
	auctionRepo.SaveAuction(auctionactive)
	auctionRepo.SaveAuction(auctionactive2)

	var auctions []*domain.Auction = []*domain.Auction{auction1, auction2, auctionactive, auctionactive2}
	return auctions

}

func main() {
	// RABBITMQ PARAMS
	// parameters below need to be manually changed with deployment
	// note: variables names below may need refactoring or deletion;
	// some are unused (e.g. publishing to an exchange does need knowledge of a binding queue),
	// and sometimes these variables are put into the arguments for routing keys. needed refactor.
	rabbitMQContainerHostName := "localhost"
	rabbitMQContainerPort := "5672"

	startSoonExchangeName := "auction.start-soon"
	startSoonQueueName := ""

	endSoonExchangeName := "auction.end-soon"
	endSoonQueueName := ""

	auctionEndExchangeName := "auction.end"
	auctionEndQueueName := ""

	// newBidExchangeName := "auction.new-bid"
	// newBidQueueName := "auction.process-bid"

	newTopBidExchangeName := "auction.new-high-bid"
	newTopBidQueueName := ""

	// ItemCounterfeitExchangeName := "item.counterfeit"
	// ItemCounterfeitQueueName := "auction.consume_counterfeit"

	// ItemInappropriateExchangeName := "item.inappropriate"
	// ItemInappropriateQueueName := "auction.consume_inappropriate"

	// UserDeleteExchangeName := "user.delete"
	// UserDeleteQueueName := "auction.consume_userdelete"

	// UserActivationExchangeName := "user.activation"
	// UserActivationQueueName := "auction.consume_useractivation"

	// UserUpdateExchangeName := "user.update"
	// UserUpdateQueueName := "auction.consume_userupdate"

	var conn *amqp.Connection

	// intialize OUTBOUND ADAPTER (AlertEngine)
	var alertEngine domain.AlertEngine

	// typical to have 1 RabbitMQ connection per application, 1 channel per thread, and only 1 thread publishing/subscribing at any moment
	connStr := fmt.Sprintf("amqp://guest:guest@%s:%s/", rabbitMQContainerHostName, rabbitMQContainerPort)
	var err error
	conn, err = amqp.Dial(connStr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close() // close connection on application exit

	fmt.Println("OUTBOUND ADAPTER = deployment [outbound msgs go to RabbitMQ / other contexts]...")

	alertEngine = domain.NewRabbitMQAlertEngine( // creates a channel internally with connection
		conn,
		rabbitMQContainerHostName,
		startSoonExchangeName,
		startSoonQueueName,
		endSoonExchangeName,
		endSoonQueueName,
		auctionEndExchangeName,
		auctionEndQueueName,
		newTopBidExchangeName,
		newTopBidQueueName,
	)
	defer alertEngine.TurnDown() // close channel(s) on application exit

	// intialize REPOSITORIES
	var bidRepo domain.BidRepository
	var auctionRepo domain.AuctionRepository

	fmt.Println("REPOSITORY TYPE = in-memory [no persistence; everything held in-memory]")
	bidRepo = domain.NewInMemoryBidRepository(false) // do not use seed; assign random uuid's to new Bids
	auctionRepo = domain.NewInMemoryAuctionRepository(alertEngine)

	// seed with data
	auctions := fillReposWDummyData(bidRepo, auctionRepo)

	auction := auctions[0] // not canceled, not finalized, has a few bids

	// SEND AUCTION.NEW-TOP-BID
	send_new_high_bid(auction, alertEngine)

	// SEND AUCTION.START-SOON
	send_start_soon(auction, alertEngine)

	// // SEND AUCTION.END-SOON
	send_end_soon(auction, alertEngine)

	// // SEND AUCTION.END
	auction2 := auctions[1] // canceled and finalized
	send_end(auction2, alertEngine)
}

func send_new_high_bid(auction *domain.Auction, alertEngine domain.AlertEngine) {
	var itemId *string = &auction.Item.ItemId
	var seller *string = &auction.Item.SellerUserId
	newTopBid := *domain.NewBid("105", "20", "Mary", time.Now(), int64(100000), true)
	var formerTopBidder *string = nil
	var newTopBidder *string = &newTopBid.BidderUserId
	alertEngine.SendNewTopBidAlert(itemId, seller, formerTopBidder, newTopBidder)
}

func send_start_soon(auction *domain.Auction, alertEngine domain.AlertEngine) {
	var itemId *string = &auction.Item.ItemId
	startTime := auction.Item.StartTime
	endTime := auction.Item.EndTime
	msg := "dummy msg"
	alertEngine.SendAuctionStartSoonAlert(msg, *itemId, startTime, endTime)
}

func send_end_soon(auction *domain.Auction, alertEngine domain.AlertEngine) {
	var itemId *string = &auction.Item.ItemId
	startTime := auction.Item.StartTime
	endTime := auction.Item.EndTime
	msg := "dummy msg"
	alertEngine.SendAuctionEndSoonAlert(msg, *itemId, startTime, endTime)
}

func send_end(auction *domain.Auction, alertEngine domain.AlertEngine) {
	alertEngine.SendAuctionEndAlert(auction.ToAuctionData())
}
