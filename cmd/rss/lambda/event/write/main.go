package main

import (
	"context"
	"encoding/json"
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
	logger.Info("Checking existence of RSS entry", "exists", exists, "source", rssEntry.Source)

	if exists {
		if existingRss.LastBuildDate.Equal(rssEntry.LastBuildDate) {
			logger.Info("RSS entry is up-to-date, no update needed", "source", rssEntry.Source)
			return existingRss, nil
		}
	}
	cleansingRss, err := cleansing(ctx, logger, rssRepository, rssEntry)
	if err != nil {
		return rss.Rss{}, err
	}

	if len(cleansingRss.Items) == 0 {
		logger.Info("No records to update, skipping update", "source", rssEntry.Source)
		return rssEntry, nil
	}

	savedRss, err := rssRepository.Save(ctx, rssEntry, metadata.UserMeta{ID: rssEntry.Source, Name: rssEntry.Source})
	return savedRss, err
}

func cleansing(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, rssEntry rss.Rss) (cleansingRss rss.Rss, err error) {
	cleansingRss = rssEntry
	cleansingRss.Items = map[rss.Guid]rss.Item{}

	for key, item := range rssEntry.Items {
		findItem, err := rss.GetItem(ctx, rssRepository, rssEntry, key)
		if err != nil {
			logger.Error("Error retrieving item", "error", err, "source", rssEntry.Source, "guid", key)
			continue
		}

		if len(findItem.Items) == 0 {
			cleansingRss.Items[key] = item
		} else {
			logger.Info("Item already exists and will not be added", "source", rssEntry.Source, "guid", key)
		}
	}

	return cleansingRss, nil
}

func getMessage(record events.SNSEventRecord) (receiveMessage message.Write, err error) {
	err = json.Unmarshal([]byte(record.SNS.Message), &receiveMessage)
	if err != nil {
		return receiveMessage, err
	}

	if receiveMessage.Compressed {
		decompressedRssData, err := message.DecodeAndDecompressData(receiveMessage.Data)
		if err != nil {
			return message.Write{}, err
		}
		receiveMessage.RssFeed = decompressedRssData
		receiveMessage.Data = nil
		receiveMessage.Compressed = false
	}

	return receiveMessage, nil
}

func main() {
	lambda.Start(handler)
}
