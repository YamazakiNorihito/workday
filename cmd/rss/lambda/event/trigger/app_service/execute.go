package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/YamazakiNorihito/workday/pkg/throttle"
)

type FeedProvider interface {
	GetFeedURLs() []string
}

func Execute(ctx context.Context, logger infrastructure.Logger, publisher publisher.SubscribeMessagePublisher, throttleConfig throttle.Config, feedProvider FeedProvider) error {
	err := Trigger(ctx, logger, publisher, throttleConfig, feedProvider)
	if err != nil {
		return err
	}

	logger.Info("Message Trigger successfully")
	return nil
}

func Trigger(ctx context.Context, logger infrastructure.Logger, publisher publisher.SubscribeMessagePublisher, throttleConfig throttle.Config, feedProvider FeedProvider) error {
	feedURLs := feedProvider.GetFeedURLs()
	for i, feedURL := range feedURLs {
		err := publisher.Publish(ctx, feedURL)
		if err != nil {
			return err
		}
		if (i+1)%throttleConfig.BatchSize == 0 {
			throttleConfig.Sleep()
		}
	}
	return nil
}
