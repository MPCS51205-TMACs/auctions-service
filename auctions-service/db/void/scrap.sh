# some random bash things I already wrote and may want to reference in future

#!/usr/bin/env bash
# Filename: create_db_and_empty_tables.sh


POSTGRES_CONTAINER_USERNAME="postgres"
# meant to be executed in postgres container
psql -U $POSTGRES_CONTAINER_USERNAME
DROP DATABASE IF EXISTS auctiondb;
CREATE DATABASE auctiondb;
\q
