package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.SNSEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	dynamodbClient := cfg.NewDynamodbClient()
	snsClient := cfg.NewSnsClient()

	snsTopicClient := awsConfig.NewSnsTopicClient(snsClient, os.Getenv("OUTPUT_TOPIC_RSS_ARN"))
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("SNS Event", "event", shared.SnsEventToJson(event))

	for _, record := range event.Records {
		recordLogger := logger.With("messageID", record.SNS.MessageID)
		err := processRecord(ctx, recordLogger, rssRepository, snsTopicClient, record)

		if err != nil {
			recordLogger.Error("Failed", "error", err)
		}
		logger.Info("finish")
	}

	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, rssWritePublisher shared.Publisher, record events.SNSEventRecord) error {
	receiveMessage, err := getMessage(record)
	if err != nil {
		return err
	}

	logger.Info("Processing command", receiveMessage.RssFeed.Source)
	entryRss, err := Core(ctx, logger, rssRepository, receiveMessage.RssFeed)
	if err != nil {
		return err
	}

	err = Publish(ctx, rssWritePublisher, entryRss)
	if err != nil {
		return err
	}
	logger.Info("Message published successfully", "feedURL", receiveMessage.RssFeed.Source)
	return nil
}

func Core(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, rssEntry rss.Rss) (rss.Rss, error) {
	exists, existingRss := rss.Exists(ctx, rssRepository, rssEntry)
	logger.Info("Checking existence of RSS entry", "exists", exists, "source", rssEntry.Source)

	if exists == false {
		return rssEntry, nil
	}

	existingRss.SetLastBuildDate(rssEntry.LastBuildDate)
	for _, item := range rssEntry.Items {
		existingRss.AddOrUpdateItem(item)
	}

	cleansingRss, err := cleansing(ctx, logger, rssRepository, existingRss)
	if err != nil {
		return rss.Rss{}, err
	}

	return cleansingRss, nil
}

func Publish(ctx context.Context, rssWritePublisher shared.Publisher, rssEntry rss.Rss) error {
	writeMessage, err := message.NewWriteMessage(rssEntry)
	if err != nil {
		return err
	}

	rssJson, _ := json.Marshal(writeMessage)
	err = rssWritePublisher.Publish(ctx, string(rssJson))
	if err != nil {
		return err
	}
	return nil
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
		return message.Write{}, err
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
