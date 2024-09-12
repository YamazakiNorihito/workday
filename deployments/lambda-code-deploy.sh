#!/bin/bash

set -euo pipefail

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

aws s3 sync binaries "s3://${BUCKET}/binaries" --exclude "deploy*" --profile "${PROFILE}"

# Update Lambda functions
# shellcheck disable=SC2154
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