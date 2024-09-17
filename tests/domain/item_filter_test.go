package domain

import (
	"encoding/json"
	"testing"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/stretchr/testify/assert"
)

func TestItemFilter_NewItemFilter(t *testing.T) {
	t.Run("should create new ItemFilter with provided keywords", func(t *testing.T) {
		// Arrange
		includeKeywords := []string{"go", "golang"}
		excludeKeywords := []string{"python", "ruby"}

		// Act
		filter := rss.NewItemFilter(includeKeywords, excludeKeywords)

		// Assert
		assert.NotNil(t, filter)
		assert.Equal(t, []string{"go", "golang"}, filter.IncludeKeywords)
		assert.Equal(t, []string{"python", "ruby"}, filter.ExcludeKeywords)
	})

	t.Run("should create new ItemFilter with empty slices when nil slices are provided", func(t *testing.T) {
		// Act
		filter := rss.NewItemFilter(nil, nil)

		// Assert
		assert.NotNil(t, filter)
		assert.Empty(t, filter.IncludeKeywords)
		assert.Empty(t, filter.ExcludeKeywords)
	})
}

func TestItemFilter_GetIncludeKeywords(t *testing.T) {
	t.Run("should return the correct include keywords", func(t *testing.T) {
		// Arrange
		includeKeywords := []string{"go", "golang"}
		filter := rss.NewItemFilter(includeKeywords, nil)

		// Act
		result := filter.GetIncludeKeywords()

		// Assert
		assert.Equal(t, []string{"go", "golang"}, result)
	})
}

func TestItemFilter_GetExcludeKeywords(t *testing.T) {
	t.Run("should return the correct exclude keywords", func(t *testing.T) {
		// Arrange
		excludeKeywords := []string{"python", "ruby"}
		filter := rss.NewItemFilter(nil, excludeKeywords)

		// Act
		result := filter.GetExcludeKeywords()

		// Assert
		assert.Equal(t, []string{"python", "ruby"}, result)
	})
}

func TestItemFilter_Serialize(t *testing.T) {
	t.Run("should serialize to JSON correctly", func(t *testing.T) {
		// Arrange
		includeKeywords := []string{"go", "golang"}
		excludeKeywords := []string{"python", "ruby"}
		filter := rss.NewItemFilter(includeKeywords, excludeKeywords)

		// Act
		jsonData, err := json.Marshal(filter)

		// Assert
		assert.NoError(t, err)
		expectedJSON := `{"include_keywords":["go","golang"],"exclude_keywords":["python","ruby"]}`
		assert.JSONEq(t, expectedJSON, string(jsonData))
	})

	t.Run("should serialize to JSON correctly when slices are empty", func(t *testing.T) {
		// Arrange
		filter := rss.NewItemFilter(nil, nil)

		// Act
		jsonData, err := json.Marshal(filter)

		// Assert
		assert.NoError(t, err)
		expectedJSON := `{"include_keywords":[],"exclude_keywords":[]}`
		assert.JSONEq(t, expectedJSON, string(jsonData))
	})
}
