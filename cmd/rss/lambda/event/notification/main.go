package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/localdebug"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

const updateTimeThreshold = 30 * time.Minute

type SlackSender interface {
	PostMessageContext(ctx context.Context, options ...slack.MsgOption) (respChannel string, respTimestamp string, err error)
}

type slackChannelClient struct {
	client    *slack.Client
	channelId string
}

func (s *slackChannelClient) PostMessageContext(ctx context.Context, options ...slack.MsgOption) (respChannel string, respTimestamp string, err error) {
	return s.client.PostMessageContext(ctx, s.channelId, options...)
}

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	dynamodbClient := cfg.NewDynamodbClient()
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("DynamoDBEvent Event", "event", shared.DynamoDBEventToJson(event))

	slackClient := slack.New(os.Getenv("SLACK_TOKEN"))
	slackChannelClient := &slackChannelClient{
		client:    slackClient,
		channelId: os.Getenv("SLACK_CHANNEL_ID"),
	}

	for _, record := range event.Records {
		recordLogger := logger.With("eventID", record.EventID)
		processRecord(ctx, recordLogger, rssRepository, slackChannelClient, record)
	}
	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, slackSender SlackSender, record events.DynamoDBEventRecord) {
	logger.Info("Processing DynamoDB", "record", record)

	if record.EventName == "REMOVE" {
		logger.Info("REMOVE event detected, skipping processing")
		return
	}

	if record.Change.NewImage["sortKey"].String() != "rss" {
		logger.Info("対象外のレコードです")
		return
	}

	source := record.Change.NewImage["source"].String()
	err := Core(ctx, logger, rssRepository, slackSender, source)
	if err != nil {
		logger.Error("Failed", "error", err)
		return
	}
	logger.Info("Success")
	return
}

func Core(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, slackSender SlackSender, source string) error {

	newImageRss, err := rssRepository.FindBySource(ctx, source)
	if err != nil {
		return err
	}

	now := time.Now()
	if now.Sub(newImageRss.LastBuildDate) > updateTimeThreshold {
		logger.Info(fmt.Sprintf("The LastBuildDate is not within the last %v. Skipping processing.", updateTimeThreshold))
		return nil
	}

	modifyRss, err := rss.GetItems(ctx, rssRepository, newImageRss)
	if err != nil {
		return err
	}

	postMessage := makeMessage(modifyRss, func(item rss.Item) bool {
		result := now.Sub(item.PubDate) <= updateTimeThreshold
		logger.Info(fmt.Sprintf("Checking item with PubDate: %s, Current time: %s, Update time threshold: %v, Result: %t",
			item.PubDate.Format(time.RFC3339), now.Format(time.RFC3339), updateTimeThreshold, result))
		return result
	})
	respChannel, respTimestamp, err := slackSender.PostMessageContext(
		ctx,
		slack.MsgOptionText(postMessage, false),
		slack.MsgOptionUsername(source),
	)

	if err != nil {
		return err
	}

	logger.Info("Successfully sent message to Slack", "source", source, "response channel", respChannel, "timestamp", respTimestamp)
	return nil
}

func makeMessage(r rss.Rss, itemFilter func(item rss.Item) bool) string {
	if len(r.Items) == 0 {
		return ""
	}

	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("*フィードタイトル:* <%s|%s>\n*フィード詳細:* %s", r.Link, r.Title, r.Description))
	if !r.LastBuildDate.IsZero() {
		messageBuilder.WriteString(fmt.Sprintf("\n*最終更新日:* %s", r.LastBuildDate.Format(time.RFC3339)))
	}
	messageBuilder.WriteString("\n\n*最新の記事:*\n")

	const viewLength = 74 * 2
	i := 1
	for _, item := range r.Items {
		if itemFilter(item) == false {
			continue
		}

		truncatedDescription := truncate(item.Description, viewLength)
		messageBuilder.WriteString(fmt.Sprintf("%d. *記事タイトル:* <%s|%s>\n    *公開日:* %s\n    *概要:* %s\n    *カテゴリ:* %s\n\n",
			i, item.Link, item.Title, item.PubDate.Format(time.RFC3339), truncatedDescription, strings.Join(item.Tags, ", ")))
		i++
	}

	return messageBuilder.String()
}

func truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}

func main() {
	if os.Getenv("ENV") == "myhost" {
		streamARN := "arn:aws:dynamodb:ddblocal:000000000000:table/Rss/stream/2024-06-04T23:55:01.203"
		ctx := context.Background()
		localdebug.PollStreamAndInvokeHandler(ctx, streamARN, handler)
	} else {
		lambda.Start(handler)
	}
}
