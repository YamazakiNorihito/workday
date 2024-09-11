#!/bin/bash

set -euo pipefail

# Directories
SRC_DIR="./../cmd/rss/lambda/"
BIN_DIR="./binaries/rss/lambda/"
LAMBDA_DIRS=("event/notification" "event/subscribe" "event/trigger" "event/write" "event/translate" "event/clean" "api/create" "api/feeds")

# Build and package functions
for dir in "${LAMBDA_DIRS[@]}"; do
    echo "Building Lambda function in $dir..."
    make -C "$SRC_DIR/$dir" build
    
    echo "Creating binary directory for $dir..."
    mkdir -p "$BIN_DIR/$dir"
    
    echo "Copying function.zip to $BIN_DIR/$dir..."
    cp "$SRC_DIR/$dir/function.zip" "$BIN_DIR/$dir/function.zip"
done

echo "Build process complete."

# AWS Settings
PROFILE="workday"
BUCKET="nybeyond-com-deploy"
REGION="us-east-1"
aws s3 sync binaries "s3://${BUCKET}/binaries" --exclude "deploy*" --profile "${PROFILE}"

# Function names and corresponding S3 keys
FUNCTIONs=("RssNotificationFunction:binaries/rss/lambda/event/notification/function.zip"
        "RssSubscribeFunction:binaries/rss/lambda/event/subscribe/function.zip"
        "RssTriggerFunction:binaries/rss/lambda/event/trigger/function.zip"
        "RssWriteFunction:binaries/rss/lambda/event/write/function.zip"
        "RssTranslateFunction:binaries/rss/lambda/event/translate/function.zip"
        "RssCleanFunction:binaries/rss/lambda/event/clean/function.zip"
        "RssCreateFunction:binaries/rss/lambda/api/create/function.zip"
        "RssFeedsFunction:binaries/rss/lambda/api/feeds/function.zip")

# Update Lambda functions
for FUNCTION in "${FUNCTIONs[@]}"; do
    function_name="${FUNCTION%%:*}"
    s3_key="${FUNCTION#*:}"

    echo "Updating Lambda function $function_name with S3 key $s3_key..."
    aws lambda update-function-code \
        --function-name "$function_name" \
        --s3-bucket "${BUCKET}" \
        --s3-key "$s3_key" \
        --profile "${PROFILE}" \
        --region "${REGION}"
done

echo "Lambda functions update complete."