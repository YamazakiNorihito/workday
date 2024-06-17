#!/bin/bash

set -euo pipefail

# Directories
SRC_DIR="./../cmd/rss/lambda/event"
BIN_DIR="./binaries/rss/lambda/event"
LAMBDA_DIRS=("notification" "subscribe" "trigger" "write")

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
FUNCTION_NAMES=("RssNotificationFunction" "RssSubscribeFunction" "RssTriggerFunction" "RssWriteFunction")
S3_KEYS=("binaries/rss/lambda/event/notification/function.zip" "binaries/rss/lambda/event/subscribe/function.zip" "binaries/rss/lambda/event/trigger/function.zip" "binaries/rss/lambda/event/write/function.zip")

# Update Lambda functions
for i in "${!FUNCTION_NAMES[@]}"; do
    function_name="${FUNCTION_NAMES[$i]}"
    s3_key="${S3_KEYS[$i]}"
    echo "Updating Lambda function $function_name with S3 key $s3_key..."
    aws lambda update-function-code \
        --function-name "$function_name" \
        --s3-bucket "${BUCKET}" \
        --s3-key "$s3_key" \
        --profile "${PROFILE}" \
        --region "${REGION}"
done

echo "Lambda functions update complete."