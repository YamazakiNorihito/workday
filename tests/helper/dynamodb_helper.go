package helper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetItems(ctx context.Context, client *dynamodb.Client, tableName string) ([]map[string]interface{}, error) {
	resp, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &tableName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}

	var items []map[string]interface{}
	for _, item := range resp.Items {
		var unmarshalItem map[string]interface{}
		if err := attributevalue.UnmarshalMap(item, &unmarshalItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		items = append(items, unmarshalItem)
	}

	return items, nil
}

func GetItem(ctx context.Context, client *dynamodb.Client, tableName string, partitionkey string, sortKey string) (map[string]interface{}, error) {
	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id":      &types.AttributeValueMemberS{Value: partitionkey},
			"sortKey": &types.AttributeValueMemberS{Value: sortKey},
		},
	})
	if err != nil {
		return nil, err
	}

	var unmarshalItem map[string]interface{}
	if err := attributevalue.UnmarshalMap(result.Item, &unmarshalItem); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return unmarshalItem, nil
}

func existsTable(client *dynamodb.Client, tableName string) (bool, error) {
	exists := true
	_, err := client.DescribeTable(
		context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist.\n", tableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", tableName, err)
		}
		exists = false
	}
	return exists, err
}

func deleteTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName)})
	if err != nil {
		log.Printf("Couldn't delete table %v. Here's why: %v\n", tableName, err)
	}
	return err
}

type SchemaProvider func() ([]types.AttributeDefinition, []types.KeySchemaElement)

func createTable(ctx context.Context, client *dynamodb.Client, tableName string, schemaProvider SchemaProvider) error {
	attributeDefinitions, keySchema := schemaProvider()

	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: attributeDefinitions,
		KeySchema:            keySchema,
		TableName:            aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", tableName, err)
		return err
	}

	waiter := dynamodb.NewTableExistsWaiter(client)
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 5*time.Minute)
	if err != nil {
		log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		return err
	}
	return nil
}

func DropTableIfNotExists(ctx context.Context, client *dynamodb.Client, tableName string) {
	exists, err := existsTable(client, tableName)
	if err != nil {
		panic(fmt.Sprintf("Error checking if table %v exists: %v", tableName, err))
	}
	if !exists {
		return
	}

	err = deleteTable(ctx, client, tableName)
	if err != nil {
		panic(fmt.Sprintf("Couldn't delete table %v: %v", tableName, err))
	}
}

func CreateTableIfNotExists(ctx context.Context, client *dynamodb.Client, tableName string, schemaProvider SchemaProvider) {
	exists, err := existsTable(client, tableName)
	if err != nil {
		panic(fmt.Sprintf("Error checking if table %v exists: %v", tableName, err))
	}
	if exists {
		return
	}

	err = createTable(ctx, client, tableName, schemaProvider)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create table %v: %v", tableName, err))
	}
}
