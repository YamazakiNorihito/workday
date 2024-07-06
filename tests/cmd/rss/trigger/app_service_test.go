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

type spyFeedProvider struct{}

func (r *spyFeedProvider) GetFeedURLAndLanguage(ctx context.Context) (map[string]string, error) {
	languageFeedMap := map[string]string{
		"https://azure.microsoft.com/ja-jp/blog/feed/":  "en",
		"https://aws.amazon.com/jp/blogs/news/feed/":    "ja",
		"https://developers-jp.googleblog.com/atom.xml": "ja",
		"https://techblog.nhn-techorus.com/feed":        "ja",
		"https://buildersbox.corp-sansan.com/rss":       "ja",
		"https://knowledge.sakura.ad.jp/rss/":           "ja",
		"https://www.oreilly.co.jp/catalog/soon.xml":    "ja",
		"https://go.dev/blog/feed.atom":                 "ja",
		"https://connpass.com/explore/ja.atom":          "ja",
		"https://www.ipa.go.jp/security/alert-rss.rdf":  "ja",
		"https://feed.infoq.com":                        "en",
		"https://techcrunch.com/feed":                   "en",
	}

	return languageFeedMap, nil
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
		feedRepository := spyFeedProvider{}

		// Act
		err := app_service.Trigger(ctx, &logger, *subscribeMessagePublisher, throttle, &feedRepository)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, messageClient.Messages, 12)
		assert.ElementsMatch(t, messageClient.Messages, []string{
			"{\"feed_url\":\"https://azure.microsoft.com/ja-jp/blog/feed/\",\"language\":\"en\"}",
			"{\"feed_url\":\"https://aws.amazon.com/jp/blogs/news/feed/\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://developers-jp.googleblog.com/atom.xml\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://techblog.nhn-techorus.com/feed\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://buildersbox.corp-sansan.com/rss\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://knowledge.sakura.ad.jp/rss/\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://www.oreilly.co.jp/catalog/soon.xml\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://go.dev/blog/feed.atom\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://connpass.com/explore/ja.atom\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://www.ipa.go.jp/security/alert-rss.rdf\",\"language\":\"ja\"}",
			"{\"feed_url\":\"https://feed.infoq.com\",\"language\":\"en\"}",
			"{\"feed_url\":\"https://techcrunch.com/feed\",\"language\":\"en\"}",
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
				feedRepository := spyFeedProvider{}

				// Act
				err := app_service.Execute(ctx, &logger, *subscribeMessagePublisher, throttle, &feedRepository)

				// Assert
				assert.NoError(t, err)

				// We are not testing the actual sleep duration. Instead, we are verifying
				// that the Sleep function is called the expected number of times based on
				// the specified batch size.
				assert.Equal(t, tc.expectedSleepCount, actSleepCount)
				assert.ElementsMatch(t, messageClient.Messages, []string{
					"{\"feed_url\":\"https://azure.microsoft.com/ja-jp/blog/feed/\",\"language\":\"en\"}",
					"{\"feed_url\":\"https://aws.amazon.com/jp/blogs/news/feed/\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://developers-jp.googleblog.com/atom.xml\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://techblog.nhn-techorus.com/feed\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://buildersbox.corp-sansan.com/rss\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://knowledge.sakura.ad.jp/rss/\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://www.oreilly.co.jp/catalog/soon.xml\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://go.dev/blog/feed.atom\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://connpass.com/explore/ja.atom\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://www.ipa.go.jp/security/alert-rss.rdf\",\"language\":\"ja\"}",
					"{\"feed_url\":\"https://feed.infoq.com\",\"language\":\"en\"}",
					"{\"feed_url\":\"https://techcrunch.com/feed\",\"language\":\"en\"}",
				})
			})
		}
	})
}
