package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.EventBridgeEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	snsClient := cfg.NewSnsClient()
	snsTopicClient := awsConfig.NewSnsTopicClient(snsClient, os.Getenv("RSS_SUBSCRIBE_ARN"))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("EventBridgeID", event.ID)
	logger.Info("EventBridgeEvent Event", "event", shared.EventBridgeEventToJson(event))

	err := Core(ctx, logger, snsTopicClient)
	if err != nil {
		logger.Error("Core function execution failed", "error", err)
		return err
	}

	logger.Info("Message published successfully")
	return nil
}

func Core(ctx context.Context, logger infrastructure.Logger, rssWritePublisher shared.Publisher) error {
	feedURLs := [12]string{
		"https://azure.microsoft.com/ja-jp/blog/feed/",
		"https://aws.amazon.com/jp/blogs/news/feed/",
		"https://developers-jp.googleblog.com/atom.xml",
		"https://techblog.nhn-techorus.com/feed",
		"https://buildersbox.corp-sansan.com/rss",
		"https://knowledge.sakura.ad.jp/rss/",
		"https://www.oreilly.co.jp/catalog/soon.xml",
		"https://go.dev/blog/feed.atom",
		"https://connpass.com/explore/ja.atom",
		"https://www.ipa.go.jp/security/alert-rss.rdf",
		"https://feed.infoq.com",
		"https://techcrunch.com/feed",
	}

	for _, feedURL := range feedURLs {
		message := message.Subscribe{FeedURL: feedURL}
		rssJson, _ := json.Marshal(message)
		err := rssWritePublisher.Publish(ctx, string(rssJson))
		if err != nil {
			logger.Error("Failed to publish RSS entry", "error", err)
			return err
		}
	}
	return nil
}

func main() {
	if os.Getenv("ENV") == "myhost" {
		event := events.EventBridgeEvent{
			ID: "test-event-id",
		}
		handler(context.Background(), event)
	} else {
		lambda.Start(handler)
	}
}
