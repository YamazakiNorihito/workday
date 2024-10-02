package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRss_New(t *testing.T) {
	t.Run("should create new Rss when all required fields are provided", func(t *testing.T) {
		// Arrange
		now := time.Now()
		title := "Test Title"
		source := "Test Source"
		link := "http://example.com"
		description := "Test Description"
		language := "en"
		lastBuildDate := now.Add(-10 * time.Minute)

		// Act
		rss, err := rss.New(title, source, link, description, language, lastBuildDate)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, rss.ID)
		assert.Equal(t, "Test Title", rss.Title)
		assert.Equal(t, "Test Source", rss.Source)
		assert.Equal(t, "http://example.com", rss.Link)
		assert.Equal(t, "Test Description", rss.Description)
		assert.Equal(t, "en", rss.Language)
		assert.Equal(t, now.Add(-10*time.Minute), rss.LastBuildDate)
		assert.Empty(t, rss.Items)
	})

	t.Run("should return error when required fields are missing", func(t *testing.T) {
		// Arrange
		var tests = []struct {
			testName      string
			title         string
			source        string
			link          string
			description   string
			language      string
			lastBuildDate time.Time
			expectedError error
		}{
			{
				testName:      "should return error when title is empty",
				title:         "",
				source:        "Test Source",
				link:          "http://example.com",
				description:   "Test Description",
				language:      "en",
				lastBuildDate: time.Now(),
			},
			{
				testName:      "should return error when source is empty",
				title:         "Test Title",
				source:        "",
				link:          "http://example.com",
				description:   "Test Description",
				language:      "en",
				lastBuildDate: time.Now(),
			},
			{
				testName:      "should return error when link is empty",
				title:         "Test Title",
				source:        "Test Source",
				link:          "",
				description:   "Test Description",
				language:      "en",
				lastBuildDate: time.Now(),
			},
			{
				testName:      "should return error when lastBuildDate is empty",
				title:         "Test Title",
				source:        "Test Source",
				link:          "http://example.com",
				description:   "Test Description",
				language:      "en",
				lastBuildDate: time.Time{},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.testName, func(t *testing.T) {
				// Act
				_, err := rss.New(tt.title, tt.source, tt.link, tt.description, tt.language, tt.lastBuildDate)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestRss_SetLastBuildDate(t *testing.T) {
	t.Run("should set LastBuildDate when valid date is provided", func(t *testing.T) {
		// Arrange
		oldDate := time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC)
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", oldDate)
			return err
		})

		newDate := time.Date(2024, time.June, 3, 14, 40, 1, 1, time.UTC)

		// Act
		err := test_rss.SetLastBuildDate(newDate)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2024, time.June, 3, 14, 40, 1, 1, time.UTC), test_rss.LastBuildDate)
	})

	t.Run("should return error when LastBuildDate is zero", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		err := test_rss.SetLastBuildDate(time.Time{})

		// Assert
		assert.Error(t, err)
	})
}

func TestRss_SetLanguage(t *testing.T) {
	t.Run("should override the current language when a new valid language code is provided", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "ja", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		err := test_rss.SetLanguage("en")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "en", test_rss.Language)
	})

	t.Run("should clear the language code when provided with an empty string", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		err := test_rss.SetLanguage("")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "", test_rss.Language)
	})
}

func TestRss_AddOrUpdateItem(t *testing.T) {
	t.Run("should add new item when item does not exist", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Test Title", "http://example.com", "Test description", "Test Author", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(test_item)

		// Assert
		assert.Len(t, test_rss.Items, 1)

		item, exists := test_rss.Items[rss.Guid{Value: "guid-12345"}]
		if !exists {
			t.Fatalf("expected item with guid 'guid-123' to be present")
		}
		assert.Equal(t, rss.Guid{Value: "guid-12345"}, item.Guid)
		assert.Equal(t, "Test Title", item.Title)
		assert.Equal(t, "http://example.com", item.Link)
		assert.Equal(t, "Test description", item.Description)
		assert.Equal(t, "Test Author", item.Author)
		assert.Equal(t, time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC), item.PubDate)
	})

	t.Run("should add new items when items do not exist", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		var test_item1, test_item2 rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item1, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Test Title 1", "http://example.com/1", "Test description 1", "Test Author 1", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC))
			if err != nil {
				return err
			}
			test_item2, err = rss.NewItem(rss.Guid{Value: "guid-67890"}, "Test Title 2", "http://example.com/2", "Test description 2", "Test Author 2", time.Date(2023, time.June, 2, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(test_item1)
		test_rss.AddOrUpdateItem(test_item2)

		// Assert
		assert.Len(t, test_rss.Items, 2)

		item1, exists1 := test_rss.Items[rss.Guid{Value: "guid-12345"}]
		if !exists1 {
			t.Fatalf("expected item with guid 'guid-12345' to be present")
		}
		assert.Equal(t, rss.Guid{Value: "guid-12345"}, item1.Guid)
		assert.Equal(t, "Test Title 1", item1.Title)
		assert.Equal(t, "http://example.com/1", item1.Link)
		assert.Equal(t, "Test description 1", item1.Description)
		assert.Equal(t, "Test Author 1", item1.Author)
		assert.Equal(t, time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC), item1.PubDate)

		item2, exists2 := test_rss.Items[rss.Guid{Value: "guid-67890"}]
		if !exists2 {
			t.Fatalf("expected item with guid 'guid-67890' to be present")
		}
		assert.Equal(t, rss.Guid{Value: "guid-67890"}, item2.Guid)
		assert.Equal(t, "Test Title 2", item2.Title)
		assert.Equal(t, "http://example.com/2", item2.Link)
		assert.Equal(t, "Test description 2", item2.Description)
		assert.Equal(t, "Test Author 2", item2.Author)
		assert.Equal(t, time.Date(2023, time.June, 2, 13, 30, 0, 0, time.UTC), item2.PubDate)
	})

	t.Run("should update existing item when item already exists", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		var original_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			original_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Original Title", "http://example.com/original", "Original description", "Original Author", time.Date(2023, time.January, 1, 13, 30, 0, 0, time.UTC))
			return err
		})
		test_rss.AddOrUpdateItem(original_item)

		var updated_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			updated_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Updated Title", "http://example.com/updated", "Updated description", "Updated Author", time.Date(2024, time.January, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(updated_item)

		// Assert
		assert.Len(t, test_rss.Items, 1)

		item, exists := test_rss.Items[rss.Guid{Value: "guid-12345"}]
		if !exists {
			t.Fatalf("expected item with guid 'guid-123' to be present")
		}
		assert.Equal(t, rss.Guid{Value: "guid-12345"}, item.Guid)
		assert.Equal(t, "Updated Title", item.Title)
		assert.Equal(t, "http://example.com/updated", item.Link)
		assert.Equal(t, "Updated description", item.Description)
		assert.Equal(t, "Updated Author", item.Author)
		assert.Equal(t, time.Date(2024, time.January, 1, 13, 30, 0, 0, time.UTC), item.PubDate)
	})

	t.Run("should add item when item matches include keywords", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})
		test_rss.SetItemFilter([]string{"include_keyword"}, nil)

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "This title contains include_keyword", "http://example.com", "Test description", "Test Author", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(test_item)

		// Assert
		assert.Len(t, test_rss.Items, 1)
		item, exists := test_rss.Items[rss.Guid{Value: "guid-12345"}]
		if !exists {
			t.Fatalf("expected item with guid 'guid-12345' to be present")
		}
		assert.Equal(t, "This title contains include_keyword", item.Title)
	})

	t.Run("should not add item when item does not match include keywords", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})
		test_rss.SetItemFilter([]string{"include_keyword"}, nil)

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "No matching keyword in title", "http://example.com", "Test description", "Test Author", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(test_item)

		// Assert
		assert.Len(t, test_rss.Items, 0)
	})

	t.Run("should not add item when item matches exclude keywords", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})
		test_rss.SetItemFilter(nil, []string{"exclude_keyword"})

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "This title contains exclude_keyword", "http://example.com", "Test description", "Test Author", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(test_item)

		// Assert
		assert.Len(t, test_rss.Items, 0)
	})

	t.Run("should add item when ItemFilter is empty", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})
		// No filters set

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Test Title", "http://example.com", "Test description", "Test Author", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(test_item)

		// Assert
		assert.Len(t, test_rss.Items, 1)
	})

	t.Run("should add item when item matches include and does not match exclude keywords", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})
		test_rss.SetItemFilter([]string{"include_keyword"}, []string{"exclude_keyword"})

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "This title contains include_keyword", "http://example.com", "Test description without exclude_keyword", "Test Author", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(test_item)

		// Assert
		assert.Len(t, test_rss.Items, 0)
	})

	t.Run("should not add item when item matches include and exclude keywords", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})
		test_rss.SetItemFilter([]string{"include_keyword"}, []string{"exclude_keyword"})

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "This title contains include_keyword and exclude_keyword", "http://example.com", "Test description", "Test Author", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC))
			return err
		})

		// Act
		test_rss.AddOrUpdateItem(test_item)

		// Assert
		assert.Len(t, test_rss.Items, 0)
	})
}

func TestRss_Serialize(t *testing.T) {
	t.Run("should serialize to JSON correctly", func(t *testing.T) {
		// Arrange
		var test_rss rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			test_rss, err = rss.New("Test Title", "Test Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))

			includeKeywords := []string{"go", "golang"}
			excludeKeywords := []string{"python", "ruby"}
			test_rss.SetItemFilter(includeKeywords, excludeKeywords)
			return err
		})

		var original_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			original_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Original Title-go", "http://example.com/original", "Original description", "Original Author", time.Date(2023, time.January, 1, 13, 30, 0, 0, time.UTC))
			if err != nil {
				return err
			}
			original_item.AddTag("tag1")
			original_item.AddTag("tag2")
			return nil
		})
		test_rss.AddOrUpdateItem(original_item)

		// Act
		jsonData, err := json.Marshal(test_rss)
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}

		// Assert
		expected := `{
			"id":"` + test_rss.ID.String() + `",
			"source":"Test Source",
			"title":"Test Title",
			"link":"http://example.com",
			"description":"Test Description",
			"language":"en",
			"last_build_date":"2024-06-01T13:30:00Z",
			"items":{
				"guid-12345":{
					"guid":"guid-12345",
					"title":"Original Title-go",
					"link":"http://example.com/original",
					"description":"Original description",
					"author":"Original Author",
					"pubDate":"2023-01-01T13:30:00Z",
					"tags":["tag1","tag2"]
				}
			},
			"item_filter":{
				"include_keywords":["go","golang"],
				"exclude_keywords":["python","ruby"]
			},
			"create_by":{"id":"","name":""},
			"create_at":"0001-01-01T00:00:00Z",
			"update_by":{"id":"","name":""},
			"update_at":"0001-01-01T00:00:00Z"
		}`
		assert.JSONEq(t, expected, string(jsonData))
	})
}

func TestRss_Deserialize(t *testing.T) {
	t.Run("should deserialize from JSON correctly", func(t *testing.T) {
		// Arrange
		jsonData := `{
			"id":"8091f704-e1c5-41a5-a305-52d4c7876d5f",
			"source":"Test Source",
			"title":"Test Title",
			"link":"http://example.com",
			"description":"Test Description",
			"language":"en",
			"last_build_date":"2020-06-01T13:30:00Z",
			"items":{
				"guid-12345":{
					"guid":"guid-12345",
					"title":"Original Title",
					"link":"http://example.com/original",
					"description":"Original description",
					"author":"Original Author",
					"pubDate":"2021-06-01T13:30:00Z",
					"tags":["tag1","tag2"]
				}
			},
			"create_by":{"id":"test-create-id","name":"test-create-name"},
			"create_at":"2022-06-01T13:30:00Z",
			"update_by":{"id":"test-update-id","name":"test-update-name"},
			"update_at":"2023-06-01T13:30:00Z"
		}`

		// Act
		var test_rss rss.Rss
		err := json.Unmarshal([]byte(jsonData), &test_rss)

		// Assert
		assert.NoError(t, err)

		lastBuildDate, _ := time.Parse(time.RFC3339, "2020-06-01T13:30:00Z")
		pubDate, _ := time.Parse(time.RFC3339, "2021-06-01T13:30:00Z")
		createdAt, _ := time.Parse(time.RFC3339, "2022-06-01T13:30:00Z")
		updatedAt, _ := time.Parse(time.RFC3339, "2023-06-01T13:30:00Z")
		id, _ := uuid.Parse("8091f704-e1c5-41a5-a305-52d4c7876d5f")

		expectedRss := rss.Rss{
			ID:            id,
			Source:        "Test Source",
			Title:         "Test Title",
			Link:          "http://example.com",
			Description:   "Test Description",
			Language:      "en",
			LastBuildDate: lastBuildDate,
			Items: map[rss.Guid]rss.Item{
				rss.Guid{Value: "guid-12345"}: {
					Guid:        rss.Guid{Value: "guid-12345"},
					Title:       "Original Title",
					Link:        "http://example.com/original",
					Description: "Original description",
					Author:      "Original Author",
					PubDate:     pubDate,
					Tags:        []string{"tag1", "tag2"},
				},
			},
			CreatedBy: metadata.CreateBy{
				ID:   "test-create-id",
				Name: "test-create-name",
			},
			CreatedAt: createdAt,
			UpdatedBy: metadata.UpdateBy{
				ID:   "test-update-id",
				Name: "test-update-name",
			},
			UpdatedAt: updatedAt,
		}

		assert.Equal(t, expectedRss, test_rss)
	})
}
