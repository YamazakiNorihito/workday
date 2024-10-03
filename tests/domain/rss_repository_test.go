package domain

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/YamazakiNorihito/workday/tests/helper/assert_helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func schemaProvider() ([]types.AttributeDefinition, []types.KeySchemaElement, []types.GlobalSecondaryIndex) {
	attributeDefinitions := []types.AttributeDefinition{
		{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeS,
		},
		{
			AttributeName: aws.String("sortKey"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}

	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		},
		{
			AttributeName: aws.String("sortKey"),
			KeyType:       types.KeyTypeRange,
		},
	}

	gsi := []types.GlobalSecondaryIndex{
		{
			IndexName: aws.String("SortKeyIndex"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("sortKey"),
					KeyType:       types.KeyTypeHash,
				},
			},
			Projection: &types.Projection{
				ProjectionType: types.ProjectionTypeAll,
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
		},
	}

	return attributeDefinitions, keySchema, gsi
}

func TestRssRepository_Save(t *testing.T) {
	t.Run("should save new rss when CreatedBy is empty", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		test_rss := getTestRss(t)

		// Act
		_, err := rssRepository.Save(ctx, test_rss, metadata.UserMeta{ID: "test-id", Name: "test-user"})

		// Assert
		assert.NoError(t, err)

		// rss
		tableName := "Rss"
		partitionkey := "Test_Source"
		sortKey := "rss"
		actual_rss, _ := helper.GetItem(ctx, client, tableName, partitionkey, sortKey)

		assert.NotEmpty(t, actual_rss)
		assert.Contains(t, actual_rss, "rss_id")
		assert.Contains(t, actual_rss, "source")
		assert.Contains(t, actual_rss, "title")
		assert.Contains(t, actual_rss, "link")
		assert.Contains(t, actual_rss, "description")
		assert.Contains(t, actual_rss, "language")
		assert.Contains(t, actual_rss, "last_build_date")
		assert.Contains(t, actual_rss, "item_filter")
		assert.Contains(t, actual_rss["item_filter"], "include_keywords")
		assert.Contains(t, actual_rss["item_filter"], "exclude_keywords")
		assert.Contains(t, actual_rss, "create_by")
		assert.Contains(t, actual_rss, "create_at")
		assert.Contains(t, actual_rss, "update_by")
		assert.Contains(t, actual_rss, "update_at")

		assert_helper.EqualUUID(t, test_rss.ID, actual_rss["rss_id"])
		assert.Equal(t, "Test_Source", actual_rss["source"])
		assert.Equal(t, "Test Title", actual_rss["title"])
		assert.Equal(t, "http://example.com", actual_rss["link"])
		assert.Equal(t, "Test Description", actual_rss["description"])
		assert.Equal(t, "en", actual_rss["language"])
		assert_helper.EqualUnixTime(t, time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC), actual_rss["last_build_date"])

		itemFilter := actual_rss["item_filter"].(map[string]interface{})
		assert.Equal(t, "include_keyword", itemFilter["include_keywords"].([]interface{})[0])
		assert.Equal(t, "exclude_keyword", itemFilter["exclude_keywords"].([]interface{})[0])
		assert_helper.EqualUserMeta(t, metadata.UserMeta{ID: "test-id", Name: "test-user"}, actual_rss["create_by"])
		assert.NotEmpty(t, actual_rss["create_at"])
		assert_helper.EqualUserMeta(t, metadata.UserMeta{ID: "test-id", Name: "test-user"}, actual_rss["update_by"])
		assert.NotEmpty(t, actual_rss["update_at"])

		// item1
		item1_sortKey := test_rss.ID.String() + "#" + "guid-12345"
		actual_item1, _ := helper.GetItem(ctx, client, tableName, partitionkey, item1_sortKey)
		assert.NotEmpty(t, actual_item1)
		assert.Equal(t, "Test Title 1", actual_item1["title"])
		assert.Equal(t, "http://example.com/1", actual_item1["link"])
		assert.Equal(t, "Test description 1-include_keyword", actual_item1["description"])
		assert.Equal(t, "Test Author 1", actual_item1["author"])
		assert_helper.EqualUnixTime(t, time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC), actual_item1["pub_date"])

		// item2
		item2_sortKey := test_rss.ID.String() + "#" + "guid-67890"
		actual_item2, _ := helper.GetItem(ctx, client, tableName, partitionkey, item2_sortKey)
		assert.NotEmpty(t, actual_item2)
		assert.Equal(t, "Test Title 2", actual_item2["title"])
		assert.Equal(t, "http://example.com/2", actual_item2["link"])
		assert.Equal(t, "Test description 2-include_keyword", actual_item2["description"])
		assert.Equal(t, "Test Author 2", actual_item2["author"])
		assert_helper.EqualUnixTime(t, time.Date(2023, time.June, 2, 13, 30, 0, 0, time.UTC), actual_item2["pub_date"])
	})

	t.Run("should update existing rss when RSS with same ID exists", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		initial_rss := getTestRss(t)

		// Save initial RSS
		test_rss, err := rssRepository.Save(ctx, initial_rss, metadata.UserMeta{ID: "test-id", Name: "test-user"})
		assert.NoError(t, err)

		// Modify RSS for update
		helper.MustSucceed(t, func() error {
			err := test_rss.SetLastBuildDate(time.Date(2024, time.August, 17, 11, 11, 11, 11, time.UTC))
			if err != nil {
				return err
			}
			rssItem, err := rss.NewItem(rss.Guid{Value: "guid-11111"}, "Test Title 3", "http://example.com/3", "Test description 3-include_keyword", "Test Author 3", time.Date(2023, time.June, 3, 13, 30, 0, 0, time.UTC))
			if err != nil {
				return err
			}
			test_rss.AddOrUpdateItem(rssItem)
			return nil
		})

		// Act
		_, err = rssRepository.Save(ctx, test_rss, metadata.UserMeta{ID: "update-id", Name: "update-user"})
		assert.NoError(t, err)

		// Assert
		tableName := "Rss"
		partitionkey := "Test_Source"
		sortKey := "rss"
		actual_rss, _ := helper.GetItem(ctx, client, tableName, partitionkey, sortKey)

		// rss
		assert.NotEmpty(t, actual_rss)

		assert_helper.EqualUUID(t, test_rss.ID, actual_rss["rss_id"])
		assert.Equal(t, "Test_Source", actual_rss["source"])
		assert.Equal(t, "Test Title", actual_rss["title"])
		assert.Equal(t, "http://example.com", actual_rss["link"])
		assert.Equal(t, "Test Description", actual_rss["description"])
		assert.Equal(t, "en", actual_rss["language"])
		assert_helper.EqualUnixTime(t, time.Date(2024, time.August, 17, 11, 11, 11, 11, time.UTC), actual_rss["last_build_date"])

		itemFilter := actual_rss["item_filter"].(map[string]interface{})
		assert.Equal(t, "include_keyword", itemFilter["include_keywords"].([]interface{})[0])
		assert.Equal(t, "exclude_keyword", itemFilter["exclude_keywords"].([]interface{})[0])
		assert_helper.EqualUserMeta(t, metadata.UserMeta{ID: "test-id", Name: "test-user"}, actual_rss["create_by"])
		assert.NotEmpty(t, actual_rss["create_at"])
		assert_helper.EqualUserMeta(t, metadata.UserMeta{ID: "update-id", Name: "update-user"}, actual_rss["update_by"])
		assert.NotEmpty(t, actual_rss["update_at"])

		// item1
		item1_sortKey := test_rss.ID.String() + "#" + "guid-12345"
		actual_item1, _ := helper.GetItem(ctx, client, tableName, partitionkey, item1_sortKey)
		assert.NotEmpty(t, actual_item1)
		assert.Equal(t, "Test Title 1", actual_item1["title"])
		assert.Equal(t, "http://example.com/1", actual_item1["link"])
		assert.Equal(t, "Test description 1-include_keyword", actual_item1["description"])
		assert.Equal(t, "Test Author 1", actual_item1["author"])
		assert_helper.EqualUnixTime(t, time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC), actual_item1["pub_date"])

		// item2
		item2_sortKey := test_rss.ID.String() + "#" + "guid-67890"
		actual_item2, _ := helper.GetItem(ctx, client, tableName, partitionkey, item2_sortKey)
		assert.NotEmpty(t, actual_item2)
		assert.Equal(t, "Test Title 2", actual_item2["title"])
		assert.Equal(t, "http://example.com/2", actual_item2["link"])
		assert.Equal(t, "Test description 2-include_keyword", actual_item2["description"])
		assert.Equal(t, "Test Author 2", actual_item2["author"])
		assert_helper.EqualUnixTime(t, time.Date(2023, time.June, 2, 13, 30, 0, 0, time.UTC), actual_item2["pub_date"])

		// item3
		item3_sortKey := test_rss.ID.String() + "#" + "guid-11111"
		actual_item3, _ := helper.GetItem(ctx, client, tableName, partitionkey, item3_sortKey)
		assert.NotEmpty(t, actual_item3)
		assert.Equal(t, "Test Title 3", actual_item3["title"])
		assert.Equal(t, "http://example.com/3", actual_item3["link"])
		assert.Equal(t, "Test description 3-include_keyword", actual_item3["description"])
		assert.Equal(t, "Test Author 3", actual_item3["author"])
		assert_helper.EqualUnixTime(t, time.Date(2023, time.June, 3, 13, 30, 0, 0, time.UTC), actual_item3["pub_date"])
	})
}

func TestRssRepository_FindBySource(t *testing.T) {
	t.Run("should return error when source is empty", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss, err := rssRepository.FindBySource(ctx, "")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, rss.Rss{}, actual_rss)
	})

	t.Run("should return Rss when source exists", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss, err := rssRepository.FindBySource(ctx, "Test_Source")

		// Assert
		assert.NoError(t, err)

		assert.NotEmpty(t, actual_rss.ID)
		assert.Equal(t, "Test_Source", actual_rss.Source)
		assert.Equal(t, "Test Title", actual_rss.Title)
		assert.Equal(t, "http://example.com", actual_rss.Link)
		assert.Equal(t, "Test Description", actual_rss.Description)
		assert.Equal(t, "en", actual_rss.Language)
		assert.Equal(t, time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC), actual_rss.LastBuildDate.UTC())
		assert.Equal(t, rss.ItemFilter{
			IncludeKeywords: []string{"Test"},
			ExcludeKeywords: []string{"unit"},
		}, actual_rss.ItemFilter)
		assert.Equal(t, metadata.CreateBy{ID: "test-id", Name: "test-user"}, actual_rss.CreatedBy)
		assert.NotEmpty(t, actual_rss.CreatedAt)
		assert.Equal(t, metadata.UpdateBy{ID: "test-id", Name: "test-user"}, actual_rss.UpdatedBy)
		assert.NotEmpty(t, actual_rss.UpdatedAt)

		// item
		assert.Len(t, actual_rss.Items, 0)
	})

	t.Run("should return empty Rss when source does not exist", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss, err := rssRepository.FindBySource(ctx, "non-existent_Source")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, actual_rss, rss.Rss{Items: make(map[rss.Guid]rss.Item)})
		assert.Equal(t, actual_rss.ID, uuid.Nil)
		assert.Len(t, actual_rss.Items, 0)
	})
}

func TestRssRepository_FindAll(t *testing.T) {
	t.Run("should return error when source is empty", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)

		// Act
		actual_rss_list, err := rssRepository.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, actual_rss_list, 0)
	})
	t.Run("should return error when source is empty1", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss_list, err := rssRepository.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, actual_rss_list, 1)

		actual_rss := actual_rss_list[0]
		assert.Equal(t, "Test_Source", actual_rss.Source)
		assert.Equal(t, "Test Title", actual_rss.Title)
		assert.Equal(t, "http://example.com", actual_rss.Link)
		assert.Equal(t, "Test Description", actual_rss.Description)
		assert.Equal(t, "en", actual_rss.Language)
		assert.Equal(t, time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC), actual_rss.LastBuildDate.UTC())
		assert.Equal(t, rss.ItemFilter{
			IncludeKeywords: []string{"Test"},
			ExcludeKeywords: []string{"unit"},
		}, actual_rss.ItemFilter)
		assert.Equal(t, metadata.CreateBy{ID: "test-id", Name: "test-user"}, actual_rss.CreatedBy)
		assert.NotEmpty(t, actual_rss.CreatedAt)
		assert.Equal(t, metadata.UpdateBy{ID: "test-id", Name: "test-user"}, actual_rss.UpdatedBy)
		assert.NotEmpty(t, actual_rss.UpdatedAt)
	})

	t.Run("should return Rss when source exists", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)

		var testRSS1 rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			testRSS1, err = rss.New("Test Title 1", "Test_Source_1", "http://example.com/1", "Test Description 1", "ja", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
			if err != nil {
				return err
			}

			testRSS1.SetItemFilter([]string{"Test"}, []string{"unit"})
			items := []struct {
				guid, title, link, description, author string
				pubDate                                time.Time
			}{
				{"guid-12345", "Test Title Item 1", "http://example.com/item1", "Test description item 1", "Test Author 1", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC)},
				{"guid-67890", "Test Title Item 2", "http://example.com/item2", "Test description item 2", "Test Author 2", time.Date(2023, time.June, 2, 13, 30, 0, 0, time.UTC)},
			}

			for _, item := range items {
				rssItem, err := rss.NewItem(rss.Guid{Value: item.guid}, item.title, item.link, item.description, item.author, item.pubDate)
				if err != nil {
					return err
				}
				testRSS1.AddOrUpdateItem(rssItem)
			}

			_, err = rssRepository.Save(ctx, testRSS1, metadata.UserMeta{ID: "test-id-1", Name: "test-user-1"})
			if err != nil {
				return err
			}
			return nil
		})

		var testRSS2 rss.Rss
		helper.MustSucceed(t, func() error {
			var err error
			testRSS2, err = rss.New("Test Title 2", "Test_Source_2", "http://example.com/2", "Test Description 2", "en", time.Date(2024, time.June, 2, 13, 30, 0, 0, time.UTC))
			if err != nil {
				return err
			}

			testRSS2.SetItemFilter([]string{"Test"}, []string{"unit"})
			items := []struct {
				guid, title, link, description, author string
				pubDate                                time.Time
			}{
				{"guid-54321", "Test Title Item 3", "http://example.com/item3", "Test description item 3", "Test Author 3", time.Date(2023, time.June, 3, 13, 30, 0, 0, time.UTC)},
				{"guid-09876", "Test Title Item 4", "http://example.com/item4", "Test description item 4", "Test Author 4", time.Date(2023, time.June, 4, 13, 30, 0, 0, time.UTC)},
			}

			for _, item := range items {
				rssItem, err := rss.NewItem(rss.Guid{Value: item.guid}, item.title, item.link, item.description, item.author, item.pubDate)
				if err != nil {
					return err
				}
				testRSS2.AddOrUpdateItem(rssItem)
			}

			_, err = rssRepository.Save(ctx, testRSS2, metadata.UserMeta{ID: "test-id-2", Name: "test-user-2"})
			if err != nil {
				return err
			}
			return nil
		})

		// Act
		actual_rss_list, err := rssRepository.FindAll(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, actual_rss_list, 2)

		expectedRssFeeds := map[string]rss.Rss{
			"Test_Source_1": {
				Title:         "Test Title 1",
				Source:        "Test_Source_1",
				Link:          "http://example.com/1",
				Description:   "Test Description 1",
				Language:      "ja",
				LastBuildDate: time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC),
				ItemFilter:    rss.ItemFilter{IncludeKeywords: []string{"Test"}, ExcludeKeywords: []string{"unit"}},
				CreatedBy:     metadata.CreateBy{ID: "test-id-1", Name: "test-user-1"},
				UpdatedBy:     metadata.UpdateBy{ID: "test-id-1", Name: "test-user-1"},
			},
			"Test_Source_2": {
				Title:         "Test Title 2",
				Source:        "Test_Source_2",
				Link:          "http://example.com/2",
				Description:   "Test Description 2",
				Language:      "en",
				LastBuildDate: time.Date(2024, time.June, 2, 13, 30, 0, 0, time.UTC),
				ItemFilter:    rss.ItemFilter{IncludeKeywords: []string{"Test"}, ExcludeKeywords: []string{"unit"}},
				CreatedBy:     metadata.CreateBy{ID: "test-id-2", Name: "test-user-2"},
				UpdatedBy:     metadata.UpdateBy{ID: "test-id-2", Name: "test-user-2"},
			},
		}

		for _, actualRss := range actual_rss_list {
			expectedRss, exists := expectedRssFeeds[actualRss.Source]
			assert.True(t, exists)
			assert.Equal(t, expectedRss.Title, actualRss.Title)
			assert.Equal(t, expectedRss.Link, actualRss.Link)
			assert.Equal(t, expectedRss.Description, actualRss.Description)
			assert.Equal(t, expectedRss.Language, actualRss.Language)
			assert.Equal(t, expectedRss.LastBuildDate.UTC(), actualRss.LastBuildDate.UTC())
			assert.Equal(t, expectedRss.ItemFilter, actualRss.ItemFilter)
			assert.Equal(t, expectedRss.CreatedBy, actualRss.CreatedBy)
			assert.NotEmpty(t, actualRss.CreatedAt)
			assert.Equal(t, expectedRss.UpdatedBy, actualRss.UpdatedBy)
			assert.NotEmpty(t, actualRss.UpdatedAt)
		}
	})
}

func TestRssRepository_FindItems(t *testing.T) {
	t.Run("should return error when rss is empty", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss, err := rssRepository.FindItems(ctx, rss.Rss{})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, rss.Rss{}, actual_rss)
		assert.Len(t, actual_rss.Items, 0)
	})

	t.Run("should return Items when rss exists", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setUpRss := setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss, err := rssRepository.FindItems(ctx, setUpRss)

		// Assert
		assert.NoError(t, err)

		// item1
		actual_item1, _ := actual_rss.Items[rss.Guid{Value: "guid-12345"}]
		assert.NotEmpty(t, actual_item1)
		assert.Equal(t, "Test Title 1", actual_item1.Title)
		assert.Equal(t, "http://example.com/1", actual_item1.Link)
		assert.Equal(t, "Test description 1", actual_item1.Description)
		assert.Equal(t, "Test Author 1", actual_item1.Author)
		assert.Equal(t, time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC), actual_item1.PubDate.UTC())

		// item2
		actual_item2, _ := actual_rss.Items[rss.Guid{Value: "guid-67890"}]
		assert.NotEmpty(t, actual_item2)
		assert.Equal(t, "Test Title 2", actual_item2.Title)
		assert.Equal(t, "http://example.com/2", actual_item2.Link)
		assert.Equal(t, "Test description 2", actual_item2.Description)
		assert.Equal(t, "Test Author 2", actual_item2.Author)
		assert.Equal(t, time.Date(2023, time.June, 2, 13, 30, 0, 0, time.UTC), actual_item2.PubDate.UTC())
	})
}

func TestRssRepository_FindItem(t *testing.T) {
	t.Run("should return error when rss is empty", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss, err := rssRepository.FindItemsByPk(ctx, rss.Rss{}, rss.Guid{})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, rss.Rss{}, actual_rss)
		assert.Len(t, actual_rss.Items, 0)
	})

	t.Run("should return Items when rss exists", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setUpRss := setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss, err := rssRepository.FindItemsByPk(ctx, setUpRss, rss.Guid{Value: "guid-12345"})

		// Assert
		assert.NoError(t, err)

		// item1
		actual_item1, _ := actual_rss.Items[rss.Guid{Value: "guid-12345"}]
		assert.NotEmpty(t, actual_item1)
		assert.Equal(t, "Test Title 1", actual_item1.Title)
		assert.Equal(t, "http://example.com/1", actual_item1.Link)
		assert.Equal(t, "Test description 1", actual_item1.Description)
		assert.Equal(t, "Test Author 1", actual_item1.Author)
		assert.Equal(t, time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC), actual_item1.PubDate.UTC())
	})

	t.Run("should return empty Items when GUID does not exist", func(t *testing.T) {
		// Arrange
		ctx, client := setUp()
		rssRepository := rss.NewDynamoDBRssRepository(client)
		setUpRss := setupExpectedRss(t, ctx, rssRepository)

		// Act
		actual_rss, err := rssRepository.FindItemsByPk(ctx, setUpRss, rss.Guid{Value: "non-existent-guid"})

		// Assert
		assert.NoError(t, err)
		actual_item, exists := actual_rss.Items[rss.Guid{Value: "non-existent-guid"}]
		assert.False(t, exists)
		assert.Empty(t, actual_item)
		assert.Len(t, actual_rss.Items, 0)
	})
}

func setUp() (ctx context.Context, client *dynamodb.Client) {
	ctx = context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Credentials = credentials.NewStaticCredentialsProvider("dummy", "dummy", "")
		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("Error loading AWS configuration: %v", err))
	}

	client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})

	helper.DropTableIfNotExists(ctx, client, "Rss")
	helper.CreateTableIfNotExists(ctx, client, "Rss", schemaProvider)

	return ctx, client
}

func getTestRss(t *testing.T) rss.Rss {
	var test_rss rss.Rss
	helper.MustSucceed(t, func() error {
		var err error
		test_rss, err = rss.New("Test Title", "Test_Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
		if err != nil {
			return err
		}

		test_rss.SetItemFilter([]string{"include_keyword"}, []string{"exclude_keyword"})

		items := []struct {
			guid, title, link, description, author string
			pubDate                                time.Time
		}{
			{"guid-12345", "Test Title 1", "http://example.com/1", "Test description 1-include_keyword", "Test Author 1", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC)},
			{"guid-67890", "Test Title 2", "http://example.com/2", "Test description 2-include_keyword", "Test Author 2", time.Date(2023, time.June, 2, 13, 30, 0, 0, time.UTC)},
		}

		for _, item := range items {
			rssItem, err := rss.NewItem(rss.Guid{Value: item.guid}, item.title, item.link, item.description, item.author, item.pubDate)
			if err != nil {
				return err
			}
			test_rss.AddOrUpdateItem(rssItem)
		}

		return nil
	})

	return test_rss
}

type rssSaver interface {
	Save(ctx context.Context, rss rss.Rss, updateBy metadata.UserMeta) (rss.Rss, error)
}

func setupExpectedRss(t *testing.T, ctx context.Context, rssSaver rssSaver) (setUpRss rss.Rss) {
	var test_rss rss.Rss
	helper.MustSucceed(t, func() error {
		var err error
		test_rss, err = rss.New("Test Title", "Test_Source", "http://example.com", "Test Description", "en", time.Date(2024, time.June, 1, 13, 30, 0, 0, time.UTC))
		if err != nil {
			return err
		}

		test_rss.SetItemFilter([]string{"Test"}, []string{"unit"})

		items := []struct {
			guid, title, link, description, author string
			pubDate                                time.Time
		}{
			{"guid-12345", "Test Title 1", "http://example.com/1", "Test description 1", "Test Author 1", time.Date(2023, time.June, 1, 13, 30, 0, 0, time.UTC)},
			{"guid-67890", "Test Title 2", "http://example.com/2", "Test description 2", "Test Author 2", time.Date(2023, time.June, 2, 13, 30, 0, 0, time.UTC)},
		}

		for _, item := range items {
			rssItem, err := rss.NewItem(rss.Guid{Value: item.guid}, item.title, item.link, item.description, item.author, item.pubDate)
			if err != nil {
				return err
			}
			test_rss.AddOrUpdateItem(rssItem)
		}

		return nil
	})

	_, err := rssSaver.Save(ctx, test_rss, metadata.UserMeta{ID: "test-id", Name: "test-user"})

	if err != nil {
		panic(err)
	}

	return test_rss
}
