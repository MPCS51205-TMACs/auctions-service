package domain

// // type AuctionRepository interface {
// // 	GetBid(bidId string) *Bid
// // 	GetBidsByUserId(userId string) *[]*Bid
// // 	GetBidsByItemId(itemId string) *[]*Bid
// // 	SaveBid(bid *Bid)
// // 	NextBidId() string
// // }

import (
	"auctions-service/common"
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/lib/pq"
	// _ "github.com/lib/pq" // postgres
)

type postgresSQLAuctionRepository struct {
	db      *sql.DB
	bidRepo BidRepository
}

func NewPostgresSQLAuctionRepository(bidRepo BidRepository) AuctionRepository {

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

	return &postgresSQLAuctionRepository{db, bidRepo}
}

type AuctionData struct {
	ItemId             string
	SellerUserId       string
	StartPriceInCents  int
	StartTime          time.Time
	EndTime            time.Time
	SentStartSoonAlert bool
	SentEndSoonAlert   bool
	TimeCanceled       pq.NullTime // might be null
	FinalizationTime   pq.NullTime // might be null
}

func (repo *postgresSQLAuctionRepository) GetAuction(itemId string) *Auction {

	var result AuctionData

	queryStr := "select auctions.*,auctionscancellations.timeCanceled,auctionsfinalizations.timeFinalized  from auctions \n" +
		"left join auctionsfinalizations \n" +
		"on auctions.itemid = auctionsfinalizations.itemId \n" +
		"left join auctionscancellations \n" +
		"on auctions.itemid = auctionscancellations.itemId \n" +
		fmt.Sprintf("where auctions.itemid = '%s';", itemId)

	rows, err := repo.db.Query(queryStr)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		debug.PrintStack()
		return nil
	}

	for rows.Next() {

		err := rows.Scan(
			&result.ItemId,
			&result.SellerUserId,
			&result.StartPriceInCents,
			&result.StartTime,
			&result.EndTime,
			&result.SentStartSoonAlert,
			&result.SentEndSoonAlert,
			&result.FinalizationTime,
			&result.TimeCanceled,
		)

		if err != nil {
			debug.PrintStack()
			return nil
		}

		itemId := result.ItemId
		sellerUserId := result.SellerUserId
		startPriceInCents := result.StartPriceInCents
		startime := result.StartTime
		endtime := result.EndTime
		sentStartSoonAlert := result.SentStartSoonAlert
		sentEndSoonAlert := result.SentEndSoonAlert
		finalizationTime := result.FinalizationTime
		cancellationTime := result.TimeCanceled

		item := NewItem(itemId, sellerUserId, startime, endtime, int64(startPriceInCents))
		bids := repo.bidRepo.GetBidsByItemId(itemId)

		var cancellation *Cancellation = nil
		if cancellationTime.Valid {
			cancellation = NewCancellation(cancellationTime.Time)
		}

		var finalization *Finalization = nil
		if finalizationTime.Valid {
			finalization = NewFinalization(finalizationTime.Time)
		}

		auction := NewAuction(item, bids, cancellation, sentStartSoonAlert, sentEndSoonAlert, finalization)

		return auction

	}
	return nil
}

func (repo *postgresSQLAuctionRepository) GetAuctions(leftBound time.Time, rightBound time.Time) []*Auction {
	var result AuctionData

	queryStr := "select auctions.*,auctionsfinalizations.timeFinalized,auctionscancellations.timeCanceled from auctions \n" +
		"left join auctionsfinalizations \n" +
		"on auctions.itemid = auctionsfinalizations.itemId \n" +
		"left join auctionscancellations \n" +
		"on auctions.itemid = auctionscancellations.itemId \n" +
		fmt.Sprintf("WHERE  not (auctions.endtime < '%s'::timestamp(6) \n", common.TimeToSQLTimestamp6(leftBound)) +
		fmt.Sprintf("OR auctions.starttime > '%s'::timestamp(6));", common.TimeToSQLTimestamp6(rightBound))

	// fmt.Println(queryStr)
	rows, err := repo.db.Query(queryStr)
	defer rows.Close()

	auctions := []*Auction{}

	if err != nil {
		log.Fatalln(err)
		debug.PrintStack()
		return auctions
	}

	for rows.Next() {

		err := rows.Scan(
			&result.ItemId,
			&result.SellerUserId,
			&result.StartPriceInCents,
			&result.StartTime,
			&result.EndTime,
			&result.SentStartSoonAlert,
			&result.SentEndSoonAlert,
			&result.FinalizationTime,
			&result.TimeCanceled,
		)

		if err != nil {
			fmt.Println(err)
			debug.PrintStack()
			return nil
		}

		itemId := result.ItemId
		sellerUserId := result.SellerUserId
		startPriceInCents := result.StartPriceInCents
		startime := result.StartTime
		endtime := result.EndTime
		sentStartSoonAlert := result.SentStartSoonAlert
		sentEndSoonAlert := result.SentEndSoonAlert
		finalizationTime := result.FinalizationTime
		cancellationTime := result.TimeCanceled

		item := NewItem(itemId, sellerUserId, startime, endtime, int64(startPriceInCents))
		bids := repo.bidRepo.GetBidsByItemId(itemId)

		var cancellation *Cancellation = nil
		if cancellationTime.Valid {
			cancellation = NewCancellation(cancellationTime.Time)
		}

		var finalization *Finalization = nil
		if finalizationTime.Valid {
			finalization = NewFinalization(finalizationTime.Time)
		}

		auction := NewAuction(item, bids, cancellation, sentStartSoonAlert, sentEndSoonAlert, finalization)
		auctions = append(auctions, auction)

	}
	return auctions
}

func (repo *postgresSQLAuctionRepository) SaveAuction(auctionToSave *Auction) {
	// USE UPSERT SYNTAX (insert if not already in db; update if already exists)
	if auctionToSave == nil {
		return
	}

	itemId := auctionToSave.Item.ItemId
	sellerUserId := auctionToSave.Item.SellerUserId
	startPriceInCents := auctionToSave.Item.StartPriceInCents
	startime := auctionToSave.Item.StartTime
	endtime := auctionToSave.Item.EndTime

	var sentStartSoonAlert string = "FALSE"
	if auctionToSave.sentStartSoonAlert {
		sentStartSoonAlert = "TRUE"
	}

	var sentEndSoonAlert string = "FALSE"
	if auctionToSave.sentEndSoonAlert {
		sentEndSoonAlert = "TRUE"
	}

	var timeCanceled pq.NullTime
	if auctionToSave.cancellation != nil {
		timeCanceled = pq.NullTime{auctionToSave.cancellation.TimeReceived, true}
	} else {
		timeCanceled = pq.NullTime{time.Now(), false} // meaningless time
	}

	var timeFinalized pq.NullTime
	if auctionToSave.finalization != nil {
		timeFinalized = pq.NullTime{auctionToSave.finalization.TimeReceived, true}
	} else {
		timeFinalized = pq.NullTime{time.Now(), false} // meaningless time
	}

	// note the following code is not *transactional*
	// save associated cancellation if exists
	if timeCanceled.Valid {
		sqlStr := "INSERT INTO auctionscancellations (itemId, timeCanceled) VALUES \n" +
			fmt.Sprintf("('%s',TIMESTAMP '%s') \n", itemId, common.TimeToSQLTimestamp6(timeCanceled.Time)) +
			"on conflict (itemId) do update \n" +
			"set itemId=excluded.itemId, \n" +
			"timeCanceled=excluded.timeCanceled;"

		_, err := repo.db.Exec(sqlStr)
		if err != nil {
			fmt.Println("got error: ")
			fmt.Println(err)
			debug.PrintStack()
		}
	}

	// save associated finalization if exists
	if timeFinalized.Valid {
		sqlStr := "INSERT INTO auctionsfinalizations (itemId, timeFinalized) VALUES \n" +
			fmt.Sprintf("('%s',TIMESTAMP '%s') \n", itemId, common.TimeToSQLTimestamp6(timeFinalized.Time)) +
			"on conflict (itemId) do update \n" +
			"set itemId=excluded.itemId, \n" +
			"timeFinalized=excluded.timeFinalized;"

		_, err := repo.db.Exec(sqlStr)
		if err != nil {
			fmt.Println("got error: ")
			fmt.Println(err)
			debug.PrintStack()
		}
	}

	// save associated auction
	sqlStr := "INSERT INTO auctions (itemId, sellerUserId, startPriceInCents, startTime, endTime, sentStartSoonAlert, sentEndSoonAlert) VALUES \n" +
		fmt.Sprintf("('%s','%s',%d,TIMESTAMP '%s',TIMESTAMP '%s',%s,%s) \n", itemId, sellerUserId, startPriceInCents, common.TimeToSQLTimestamp6(startime), common.TimeToSQLTimestamp6(endtime), sentStartSoonAlert, sentEndSoonAlert) +
		"on conflict (itemId) do update \n" +
		"set itemId=excluded.itemId, \n" +
		"sellerUserId=excluded.sellerUserId, \n" +
		"startPriceInCents=excluded.startPriceInCents, \n" +
		"startTime=excluded.startTime, \n" +
		"endTime=excluded.endTime, \n" +
		"sentStartSoonAlert=excluded.sentStartSoonAlert, \n" +
		"sentEndSoonAlert=excluded.sentEndSoonAlert;"

	_, err := repo.db.Exec(sqlStr)
	if err != nil {
		fmt.Println("got error: ")
		fmt.Println(err)
		debug.PrintStack()
	}

}

func (repo *postgresSQLAuctionRepository) NumAuctionsSaved() int {
	queryStr := "select count(*) from auctions;"

	var count int
	row := repo.db.QueryRow(queryStr)
	switch err := row.Scan(&count); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		return count
	default:
		debug.PrintStack()
		panic(err)
	}

	return count

}
