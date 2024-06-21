package app_service

import (
	"context"
	"encoding/json"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/YamazakiNorihito/workday/pkg/throttle"
)

type FeedProvider interface {
	GetFeedURLs() []string
}

func Execute(ctx context.Context, logger infrastructure.Logger, rssWritePublisher shared.Publisher, throttleConfig throttle.Config, feedProvider FeedProvider) error {
	feedURLs := feedProvider.GetFeedURLs()
	for i, feedURL := range feedURLs {
		message := message.Subscribe{FeedURL: feedURL}
		rssJson, _ := json.Marshal(message)
		err := rssWritePublisher.Publish(ctx, string(rssJson))
		if err != nil {
			return err
		}
		if (i+1)%throttleConfig.BatchSize == 0 {
			throttleConfig.Sleep()
		}
	}
	return nil
}
