#!/bin/bash

set -euo pipefail

# shellcheck disable=SC1091
source ./config.sh

NO_BUILD=false

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --no-build) NO_BUILD=true ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# -----------------------------
# Create execute file
# -----------------------------
if [ "$NO_BUILD" = false ]; then
    for dir in "${LAMBDA_DIRS[@]}"; do
        echo "Building Lambda function in $dir..."
        make -C "$SRC_DIR/$dir" build
        
        echo "Creating binary directory for $dir..."
        mkdir -p "$BIN_DIR/$dir"
        
        echo "Copying function.zip to $BIN_DIR/$dir..."
        cp "$SRC_DIR/$dir/function.zip" "$BIN_DIR/$dir/function.zip"
    done
else
    echo "Skipping build process."
fi

echo "Build process complete."

# -----------------------------
# Create S3 Bucket for Deployment
# -----------------------------
./create_s3_bucket_if_not_exists.sh "${BUCKET}" "${REGION}" "${PROFILE}"

# テンプレートをS3にアップロード
aws s3 sync . "s3://${BUCKET}/" --exclude "deploy*" --profile "${PROFILE}"

# -----------------------------
# Deploy CloudFormation
# -----------------------------
# AWS CloudFormationスタックを作成
aws cloudformation deploy \
  --stack-name "${STACK_NAME}" \
  --template-file "template.yaml" \
  --s3-bucket "${BUCKET}" \
  --capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND \
  --parameter-overrides TemplateBucket="${BUCKET}" TranslateApiUrl="${GAS_TRANALATE_API}" \
  --region "${REGION}" \
  --profile "${PROFILE}"

echo "CloudFormationスタック ${STACK_NAME} が正常にデプロイされました。"

# aws cloudformation package --template-file template.yaml --s3-bucket nybeyond-com-deploy --output-template-file out.yaml --profile workday --region us-east-1
# delete stack command
# aws cloudformation delete-stack --stack-name nybeyond-com-workday --profile workday --region us-east-1
# aws cloudformation describe-stack-events --stack-name nybeyond-com-workday --profile workday --region us-east-1
# aws lambda get-function-configuration --function-name test-go --profile medcom.ne.jp --region us-east-1
# aws lambda get-function --function-name test-go --profile workday --region us-east-1
# aws cloudformation deploy --stack-name nybeyond-com-workday --profile workday --region us-east-1--s3-bucket nybeyond-com-deploy --capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND --template-file lambda/eventSourceMapping/dynamoDB.yaml --parameter-overrides FunctionArn=arn:aws:lambda:us-east-1:155345814070:function:RssNotificationFunction DynamoDBStreamArn=arn:aws:dynamodb:us-east-1:155345814070:table/Rss
