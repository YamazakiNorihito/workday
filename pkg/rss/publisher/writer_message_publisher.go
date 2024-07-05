package publisher

import (
	"context"
	"encoding/json"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
)

type WriterMessagePublisher struct {
	publisher MessagePublisher
}

func NewWriterMessagePublisher(publisher MessagePublisher) *WriterMessagePublisher {
	return &WriterMessagePublisher{publisher: publisher}
}

func (p *WriterMessagePublisher) Publish(ctx context.Context, rssEntry rss.Rss) error {
	message, err := message.NewWriteMessage(rssEntry)
	if err != nil {
		return err
	}

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
