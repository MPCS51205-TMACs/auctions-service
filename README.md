# auctions-service

text / call Austin Clark for details 703-785-1134.

## Overview

This microservice (context) deals with bid capture and auction lifecycle management.

The context offers life-cycle management operations for Auctions (create, stop, cancel, ...) that are exposed through a RESTful API (reached through HTTP requests). It also accepts new bids for auctions (which tentatively come in as RabbitMQ messages). When new bids come in for an auction, the service determines whether it is a new top bid for the auction. If so, it handles sending out alerts to the involved persons, persisting the bid, etc.. Internally, the `auctions-service` also introspects and performs actions periodically. For example, it introspects to determine whether to send out alerts about Auctions that are soon to start / soon to end, whether to load new Auction objects into memory from the database (i.e. all auctions that are active and all auctions that are pending and going to start within the hour), and whether to finalize any Auctions that are over or canceled (described later).

It is written in `Golang` and uses `postgres` (SQL) for data persistence. `Golang` was chosen as a language because it works well for implementing parallelism. This context must perform many activities concurrently; `Golang`'s lightweight goroutine parallelism construct allows us to develop efficient parallel code. 

## Description

Clarifying Note: there are various terms used to describe the lifecycle of an auction. In the `auctions-service`, an auction advances through states based on time elapsation: `PENDING -> ACTIVE -> OVER`. When an auction is in the `PENDING` or `ACTIVE` state, circumstances permitting the auction may be canceled or stopped by a client. In either case, the Auction goes to a `CANCELED` state. In periodic intervals, the `auctions-service` introspects and looks for auctions that are `CANCELED` or `OVER`. If it finds either of these, it will process the end result of that auction (e.g. declare and announce a winner if an auction is `OVER`), and it will advance the Auction to a `FINALIZED` state. When an auction is `FINALIZED`, it means that its result has been determined, and it has been "archived". Its contents are frozen and cannot be mutated.

An auction's state is determined by an internal function, which accepts a `time.Time` object. This allows us to perform retroactive analysis; i.e., we can determine an Auction's state at a given time and interact with it as if the time of the interaction was at the injected `time.Time` time. This leads to flexible behavior. For example, we can stamp a bid with a time when it is received, and we can at a later point in time choose to process the bid as if we were processing it when it was received (rather than at the time it is actually analyzed).

When the `auctions-service` finalizes an auction, it wraps up all data related to the auction and publishes the information to RabbitMQ queue (as one or more messages). Consumers of this data include `closed-auction-metrics`, a context providing functionality to analyze/visualize past auctions. Other consumers way include +contexts like `ShoppingCart` or `Items`.

`postgres` (SQL) was chosen for persistence because Auctions has a very structured schema unlikely to change. For increased performance, `auctions-service` maintains a hash table of `item_id`'s (string) that map to `Auction` objects; this acts as a cache (with no eviction) to reduce the number of reads to the database. Every state change involve a bid or auction involves a write to the database, however, to ensure no lost information if the service shuts down; so, it is unclear how helpful maintaining the inmemory hashtable of `Auction`s is.

## Developer notes

* currently, new bids will come in through the REST API, but we plan to have these come in through RabbitMQ messages.
* Time is stored in Postgres with the `timestamp(6)` datatime, which allows us to store and work with microsecond level precision. All times are interpretted and marked as UTC time when parsed, created, etc.
* currently, we do not make immense use of parallelism. The top-level `AuctionsService` object has a coarse-grained implementation. That is, it has an internal lock that surrounds each of its methods

## File structure / Architecture

This context uses a layered architecture with code written golang. postgres for storage.

```
|-- auction-service/
|   |-- auction-service/
|   |   |-- Dockerfile
|   |   |-- main/
|   |   |   |-- main.go
|   |   |   |-- auction-service.go
|   |   |   |-- ...
|   |   |-- domain/
|   |   |   |-- auction.go
|   |   |   |-- bid.go
|   |   |   |-- ...
|   |   |-- common/
|   |   |   |-- utils.go
|   |   |-- db/
|   |   |   |-- create_db_and_empty_tables.sql
|   |   |   |-- fill_tables_w_data.sql
```

A `main` package holds code dealing with the interface for the service and includes the entrypoint for the microservice: `main.go`. The `domain` package contains domain logic code (e.g. structs like `Auction`,`Bid`,...). `common` contains commonly used functions, and `db` contains files to create a sql database and optionally generate seed data.

`main.go` is the entrypoint for the microservice (i.e. is the "main" function). This script creates an instance of an `AuctionService` and passes a pointer to the `AuctionService` to various goroutines that it spawns. Each of these goroutines generally handle a way of invoking the `auctions-service` microservice. In total, the `AuctionService` can be invoked 3 different ways. It can be invoked through a RESTful API (i.e. a client sends GET/POST requests to specific URL endpoints), it can be invoked through RabbitMQ messages that get published and this context is listening to, and it can be invoked through an `AuctionSessionManager` object that is created internally. An `AuctionSessionManager` object periodically invokes the `AuctionService` to do activities. Currently, the `AuctionSessionManager` periodically invokes the `AuctionService` to 

1. load new Auctions into memory from the SQL database (i.e. that are going to start within an hour from now)
2. peruse the Auctions in memory to see if any need to send out their start soon or end soon alerts
3. finalize any Auctions that are `OVER` or `CANCELED` (i.e. mark the Auction's state as finalized and frozen, package up the information, and publish the information as an Auction-End event)

![alt text](AuctionsDomainModel.png "Title")


## persistence

`create_db_and_empty_tables.sql` and `fill_tables_w_data.sql` are sql scripts that define how the SQL database should be created and what seed data should be used when the database is created. These sql scripts are mounted to `/docker-entrypoint-initdb.d` from a docker-compose file. These SQL scripts are run when the first `docker-compose up -d` is called. 

. `src/db/init` folder into the mongo intitialization script directory for the mongo container (`/docker-entrypoint-initdb.d`). However, I could not manage to generate json time data correctly with this approach, so I have been using `insert_starter_auction_data_into_mongo.py` instead.

a docker-volume `pgdata` is created and mounted to `/var/lib/postgresql/data` (where the postgres db container stores data). This enables persistence of data between `docker-compose up -d` and `docker-compose down` calls. If the user wishes to clear out the database and start from an empty database, they do a `docker volume rm project-dir_pgdata` call, which deletes the local docker-volume storing the persisted data on the host-system. The next call to  `docker-compose up -d` will create the docker-volume again from scratch, so the postgres database container will run the initialization scripts: it will create the database and tables, and it will insert seed data (it will run all sql scripts that have been mounted to `/docker-entrypoint-initdb.d`). When the system is brought down with `docker-compose down`, the data in the postgres container will be saved to the docker-volume `pgdata` on the host machine. Upon the next call to `docker-compose up -d`, the postgres container will recognize data from the host's system, and it will load that data into the database (and it will not perform the initialization scripts).

## deployment

The service involves 3 containers: a container for the actual service `auctions-service`, a container for the persistence with mongo db (`mongo-server`), and a container for the RabbitMQ instance (`rabbitmq-server`) (in full deployment, this service would share a RabbitMQ container with other microservices).

If firing up just this microservice, perform a docker-compose call using `docker-compose.yml`
```
$ docker-compose up -d
```

turn down with
```
$ docker-compose down
```

## to run tests 

There are a fair amount of tests for domain logic, but not a lot of integration tests.

to run Golang tests for a specific package: cd to the package (e.g. `domain`) and run
```
$ go test -v
```

To run a specific test
```
$ go test -v -run REGEX_EXPRESSION_REFERENCING_TEST_METHOD_NAME
$ go test -v -run ^TestProcessBid$
```

## starting the application

If firing up just this microservice, perform a docker-compose call using `docker-compose.yml`
```
$ docker-compose up -d
```

turn down with
```
$ docker-compose down
```