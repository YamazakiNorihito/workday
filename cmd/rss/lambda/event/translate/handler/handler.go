package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/translate/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/aws/aws-lambda-go/events"
)

type feedProvider struct {
	repo *rss.DynamoDBRssRepository
}

func (r *feedProvider) GetSourceLanguage(ctx context.Context, logger infrastructure.Logger, source string) (sourceLanguageCode string, ok bool) {
	feeds, err := r.repo.FindAll(ctx)
	if err != nil {
		logger.Warn("Failed to retrieve feeds from repository", "error", err)
		return "", false
	}
	languageFeedMap := make(map[string]string)
	for _, feed := range feeds {
		languageFeedMap[feed.Source] = feed.Language
	}

	sourceLanguageCode, ok = languageFeedMap[source]
	return
}

type executer func(ctx context.Context, logger infrastructure.Logger, rssEntry rss.Rss) error

func Handler(ctx context.Context, event events.SNSEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	snsClient := cfg.NewSnsClient()
	translateClient := cfg.NewTranslateClient()

	snsTopicClient := awsConfig.NewSnsTopicClient(snsClient, os.Getenv("OUTPUT_TOPIC_RSS_ARN"))
	publisher := publisher.NewWriterMessagePublisher(snsTopicClient)
	easyTranslateClient := awsConfig.NewTranslateClient(translateClient)

	dynamodbClient := cfg.NewDynamodbClient()
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)
	feedProvider := feedProvider{repo: rssRepository}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("SNS Event", "event", shared.SnsEventToJson(event))

	executer := func(ctx context.Context, logger infrastructure.Logger, rssEntry rss.Rss) error {
		return app_service.Execute(ctx, logger, easyTranslateClient, &feedProvider, *publisher, rssEntry)
	}

	for _, record := range event.Records {
		recordLogger := logger.With("messageID", record.SNS.MessageID)
		err := processRecord(ctx, recordLogger, executer, record)

		if err != nil {
			recordLogger.Error("Failed", "error", err)
		}
		logger.Info("finish")
	}

	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, executer executer, record events.SNSEventRecord) error {
	receiveMessage, err := getMessage(record)
	if err != nil {
		return err
	}
	return executer(ctx, logger, receiveMessage.RssFeed)
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
