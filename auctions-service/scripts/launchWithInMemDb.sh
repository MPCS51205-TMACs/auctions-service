#!/usr/bin/env bash
# Filename: launchWithInMemDb.sh

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)
PROJECT_DIR="$(dirname "${SCRIPT_DIR}")"
MAIN_DIR=$PROJECT_DIR/main/

echo SCRIPT_DIR=$SCRIPT_DIR
echo PROJECT_DIR=$PROJECT_DIR
echo MAIN_DIR=$MAIN_DIR

# launch application with in memory repositories flag
go run $MAIN_DIR/ inmemory