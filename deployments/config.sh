#!/bin/bash

# shellcheck disable=SC2034
PROFILE="workday"
BUCKET="nybeyond-com-deploy"
REGION="us-east-1"
STACK_NAME="nybeyond-com-workday"
GAS_TRANALATE_API="https://script.google.com/macros/s/AKfycbwrnNBNPJ94-HGK-Ske-aIjfI_bGuRQ37tg3MsI6Fqsb3n9psq_Z02znIwUjpMaLRudow/exec"


SRC_DIR="./../cmd/rss/lambda/"
BIN_DIR="./binaries/rss/lambda/"
LAMBDA_DIRS=("event/notification" "event/subscribe" "event/trigger" "event/write" "event/translate" "event/clean" "api/create" "api/feeds")

# Function names and corresponding S3 keys
FUNCTIONs=("RssNotificationFunction:binaries/rss/lambda/event/notification/function.zip"
        "RssSubscribeFunction:binaries/rss/lambda/event/subscribe/function.zip"
        "RssTriggerFunction:binaries/rss/lambda/event/trigger/function.zip"
        "RssWriteFunction:binaries/rss/lambda/event/write/function.zip"
        "RssTranslateFunction:binaries/rss/lambda/event/translate/function.zip"
        "RssCleanFunction:binaries/rss/lambda/event/clean/function.zip"
        "RssCreateFunction:binaries/rss/lambda/api/create/function.zip"
        "RssFeedsFunction:binaries/rss/lambda/api/feeds/function.zip")
