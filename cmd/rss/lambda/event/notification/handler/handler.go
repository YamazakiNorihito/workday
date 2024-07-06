package handler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/notification/app_service"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

const updateTimeThreshold = 30 * time.Minute

type executer func(ctx context.Context, logger infrastructure.Logger, isNew bool, source string) error

type slackChannelClient struct {
	client    *slack.Client
	channelId string
}

func (s *slackChannelClient) PostMessageContext(ctx context.Context, text string, username string) (respChannel string, respTimestamp string, err error) {
	return s.client.PostMessageContext(ctx, s.channelId, slack.MsgOptionText(text, false), slack.MsgOptionUsername(username))
}

func Handler(ctx context.Context, event events.DynamoDBEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	dynamodbClient := cfg.NewDynamodbClient()
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)

	slackClient := slack.New(os.Getenv("SLACK_TOKEN"))
	slackChannelClient := &slackChannelClient{
		client:    slackClient,
		channelId: os.Getenv("SLACK_CHANNEL_ID"),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("DynamoDBEvent Event", "event", shared.DynamoDBEventToJson(event))

	executer := func(ctx context.Context, logger infrastructure.Logger, isNew bool, source string) error {
		now := time.Now()
		conditions := app_service.RssConditions{
			Target: func(r rss.Rss) bool {
				if isNew {
					return true
				}
				shouldProcess := now.Sub(r.LastBuildDate) <= updateTimeThreshold
				logger.Info("The LastBuildDate is not within the last update time threshold. Skipping processing.", "ID", r.ID, "UpdateTimeThreshold", updateTimeThreshold, "isOutdated", shouldProcess)
				return shouldProcess
			},
			ItemFilter: func(item rss.Item) bool {
				if isNew {
					return true
				}
				result := now.Sub(item.PubDate) <= updateTimeThreshold
				logger.Info(fmt.Sprintf("Checking item with GUID: %s, PubDate: %s, Current time: %s, Update time threshold: %v, Result: %t",
					item.Guid.Value, item.PubDate.Format(time.RFC3339), now.Format(time.RFC3339), updateTimeThreshold, result))
				return result
			},
		}

		return app_service.Execute(ctx, logger, rssRepository, slackChannelClient, conditions, source)
	}

	for _, record := range event.Records {
		recordLogger := logger.With("eventID", record.EventID)
		err := processRecord(ctx, recordLogger, executer, record)

		if err != nil {
			recordLogger.Error("Failed", "error", err)
		}
		logger.Info("finish")
	}

	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, executer executer, record events.DynamoDBEventRecord) error {
	logger.Info("Processing DynamoDB", "record", record)

	if record.EventName == "REMOVE" {
		logger.Info("REMOVE event detected, skipping processing")
		return nil
	}

	if record.Change.NewImage["sortKey"].String() != "rss" {
		logger.Info("対象外のレコードです")
		return nil
	}

	source := record.Change.NewImage["source"].String()
	return executer(ctx, logger, record.EventName == "INSERT", source)
}
