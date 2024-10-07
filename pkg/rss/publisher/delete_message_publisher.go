package publisher

import (
	"context"
	"encoding/json"

	"github.com/YamazakiNorihito/workday/pkg/rss/message"
)

type DeleteMessagePublisher struct {
	publisher MessagePublisher
}

func NewDeleteMessagePublisher(publisher MessagePublisher) *DeleteMessagePublisher {
	return &DeleteMessagePublisher{publisher: publisher}
}

func (p *DeleteMessagePublisher) Publish(ctx context.Context, source string) error {
	message := message.Delete{Source: source}

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
