package subscribe

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/subscribe/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

type spyRssWritePublisher struct{ Messages []string }

func (r *spyRssWritePublisher) Publish(ctx context.Context, message string) error {
	r.Messages = append(r.Messages, message)
	return nil
}

func TestAppService_Subscribe(t *testing.T) {
	t.Run("should successfully return RSS feed", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mockFeed := `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
<channel>
  <title>ダミーニュースのフィード</title>
  <link>http://www.example.com/</link>
  <description>このフィードはダミーニュースを提供します。</description>
  <language>ja</language>

  <item>
    <title>ダミー記事1</title>
    <guid>http://www.example.com/dummy-guid1</guid>
    <link>http://www.example.com/dummy-article1</link>
    <description>これはダミー記事1の概要です。詳細はリンクをクリックしてください。</description>
    <pubDate>Mon, 03 Jul 2024 12:00:00 GMT</pubDate>
    <author>item1@dummy.com</author>
  </item>

  <item>
    <title>ダミー記事2</title>
    <guid>http://www.example.com/dummy-guid2</guid>
    <link>http://www.example.com/dummy-article2</link>
    <description>これはダミー記事2の概要です。詳細はリンクをクリックしてください。</description>
    <pubDate>Mon, 03 Jul 2024 12:30:00 GMT</pubDate>
    <author>item2@dummy.com</author>
  </item>

  <item>
    <title>ダミー記事3</title>
    <guid>http://www.example.com/dummy-guid3</guid>
    <link>http://www.example.com/dummy-article3</link>
    <description>これはダミー記事3の概要です。詳細はリンクをクリックしてください。</description>
    <pubDate>Mon, 03 Jul 2024 13:00:00 GMT</pubDate>
    <author>item3@dummy.com</author>
  </item>

</channel>
</rss>`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockFeed))
		}))
		defer server.Close()

		ctx := context.Background()
		logger := helper.MockLogger{}

		client := server.Client()
		repo := app_service.NewFeedRepository(client, server.URL)

		// Act
		act_rss, err := app_service.Subscribe(ctx, &logger, &repo)

		// Assert
		assert.NoError(t, err)

		assert.NotEmpty(t, act_rss.ID)
		assert.Equal(t, "ダミーニュースのフィード", act_rss.Title)
		port := getPort(server.URL)
		assert.Equal(t, "127.0.0.1:"+port, act_rss.Source)
		assert.Equal(t, server.URL, act_rss.Link)
		assert.Equal(t, "このフィードはダミーニュースを提供します。", act_rss.Description)
		assert.Equal(t, "ja", act_rss.Language)
		assert.Equal(t, time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC), act_rss.LastBuildDate)

		assert.Len(t, act_rss.Items, 3)
		// item1
		{
			item1 := act_rss.Items[rss.Guid{Value: "http://www.example.com/dummy-guid1"}]
			assert.Equal(t, rss.Guid{Value: "http://www.example.com/dummy-guid1"}, item1.Guid)
			assert.Equal(t, "ダミー記事1", item1.Title)
			assert.Equal(t, "http://www.example.com/dummy-article1", item1.Link)
			assert.Equal(t, "これはダミー記事1の概要です。詳細はリンクをクリックしてください。", item1.Description)
			assert.Equal(t, time.Date(2024, time.July, 3, 12, 0, 0, 0, time.UTC), item1.PubDate)
			assert.Equal(t, "item1@dummy.com", item1.Author)
		}

		// item2
		{
			item2 := act_rss.Items[rss.Guid{Value: "http://www.example.com/dummy-guid2"}]
			assert.Equal(t, rss.Guid{Value: "http://www.example.com/dummy-guid2"}, item2.Guid)
			assert.Equal(t, "ダミー記事2", item2.Title)
			assert.Equal(t, "http://www.example.com/dummy-article2", item2.Link)
			assert.Equal(t, "これはダミー記事2の概要です。詳細はリンクをクリックしてください。", item2.Description)
			assert.Equal(t, time.Date(2024, time.July, 3, 12, 30, 0, 0, time.UTC), item2.PubDate)
			assert.Equal(t, "item2@dummy.com", item2.Author)
		}

		// item3
		{
			item3 := act_rss.Items[rss.Guid{Value: "http://www.example.com/dummy-guid3"}]
			assert.Equal(t, rss.Guid{Value: "http://www.example.com/dummy-guid3"}, item3.Guid)
			assert.Equal(t, "ダミー記事3", item3.Title)
			assert.Equal(t, "http://www.example.com/dummy-article3", item3.Link)
			assert.Equal(t, "これはダミー記事3の概要です。詳細はリンクをクリックしてください。", item3.Description)
			assert.Equal(t, time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC), item3.PubDate)
			assert.Equal(t, "item3@dummy.com", item3.Author)
		}
	})
}

func TestAppService_Publish(t *testing.T) {
	t.Run("should be published as WriteMessage", func(t *testing.T) {
		// Arrange
		var dummy_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			dummy_rss, err = rss.New("ダミーニュースのフィード", "127.0.0.1:8080", "http://127.0.0.1:8080", "このフィードはダミーニュースを提供します。", "ja", time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC))
			if err != nil {
				return err
			}

			dummy_item, err := rss.NewItem(rss.Guid{Value: "http://www.example.com/dummy-guid1"}, "ダミー記事1", "http://www.example.com/dummy-article1", "これはダミー記事1の概要です。詳細はリンクをクリックしてください。", "item1@dummy.com", time.Date(2024, time.July, 3, 12, 0, 0, 0, time.UTC))
			if err != nil {
				return err
			}
			dummy_rss.AddOrUpdateItem(dummy_item)
			return err
		})
		ctx := context.Background()
		publisher := spyRssWritePublisher{}

		// Act
		err := app_service.Publish(ctx, &publisher, dummy_rss)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, publisher.Messages, 1)

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
		  }`, dummy_rss.ID.String())

		assert.JSONEq(t, expectedJSON, publisher.Messages[0])
	})
}

func getPort(rawURL string) (port string) {
	u, _ := url.Parse(rawURL)
	return u.Port()
}
