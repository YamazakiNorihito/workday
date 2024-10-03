package trigger

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/trigger/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
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

type testRssFeed struct {
	URL      string
	Language string

	IncludeKeywords []string
	ExcludeKeywords []string
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

		rssRepository := helper.SpyRssRepository{
			FindAllFunc: func(ctx context.Context) ([]rss.Rss, error) {
				rssFeeds := []testRssFeed{
					{
						URL:             "https://azure.microsoft.com/ja-jp/blog/feed/",
						Language:        "en",
						IncludeKeywords: []string{"Azure", "Cloud", "Microsoft"},
						ExcludeKeywords: []string{"AWS", "Google Cloud"},
					},
					{
						URL:             "https://aws.amazon.com/jp/blogs/news/feed/",
						Language:        "ja",
						IncludeKeywords: []string{"AWS", "Lambda", "S3"},
						ExcludeKeywords: nil, // Only IncludeKeywords are set
					},
					{
						URL:             "https://developers-jp.googleblog.com/atom.xml",
						Language:        "ja",
						IncludeKeywords: nil, // Only ExcludeKeywords are set
						ExcludeKeywords: []string{"Azure", "AWS"},
					},
					{
						URL:             "https://techblog.nhn-techorus.com/feed",
						Language:        "ja",
						IncludeKeywords: nil, // Neither IncludeKeywords nor ExcludeKeywords are set
						ExcludeKeywords: nil,
					},
					{
						URL:             "https://buildersbox.corp-sansan.com/rss",
						Language:        "ja",
						IncludeKeywords: []string{"Sansan", "Cloud", "API"},
						ExcludeKeywords: []string{"AWS", "Google"},
					},
					{
						URL:             "https://knowledge.sakura.ad.jp/rss/",
						Language:        "ja",
						IncludeKeywords: nil, // Only ExcludeKeywords are set
						ExcludeKeywords: []string{"Microsoft", "Amazon"},
					},
					{
						URL:             "https://www.oreilly.co.jp/catalog/soon.xml",
						Language:        "ja",
						IncludeKeywords: []string{"O'Reilly", "Books", "Technology"},
						ExcludeKeywords: nil, // Only IncludeKeywords are set
					},
					{
						URL:             "https://go.dev/blog/feed.atom",
						Language:        "ja",
						IncludeKeywords: nil, // Neither IncludeKeywords nor ExcludeKeywords are set
						ExcludeKeywords: nil,
					},
					{
						URL:             "https://connpass.com/explore/ja.atom",
						Language:        "ja",
						IncludeKeywords: []string{"Connpass", "Events", "Tech"},
						ExcludeKeywords: []string{"Non-Tech", "Marketing"},
					},
					{
						URL:             "https://www.ipa.go.jp/security/alert-rss.rdf",
						Language:        "ja",
						IncludeKeywords: []string{"Security", "IPA", "Alerts"},
						ExcludeKeywords: []string{"Old Versions"},
					},
					{
						URL:             "https://feed.infoq.com",
						Language:        "en",
						IncludeKeywords: []string{"InfoQ", "Technology", "Software"},
						ExcludeKeywords: nil, // Only IncludeKeywords are set
					},
					{
						URL:             "https://techcrunch.com/feed",
						Language:        "en",
						IncludeKeywords: nil, // Only ExcludeKeywords are set
						ExcludeKeywords: []string{"Non-Tech", "Startups"},
					},
				}

				var rssList []rss.Rss
				counter := 1

				for _, feed := range rssFeeds {
					dummy_rss, err := rss.New(
						fmt.Sprintf("ダミーニュースのフィード%d", counter),
						fmt.Sprintf("127.0.0.1:808%d", counter),
						feed.URL,
						fmt.Sprintf("このフィードはダミーニュース%dを提供します。", counter),
						feed.Language,
						time.Date(2024, time.July, 3+counter, 13+counter, 0, 0, 0, time.UTC),
					)
					if err != nil {
						return nil, err
					}

					dummy_rss.SetItemFilter(feed.IncludeKeywords, feed.ExcludeKeywords)
					rssList = append(rssList, dummy_rss)
					counter++
				}

				return rssList, nil
			},
		}

		// Act
		err := app_service.Trigger(ctx, &logger, *subscribeMessagePublisher, throttle, &rssRepository)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, messageClient.Messages, 12)
		assert.ElementsMatch(t, messageClient.Messages, []string{
			"{\"feed_url\":\"https://azure.microsoft.com/ja-jp/blog/feed/\",\"language\":\"en\",\"item_filter\":{\"include_keywords\":[\"Azure\",\"Cloud\",\"Microsoft\"],\"exclude_keywords\":[\"AWS\",\"Google Cloud\"]}}",
			"{\"feed_url\":\"https://aws.amazon.com/jp/blogs/news/feed/\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"AWS\",\"Lambda\",\"S3\"],\"exclude_keywords\":[]}}",
			"{\"feed_url\":\"https://developers-jp.googleblog.com/atom.xml\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[\"Azure\",\"AWS\"]}}",
			"{\"feed_url\":\"https://techblog.nhn-techorus.com/feed\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[]}}",
			"{\"feed_url\":\"https://buildersbox.corp-sansan.com/rss\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"Sansan\",\"Cloud\",\"API\"],\"exclude_keywords\":[\"AWS\",\"Google\"]}}",
			"{\"feed_url\":\"https://knowledge.sakura.ad.jp/rss/\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[\"Microsoft\",\"Amazon\"]}}",
			"{\"feed_url\":\"https://www.oreilly.co.jp/catalog/soon.xml\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"O'Reilly\",\"Books\",\"Technology\"],\"exclude_keywords\":[]}}",
			"{\"feed_url\":\"https://go.dev/blog/feed.atom\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[]}}",
			"{\"feed_url\":\"https://connpass.com/explore/ja.atom\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"Connpass\",\"Events\",\"Tech\"],\"exclude_keywords\":[\"Non-Tech\",\"Marketing\"]}}",
			"{\"feed_url\":\"https://www.ipa.go.jp/security/alert-rss.rdf\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"Security\",\"IPA\",\"Alerts\"],\"exclude_keywords\":[\"Old Versions\"]}}",
			"{\"feed_url\":\"https://feed.infoq.com\",\"language\":\"en\",\"item_filter\":{\"include_keywords\":[\"InfoQ\",\"Technology\",\"Software\"],\"exclude_keywords\":[]}}",
			"{\"feed_url\":\"https://techcrunch.com/feed\",\"language\":\"en\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[\"Non-Tech\",\"Startups\"]}}",
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

		rssRepository := helper.SpyRssRepository{
			FindAllFunc: func(ctx context.Context) ([]rss.Rss, error) {
				rssFeeds := []testRssFeed{
					{
						URL:             "https://azure.microsoft.com/ja-jp/blog/feed/",
						Language:        "en",
						IncludeKeywords: []string{"Azure", "Cloud", "Microsoft"},
						ExcludeKeywords: []string{"AWS", "Google Cloud"},
					},
					{
						URL:             "https://aws.amazon.com/jp/blogs/news/feed/",
						Language:        "ja",
						IncludeKeywords: []string{"AWS", "Lambda", "S3"},
						ExcludeKeywords: nil, // Only IncludeKeywords are set
					},
					{
						URL:             "https://developers-jp.googleblog.com/atom.xml",
						Language:        "ja",
						IncludeKeywords: nil, // Only ExcludeKeywords are set
						ExcludeKeywords: []string{"Azure", "AWS"},
					},
					{
						URL:             "https://techblog.nhn-techorus.com/feed",
						Language:        "ja",
						IncludeKeywords: nil, // Neither IncludeKeywords nor ExcludeKeywords are set
						ExcludeKeywords: nil,
					},
					{
						URL:             "https://buildersbox.corp-sansan.com/rss",
						Language:        "ja",
						IncludeKeywords: []string{"Sansan", "Cloud", "API"},
						ExcludeKeywords: []string{"AWS", "Google"},
					},
					{
						URL:             "https://knowledge.sakura.ad.jp/rss/",
						Language:        "ja",
						IncludeKeywords: nil, // Only ExcludeKeywords are set
						ExcludeKeywords: []string{"Microsoft", "Amazon"},
					},
					{
						URL:             "https://www.oreilly.co.jp/catalog/soon.xml",
						Language:        "ja",
						IncludeKeywords: []string{"O'Reilly", "Books", "Technology"},
						ExcludeKeywords: nil, // Only IncludeKeywords are set
					},
					{
						URL:             "https://go.dev/blog/feed.atom",
						Language:        "ja",
						IncludeKeywords: nil, // Neither IncludeKeywords nor ExcludeKeywords are set
						ExcludeKeywords: nil,
					},
					{
						URL:             "https://connpass.com/explore/ja.atom",
						Language:        "ja",
						IncludeKeywords: []string{"Connpass", "Events", "Tech"},
						ExcludeKeywords: []string{"Non-Tech", "Marketing"},
					},
					{
						URL:             "https://www.ipa.go.jp/security/alert-rss.rdf",
						Language:        "ja",
						IncludeKeywords: []string{"Security", "IPA", "Alerts"},
						ExcludeKeywords: []string{"Old Versions"},
					},
					{
						URL:             "https://feed.infoq.com",
						Language:        "en",
						IncludeKeywords: []string{"InfoQ", "Technology", "Software"},
						ExcludeKeywords: nil, // Only IncludeKeywords are set
					},
					{
						URL:             "https://techcrunch.com/feed",
						Language:        "en",
						IncludeKeywords: nil, // Only ExcludeKeywords are set
						ExcludeKeywords: []string{"Non-Tech", "Startups"},
					},
				}

				var rssList []rss.Rss
				counter := 1

				for _, feed := range rssFeeds {
					dummy_rss, err := rss.New(
						fmt.Sprintf("ダミーニュースのフィード%d", counter),
						fmt.Sprintf("127.0.0.1:808%d", counter),
						feed.URL,
						fmt.Sprintf("このフィードはダミーニュース%dを提供します。", counter),
						feed.Language,
						time.Date(2024, time.July, 3+counter, 13+counter, 0, 0, 0, time.UTC),
					)
					if err != nil {
						return nil, err
					}

					dummy_rss.SetItemFilter(feed.IncludeKeywords, feed.ExcludeKeywords)

					rssList = append(rssList, dummy_rss)
					counter++
				}

				return rssList, nil
			},
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

				// Act
				err := app_service.Execute(ctx, &logger, *subscribeMessagePublisher, throttle, &rssRepository)

				// Assert
				assert.NoError(t, err)

				// We are not testing the actual sleep duration. Instead, we are verifying
				// that the Sleep function is called the expected number of times based on
				// the specified batch size.
				assert.Equal(t, tc.expectedSleepCount, actSleepCount)
				assert.ElementsMatch(t, messageClient.Messages, []string{
					"{\"feed_url\":\"https://azure.microsoft.com/ja-jp/blog/feed/\",\"language\":\"en\",\"item_filter\":{\"include_keywords\":[\"Azure\",\"Cloud\",\"Microsoft\"],\"exclude_keywords\":[\"AWS\",\"Google Cloud\"]}}",
					"{\"feed_url\":\"https://aws.amazon.com/jp/blogs/news/feed/\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"AWS\",\"Lambda\",\"S3\"],\"exclude_keywords\":[]}}",
					"{\"feed_url\":\"https://developers-jp.googleblog.com/atom.xml\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[\"Azure\",\"AWS\"]}}",
					"{\"feed_url\":\"https://techblog.nhn-techorus.com/feed\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[]}}",
					"{\"feed_url\":\"https://buildersbox.corp-sansan.com/rss\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"Sansan\",\"Cloud\",\"API\"],\"exclude_keywords\":[\"AWS\",\"Google\"]}}",
					"{\"feed_url\":\"https://knowledge.sakura.ad.jp/rss/\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[\"Microsoft\",\"Amazon\"]}}",
					"{\"feed_url\":\"https://www.oreilly.co.jp/catalog/soon.xml\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"O'Reilly\",\"Books\",\"Technology\"],\"exclude_keywords\":[]}}",
					"{\"feed_url\":\"https://go.dev/blog/feed.atom\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[]}}",
					"{\"feed_url\":\"https://connpass.com/explore/ja.atom\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"Connpass\",\"Events\",\"Tech\"],\"exclude_keywords\":[\"Non-Tech\",\"Marketing\"]}}",
					"{\"feed_url\":\"https://www.ipa.go.jp/security/alert-rss.rdf\",\"language\":\"ja\",\"item_filter\":{\"include_keywords\":[\"Security\",\"IPA\",\"Alerts\"],\"exclude_keywords\":[\"Old Versions\"]}}",
					"{\"feed_url\":\"https://feed.infoq.com\",\"language\":\"en\",\"item_filter\":{\"include_keywords\":[\"InfoQ\",\"Technology\",\"Software\"],\"exclude_keywords\":[]}}",
					"{\"feed_url\":\"https://techcrunch.com/feed\",\"language\":\"en\",\"item_filter\":{\"include_keywords\":[],\"exclude_keywords\":[\"Non-Tech\",\"Startups\"]}}",
				})
			})
		}
	})
}
