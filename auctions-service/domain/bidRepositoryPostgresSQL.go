package domain

import (
	"auctions-service/common"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // postgres
)

// type BidRepository interface {
// 	GetBid(bidId string) *Bid
// 	GetBidsByUserId(userId string) *[]*Bid
// 	GetBidsByItemId(itemId string) *[]*Bid
// 	SaveBid(bid *Bid)
// 	SaveBids(bid *Bid)
// 	NextBidId() string
// }

type postgresSQLBidRepository struct {
	db *sql.DB
}

func NewPostgresSQLBidRepository(useDeterministicSeed bool) BidRepository {

	postgresUsername := "postgres"
	postgresPassword := "mysecret"
	postgresContainerhost := "postgres-server"
	postgresContainerport := "5432"
	postgresDbName := "auctiondb"
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", postgresUsername, postgresPassword, postgresContainerhost, postgresContainerport, postgresDbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if useDeterministicSeed {
		seed := int64(1)
		rnd := rand.New(rand.NewSource(seed))
		uuid.SetRand(rnd)
	}

	return &postgresSQLBidRepository{db}
}

type BidData struct {
	BidId            string
	ItemId           string
	BidderUserId     string
	AmountInCents    int
	TimeBidProcessed time.Time
	Active           bool
}

func (repo *postgresSQLBidRepository) GetBid(bidId string) *Bid {

	var result BidData
	queryStr := fmt.Sprintf("SELECT * FROM bids WHERE bidid = '%s'", bidId)
	rows, err := repo.db.Query(queryStr)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	for rows.Next() {
		// bidId varchar(255) PRIMARY KEY,
		// itemId varchar(255) NOT NULL,
		// bidderUserId varchar(255) NOT NULL,
		// amountInCents BIGINT NOT NULL,
		// timeBidProcessed timestamp(6) NOT NULL,
		// active boolean NOT NULL

		err := rows.Scan(
			&result.BidId,
			&result.ItemId,
			&result.BidderUserId,
			&result.AmountInCents,
			&result.TimeBidProcessed,
			&result.Active,
		)

		if err != nil {
			return nil
		}

		bid := NewBid(result.BidId, result.ItemId, result.BidderUserId, result.TimeBidProcessed, int64(result.AmountInCents), result.Active)
		return bid // returns first bid found matching
	}
	return nil
}

func (repo *postgresSQLBidRepository) GetBidsByUserId(biddeUserId string) *[]*Bid {
	var result BidData
	queryStr := fmt.Sprintf("SELECT * FROM bids WHERE bidderuserid = '%s'", biddeUserId)
	rows, err := repo.db.Query(queryStr)
	defer rows.Close()

	bids := []*Bid{}

	if err != nil {
		log.Fatalln(err)
		return &bids
	}

	for rows.Next() {
		// bidId varchar(255) PRIMARY KEY,
		// itemId varchar(255) NOT NULL,
		// bidderUserId varchar(255) NOT NULL,
		// amountInCents BIGINT NOT NULL,
		// timeBidProcessed timestamp(6) NOT NULL,
		// active boolean NOT NULL

		err := rows.Scan(
			&result.BidId,
			&result.ItemId,
			&result.BidderUserId,
			&result.AmountInCents,
			&result.TimeBidProcessed,
			&result.Active,
		)

		if err != nil {
			// http.NotFound(w, r)
			return &bids
		}

		bid := NewBid(result.BidId, result.ItemId, result.BidderUserId, result.TimeBidProcessed, int64(result.AmountInCents), result.Active)
		bids = append(bids, bid)

	}
	return &bids
}

func (repo *postgresSQLBidRepository) GetBidsByItemId(itemId string) *[]*Bid {
	var result BidData
	queryStr := fmt.Sprintf("SELECT * FROM bids WHERE itemid = '%s'", itemId)
	rows, err := repo.db.Query(queryStr)
	defer rows.Close()

	bids := []*Bid{}

	if err != nil {
		log.Fatalln(err)
		return &bids
	}

	for rows.Next() {
		// bidId varchar(255) PRIMARY KEY,
		// itemId varchar(255) NOT NULL,
		// bidderUserId varchar(255) NOT NULL,
		// amountInCents BIGINT NOT NULL,
		// timeBidProcessed timestamp(6) NOT NULL,
		// active boolean NOT NULL

		err := rows.Scan(
			&result.BidId,
			&result.ItemId,
			&result.BidderUserId,
			&result.AmountInCents,
			&result.TimeBidProcessed,
			&result.Active,
		)

		if err != nil {
			return &bids
		}

		bid := NewBid(result.BidId, result.ItemId, result.BidderUserId, result.TimeBidProcessed, int64(result.AmountInCents), result.Active)
		bids = append(bids, bid)

	}
	return &bids
}

func (repo *postgresSQLBidRepository) SaveBid(bidToSave *Bid) {
	// USE UPSERT SYNTAX (insert if not already in db; update if already exists)
	if bidToSave == nil {
		return
	}

	bidId := bidToSave.BidId
	itemId := bidToSave.ItemId
	bidderUserId := bidToSave.BidderUserId
	amountInCents := bidToSave.AmountInCents
	timeBidProcessed := common.TimeToSQLTimestamp6(bidToSave.TimeReceived)
	var active string
	if bidToSave.active {
		active = "TRUE"
	} else {
		active = "FALSE"
	}

	sqlStr := "INSERT INTO bids (bidId, itemId, bidderUserId, amountInCents, timeBidProcessed, active)\n" +
		fmt.Sprintf("VALUES ('%s','%s','%s',%d,TIMESTAMP '%s',%s)\n", bidId, itemId, bidderUserId, amountInCents, timeBidProcessed, active) +
		"on conflict (bidId) do update\n" +
		"set itemId=excluded.itemId,\n" +
		"bidderUserId=excluded.bidderUserId,\n" +
		"amountInCents=excluded.amountInCents,\n" +
		"timeBidProcessed=excluded.timeBidProcessed,\n" +
		"active=excluded.active;"

	_, err := repo.db.Exec(sqlStr)
	if err != nil {
		fmt.Println("got error: ")
		fmt.Println(err)
	}
}

func (repo *postgresSQLBidRepository) SaveBids(bidsToSave *[]*Bid) {
	// USE UPSERT SYNTAX (insert if not already in db; update if already exists)
	if len(*bidsToSave) == 0 {
		return
	}

	sqlStr := "INSERT INTO bids (bidId, itemId, bidderUserId, amountInCents, timeBidProcessed, active)\n"

	for idx, bidToSave := range *bidsToSave {
		bidId := bidToSave.BidId
		itemId := bidToSave.ItemId
		bidderUserId := bidToSave.BidderUserId
		amountInCents := bidToSave.AmountInCents
		timeBidProcessed := common.TimeToSQLTimestamp6(bidToSave.TimeReceived)
		var active string
		if bidToSave.active {
			active = "TRUE"
		} else {
			active = "FALSE"
		}
		if idx == 0 {
			sqlStr += fmt.Sprintf("VALUES ")
		}
		sqlStr += fmt.Sprintf("('%s','%s','%s',%d,TIMESTAMP '%s',%s)", bidId, itemId, bidderUserId, amountInCents, timeBidProcessed, active)
		if idx != len(*bidsToSave)-1 { // if not last idx
			sqlStr += ",\n"
		} else {
			sqlStr += "\n"
		}

	}

	sqlStr += "on conflict (bidId) do update\n" +
		"set itemId=excluded.itemId,\n" +
		"bidderUserId=excluded.bidderUserId,\n" +
		"amountInCents=excluded.amountInCents,\n" +
		"timeBidProcessed=excluded.timeBidProcessed,\n" +
		"active=excluded.active;"

	_, err := repo.db.Exec(sqlStr)
	if err != nil {
		fmt.Println("got error: ")
		fmt.Println(err)
	}
}

func (repo *postgresSQLBidRepository) DeleteBid(bidId string) {

	sqlStr := fmt.Sprintf("DELETE FROM bids WHERE bidId = '%s';", bidId)

	_, err := repo.db.Exec(sqlStr)
	if err != nil {
		fmt.Println("got error: ")
		fmt.Println(err)
	}
}

func (repo *postgresSQLBidRepository) NextBidId() string {
	return uuid.New().String()
}
