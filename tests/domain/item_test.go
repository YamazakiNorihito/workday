package domain

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

func TestItem_NewItem(t *testing.T) {
	t.Run("should create new Item when all required fields are provided", func(t *testing.T) {
		// Arrange
		now := time.Now()

		// Act
		item, err := rss.NewItem(rss.Guid{Value: "guid-12345"}, "Test Title", "http://example.com", "Test Description", "Test Author", now)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "Test Title", item.Title)
		assert.Equal(t, "http://example.com", item.Link)
		assert.Equal(t, "Test Description", item.Description)
		assert.Equal(t, "Test Author", item.Author)
		assert.Equal(t, rss.Guid{Value: "guid-12345"}, item.Guid)
		assert.Equal(t, now, item.PubDate)
		assert.Empty(t, item.Tags)
	})

	t.Run("should return error when required fields are missing", func(t *testing.T) {
		// Arrange
		var tests = []struct {
			testName      string
			title         string
			link          string
			description   string
			author        string
			guid          rss.Guid
			pubDate       time.Time
			expectedError error
		}{
			{
				testName:      "should return error when title is empty",
				title:         "",
				link:          "http://example.com",
				description:   "Test Description",
				author:        "Test Author",
				guid:          rss.Guid{Value: "guid-12345"},
				pubDate:       time.Now(),
				expectedError: errors.New("title, link, and guid cannot be empty"),
			},
			{
				testName:      "should return error when link is empty",
				title:         "Test Title",
				link:          "",
				description:   "Test Description",
				author:        "Test Author",
				guid:          rss.Guid{Value: "guid-12345"},
				pubDate:       time.Now(),
				expectedError: errors.New("title, link, and guid cannot be empty"),
			},
			{
				testName:      "should return error when guid is empty",
				title:         "Test Title",
				link:          "http://example.com",
				description:   "Test Description",
				author:        "Test Author",
				guid:          rss.Guid{},
				pubDate:       time.Now(),
				expectedError: errors.New("title, link, and guid cannot be empty"),
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.testName, func(t *testing.T) {
				// Act
				_, err := rss.NewItem(tt.guid, tt.title, tt.link, tt.description, tt.author, tt.pubDate)

				// Assert
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			})
		}
	})

	t.Run("should create new Item when optional fields are missing", func(t *testing.T) {
		// Arrange
		var tests = []struct {
			testName    string
			title       string
			link        string
			description string
			author      string
			guid        rss.Guid
			pubDate     time.Time
		}{
			{
				testName:    "should create Item with empty author",
				title:       "Test Title",
				link:        "http://example.com",
				description: "Test Description",
				author:      "",
				guid:        rss.Guid{Value: "guid-12345"},
				pubDate:     time.Now(),
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.testName, func(t *testing.T) {
				// Act
				item, err := rss.NewItem(tt.guid, tt.title, tt.link, tt.description, tt.author, tt.pubDate)

				// Assert
				assert.NoError(t, err)
				assert.NotEmpty(t, item.Guid)
				assert.Equal(t, tt.title, item.Title)
				assert.Equal(t, tt.link, item.Link)
				assert.Equal(t, tt.description, item.Description)
				assert.Equal(t, tt.author, item.Author)
				assert.Equal(t, tt.guid, item.Guid)
				assert.Equal(t, tt.pubDate, item.PubDate)
				assert.Empty(t, item.Tags)
			})
		}
	})
}

func TestItem_AddTag(t *testing.T) {
	t.Run("should add a unique tag successfully", func(t *testing.T) {
		// Arrange
		now := time.Now()

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Test Title", "http://example.com", "Test description", "Test Author", now)
			return err
		})

		// Act
		test_item.AddTag("new-tag")

		// Assert
		assert.Equal(t, "Test Title", test_item.Title)
		assert.Equal(t, "http://example.com", test_item.Link)
		assert.Equal(t, "Test description", test_item.Description)
		assert.Equal(t, "Test Author", test_item.Author)
		assert.Equal(t, rss.Guid{Value: "guid-12345"}, test_item.Guid)
		assert.Equal(t, now, test_item.PubDate)
		assert.Equal(t, test_item.Tags, []string{"new-tag"})
		assert.Len(t, test_item.Tags, 1)
	})

	t.Run("should add multiple unique tags successfully", func(t *testing.T) {
		// Arrange
		now := time.Now()

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Test Title", "http://example.com", "Test description", "Test Author", now)
			return err
		})

		// Act
		test_item.AddTag("first-tag")
		test_item.AddTag("second-tag")

		// Assert
		assert.Equal(t, "Test Title", test_item.Title)
		assert.Equal(t, "http://example.com", test_item.Link)
		assert.Equal(t, "Test description", test_item.Description)
		assert.Equal(t, "Test Author", test_item.Author)
		assert.Equal(t, rss.Guid{Value: "guid-12345"}, test_item.Guid)
		assert.Equal(t, now, test_item.PubDate)
		assert.Equal(t, test_item.Tags, []string{"first-tag", "second-tag"})
		assert.Len(t, test_item.Tags, 2)
	})

	t.Run("should maintain unique tags when duplicate tag is added", func(t *testing.T) {
		// Arrange
		now := time.Now()

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Test Title", "http://example.com", "Test description", "Test Author", now)
			return err
		})

		// Act
		test_item.AddTag("duplicate-tag")
		test_item.AddTag("duplicate-tag")

		// Assert
		assert.Equal(t, "Test Title", test_item.Title)
		assert.Equal(t, "http://example.com", test_item.Link)
		assert.Equal(t, "Test description", test_item.Description)
		assert.Equal(t, "Test Author", test_item.Author)
		assert.Equal(t, rss.Guid{Value: "guid-12345"}, test_item.Guid)
		assert.Equal(t, now, test_item.PubDate)
		assert.Equal(t, test_item.Tags, []string{"duplicate-tag"})
		assert.Len(t, test_item.Tags, 1)
	})
}

func TestItem_Serialize(t *testing.T) {
	t.Run("should serialize to JSON correctly", func(t *testing.T) {
		// Arrange
		now := time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC)

		var test_item rss.Item
		helper.MustSucceed(t, func() error {
			var err error
			test_item, err = rss.NewItem(rss.Guid{Value: "guid-12345"}, "Test Title", "http://example.com", "Test description", "Test Author", now)
			return err
		})
		test_item.AddTag("tag1")
		test_item.AddTag("tag2")

		// Act
		jsonData, err := json.Marshal(test_item)
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}

		// Assert
		assert.Equal(t,
			`{"guid":"guid-12345","title":"Test Title","link":"http://example.com","description":"Test description","author":"Test Author","pubDate":"2024-06-01T13:30:00Z","tags":["tag1","tag2"]}`,
			string(jsonData))
	})
}

func TestItem_Deserialize(t *testing.T) {
	t.Run("should deserialize from JSON correctly", func(t *testing.T) {
		// Arrange
		jsonData := `{"guid":"guid-12345","title":"Test Title","link":"http://example.com","description":"Test description","author":"Test Author","pubDate":"2024-06-01T13:30:00Z","tags":["tag1","tag2"]}`

		// Act
		var test_item rss.Item
		err := json.Unmarshal([]byte(jsonData), &test_item)

		// Assert
		assert.NoError(t, err)

		expectedTime := time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC)
		expectedItem := rss.Item{
			Guid:        rss.Guid{Value: "guid-12345"},
			Title:       "Test Title",
			Link:        "http://example.com",
			Description: "Test description",
			Author:      "Test Author",
			PubDate:     expectedTime,
			Tags:        []string{"tag1", "tag2"},
		}

		assert.Equal(t, expectedItem, test_item)
	})
}
