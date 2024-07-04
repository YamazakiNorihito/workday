package clean

import (
	"context"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/notification/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

type call struct {
	Text     string
	Username string
}

type spySlackChannelClient struct {
	Calls []call
}

func (s *spySlackChannelClient) PostMessageContext(ctx context.Context, text string, username string) (respChannel string, respTimestamp string, err error) {
	s.Calls = append(s.Calls, call{Text: text, Username: username})
	return "C1234567890", "1234567890.123456", nil
}

func TestAppService_Notification(t *testing.T) {
	t.Run("should notify Slack when new articles are added", func(t *testing.T) {
		// Arrange
		test_rss := generatorTestRss(t)
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := helper.SpyRssRepository{
			FindBySourceFunc: func(ctx context.Context, source string) (rss.Rss, error) {
				if source != "127.0.0.1:8080" {
					panic("sourceLanguageCode is not 'en' as expected")
				}
				return test_rss, nil
			},
			FindItemsFunc: func(ctx context.Context, rss rss.Rss) (rss.Rss, error) {
				if rss.ID.String() != test_rss.ID.String() {
					panic("rss is not 'en' as expected")
				}
				return test_rss, nil
			},
		}
		slackChannelClient := spySlackChannelClient{}

		conditions := app_service.RssConditions{
			Target:     func(rss.Rss) bool { return true },
			ItemFilter: func(item rss.Item) bool { return true },
		}
		source := "127.0.0.1:8080"

		// Act
		err := app_service.Notification(ctx, &logger, &repo, &slackChannelClient, conditions, source)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, slackChannelClient.Calls, 1)
		act_call := slackChannelClient.Calls[0]
		assert.Equal(t, "127.0.0.1:8080", act_call.Username)
		assert.Equal(t, `*フィードタイトル:* <http://127.0.0.1:8080|ダミーニュースのフィード>
*フィード詳細:* このフィードはダミーニュースを提供します。
*最終更新日:* 2024-07-03T13:00:00Z

*最新の記事:*
1. *記事タイトル:* <http://www.example.com/dummy-article1|ダミー記事1>
    *公開日:* 2024-07-03T12:00:00Z
    *概要:* これはダミー記事1の概要です。詳細はリンクをクリックしてください。
    *カテゴリ:* 

2. *記事タイトル:* <http://www.example.com/dummy-article2|ダミー記事2>
    *公開日:* 2024-07-03T12:30:00Z
    *概要:* これはダミー記事2の概要です。詳細はリンクをクリックしてください。
    *カテゴリ:* 

`, act_call.Text)
	})
	t.Run("should notify Slack only for RSS feeds matching specified conditions", func(t *testing.T) {
		// Arrange
		dummy_rss := generatorTestRss(t)
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := helper.SpyRssRepository{
			FindBySourceFunc: func(ctx context.Context, source string) (rss.Rss, error) {
				if source != "127.0.0.1:8080" {
					panic("sourceLanguageCode is not 'en' as expected")
				}
				return dummy_rss, nil
			},
		}
		slackChannelClient := spySlackChannelClient{}

		conditions := app_service.RssConditions{
			Target:     func(rss.Rss) bool { return false },
			ItemFilter: func(item rss.Item) bool { return true },
		}
		source := "127.0.0.1:8080"

		// Act
		err := app_service.Notification(ctx, &logger, &repo, &slackChannelClient, conditions, source)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, slackChannelClient.Calls, 0)
	})
	t.Run("should notify Slack only for items matching specified conditions", func(t *testing.T) {
		// Arrange
		dummy_rss := generatorTestRss(t)
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := helper.SpyRssRepository{
			FindBySourceFunc: func(ctx context.Context, source string) (rss.Rss, error) {
				if source != "127.0.0.1:8080" {
					panic("sourceLanguageCode is not 'en' as expected")
				}
				return dummy_rss, nil
			},
			FindItemsFunc: func(ctx context.Context, rss rss.Rss) (rss.Rss, error) {
				if rss.ID.String() != dummy_rss.ID.String() {
					panic("rss is not 'en' as expected")
				}
				return dummy_rss, nil
			},
		}
		slackChannelClient := spySlackChannelClient{}

		conditions := app_service.RssConditions{
			Target: func(rss.Rss) bool { return true },
			ItemFilter: func(item rss.Item) bool {
				return item.Guid == rss.Guid{Value: "http://www.example.com/dummy-guid2"}
			},
		}
		source := "127.0.0.1:8080"

		// Act
		err := app_service.Notification(ctx, &logger, &repo, &slackChannelClient, conditions, source)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, slackChannelClient.Calls, 1)
		act_call := slackChannelClient.Calls[0]
		assert.Equal(t, "127.0.0.1:8080", act_call.Username)
		assert.Equal(t, `*フィードタイトル:* <http://127.0.0.1:8080|ダミーニュースのフィード>
*フィード詳細:* このフィードはダミーニュースを提供します。
*最終更新日:* 2024-07-03T13:00:00Z

*最新の記事:*
1. *記事タイトル:* <http://www.example.com/dummy-article2|ダミー記事2>
    *公開日:* 2024-07-03T12:30:00Z
    *概要:* これはダミー記事2の概要です。詳細はリンクをクリックしてください。
    *カテゴリ:* 

`, act_call.Text)
	})
}

func generatorTestRss(t *testing.T) rss.Rss {
	var dummy_rss rss.Rss
	helper.MustSucceed(t, func() error {
		var err error
		dummy_rss, err = rss.New("ダミーニュースのフィード", "127.0.0.1:8080", "http://127.0.0.1:8080", "このフィードはダミーニュースを提供します。", "ja", time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC))
		if err != nil {
			return err
		}

		dummy_item1, err := rss.NewItem(rss.Guid{Value: "http://www.example.com/dummy-guid1"}, "ダミー記事1", "http://www.example.com/dummy-article1", "これはダミー記事1の概要です。詳細はリンクをクリックしてください。", "item1@dummy.com", time.Date(2024, time.July, 3, 12, 0, 0, 0, time.UTC))
		if err != nil {
			return err
		}
		dummy_rss.AddOrUpdateItem(dummy_item1)

		dummy_item2, err := rss.NewItem(rss.Guid{Value: "http://www.example.com/dummy-guid2"}, "ダミー記事2", "http://www.example.com/dummy-article2", "これはダミー記事2の概要です。詳細はリンクをクリックしてください。", "item2@dummy.com", time.Date(2024, time.July, 3, 12, 30, 0, 0, time.UTC))
		if err != nil {
			return err
		}
		dummy_rss.AddOrUpdateItem(dummy_item2)
		return err
	})

	return dummy_rss
}
