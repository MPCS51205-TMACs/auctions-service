#!/usr/bin/env bash
# Filename: build-all-images.sh

# stop script upon encountering the first error
# without closing terminal window
set -e


#####
# this script builds all images for the project from
# dockerfiles
#####

# RUN set-env-vars.sh first!

# need: location of docker file for the image you want to build
# need: a name for the image you are going to create

export POSTGRES_IMG_DOCKERFILE="$PROJECT_DIR_PATH/postgres/"
export POSTGRES_IMG_NAME="postgres-server"

export RABBITMQ_IMG_DOCKERFILE="$PROJECT_DIR_PATH/rabbitmq/"
export RABBITMQ_IMG_NAME="rabbitmq-server"

export AUCTIONS_IMG_DOCKERFILE="$PROJECT_DIR_PATH/auctions-service/"
export AUCTIONS_IMG_NAME="auctions-service"

echo
echo RUNNING SCRIPT TO BUILD PROJECT IMAGES...
echo SCRIPT_DIR=$SCRIPT_DIR
echo
echo POSTGRES_IMG_DOCKERFILE=$POSTGRES_IMG_DOCKERFILE
echo POSTGRES_IMG_NAME=$POSTGRES_IMG_NAME
echo
echo RABBITMQ_IMG_DOCKERFILE=$RABBITMQ_IMG_DOCKERFILE
echo RABBITMQ_IMG_NAME=$RABBITMQ_IMG_NAME
echo
echo AUCTIONS_IMG_DOCKERFILE=$AUCTIONS_IMG_DOCKERFILE
echo AUCTIONS_IMG_NAME=$AUCTIONS_IMG_NAME
echo

# we create each image from a docker file, using the docker command:
# 'docker build -t IMAGETAGNAME:VERSION PATHTODOCKERFILE'
# RABBIT MQ
echo "building RabbitMQ image ($RABBITMQ_IMG_NAME) from dockerfile..."
docker build -t "$RABBITMQ_IMG_NAME:latest" $RABBITMQ_IMG_DOCKERFILE

# Postgres SQL
echo "building Postgres image ($POSTGRES_IMG_NAME) from dockerfile..."
docker build -t "$POSTGRES_IMG_NAME:latest" $POSTGRES_IMG_DOCKERFILE

# RABBIT MQ
echo "building AuctionsService ($AUCTIONS_IMG_NAME) image from dockerfile..."
docker build -t "$AUCTIONS_IMG_NAME:latest" $AUCTIONS_IMG_DOCKERFILE

echo "done"