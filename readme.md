# Lambda 関数のローカルでのステップスルーデバッグ

AWSの公式ドキュメントを参考にして、以下の手順でデバッグを行います：

## SAM CLIを使ったデバッグ

```bash
sam local invoke -d 5858 RssWriteFunction --template sam/template.yaml --event tests/cmd/rss/write/event.json
```

```bash
sudo ifconfig lo0 alias 172.16.123.1
```

## ローカルSNSの設定

1. Docker Composeを起動します：

   ```bash
   docker compose up
   ```

2. トピックを作成します：

   ```bash
   aws sns create-topic --name rss-write --endpoint-url http://localhost:4566
   ```

## ローカルDynamoDBの設定

以下のコマンドを使用して、対象のテーブルのStreamを有効にします。`table-name`は適切に設定してください。

### Streamを有効にするコマンド

```bash
# Streamを有効にする
aws dynamodb update-table --table-name User --stream-specification StreamEnabled=true,StreamViewType=NEW_IMAGE --endpoint-url http://localhost:8000 --region us-west-2

# 設定確認
aws dynamodb describe-table --table-name User --endpoint-url http://localhost:8000 --region us-west-2
```

### Streamを無効にするコマンド

StreamViewTypeを変更する場合は、一度無効にしてから再度有効にします。

```bash
# Streamを無効にする
aws dynamodb update-table --table-name User --stream-specification StreamEnabled=false --endpoint-url http://localhost:8000 --region us-west-2
```

### streamArnの取得

以下のコマンドを実行し、出力された`LatestStreamArn`フィールドの値をデバッグ対象の`main.go`の変数`streamArn`に設定します。

```bash
aws dynamodb describe-table --table-name User --endpoint-url http://localhost:8000 --region us-west-2
```

## UT

### 全てのUTを実行コマンド

```bash
go test -v ./...
```
