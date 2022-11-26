package domain

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go" // acquired by doing 'go get github.com/rabbitmq/amqp091-go'
)

// sends messages to RabbitMQ; uses a provided connection; establishes its own channel(s).
// make sure that within the main function that uses this Engine that the following command(s)
// are added to close connections and channels: defer alertEngine.TurnDown(); defer conn.Close();
// this will respectively close channels internally created by alertEngine, and the latter will
// close the connection to RabbitMQ
type RabbitMQAlertEngine struct {
	conn                      *amqp.Connection
	ch                        *amqp.Channel
	rabbitMQContainerHostName string
	startSoonExchangeName     string
	startSoonQueueName        string
	endSoonExchangeName       string
	endSoonQueueName          string
	auctionEndExchangeName    string
	auctionEndQueueName       string
	newTopBidExchangeName     string
	newTopBidQueueName        string
}

func NewRabbitMQAlertEngine(
	conn *amqp.Connection,
	rabbitMQContainerHostName,
	startSoonExchangeName,
	startSoonQueueName,
	endSoonExchangeName,
	endSoonQueueName,
	auctionEndExchangeName,
	auctionEndQueueName,
	newTopBidExchangeName,
	newTopBidQueueName string) *RabbitMQAlertEngine {

	// rabbitMqContainerHostName := "rabbitmq-server" // e.g. "localhost"
	// exchangeName := "auctionfinalizations"
	// queueName := ""

	// make connection
	// connStr := fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitMqContainerHostName)
	// conn, err := amqp.Dial(connStr)
	// failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	// create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	// declare exchanges / queues

	// AuctionStartSoon Alerts
	err = ch.ExchangeDeclare(
		startSoonExchangeName, // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to declare exchange: "+startSoonExchangeName)

	// AuctionEndSoon Alerts
	err = ch.ExchangeDeclare(
		endSoonExchangeName, // name
		"direct",            // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare exchange: "+endSoonExchangeName)

	// AuctionEnd Alerts
	err = ch.ExchangeDeclare(
		auctionEndExchangeName, // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare exchange: "+auctionEndExchangeName)

	return &RabbitMQAlertEngine{
		conn,
		ch,
		rabbitMQContainerHostName,
		endSoonExchangeName,
		endSoonQueueName,
		startSoonExchangeName,
		startSoonQueueName,
		auctionEndExchangeName,
		auctionEndQueueName,
		newTopBidExchangeName,
		newTopBidQueueName}
}

// func (auction *RabbitMQAlertEngine) ConnectToRabbit

func (alertEngine *RabbitMQAlertEngine) SendAuctionStartSoonAlert(msg string, itemId string, startTime time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(msg)
	// fmt.Println(body)
	failOnError(err, "[AlertEngine] Error encoding JSON")

	log.Printf("[AlertEngine] Sending AuctionStartingSoon Alert (item_id=%s) to RabbitMQ exchange '%s'\n", itemId, alertEngine.startSoonExchangeName)
	// log.Printf("[AlertEngine] Sending Auction data to RabbitMQ: %s\n", body)
	err = alertEngine.ch.PublishWithContext(ctx,
		alertEngine.startSoonExchangeName, // exchange
		"",                                // routing key WITH QUEUE q.Name
		false,                             // mandatory
		false,                             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	// amqp.Publishing{
	// 	ContentType: "text/plain",
	// 	Body:        []byte(body),
	// })
	failOnError(err, "[AlertEngine] Failed to publish a message")
	log.Printf("[AlertEngine] [x] Sent AuctionStartingSoon Alert (item_id=%s)\n", itemId)
}

func (alertEngine *RabbitMQAlertEngine) SendAuctionEndSoonAlert(msg string, itemId string, endTime time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(msg)
	// fmt.Println(body)
	failOnError(err, "[AlertEngine] Error encoding JSON")

	log.Printf("[AlertEngine] Sending AuctionEndingSoon Alert (item_id=%s) to RabbitMQ exchange '%s'\n", itemId, alertEngine.endSoonExchangeName)
	// log.Printf("[AlertEngine] Sending Auction data to RabbitMQ: %s\n", body)
	err = alertEngine.ch.PublishWithContext(ctx,
		alertEngine.startSoonExchangeName, // exchange
		"",                                // routing key WITH QUEUE q.Name
		false,                             // mandatory
		false,                             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	// amqp.Publishing{
	// 	ContentType: "text/plain",
	// 	Body:        []byte(body),
	// })
	failOnError(err, "[AlertEngine] Failed to publish a message")
	log.Printf("[AlertEngine] [x] Sent AuctionEndingSoon Alert (item_id=%s)\n", itemId)
}

type NewTopBidData struct {
	ItemId          string  `json:"itemid"`
	SellerUserId    string  `json:"selleruserid"`
	FormerTopBidder *string `json:"starttime"`
	NewTopBidder    string  `json:"endtime"`
}

func NewNewTopBidData(itemId, sellerUserId, formerTopBidderUserId, newTopBidderUserId *string) *NewTopBidData {
	return &NewTopBidData{
		ItemId:          *itemId,
		SellerUserId:    *sellerUserId,
		FormerTopBidder: formerTopBidderUserId,
		NewTopBidder:    *newTopBidderUserId,
	}
}

func (alertEngine *RabbitMQAlertEngine) SendNewTopBidAlert(itemId, sellerUserId, formerTopBidderUserId, newTopBidderUserId *string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newTopBidData := NewNewTopBidData(itemId, sellerUserId, formerTopBidderUserId, newTopBidderUserId)

	body, err := json.Marshal(*newTopBidData)
	// fmt.Println(body)
	failOnError(err, "[AlertEngine] Error encoding JSON")

	log.Printf("[AlertEngine] Sending Auction data to RabbitMQ (item_id=%s) to RabbitMQ exchange '%s'\n", *itemId, alertEngine.newTopBidExchangeName)
	// log.Printf("[AlertEngine] Sending Auction data to RabbitMQ: %s\n", body)
	err = alertEngine.ch.PublishWithContext(ctx,
		alertEngine.newTopBidExchangeName, // exchange
		"",                                // routing key WITH QUEUE q.Name
		false,                             // mandatory
		false,                             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	// amqp.Publishing{
	// 	ContentType: "text/plain",
	// 	Body:        []byte(body),
	// })
	failOnError(err, "[AlertEngine] Failed to publish a message")
	log.Printf("[AlertEngine] [x] Sent (item_id=%s)\n", *itemId)
}

func (alertEngine *RabbitMQAlertEngine) sendAuctionEndAlert(auctionData *AuctionData) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(*auctionData)
	// fmt.Println(body)
	failOnError(err, "[AlertEngine] Error encoding JSON")

	log.Printf("[AlertEngine] Sending Auction data (item_id=%s) to RabbitMQ exchange '%s'\n", auctionData.Item.ItemId, alertEngine.auctionEndExchangeName)
	// log.Printf("[AlertEngine] Sending Auction data to RabbitMQ: %s\n", body)
	err = alertEngine.ch.PublishWithContext(ctx,
		alertEngine.auctionEndExchangeName, // exchange
		alertEngine.auctionEndQueueName,    // routing key WITH QUEUE q.Name
		false,                              // mandatory
		false,                              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	// amqp.Publishing{
	// 	ContentType: "text/plain",
	// 	Body:        []byte(body),
	// })
	failOnError(err, "[AlertEngine] Failed to publish a message")
	log.Printf("[AlertEngine] [x] Sent (item_id=%s)\n", auctionData.Item.ItemId)
}

func (alertEngine *RabbitMQAlertEngine) TurnDown() {
	alertEngine.ch.Close()
}

// sendAuctionEndAlert(finalizedAuction *AuctionData)

// SendNewTopBidAlert(itemId, sellerUserId, formerTopBidderUserId, newTopBidderUserId *string)

// func (alertEngine *ConsoleAlertEngine) AlertSeller(msg string, sellerUserId string) {
// 	log.Printf("[AlertEngine] Sending AlertSeller Alert to Console (user_id=%s)\n", sellerUserId)
// 	alertEngine.sendToConsole(msg)
// 	log.Printf("[AlertEngine] [x] Sent AlertSeller Alert (user_id=%s)\n", sellerUserId)
// }

// func (alertEngine *ConsoleAlertEngine) AlertBidder(msg string, bidderUserId string) {
// 	log.Printf("[AlertEngine] Sending AlertBidder Alert to Console (user_id=%s)\n", bidderUserId)
// 	alertEngine.sendToConsole(msg)
// 	log.Printf("[AlertEngine] [x] Sent AlertBidder Alert (user_id=%s)\n", bidderUserId)
// }

// func (alertEngine *RabbitMQAlertEngine) TurnDown() {
// 	log.Printf("[AlertEngine] shutting down...")
// 	alertEngine.ch.Close()
// }

// func sendAuctionDataToRabbitMQ(auction *Auction) {

// 	rabbitMqContainerHostName := "rabbitmq-server" // e.g. "localhost"
// 	exchangeName := "auctionfinalizations"
// 	queueName := ""

// 	// make connection
// 	connStr := fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitMqContainerHostName)
// 	conn, err := amqp.Dial(connStr)
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	// create a channel
// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	// declare exchange
// 	err = ch.ExchangeDeclare(
// 		exchangeName, // name
// 		"fanout",     // type
// 		true,         // durable
// 		false,        // auto-deleted
// 		false,        // internal
// 		false,        // no-wait
// 		nil,          // arguments
// 	)
// 	failOnError(err, "Failed to declare an exchange")

// 	// // declare queue for us to send messages to
// 	// q, err := ch.QueueDeclare(
// 	// 	queueName, // name
// 	// 	true,      // durable
// 	// 	false,     // delete when unused
// 	// 	false,     // exclusive
// 	// 	false,     // no-wait
// 	// 	nil,       // arguments
// 	// )
// 	// failOnError(err, "Failed to declare a queue")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	body, err := json.Marshal(*(auction.ToAuctionData()))
// 	// fmt.Println(body)
// 	failOnError(err, "Error encoding JSON")

// 	err = ch.PublishWithContext(ctx,
// 		exchangeName, // exchange
// 		queueName,    // routing key WITH QUEUE q.Name
// 		false,        // mandatory
// 		false,        // immediate
// 		amqp.Publishing{
// 			ContentType:  "application/json",
// 			Body:         body,
// 			DeliveryMode: amqp.Persistent,
// 		})
// 	// amqp.Publishing{
// 	// 	ContentType: "text/plain",
// 	// 	Body:        []byte(body),
// 	// })
// 	failOnError(err, "Failed to publish a message")
// 	log.Printf(" [x] Sent %s\n", body)
// }
