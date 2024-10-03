#!/bin/bash

# shellcheck disable=SC2034
PROFILE="workday"
BUCKET="nybeyond-com-deploy"
REGION="us-east-1"
STACK_NAME="nybeyond-com-workday"
GAS_TRANALATE_API="https://script.google.com/macros/s/AKfycbwrnNBNPJ94-HGK-Ske-aIjfI_bGuRQ37tg3MsI6Fqsb3n9psq_Z02znIwUjpMaLRudow/exec"


SRC_DIR="./../cmd/rss/lambda/"
BIN_DIR="./binaries/rss/lambda/"

# Function names and corresponding S3 keys
FUNCTIONs=("RssNotificationFunction:event/notification"
        "RssSubscribeFunction:event/subscribe"
        "RssTriggerFunction:event/trigger"
        "RssWriteFunction:event/write"
        "RssTranslateFunction:event/translate"
        "RssCleanFunction:event/clean"
        "RssCreateFunction:api/create"
        "RssFeedsFunction:api/feeds"
        "RssFeedIdFunction:api/feed_id"
        "RssPatchFunction:api/patch")
