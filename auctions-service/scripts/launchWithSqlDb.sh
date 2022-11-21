#!/usr/bin/env bash
# Filename: launchWithSqlDb.sh

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)
PROJECT_DIR="$(dirname "${SCRIPT_DIR}")"
MAIN_DIR=$PROJECT_DIR/main/
DB_DIR=$PROJECT_DIR/db/

echo SCRIPT_DIR=$SCRIPT_DIR
echo PROJECT_DIR=$PROJECT_DIR
echo MAIN_DIR=$MAIN_DIR
echo DB_DIR=$MAIN_DIR

POSTGRES_CONTAINER_HOSTNAME="postgres-server"
POSTGRES_CONTAINER_USERNAME="postgres"
POSTGRES_CONTAINER_PASSWORD="mysecret"
AUCTION_SERVICE_DB_NAME="auctionsservicedb"

psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$AUCTION_SERVICE_DB_NAME'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE $AUCTION_SERVICE_DB_NAME"
psql -h $POSTGRES_CONTAINER_HOSTNAME -U $POSTGRES_CONTAINER_USERNAME -d myDataBase -a -f myInsertFile

# launch application with in memory repositories flag
# go run $MAIN_DIR/ sql