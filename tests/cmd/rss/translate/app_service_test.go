package subscribe

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/translate/app_service"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

type spyFeedRepository struct {
	getSourceLanguageFunc func(source string) (sourceLanguageCode string, ok bool)
}

func (r *spyFeedRepository) GetSourceLanguage(source string) (sourceLanguageCode string, ok bool) {
	if r.getSourceLanguageFunc != nil {
		return r.getSourceLanguageFunc(source)
	}
	panic("getSourceLanguageFunc is not implemented")
}

type spyTranslator struct {
	translateTextFunc func(ctx context.Context, sourceLanguageCode string, targetLanguageCode string, text string) (translatedText string, err error)
}

func (r *spyTranslator) TranslateText(ctx context.Context, sourceLanguageCode string, targetLanguageCode string, text string) (translatedText string, err error) {
	if r.translateTextFunc != nil {
		return r.translateTextFunc(ctx, sourceLanguageCode, targetLanguageCode, text)
	}
	panic("translateTextFunc is not implemented")
}

func TestAppService_Clean(t *testing.T) {
	t.Run("Should return original description when item is in Japanese", func(t *testing.T) {
		// Arrange
		test_rss := GenerateJapaneseTestRss(t)
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := spyFeedRepository{
			getSourceLanguageFunc: func(source string) (sourceLanguageCode string, ok bool) {
				return "ja", true
			},
		}
		translator := spyTranslator{}

		// Act
		act_rss, err := app_service.Translate(ctx, &logger, &translator, &repo, test_rss)

		// Assert
		assert.NoError(t, err)

		assert.NotEmpty(t, act_rss.ID)
		assert.Equal(t, test_rss.ID, act_rss.ID)
		assert.Equal(t, "ダミーニュースのフィード", act_rss.Title)
		assert.Equal(t, "127.0.0.1:8080", act_rss.Source)
		assert.Equal(t, "http://127.0.0.1:8080", act_rss.Link)
		assert.Equal(t, "このフィードはダミーニュースを提供します。", act_rss.Description)
		assert.Equal(t, "ja", act_rss.Language)
		assert.Equal(t, time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC), act_rss.LastBuildDate)

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
	t.Run("Should return original description when translation setting is unset", func(t *testing.T) {
		// Arrange
		test_rss := GenerateJapaneseTestRss(t)
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := spyFeedRepository{
			getSourceLanguageFunc: func(source string) (sourceLanguageCode string, ok bool) {
				return "", false
			},
		}
		translator := spyTranslator{}

		// Act
		act_rss, err := app_service.Translate(ctx, &logger, &translator, &repo, test_rss)

		// Assert
		assert.NoError(t, err)

		assert.NotEmpty(t, act_rss.ID)
		assert.Equal(t, test_rss.ID, act_rss.ID)
		assert.Equal(t, "ダミーニュースのフィード", act_rss.Title)
		assert.Equal(t, "127.0.0.1:8080", act_rss.Source)
		assert.Equal(t, "http://127.0.0.1:8080", act_rss.Link)
		assert.Equal(t, "このフィードはダミーニュースを提供します。", act_rss.Description)
		assert.Equal(t, "ja", act_rss.Language)
		assert.Equal(t, time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC), act_rss.LastBuildDate)

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
	t.Run("Should translate and return Japanese description when original is in a foreign language", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Dummy News Feed", "127.0.0.1:8080", "http://127.0.0.1:8080", "This feed provides dummy news.", "en", time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC))
			if err != nil {
				return err
			}

			dummy_item1, err := rss.NewItem(rss.Guid{Value: "http://www.example.com/dummy-guid1"}, "Dummy Article 1", "http://www.example.com/dummy-article1", "Here is a summary of dummy article 1. Please click the link for more details.", "item1@dummy.com", time.Date(2024, time.July, 3, 12, 0, 0, 0, time.UTC))
			if err != nil {
				return err
			}
			test_rss.AddOrUpdateItem(dummy_item1)

			dummy_item2, err := rss.NewItem(rss.Guid{Value: "http://www.example.com/dummy-guid2"}, "Dummy Article 2", "http://www.example.com/dummy-article2", "Here is a summary of dummy article 2. Please click the link for more details.", "item2@dummy.com", time.Date(2024, time.July, 3, 12, 30, 0, 0, time.UTC))
			if err != nil {
				return err
			}
			test_rss.AddOrUpdateItem(dummy_item2)
			return err
		})

		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := spyFeedRepository{
			getSourceLanguageFunc: func(source string) (sourceLanguageCode string, ok bool) {
				return "en", true
			},
		}

		act_translateTextFunc_call_count := 0
		translator := spyTranslator{
			translateTextFunc: func(ctx context.Context, sourceLanguageCode string, targetLanguageCode string, text string) (translatedText string, err error) {
				act_translateTextFunc_call_count++
				if sourceLanguageCode != "en" {
					panic("sourceLanguageCode is not 'en' as expected")
				}
				if targetLanguageCode != "ja" {
					panic("targetLanguageCode is not 'ja' as expected")
				}

				if text == "Here is a summary of dummy article 1. Please click the link for more details." {
					return "これはダミー記事1の概要です。詳細はリンクをクリックしてください。", nil
				}

				if text == "Here is a summary of dummy article 2. Please click the link for more details." {
					return "これはダミー記事2の概要です。詳細はリンクをクリックしてください。", nil
				}

				return "", nil
			},
		}

		// Act
		act_rss, err := app_service.Translate(ctx, &logger, &translator, &repo, test_rss)

		// Assert
		assert.NoError(t, err)

		assert.Equal(t, test_rss.ID, act_rss.ID)
		assert.Equal(t, "Dummy News Feed", act_rss.Title)
		assert.Equal(t, "127.0.0.1:8080", act_rss.Source)
		assert.Equal(t, "http://127.0.0.1:8080", act_rss.Link)
		assert.Equal(t, "This feed provides dummy news.", act_rss.Description)
		assert.Equal(t, "en", act_rss.Language)
		assert.Equal(t, time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC), act_rss.LastBuildDate)

		assert.Len(t, act_rss.Items, 2)
		// item1
		{
			item1 := act_rss.Items[rss.Guid{Value: "http://www.example.com/dummy-guid1"}]
			assert.Equal(t, rss.Guid{Value: "http://www.example.com/dummy-guid1"}, item1.Guid)
			assert.Equal(t, "Dummy Article 1", item1.Title)
			assert.Equal(t, "http://www.example.com/dummy-article1", item1.Link)
			assert.Equal(t, "これはダミー記事1の概要です。詳細はリンクをクリックしてください。", item1.Description)
			assert.Equal(t, time.Date(2024, time.July, 3, 12, 0, 0, 0, time.UTC), item1.PubDate)
			assert.Equal(t, "item1@dummy.com", item1.Author)
		}
		// item2
		{
			item2 := act_rss.Items[rss.Guid{Value: "http://www.example.com/dummy-guid2"}]
			assert.Equal(t, rss.Guid{Value: "http://www.example.com/dummy-guid2"}, item2.Guid)
			assert.Equal(t, "Dummy Article 2", item2.Title)
			assert.Equal(t, "http://www.example.com/dummy-article2", item2.Link)
			assert.Equal(t, "これはダミー記事2の概要です。詳細はリンクをクリックしてください。", item2.Description)
			assert.Equal(t, time.Date(2024, time.July, 3, 12, 30, 0, 0, time.UTC), item2.PubDate)
			assert.Equal(t, "item2@dummy.com", item2.Author)
		}

		assert.Equal(t, act_translateTextFunc_call_count, 2)
	})
}

func GenerateJapaneseTestRss(t *testing.T) rss.Rss {
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
