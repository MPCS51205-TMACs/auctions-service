#!/usr/bin/env bash
# Filename: shutdown-cleardb-rebuild-deploy.sh

# get env vars

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)
. $SCRIPT_DIR/set-env-vars.sh # gives us some env vars

# shutdown
docker-compose -f $PROJECT_DIR_PATH/docker-compose-sql.yml down

# clear sql database
. $SCRIPT_DIR/clear-sql-db.sh

# rebuild images
. $SCRIPT_DIR/build-all-images.sh

# startup
docker-compose -f $PROJECT_DIR_PATH/docker-compose-sql.yml up -d

# # optional: go into auctions-service docker container. though the app should already be running so this not needed
# docker exec -it auctions-service /bin/bash
# ./main