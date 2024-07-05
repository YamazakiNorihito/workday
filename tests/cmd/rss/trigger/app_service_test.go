package trigger

import (
	"context"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/trigger/app_service"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/YamazakiNorihito/workday/pkg/throttle"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

type spyMessageClient struct{ Messages []string }

func (r *spyMessageClient) Publish(ctx context.Context, message string) error {
	r.Messages = append(r.Messages, message)
	return nil
}

type spyFeedRepository struct{}

func (r *spyFeedRepository) GetFeedURLs() []string {
	return []string{
		"https://azure.microsoft.com/ja-jp/blog/feed/",
		"https://aws.amazon.com/jp/blogs/news/feed/",
		"https://developers-jp.googleblog.com/atom.xml",
		"https://techblog.nhn-techorus.com/feed",
		"https://buildersbox.corp-sansan.com/rss",
		"https://knowledge.sakura.ad.jp/rss/",
		"https://www.oreilly.co.jp/catalog/soon.xml",
		"https://go.dev/blog/feed.atom",
		"https://connpass.com/explore/ja.atom",
		"https://www.ipa.go.jp/security/alert-rss.rdf",
		"https://feed.infoq.com",
		"https://techcrunch.com/feed",
	}
}

func TestAppService_Execute(t *testing.T) {
	t.Run("should handle event successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		logger := helper.MockLogger{}
		messageClient := spyMessageClient{}
		subscribeMessagePublisher := publisher.NewSubscribeMessagePublisher(&messageClient)
		throttle := throttle.Config{
			BatchSize: 1,
			Sleep:     func() { time.Sleep(1 * time.Microsecond) },
		}
		feedRepository := spyFeedRepository{}

		// Act
		err := app_service.Trigger(ctx, &logger, *subscribeMessagePublisher, throttle, &feedRepository)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, messageClient.Messages, 12)
		assert.ElementsMatch(t, messageClient.Messages, []string{
			"{\"feed_url\":\"https://azure.microsoft.com/ja-jp/blog/feed/\"}",
			"{\"feed_url\":\"https://aws.amazon.com/jp/blogs/news/feed/\"}",
			"{\"feed_url\":\"https://developers-jp.googleblog.com/atom.xml\"}",
			"{\"feed_url\":\"https://techblog.nhn-techorus.com/feed\"}",
			"{\"feed_url\":\"https://buildersbox.corp-sansan.com/rss\"}",
			"{\"feed_url\":\"https://knowledge.sakura.ad.jp/rss/\"}",
			"{\"feed_url\":\"https://www.oreilly.co.jp/catalog/soon.xml\"}",
			"{\"feed_url\":\"https://go.dev/blog/feed.atom\"}",
			"{\"feed_url\":\"https://connpass.com/explore/ja.atom\"}",
			"{\"feed_url\":\"https://www.ipa.go.jp/security/alert-rss.rdf\"}",
			"{\"feed_url\":\"https://feed.infoq.com\"}",
			"{\"feed_url\":\"https://techcrunch.com/feed\"}",
		})
	})

	t.Run("should handle event with throttle successfully", func(t *testing.T) {
		testCases := []struct {
			name               string
			batchSize          int
			expectedSleepCount int
		}{
			{name: "BatchSize 1", batchSize: 1, expectedSleepCount: 12},
			{name: "BatchSize 3", batchSize: 3, expectedSleepCount: 4},
			{name: "BatchSize 5", batchSize: 5, expectedSleepCount: 2},
			{name: "BatchSize 10", batchSize: 10, expectedSleepCount: 1},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				ctx := context.Background()
				logger := helper.MockLogger{}
				messageClient := spyMessageClient{}
				subscribeMessagePublisher := publisher.NewSubscribeMessagePublisher(&messageClient)
				actSleepCount := 0
				throttle := throttle.Config{
					BatchSize: tc.batchSize,
					Sleep: func() {
						actSleepCount++
						time.Sleep(1 * time.Microsecond)
					},
				}
				feedRepository := spyFeedRepository{}

				// Act
				err := app_service.Execute(ctx, &logger, *subscribeMessagePublisher, throttle, &feedRepository)

				// Assert
				assert.NoError(t, err)

				// We are not testing the actual sleep duration. Instead, we are verifying
				// that the Sleep function is called the expected number of times based on
				// the specified batch size.
				assert.Equal(t, tc.expectedSleepCount, actSleepCount)
				assert.ElementsMatch(t, messageClient.Messages, []string{
					"{\"feed_url\":\"https://azure.microsoft.com/ja-jp/blog/feed/\"}",
					"{\"feed_url\":\"https://aws.amazon.com/jp/blogs/news/feed/\"}",
					"{\"feed_url\":\"https://developers-jp.googleblog.com/atom.xml\"}",
					"{\"feed_url\":\"https://techblog.nhn-techorus.com/feed\"}",
					"{\"feed_url\":\"https://buildersbox.corp-sansan.com/rss\"}",
					"{\"feed_url\":\"https://knowledge.sakura.ad.jp/rss/\"}",
					"{\"feed_url\":\"https://www.oreilly.co.jp/catalog/soon.xml\"}",
					"{\"feed_url\":\"https://go.dev/blog/feed.atom\"}",
					"{\"feed_url\":\"https://connpass.com/explore/ja.atom\"}",
					"{\"feed_url\":\"https://www.ipa.go.jp/security/alert-rss.rdf\"}",
					"{\"feed_url\":\"https://feed.infoq.com\"}",
					"{\"feed_url\":\"https://techcrunch.com/feed\"}",
				})
			})
		}
	})
}
