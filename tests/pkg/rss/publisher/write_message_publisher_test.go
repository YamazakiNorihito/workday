package publisher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

type spyMessageClient struct{ Messages []string }

func (r *spyMessageClient) Publish(ctx context.Context, message string) error {
	r.Messages = append(r.Messages, message)
	return nil
}

func TestWriterMessagePublisher_Publish(t *testing.T) {
	t.Run("should be published as WriteMessage", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("ダミーニュースのフィード", "127.0.0.1:8080", "http://127.0.0.1:8080", "このフィードはダミーニュースを提供します。", "ja", time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC))
			if err != nil {
				return err
			}

			dummy_item, err := rss.NewItem(rss.Guid{Value: "http://www.example.com/dummy-guid1"}, "ダミー記事1", "http://www.example.com/dummy-article1", "これはダミー記事1の概要です。詳細はリンクをクリックしてください。", "item1@dummy.com", time.Date(2024, time.July, 3, 12, 0, 0, 0, time.UTC))
			if err != nil {
				return err
			}
			test_rss.AddOrUpdateItem(dummy_item)
			return err
		})
		ctx := context.Background()
		messageClient := spyMessageClient{}
		writerMessagePublisher := publisher.NewWriterMessagePublisher(&messageClient)
		// Act
		err := writerMessagePublisher.Publish(ctx, test_rss)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, messageClient.Messages, 1)

		expectedJSON := fmt.Sprintf(`{
			"rss": {
			  "id": "%s",
			  "source": "127.0.0.1:8080",
			  "title": "ダミーニュースのフィード",
			  "link": "http://127.0.0.1:8080",
			  "description": "このフィードはダミーニュースを提供します。",
			  "language": "ja",
			  "last_build_date": "2024-07-03T13:00:00Z",
			  "items": {
				"http://www.example.com/dummy-guid1": {
				  "guid": "http://www.example.com/dummy-guid1",
				  "title": "ダミー記事1",
				  "link": "http://www.example.com/dummy-article1",
				  "description": "これはダミー記事1の概要です。詳細はリンクをクリックしてください。",
				  "author": "item1@dummy.com",
				  "pubDate": "2024-07-03T12:00:00Z",
				  "tags": []
				}
			  },
			  "create_by": {
				"id": "",
				"name": ""
			  },
			  "create_at": "0001-01-01T00:00:00Z",
			  "update_by": {
				"id": "",
				"name": ""
			  },
			  "update_at": "0001-01-01T00:00:00Z"
			},
			"compressed": false
		  }`, test_rss.ID.String())

		assert.JSONEq(t, expectedJSON, messageClient.Messages[0])
	})
}
