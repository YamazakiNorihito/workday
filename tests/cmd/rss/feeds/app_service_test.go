package clean

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/feeds/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

func TestAppService_AllRssFeeds(t *testing.T) {
	t.Run("should successfully return the original RSS feed for new RSS feed", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := helper.SpyRssRepository{
			FindAllFunc: func(ctx context.Context) ([]rss.Rss, error) {
				dummy_rss1, err := rss.New(
					"ダミーニュースのフィード1",
					"127.0.0.1:8080",
					"http://127.0.0.1:8080/1",
					"このフィードはダミーニュース1を提供します。",
					"ja",
					time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC),
				)
				if err != nil {
					return nil, err
				}

				dummy_rss2, err := rss.New(
					"ダミーニュースのフィード2",
					"127.0.0.1:8081",
					"http://127.0.0.1:8081/2",
					"このフィードはダミーニュース2を提供します。",
					"ja",
					time.Date(2024, time.July, 4, 14, 0, 0, 0, time.UTC),
				)
				if err != nil {
					return nil, err
				}

				return []rss.Rss{dummy_rss1, dummy_rss2}, nil
			},
		}

		// Act
		act_rssFeeds, err := app_service.AllRssFeeds(ctx, &logger, &repo)

		// Assert
		assert.NoError(t, err)

		assert.Len(t, act_rssFeeds, 2)

		// rssFeed1
		assert.Equal(t, "ダミーニュースのフィード1", act_rssFeeds[0].Title)
		assert.Equal(t, "127.0.0.1:8080", act_rssFeeds[0].Source)
		assert.Equal(t, "http://127.0.0.1:8080/1", act_rssFeeds[0].Link)
		assert.Equal(t, time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC), act_rssFeeds[0].LastBuildDate)

		// rssFeed2
		assert.Equal(t, "ダミーニュースのフィード2", act_rssFeeds[1].Title)
		assert.Equal(t, "127.0.0.1:8081", act_rssFeeds[1].Source)
		assert.Equal(t, "http://127.0.0.1:8081/2", act_rssFeeds[1].Link)
		assert.Equal(t, time.Date(2024, time.July, 4, 14, 0, 0, 0, time.UTC), act_rssFeeds[1].LastBuildDate)
	})
}

func TestAppService_AllRssFeeds_Error(t *testing.T) {
	t.Run("should return an error when repository returns an error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := helper.SpyRssRepository{
			FindAllFunc: func(ctx context.Context) ([]rss.Rss, error) {
				return nil, fmt.Errorf("repository error")
			},
		}

		// Act
		act_rssFeeds, err := app_service.AllRssFeeds(ctx, &logger, &repo)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, act_rssFeeds)
		assert.Equal(t, "repository error", err.Error())
	})
}
