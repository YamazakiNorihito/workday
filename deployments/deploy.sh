#!/bin/bash

set -euo pipefail

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
SRC_DIR="./../cmd/rss/lambda/event"
BIN_DIR="./binaries/rss/lambda/event"

LAMBDA_DIRS=("notification" "subscribe" "trigger" "write" "translate")

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


# プロファイル名の設定
PROFILE="workday"
BUCKET="nybeyond-com-deploy"
REGION="us-east-1"
# -----------------------------
# Create Deploy S3
# -----------------------------
if ! aws s3 ls "s3://${BUCKET}" --profile "${PROFILE}" ; then
  echo "バケットが存在しません。新しいバケットを作成します: ${BUCKET}"
  aws s3 mb "s3://${BUCKET}" --region "${REGION}" --profile "${PROFILE}"
  echo "バケットが完全に動作するのを待っています..."
  sleep 10 

  max_retries=5
  count=0
  until aws s3api put-bucket-policy --bucket "${BUCKET}" --policy file://deploy-bucket-policy.json --profile "${PROFILE}"
  do
    count=$((count+1))
    if [ "${count}" -eq "${max_retries}" ]; then
      echo "Failed to apply policy after ${max_retries} attempts."
      exit 1
    fi
    echo "Retrying to apply policy...attempt ${count}"
    sleep 10
  done
fi
# テンプレートをS3にアップロード
aws s3 sync . "s3://${BUCKET}/" --exclude "deploy*" --profile "${PROFILE}"

# -----------------------------
# Deploy CloudFormation
# -----------------------------
stack_name="nybeyond-com-workday"

# AWS CloudFormationスタックを作成
aws cloudformation deploy \
  --stack-name "${stack_name}" \
  --template-file "template.yaml" \
  --s3-bucket "${BUCKET}" \
  --capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND \
  --parameter-overrides TemplateBucket=${BUCKET} \
  --region "${REGION}" \
  --profile "${PROFILE}"

echo "CloudFormationスタック ${stack_name} が正常にデプロイされました。"

# aws cloudformation package --template-file template.yaml --s3-bucket nybeyond-com-deploy --output-template-file out.yaml --profile workday --region us-east-1
# delete stack command
# aws cloudformation delete-stack --stack-name nybeyond-com-workday --profile workday --region us-east-1
# aws cloudformation describe-stack-events --stack-name nybeyond-com-workday --profile workday --region us-east-1
# aws lambda get-function-configuration --function-name test-go --profile medcom.ne.jp --region us-east-1
# aws lambda get-function --function-name test-go --profile workday --region us-east-1
# aws cloudformation deploy --stack-name nybeyond-com-workday --profile workday --region us-east-1--s3-bucket nybeyond-com-deploy --capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND --template-file lambda/eventSourceMapping/dynamoDB.yaml --parameter-overrides FunctionArn=arn:aws:lambda:us-east-1:155345814070:function:RssNotificationFunction DynamoDBStreamArn=arn:aws:dynamodb:us-east-1:155345814070:table/Rss
