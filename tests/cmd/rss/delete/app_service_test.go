package clean

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/delete/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

func TestAppService_Clean(t *testing.T) {
	t.Run("should delete RSS feed when found by source", func(t *testing.T) {
		// Arrange
		test_rss := generatorTestRss(t)
		ctx := context.Background()
		logger := helper.MockLogger{}

		act_rss := rss.Rss{}
		repo := helper.SpyRssRepository{
			FindBySourceFunc: func(ctx context.Context, source string) (rss.Rss, error) {
				if source == "127.0.0.1:8080" {
					return test_rss, nil
				}
				return rss.Rss{}, errors.New("RSS feed not found")
			},
			DeleteFunc: func(ctx context.Context, rss rss.Rss) error {
				act_rss = rss
				return nil
			},
		}

		// Act
		err := app_service.Delete(ctx, &logger, &repo, "127.0.0.1:8080")

		// Assert
		assert.NoError(t, err)

		assert.NotEmpty(t, act_rss.ID)
		assert.Equal(t, test_rss, act_rss)
	})

	t.Run("should return error when RSS feed not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := helper.SpyRssRepository{
			FindBySourceFunc: func(ctx context.Context, source string) (rss.Rss, error) {
				return rss.Rss{}, errors.New("RSS feed not found")
			},
			DeleteFunc: func(ctx context.Context, rss rss.Rss) error {
				return nil
			},
		}

		// Act
		err := app_service.Delete(ctx, &logger, &repo, "127.0.0.1:8080")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "RSS feed not found", err.Error())
	})

	t.Run("should return error when Delete fails", func(t *testing.T) {
		// Arrange
		test_rss := generatorTestRss(t)
		ctx := context.Background()
		logger := helper.MockLogger{}

		repo := helper.SpyRssRepository{
			FindBySourceFunc: func(ctx context.Context, source string) (rss.Rss, error) {
				return test_rss, nil
			},
			DeleteFunc: func(ctx context.Context, rss rss.Rss) error {
				return errors.New("failed to delete RSS feed")
			},
		}

		// Act
		err := app_service.Delete(ctx, &logger, &repo, "127.0.0.1:8080")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "failed to delete RSS feed", err.Error())
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

func deepCopy(src, dest interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}
