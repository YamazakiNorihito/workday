package rss

import (
	"context"
	"errors"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type rssManager struct {
	rss   rssModel
	items []itemModel
}

type rssModel struct {
	PartitionKey string `dynamodbav:"id"`
	SortKey      string `dynamodbav:"sortKey"`

	RssId         string            `dynamodbav:"rss_id"`
	Source        string            `dynamodbav:"source"`
	Title         string            `dynamodbav:"title"`
	Link          string            `dynamodbav:"link"`
	Description   string            `dynamodbav:"description"`
	Language      string            `dynamodbav:"language"`
	LastBuildDate int64             `dynamodbav:"last_build_date"`
	ItemFilter    itemFilterModel   `dynamodbav:"item_filter"`
	CreatedBy     metadata.CreateBy `dynamodbav:"create_by"`
	CreatedAt     int64             `dynamodbav:"create_at"`
	UpdatedBy     metadata.UpdateBy `dynamodbav:"update_by"`
	UpdatedAt     int64             `dynamodbav:"update_at"`
}

type itemModel struct {
	PartitionKey string   `dynamodbav:"id"`
	SortKey      string   `dynamodbav:"sortKey"`
	RssId        string   `dynamodbav:"rss_id"`
	GuId         string   `dynamodbav:"guid"`
	Title        string   `dynamodbav:"title"`
	Link         string   `dynamodbav:"link"`
	Description  string   `dynamodbav:"description"`
	Author       string   `dynamodbav:"author"`
	PubDate      int64    `dynamodbav:"pub_date"`
	Tags         []string `dynamodbav:"tags"`
}

type itemFilterModel struct {
	IncludeKeywords []string `dynamodbav:"include_keywords"`
	ExcludeKeywords []string `dynamodbav:"exclude_keywords"`
}

func (r *rssModel) NewItemModel(item Item) itemModel {
	return itemModel{
		PartitionKey: r.PartitionKey,
		SortKey:      r.RssId + "#" + item.Guid.Value,
		RssId:        r.RssId,
		GuId:         item.Guid.Value,
		Title:        item.Title,
		Link:         item.Link,
		Description:  item.Description,
		Author:       item.Author,
		PubDate:      item.PubDate.Unix(),
		Tags:         item.Tags,
	}
}

type IRssRepository interface {
	FindBySource(ctx context.Context, source string) (Rss, error)
	FindAll(ctx context.Context) ([]Rss, error)
	FindItems(ctx context.Context, rss Rss) (Rss, error)
	FindItemsByPk(ctx context.Context, rss Rss, guid Guid) (Rss, error)
	Save(ctx context.Context, rss Rss, updateBy metadata.UserMeta) (Rss, error)
	Delete(ctx context.Context, rss Rss) error
}

type DynamoDBRssRepository struct {
	dynamoDBStore infrastructure.DynamoDBStore
}

func NewDynamoDBRssRepository(client *dynamodb.Client) *DynamoDBRssRepository {
	return &DynamoDBRssRepository{dynamoDBStore: *infrastructure.NewDynamoDBStore(client, "Rss")}
}

// FindBySource retrieves an Rss instance from the repository based on the given source.
// Note:
// - Due to data volume concerns, `Item` data is not returned by this function.
// - If `Item` data is needed, use `FindItems` or `FindItemsByPk` functions instead.
func (r *DynamoDBRssRepository) FindBySource(ctx context.Context, source string) (Rss, error) {
	if source == "" {
		return Rss{}, errors.New("invalid source")
	}

	result, err := r.dynamoDBStore.GetItemById(ctx, source, "rss")
	if err != nil {
		return Rss{}, err
	}

	var rssModel rssModel
	err = attributevalue.UnmarshalMap(result.Item, &rssModel)
	if err != nil {
		return Rss{}, err
	}

	manager := rssManager{
		rss:   rssModel,
		items: []itemModel{},
	}

	return buildRss(manager), nil
}

func (r *DynamoDBRssRepository) FindAll(ctx context.Context) ([]Rss, error) {
	results, err := r.dynamoDBStore.QueryItemsBySortKey(ctx, "rss")
	if err != nil {
		return []Rss{}, err
	}

	var rssFeeds []Rss
	for _, item := range results.Items {
		var rssModel rssModel
		err = attributevalue.UnmarshalMap(item, &rssModel)
		if err != nil {
			return []Rss{}, err
		}

		manager := rssManager{
			rss:   rssModel,
			items: []itemModel{},
		}

		rss := buildRss(manager)
		rssFeeds = append(rssFeeds, rss)
	}

	return rssFeeds, nil
}

func (r *DynamoDBRssRepository) FindItems(ctx context.Context, rss Rss) (Rss, error) {
	if rss.Source == "" {
		return Rss{}, errors.New("invalid source")
	}

	manager := buildRssManager(rss)
	itemModels, err := r.getItemModels(ctx, manager.rss.PartitionKey, manager.rss.RssId)
	if err != nil {
		return Rss{}, err
	}

	finalManager := rssManager{
		rss:   manager.rss,
		items: itemModels,
	}

	return buildRss(finalManager), nil
}

func (r *DynamoDBRssRepository) FindItemsByPk(ctx context.Context, rss Rss, guid Guid) (Rss, error) {
	if rss.ID.String() == "" || guid.Value == "" {
		return Rss{}, errors.New("invalid source")
	}

	manager := buildRssManager(rss)
	searchItemModel := manager.rss.NewItemModel(Item{Guid: guid})

	findItemModel, err := r.getItemModel(ctx, searchItemModel.PartitionKey, searchItemModel.SortKey)
	if err != nil {
		return Rss{}, err
	}

	var itemModels []itemModel
	if findItemModel.PartitionKey != "" {
		itemModels = []itemModel{findItemModel}
	} else {
		itemModels = []itemModel{}
	}

	finalManager := rssManager{
		rss:   manager.rss,
		items: itemModels,
	}

	return buildRss(finalManager), nil
}

func (r *DynamoDBRssRepository) getItemModel(ctx context.Context, partitionKey string, sortKey string) (itemModel, error) {
	result, err := r.dynamoDBStore.GetItemById(ctx, partitionKey, sortKey)
	if err != nil {
		return itemModel{}, err
	}

	var model itemModel
	err = attributevalue.UnmarshalMap(result.Item, &model)
	if err != nil {
		return itemModel{}, err
	}

	return model, nil
}

func (r *DynamoDBRssRepository) getItemModels(ctx context.Context, source string, rssId string) ([]itemModel, error) {
	result, err := r.dynamoDBStore.QueryItemsBySortPrefix(ctx, source, rssId)
	if err != nil {
		return nil, err
	}

	var model []itemModel
	err = attributevalue.UnmarshalListOfMaps(result.Items, &model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *DynamoDBRssRepository) Save(ctx context.Context, rss Rss, updateBy metadata.UserMeta) (Rss, error) {

	if rss.ID == uuid.Nil {
		return rss, errors.New("invalid rss ID")
	}

	now := time.Now()

	if rss.CreatedBy.ID == "" {
		rss.CreatedAt = metadata.CreateAt(now)
		rss.CreatedBy = metadata.CreateBy(updateBy)
	}
	rss.UpdatedAt = metadata.UpdateAt(now)
	rss.UpdatedBy = metadata.UpdateBy(updateBy)

	rssManager := buildRssManager(rss)

	// rss
	err := r.dynamoDBStore.PutItem(ctx, rssManager.rss)
	if err != nil {
		return rss, err
	}

	// item
	for _, item := range rssManager.items {
		err = r.dynamoDBStore.PutItem(ctx, item)
		if err != nil {
			return rss, err
		}
	}

	return rss, nil
}

func (r *DynamoDBRssRepository) Delete(ctx context.Context, rss Rss) error {
	if rss.ID == uuid.Nil {
		return errors.New("invalid rss ID")
	}
	if rss.Source == "" {
		return errors.New("invalid the Source")
	}

	manager := buildRssManager(rss)
	itemModels, err := r.getItemModels(ctx, manager.rss.PartitionKey, manager.rss.RssId)
	if err != nil {
		return err
	}

	var deleteInputs []dynamodb.DeleteItemInput
	for _, item := range itemModels {
		//r.dynamoDBStore.DeleteItem(ctx, item.PartitionKey, item.SortKey)
		deleteInputs = append(deleteInputs, dynamodb.DeleteItemInput{
			TableName: aws.String(r.dynamoDBStore.TableName),
			Key: map[string]types.AttributeValue{
				"id":      &types.AttributeValueMemberS{Value: item.PartitionKey},
				"sortKey": &types.AttributeValueMemberS{Value: item.SortKey},
			},
		})
	}

	err = r.dynamoDBStore.BatchDeleteItems(ctx, deleteInputs)
	if err != nil {
		return err
	}

	_, err = r.dynamoDBStore.DeleteItem(ctx, manager.rss.PartitionKey, manager.rss.SortKey)
	if err != nil {
		return err
	}

	return nil
}

func buildRss(manager rssManager) Rss {
	if manager.rss.RssId == "" {
		return Rss{
			Items: make(map[Guid]Item),
		}
	}

	itemsMap := make(map[Guid]Item)

	for _, item := range manager.items {
		itemsMap[Guid{Value: item.GuId}] = Item{
			Guid:        Guid{Value: item.GuId},
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Author:      item.Author,
			PubDate:     time.Unix(item.PubDate, 0).UTC(),
			Tags:        item.Tags,
		}
	}

	rss := Rss{
		ID:            uuid.MustParse(manager.rss.RssId),
		Source:        manager.rss.Source,
		Title:         manager.rss.Title,
		Link:          manager.rss.Link,
		Description:   manager.rss.Description,
		Language:      manager.rss.Language,
		LastBuildDate: time.Unix(manager.rss.LastBuildDate, 0),
		ItemFilter:    ItemFilter(manager.rss.ItemFilter),
		Items:         itemsMap,
		CreatedBy:     manager.rss.CreatedBy,
		CreatedAt:     time.Unix(manager.rss.CreatedAt, 0).UTC(),
		UpdatedBy:     manager.rss.UpdatedBy,
		UpdatedAt:     time.Unix(manager.rss.UpdatedAt, 0).UTC(),
	}
	return rss
}

func buildRssManager(rss Rss) rssManager {

	rssModel := rssModel{
		PartitionKey:  rss.Source,
		SortKey:       "rss",
		RssId:         rss.ID.String(),
		Source:        rss.Source,
		Title:         rss.Title,
		Link:          rss.Link,
		Description:   rss.Description,
		Language:      rss.Language,
		LastBuildDate: rss.LastBuildDate.Unix(),
		ItemFilter:    itemFilterModel(rss.ItemFilter),
		CreatedBy:     rss.CreatedBy,
		CreatedAt:     rss.CreatedAt.Unix(),
		UpdatedBy:     rss.UpdatedBy,
		UpdatedAt:     rss.UpdatedAt.Unix(),
	}

	itemModels := []itemModel{}
	for _, item := range rss.Items {
		itemModel := rssModel.NewItemModel(item)
		itemModels = append(itemModels, itemModel)
	}

	manager := rssManager{
		rss:   rssModel,
		items: itemModels,
	}

	return manager
}
