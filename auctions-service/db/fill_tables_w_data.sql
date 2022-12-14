-- two calls below are basically saying "trunicate table if table exists"
-- make sure we are using auctiondb
\c auctiondb;


CREATE TABLE IF NOT exists auctions (
    itemId varchar(255) PRIMARY KEY,
    sellerUserId varchar(255) NOT NULL,
	startPriceInCents BIGINT NOT NULL,
    startTime timestamp(6) NOT NULL,
    endTime timestamp(6) NOT NULL,
	sentStartSoonAlert boolean NOT NULL,
	sentEndSoonAlert boolean NOT NULL
);
TRUNCATE TABLE auctions;

CREATE TABLE IF NOT exists auctionsCancellations (
    itemId varchar(255) PRIMARY KEY,
    timeCanceled timestamp(6) NOT NULL
);
TRUNCATE TABLE auctionsCancellations;

CREATE TABLE IF NOT exists auctionsFinalizations (
    itemId varchar(255) PRIMARY KEY,
    timeFinalized timestamp(6) NOT NULL
);
TRUNCATE TABLE auctionsFinalizations;


CREATE TABLE IF NOT exists bids (
    bidId varchar(255) PRIMARY KEY,
    itemId varchar(255) NOT NULL,
    bidderUserId varchar(255) NOT NULL,
    amountInCents BIGINT NOT NULL,
    timeBidProcessed timestamp(6) NOT NULL,
	active boolean NOT NULL
);
TRUNCATE TABLE bids;

-- insert some starter data

INSERT INTO auctions (itemId, sellerUserId, startPriceInCents, startTime, endTime, sentStartSoonAlert, sentEndSoonAlert) VALUES
	('200','270',7900,TIMESTAMP '2022-12-01 15:00:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('201','336',100,TIMESTAMP '2022-12-01 15:00:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('202','203',9500,TIMESTAMP '2022-12-01 15:00:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('203','247',4600,TIMESTAMP '2022-12-01 15:00:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('204','281',8900,TIMESTAMP '2022-12-01 15:00:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('205','248',2600,TIMESTAMP '2022-12-01 15:00:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('206','371',7500,TIMESTAMP '2022-12-01 15:00:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('207','215',8300,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('208','242',2900,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('209','226',11400,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('210','312',7300,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('211','314',6500,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('212','292',1200,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('213','254',9900,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('214','301',11000,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('215','245',10500,TIMESTAMP '2022-12-01 15:30:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('216','256',900,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('217','364',11200,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('218','256',8300,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('219','361',2900,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('220','256',7500,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('221','222',5900,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('222','269',300,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('223','367',1300,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('224','241',5900,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('225','227',2200,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('226','225',8200,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('227','263',2200,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('228','323',10000,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 17:00:00.000000',FALSE,FALSE),
	('229','363',10300,TIMESTAMP '2022-12-01 15:45:00.000000',TIMESTAMP '2022-12-01 16:00:00.000000',FALSE,FALSE),
	('230','385',5000,TIMESTAMP '2020-03-04 15:00:00.000000',TIMESTAMP '2020-03-05 15:00:00.000000',TRUE,TRUE);


INSERT INTO auctionsFinalizations (itemId, timeFinalized) VALUES
	('230',TIMESTAMP '2020-03-05 15:30:00.000000');

INSERT INTO auctionsCancellations (itemId, timeCanceled) VALUES
	('209',TIMESTAMP '2022-11-01 16:00:00.000000'),
	('212',TIMESTAMP '2022-11-07 16:00:00.000000');

INSERT INTO bids (bidId, itemId, bidderUserId, amountInCents, timeBidProcessed, active) VALUES
	('100','200','431',7914,TIMESTAMP '2022-12-01 15:15:00.000000',TRUE),
	('101','201','402',126,TIMESTAMP '2022-12-01 15:16:00.000000',TRUE),
	('102','202','406',9503,TIMESTAMP '2022-12-01 15:15:00.000000',TRUE),
	('103','203','484',4630,TIMESTAMP '2022-12-01 15:18:00.000000',TRUE),
	('104','204','431',8928,TIMESTAMP '2022-12-01 15:19:00.000000',TRUE),
	('105','205','468',2619,TIMESTAMP '2022-12-01 15:16:00.000000',TRUE),
	('106','206','443',7529,TIMESTAMP '2022-12-01 15:18:00.000000',TRUE),
	('107','207','426',8309,TIMESTAMP '2022-12-01 15:31:00.000000',TRUE),
	('108','208','433',2905,TIMESTAMP '2022-12-01 15:32:00.000000',TRUE),
	-- ('109','209','404',11410,TIMESTAMP '2022-12-01 15:30:00.000000',TRUE),
	('110','210','435',7314,TIMESTAMP '2022-12-01 15:40:00.000000',TRUE),
	('111','211','430',6523,TIMESTAMP '2022-12-01 15:41:00.000000',TRUE),
	-- ('112','212','480',1205,TIMESTAMP '2022-12-01 15:38:00.000000',TRUE),
	('113','213','440',9926,TIMESTAMP '2022-12-01 15:32:00.000000',TRUE),
	('114','214','427',11011,TIMESTAMP '2022-12-01 15:32:00.000000',TRUE),
	('115','215','482',10500,TIMESTAMP '2022-12-01 15:32:00.000000',TRUE),
	('116','216','413',900,TIMESTAMP '2022-12-01 15:45:00.000000',TRUE),
	('117','217','496',11219,TIMESTAMP '2022-12-01 15:46:00.000000',TRUE),
	('118','218','438',8318,TIMESTAMP '2022-12-01 15:50:00.000000',TRUE),
	('119','219','418',2930,TIMESTAMP '2022-12-01 15:42:00.000000',TRUE),
	('120','220','429',7501,TIMESTAMP '2022-12-01 15:45:00.000000',TRUE),
	('121','221','478',5920,TIMESTAMP '2022-12-01 15:46:00.000000',TRUE),
	('122','222','499',308,TIMESTAMP '2022-12-01 15:50:00.000000',TRUE),
	('123','223','406',1320,TIMESTAMP '2022-12-01 15:45:00.000000',TRUE),
	('124','224','488',5921,TIMESTAMP '2022-12-01 15:46:00.000000',TRUE),
	('125','225','427',2224,TIMESTAMP '2022-12-01 15:45:00.000000',TRUE),
	('126','226','451',8200,TIMESTAMP '2022-12-01 15:45:00.000000',TRUE),
	('127','227','459',2207,TIMESTAMP '2022-12-01 15:46:00.000000',TRUE),
	('128','228','462',10019,TIMESTAMP '2022-12-01 15:50:00.000000',TRUE),
	('129','229','475',10323,TIMESTAMP '2022-12-01 15:42:00.000000',TRUE),
	('130','200','431',7926,TIMESTAMP '2022-12-01 15:17:00.000000',TRUE),
	('131','201','463',137,TIMESTAMP '2022-12-01 15:18:00.000000',TRUE),
	('132','202','496',9528,TIMESTAMP '2022-12-01 15:17:00.000000',TRUE),
	('133','203','420',4648,TIMESTAMP '2022-12-01 15:20:00.000000',TRUE),
	('134','204','409',8946,TIMESTAMP '2022-12-01 15:21:00.000000',TRUE),
	('135','205','446',2635,TIMESTAMP '2022-12-01 15:18:00.000000',TRUE),
	('136','206','447',7540,TIMESTAMP '2022-12-01 15:20:00.000000',TRUE),
	('137','207','454',8337,TIMESTAMP '2022-12-01 15:33:00.000000',FALSE),
	('138','208','472',2920,TIMESTAMP '2022-12-01 15:34:00.000000',TRUE),
	-- ('139','209','490',11414,TIMESTAMP '2022-12-01 15:32:00.000000',TRUE),
	('140','210','469',7330,TIMESTAMP '2022-12-01 15:42:00.000000',TRUE),
	('141','211','493',6529,TIMESTAMP '2022-12-01 15:43:00.000000',TRUE),
	-- ('142','212','419',1232,TIMESTAMP '2022-12-01 15:40:00.000000',TRUE),
	('143','213','488',9948,TIMESTAMP '2022-12-01 15:34:00.000000',TRUE),
	('144','214','464',11035,TIMESTAMP '2022-12-01 15:34:00.000000',TRUE),
	('145','215','471',10501,TIMESTAMP '2022-12-01 15:34:00.000000',TRUE),
	('146','216','405',911,TIMESTAMP '2022-12-01 15:47:00.000000',TRUE),
	('147','217','410',11245,TIMESTAMP '2022-12-01 15:48:00.000000',TRUE),
	('148','218','430',8320,TIMESTAMP '2022-12-01 15:52:00.000000',TRUE),
	('149','219','489',2947,TIMESTAMP '2022-12-01 15:44:00.000000',TRUE),
	('150','220','499',7521,TIMESTAMP '2022-12-01 15:47:00.000000',TRUE),
	('151','221','407',5926,TIMESTAMP '2022-12-01 15:48:00.000000',TRUE),
	('152','222','406',337,TIMESTAMP '2022-12-01 15:52:00.000000',TRUE),
	('153','223','406',1321,TIMESTAMP '2022-12-01 15:47:00.000000',TRUE),
	('154','224','459',5937,TIMESTAMP '2022-12-01 15:48:00.000000',TRUE),
	('155','225','485',2239,TIMESTAMP '2022-12-01 15:47:00.000000',TRUE),
	('156','226','438',8209,TIMESTAMP '2022-12-01 15:47:00.000000',TRUE),
	('157','227','500',2219,TIMESTAMP '2022-12-01 15:48:00.000000',TRUE),
	('158','228','436',10026,TIMESTAMP '2022-12-01 15:52:00.000000',TRUE),
	('159','229','486',10331,TIMESTAMP '2022-12-01 15:44:00.000000',TRUE);


select * from auctions a;
select * from auctionsFinalizations aF;
select * from auctionsCancellations aC;
select * from bids b;
