package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/write/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/aws/aws-lambda-go/events"
)

type executer func(ctx context.Context, logger infrastructure.Logger, rssEntry rss.Rss) error

func Handler(ctx context.Context, event events.SNSEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	dynamodbClient := cfg.NewDynamodbClient()
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("SNS Event", "event", shared.SnsEventToJson(event))

	executer := func(ctx context.Context, logger infrastructure.Logger, rssEntry rss.Rss) error {
		return app_service.Execute(ctx, logger, rssRepository, rssEntry)
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
