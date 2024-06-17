package rss

import (
	"context"
	"errors"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	Save(ctx context.Context, rss Rss, updateBy metadata.UserMeta) (Rss, error)
}

type DynamoDBRssRepository struct {
	dynamoDBStore infrastructure.DynamoDBStore
}

func NewDynamoDBRssRepository(client *dynamodb.Client) *DynamoDBRssRepository {
	return &DynamoDBRssRepository{dynamoDBStore: *infrastructure.NewDynamoDBStore(client, "Rss")}
}

func (r *DynamoDBRssRepository) FindBySource(ctx context.Context, source string) (Rss, error) {

	if source == "" {
		return Rss{}, errors.New("invalid source")
	}

	rssModel, err := r.getRssModel(ctx, source)
	if err != nil {
		return Rss{}, err
	}

	itemModels, err := r.getItemModels(ctx, source, rssModel.RssId)
	if err != nil {
		return Rss{}, err
	}

	manager := rssManager{
		rss:   rssModel,
		items: itemModels,
	}

	return buildRss(manager), nil
}

func (r *DynamoDBRssRepository) getRssModel(ctx context.Context, source string) (rssModel, error) {
	result, err := r.dynamoDBStore.GetItemById(ctx, source, "rss")
	if err != nil {
		return rssModel{}, err
	}

	var model rssModel
	err = attributevalue.UnmarshalMap(result.Item, &model)
	if err != nil {
		return rssModel{}, err
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

func buildRss(manager rssManager) Rss {
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
