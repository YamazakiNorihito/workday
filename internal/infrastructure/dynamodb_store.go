package infrastructure

import (
	"context"

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
