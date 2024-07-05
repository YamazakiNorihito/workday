package publisher

import (
	"context"
	"encoding/json"

	"github.com/YamazakiNorihito/workday/pkg/rss/message"
)

type SubscribeMessagePublisher struct {
	publisher MessagePublisher
}

func NewSubscribeMessagePublisher(publisher MessagePublisher) *SubscribeMessagePublisher {
	return &SubscribeMessagePublisher{publisher: publisher}
}

func (p *SubscribeMessagePublisher) Publish(ctx context.Context, feedURL string) error {
	message := message.Subscribe{FeedURL: feedURL}

	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.publisher.Publish(ctx, string(messageJson))
	if err != nil {
		return err
	}
	return nil
}
