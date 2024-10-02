#!/bin/bash

# 外部から変数を受け取る
BUCKET=$1
REGION=$2
PROFILE=$3

# S3バケットが存在するか確認
if ! aws s3 ls "s3://${BUCKET}" --profile "${PROFILE}" ; then
  echo "バケットが存在しません。新しいバケットを作成します: ${BUCKET}"

  # バケットを作成
  aws s3 mb "s3://${BUCKET}" --region "${REGION}" --profile "${PROFILE}"

  echo "バケットが完全に動作するのを待っています..."
  sleep 20

  # リトライ設定
  max_retries=5
  count=0

  # バケットポリシーを適用する処理
  until aws s3api put-bucket-policy --bucket "${BUCKET}" --policy "$(cat <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudformation.amazonaws.com"
      },
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::${BUCKET}/*"
    }
  ]
}
EOF
)" --profile "${PROFILE}"
  do
    count=$((count+1))
    if [ "${count}" -eq "${max_retries}" ]; then
      echo "Failed to apply policy after ${max_retries} attempts."
      exit 1
    fi
    echo "Retrying to apply policy...attempt ${count}"
    sleep $((10 * count))
  done
fi