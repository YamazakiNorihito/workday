package handler

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/trigger/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/YamazakiNorihito/workday/pkg/throttle"
	"github.com/aws/aws-lambda-go/events"
)

type feedProvider struct {
	repo *rss.DynamoDBRssRepository
}

func (r *feedProvider) GetFeedURLAndLanguage(ctx context.Context) (map[string]string, error) {
	feeds, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	feedMap := make(map[string]string)
	for _, feed := range feeds {
		feedMap[feed.Link] = feed.Language
	}
	return feedMap, nil
}

type executer func(ctx context.Context, logger infrastructure.Logger) error

func Handler(ctx context.Context, event events.EventBridgeEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	snsClient := cfg.NewSnsClient()
	snsTopicClient := awsConfig.NewSnsTopicClient(snsClient, os.Getenv("OUTPUT_TOPIC_RSS_ARN"))
	publisher := publisher.NewSubscribeMessagePublisher(snsTopicClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("EventBridgeID", event.ID)
	logger.Info("EventBridgeEvent Event", "event", shared.EventBridgeEventToJson(event))

	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		panic("Invalid BATCH_SIZE value: must be an integer")
	}

	throttleConfig := throttle.Config{
		BatchSize: batchSize,
		Sleep:     func() { time.Sleep(2 * time.Second) },
	}

	dynamodbClient := cfg.NewDynamodbClient()
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)
	feedProvider := feedProvider{repo: rssRepository}

	executer := func(ctx context.Context, logger infrastructure.Logger) error {
		return app_service.Execute(ctx, logger, *publisher, throttleConfig, &feedProvider)
	}

	err = processRecord(ctx, logger, event, executer)
	if err != nil {
		logger.Error("ProcessRecord function execution failed", "error", err)
		return err
	}

	logger.Info("Message published successfully")
	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, _ events.EventBridgeEvent, executer executer) error {
	return executer(ctx, logger)
}
