package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/YamazakiNorihito/workday/pkg/throttle"
)

type FeedProvider interface {
	GetFeedURLAndLanguage(ctx context.Context) (map[string]string, error)
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
	feedURLAndLanguageMap, err := feedProvider.GetFeedURLAndLanguage(ctx)
	if err != nil {
		return err
	}
	i := 0
	for feedURL, language := range feedURLAndLanguageMap {
		err := publisher.Publish(ctx, feedURL, language)
		if err != nil {
			return err
		}
		if (i+1)%throttleConfig.BatchSize == 0 {
			throttleConfig.Sleep()
		}
		i++
	}
	return nil
}
