package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/YamazakiNorihito/workday/pkg/throttle"
)

func Execute(ctx context.Context, logger infrastructure.Logger, publisher publisher.SubscribeMessagePublisher, throttleConfig throttle.Config, rssRepository rss.IRssRepository) error {
	err := Trigger(ctx, logger, publisher, throttleConfig, rssRepository)
	if err != nil {
		return err
	}

	logger.Info("Message Trigger successfully")
	return nil
}

func Trigger(ctx context.Context, logger infrastructure.Logger, publisher publisher.SubscribeMessagePublisher, throttleConfig throttle.Config, rssRepository rss.IRssRepository) error {
	messages, err := getMessages(ctx, rssRepository)
	if err != nil {
		return err
	}
	i := 0
	for _, message := range messages {
		err := publisher.Publish(ctx, message)
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

func getMessages(ctx context.Context, rssRepository rss.IRssRepository) ([]message.Subscribe, error) {
	feeds, err := rssRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var messages []message.Subscribe
	for _, feed := range feeds {
		message := message.Subscribe{
			FeedURL:    feed.Link,
			Language:   feed.Language,
			ItemFilter: feed.ItemFilter,
		}
		messages = append(messages, message)
	}

	return messages, nil
}
