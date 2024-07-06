package subscribe

import (
	"context"
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
		repo := app_service.NewFeedRepository(client, server.URL, "ja")

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

func getPort(rawURL string) (port string) {
	u, _ := url.Parse(rawURL)
	return u.Port()
}
