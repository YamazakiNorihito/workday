package infrastructure

import (
	"context"

	"github.com/YamazakiNorihito/workday/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStore struct {
	client    *dynamodb.Client
	TableName string
}

func NewDynamoDBStore(client *dynamodb.Client, tableName string) *DynamoDBStore {
	return &DynamoDBStore{
		client:    client,
		TableName: tableName,
	}
}

func (r *DynamoDBStore) GetItemById(ctx context.Context, partitionkey string, sortKey string) (*dynamodb.GetItemOutput, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"id":      &types.AttributeValueMemberS{Value: partitionkey},
			"sortKey": &types.AttributeValueMemberS{Value: sortKey},
		},
	}
	optFns := func(o *dynamodb.Options) {
		o.RetryMaxAttempts = 1
		o.RetryMode = aws.RetryModeStandard
	}
	result, err := r.client.GetItem(ctx, input, optFns)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *DynamoDBStore) QueryItemsBySortPrefix(ctx context.Context, partitionkey string, sortKeyPrefix string) (*dynamodb.QueryOutput, error) {
	// https://docs.aws.amazon.com/ja_jp/amazondynamodb/latest/developerguide/LegacyConditionalParameters.KeyConditions.html#KeyConditionExpression.instead
	input := &dynamodb.QueryInput{
		TableName:              &r.TableName,
		KeyConditionExpression: aws.String("id = :id AND begins_with(sortKey, :sortKeyPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id":            &types.AttributeValueMemberS{Value: partitionkey},
			":sortKeyPrefix": &types.AttributeValueMemberS{Value: sortKeyPrefix},
		},
	}
	optFns := func(o *dynamodb.Options) {
		o.RetryMaxAttempts = 1
		o.RetryMode = aws.RetryModeStandard
	}
	result, err := r.client.Query(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *DynamoDBStore) QueryItemsBySortKey(ctx context.Context, sortKey string) (*dynamodb.QueryOutput, error) {
	input := &dynamodb.QueryInput{
		TableName:              &r.TableName,
		IndexName:              aws.String("SortKeyIndex"),
		KeyConditionExpression: aws.String("sortKey = :sortKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sortKey": &types.AttributeValueMemberS{Value: sortKey},
		},
	}
	optFns := func(o *dynamodb.Options) {
		o.RetryMaxAttempts = 1
		o.RetryMode = aws.RetryModeStandard
	}
	result, err := r.client.Query(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *DynamoDBStore) PutItem(ctx context.Context, item interface{}) error {

	mapItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: &r.TableName,
		Item:      mapItem,
	}
	optFns := func(o *dynamodb.Options) {
		o.RetryMaxAttempts = 3
		o.RetryMode = aws.RetryModeStandard
	}

	_, err = r.client.PutItem(ctx, input, optFns)
	if err != nil {
		return err
	}
	return nil
}

func (r *DynamoDBStore) DeleteItem(ctx context.Context, partitionKey string, sortKey string) (*dynamodb.DeleteItemOutput, error) {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"id":      &types.AttributeValueMemberS{Value: partitionKey},
			"sortKey": &types.AttributeValueMemberS{Value: sortKey},
		},
	}
	optFns := func(o *dynamodb.Options) {
		o.RetryMaxAttempts = 1
		o.RetryMode = aws.RetryModeStandard
	}

	result, err := r.client.DeleteItem(ctx, input, optFns)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *DynamoDBStore) BatchDeleteItems(ctx context.Context, deleteInputs []dynamodb.DeleteItemInput) error {
	chunks := utils.ChunkSlice(deleteInputs, 25)

	for _, chunk := range chunks {
		writeRequests := make([]types.WriteRequest, len(chunk))
		for j, input := range chunk {
			writeRequests[j] = types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: input.Key,
				},
			}
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.TableName: writeRequests,
			},
		}

		result, err := r.client.BatchWriteItem(ctx, input)
		if err != nil {
			return err
		}

		for len(result.UnprocessedItems) > 0 {
			input.RequestItems = result.UnprocessedItems
			result, err = r.client.BatchWriteItem(ctx, input)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
