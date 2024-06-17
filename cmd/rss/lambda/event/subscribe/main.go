package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcdole/gofeed"
)

func handler(ctx context.Context, event events.SNSEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	snsClient := cfg.NewSnsClient()
	dynamodbClient := cfg.NewDynamodbClient()

	snsTopicClient := awsConfig.NewSnsTopicClient(snsClient, os.Getenv("RSS_WRITE_ARN"))
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("SNS Event", "event", shared.SnsEventToJson(event))

	for _, record := range event.Records {
		recordLogger := logger.With("messageID", record.SNS.MessageID)
		err := processRecord(ctx, recordLogger, rssRepository, snsTopicClient, record)

		if err != nil {
			recordLogger.Error("Failed", "error", err)
		}
	}

	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, rssWritePublisher shared.Publisher, record events.SNSEventRecord) error {
	subscribeMessage, err := getMessage(record)
	if err != nil {
		return err
	}

	entryRss, err := Core(ctx, logger, rssRepository, subscribeMessage.FeedURL)
	if err != nil {
		return err
	}

	writeMessage := message.Write{RssEntry: entryRss}
	rssJson, _ := json.Marshal(writeMessage)
	err = rssWritePublisher.Publish(ctx, string(rssJson))
	if err != nil {
		return err
	}
	logger.Info("Message published successfully", "feedURL", subscribeMessage.FeedURL)
	return nil
}

func Core(ctx context.Context, logger infrastructure.Logger, repository rss.IRssRepository, feedURL string) (rssEntry rss.Rss, err error) {
	source := getFQDN(feedURL)
	if source == "" {
		return rss.Rss{}, fmt.Errorf("invalid Feed URL: %s", feedURL)
	}

	feed, err := getFeed(ctx, feedURL)
	if err != nil {
		return rss.Rss{}, err
	}

	feedLink := feed.Link
	if feedLink == "" {
		feedLink = feed.FeedLink
	}

	rssEntry, err = rss.New(feed.Title, source, feedLink, feed.Description, feed.Language, feed.UpdatedParsed.UTC())
	if err != nil {
		return rss.Rss{}, err
	}

	exists, existingRss := rss.Exists(ctx, repository, rssEntry)
	if exists {
		existingRss.SetLastBuildDate(rssEntry.LastBuildDate)
		rssEntry = existingRss
	}

	for _, item := range feed.Items {
		guid := rss.Guid{Value: item.GUID}

		authorEmail := ""
		if item.Author != nil {
			authorEmail = item.Author.Email
		}
		entryItem, err := rss.NewItem(guid, item.Title, item.Link, item.Description, authorEmail, *item.PublishedParsed)
		if err != nil {
			logger.Error("Validation error when creating RSS item", "error", err, "item", item.Title)
			continue
		}
		rssEntry.AddOrUpdateItem(entryItem)
	}

	return rssEntry, nil
}

func getFQDN(uri string) string {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}

func getFeed(ctx context.Context, feedURL string) (feed *gofeed.Feed, err error) {
	fp := gofeed.NewParser()
	return fp.ParseURLWithContext(feedURL, ctx)
}

func getMessage(record events.SNSEventRecord) (message message.Subscribe, err error) {
	err = json.Unmarshal([]byte(record.SNS.Message), &message)
	return message, err
}

func main() {
	if os.Getenv("ENV") == "myhost" {
		event := events.SNSEvent{
			Records: []events.SNSEventRecord{
				{
					SNS: events.SNSEntity{
						MessageID: "12345",
						Message:   `{"feed_url": "https://buildersbox.corp-sansan.com/rss"}`,
					},
				},
			},
		}
		handler(context.Background(), event)
	} else {
		lambda.Start(handler)
	}
}
