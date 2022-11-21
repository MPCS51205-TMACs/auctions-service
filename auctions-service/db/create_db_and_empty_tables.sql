-- create_db_and_empty_tables.sql
-- code borrowed from https://hashinteractive.com/blog/docker-compose-up-with-postgres-quick-tips/
DROP DATABASE IF EXISTS auctiondb;

CREATE DATABASE auctiondb;

-- Make sure we're using our `auctiondb` database
\c auctiondb;

DROP TABLE IF EXISTS auctions;
DROP TABLE IF EXISTS auctionsCancellations;
DROP TABLE IF EXISTS auctionsFinalizations;
DROP TABLE IF EXISTS bids;

CREATE TABLE auctions (
    itemId varchar(255) PRIMARY KEY,
    sellerUserId varchar(255) NOT NULL,
    startPriceInCents BIGINT NOT NULL,
    startTime timestamp(6) NOT NULL,
    endTime timestamp(6) NOT NULL,
    sentStartSoonAlert boolean NOT NULL,
	sentEndSoonAlert boolean NOT NULL
);

CREATE TABLE auctionsCancellations (
    itemId varchar(255) PRIMARY KEY,
    timeCanceled timestamp(6) NOT NULL
);

CREATE TABLE auctionsFinalizations (
    itemId varchar(255) PRIMARY KEY,
    timeFinalized timestamp(6) NOT NULL
);

CREATE TABLE bids (
    bidId varchar(255) PRIMARY KEY,
    itemId varchar(255) NOT NULL,
    bidderUserId varchar(255) NOT NULL,
    amountInCents BIGINT NOT NULL,
    timeBidProcessed timestamp(6) NOT NULL,
    active boolean NOT NULL
);