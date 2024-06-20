package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.SNSEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	dynamodbClient := cfg.NewDynamodbClient()
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("SNS Event", "event", shared.SnsEventToJson(event))

	for _, record := range event.Records {
		recordLogger := logger.With("messageID", record.SNS.MessageID)
		err := processRecord(ctx, recordLogger, rssRepository, record)

		if err != nil {
			recordLogger.Error("Failed", "error", err)
		}
		logger.Info("finish")
	}

	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, record events.SNSEventRecord) error {
	message, err := getMessage(record)
	if err != nil {
		return err
	}

	logger.Info("Processing command", message.RssFeed.Source)
	_, err = Core(ctx, logger, rssRepository, message.RssFeed)
	return err
}

func Core(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, rssEntry rss.Rss) (rss.Rss, error) {
	exists, existingRss := rss.Exists(ctx, rssRepository, rssEntry)
	logger.Info("save", "isExist", exists, "rss", rssEntry.Source)

	if exists {
		if existingRss.LastBuildDate.Equal(rssEntry.LastBuildDate) {
			logger.Info("No update needed", "rss", rssEntry.Source)
			return existingRss, nil
		}
	}

	savedRss, err := rssRepository.Save(ctx, rssEntry, metadata.UserMeta{ID: rssEntry.Source, Name: rssEntry.Source})
	return savedRss, err
}

func getMessage(record events.SNSEventRecord) (message message.Write, err error) {
	err = json.Unmarshal([]byte(record.SNS.Message), &message)
	if err != nil {
		return message, err
	}

	if message.Compressed {
		decompressedRssData, err := decompressAndDecodeData(message.Data)
		if err != nil {
			return message, err
		}

		err = json.Unmarshal(decompressedRssData, &message.RssFeed)
		if err != nil {
			return message, err
		}
		message.Data = nil
		message.Compressed = false
	}

	return message, nil
}

func decompressAndDecodeData(data []byte) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewReader(decodedData)
	gzipReader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	decompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}

func main() {
	lambda.Start(handler)
}
