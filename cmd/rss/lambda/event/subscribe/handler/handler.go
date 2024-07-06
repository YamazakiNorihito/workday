package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/subscribe/app_service"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/aws/aws-lambda-go/events"
)

type executer func(ctx context.Context, logger infrastructure.Logger, feedURL string, language string) error

func Handler(ctx context.Context, event events.SNSEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	snsClient := cfg.NewSnsClient()

	snsTopicClient := awsConfig.NewSnsTopicClient(snsClient, os.Getenv("OUTPUT_TOPIC_RSS_ARN"))
	publisher := publisher.NewWriterMessagePublisher(snsTopicClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("SNS Event", "event", shared.SnsEventToJson(event))

	httpClient := &http.Client{}
	executer := func(ctx context.Context, logger infrastructure.Logger, feedURL string, language string) error {
		repository := app_service.NewFeedRepository(httpClient, feedURL, language)
		return app_service.Execute(ctx, logger, &repository, *publisher)
	}

	for _, record := range event.Records {
		recordLogger := logger.With("messageID", record.SNS.MessageID)
		err := processRecord(ctx, recordLogger, record, executer)

		if err != nil {
			recordLogger.Error("Failed", "error", err)
		}
	}

	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, record events.SNSEventRecord, executer executer) error {
	receiveMessage, err := getMessage(record)
	if err != nil {
		return err
	}

	return executer(ctx, logger, receiveMessage.FeedURL, receiveMessage.Language)
}

func getMessage(record events.SNSEventRecord) (receiveMessage message.Subscribe, err error) {
	err = json.Unmarshal([]byte(record.SNS.Message), &receiveMessage)
	return receiveMessage, err
}
