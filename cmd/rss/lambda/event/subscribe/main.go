package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

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

	err = Publish(ctx, rssWritePublisher, entryRss)
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
		logger.Error("Failed to retrieve RSS feed", "URL", feedURL, "error", err)
		return rss.Rss{}, err
	}

	lastBuildDate := getLastBuildDate(*feed)
	rssEntry, err = rss.New(feed.Title, source, feedURL, feed.Description, feed.Language, lastBuildDate.UTC())
	if err != nil {
		return rss.Rss{}, err
	}

	for _, item := range feed.Items {
		guid, err := getGuid(*item)
		if err != nil {
			logger.Error("Failed to create GUID from RSS item link", "error", err, "item", item.Title, "link", item.Link)
			continue
		}

		author := getAuthor(*item)
		textDescription, err := getDescription(*item)
		if err != nil {
			logger.Error("Failed to extract text from HTML content in RSS item description", "error", err, "itemTitle", item.Title)
			continue
		}

		entryItem, err := rss.NewItem(guid, item.Title, item.Link, textDescription, author, *item.PublishedParsed)
		if err != nil {
			logger.Error("Validation error when creating RSS item", "error", err, "item", item.Title)
			continue
		}
		rssEntry.AddOrUpdateItem(entryItem)
	}

	return rssEntry, nil
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

func getLastBuildDate(feed gofeed.Feed) (lastBuildDate time.Time) {
	for _, item := range feed.Items {
		if item.PublishedParsed != nil && lastBuildDate.Before(*item.PublishedParsed) {
			lastBuildDate = *item.PublishedParsed
		}
	}

	if lastBuildDate.IsZero() && feed.UpdatedParsed != nil {
		lastBuildDate = *feed.UpdatedParsed
	}
	return lastBuildDate
}

func getGuid(item gofeed.Item) (rss.Guid, error) {
	guid := rss.Guid{Value: item.GUID}

	if guid.Value == "" {
		link, err := url.Parse(item.Link)
		if err != nil {
			return rss.Guid{}, err
		}
		link.RawQuery = ""
		guid = rss.Guid{Value: link.String()}
	}

	return guid, nil
}

func getAuthor(item gofeed.Item) string {
	if item.Author == nil {
		return ""
	}

	if email := item.Author.Email; email != "" {
		return email
	}
	return item.Author.Name
}

func getDescription(item gofeed.Item) (string, error) {
	description := item.Description
	if description == "" {
		description = item.Content
	}
	return ExtractTextFromHTML(description)
}

func main() {
	if os.Getenv("ENV") == "myhost" {
		event := events.SNSEvent{
			Records: []events.SNSEventRecord{
				{
					SNS: events.SNSEntity{
						MessageID: "12345",
						Message:   `{"feed_url": "https://www.ipa.go.jp/security/alert-rss.rdf"}`,
					},
				},
			},
		}
		handler(context.Background(), event)
	} else {
		lambda.Start(handler)
	}
}
