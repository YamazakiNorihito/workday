package clean

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/write/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

func TestAppService_Clean(t *testing.T) {
	t.Run("should successfully return the original RSS feed for new RSS feed", func(t *testing.T) {
		// Arrange
		test_rss := generatorTestRss(t)
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := helper.SpyRssRepository{
			FindBySourceFunc: func(ctx context.Context, source string) (rss.Rss, error) {
				return rss.Rss{}, nil
			},
			SaveFunc: func(ctx context.Context, entryRss rss.Rss, updateBy metadata.UserMeta) (rss.Rss, error) {
				var copy rss.Rss
				helper.MustSucceed(t, func() error { return deepCopy(entryRss, &copy) })

				copy.CreatedBy = metadata.CreateBy(updateBy)
				copy.CreatedAt = metadata.CreateAt(time.Date(2024, time.July, 3, 14, 0, 0, 0, time.UTC))
				copy.UpdatedBy = metadata.UpdateBy(updateBy)
				copy.UpdatedAt = metadata.UpdateAt(time.Date(2024, time.July, 3, 14, 0, 0, 0, time.UTC))
				return copy, nil
			},
		}

		// Act
		act_rss, err := app_service.Write(ctx, &logger, &repo, test_rss)

		// Assert
		assert.NoError(t, err)

		assert.NotEmpty(t, act_rss.ID)
		assert.Equal(t, "ダミーニュースのフィード", act_rss.Title)
		assert.Equal(t, "127.0.0.1:8080", act_rss.Source)
		assert.Equal(t, "http://127.0.0.1:8080", act_rss.Link)
		assert.Equal(t, "このフィードはダミーニュースを提供します。", act_rss.Description)
		assert.Equal(t, "ja", act_rss.Language)
		assert.Equal(t, time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC), act_rss.LastBuildDate)

		assert.Equal(t, metadata.CreateBy(metadata.UserMeta{ID: "127.0.0.1:8080", Name: "127.0.0.1:8080"}), act_rss.CreatedBy)
		assert.Equal(t, time.Date(2024, time.July, 3, 14, 0, 0, 0, time.UTC), act_rss.CreatedAt)
		assert.Equal(t, metadata.UpdateBy(metadata.UserMeta{ID: "127.0.0.1:8080", Name: "127.0.0.1:8080"}), act_rss.UpdatedBy)
		assert.Equal(t, time.Date(2024, time.July, 3, 14, 0, 0, 0, time.UTC), act_rss.UpdatedAt)

		assert.Len(t, act_rss.Items, 2)
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
