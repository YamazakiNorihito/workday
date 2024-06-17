package localdebug

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
)

func PollStreamAndInvokeHandler(ctx context.Context, streamArn string, debugHandle func(ctx context.Context, e events.DynamoDBEvent) error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodbstreams.NewFromConfig(cfg, func(o *dynamodbstreams.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})

	describeStreamOutput, err := client.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{
		StreamArn: &streamArn,
	})
	if err != nil {
		log.Fatalf("failed to describe stream, %v", err)
	}

	shardIteratorType := types.ShardIteratorTypeTrimHorizon
	shardIteratorOutput, err := client.GetShardIterator(ctx, &dynamodbstreams.GetShardIteratorInput{
		StreamArn:         &streamArn,
		ShardId:           describeStreamOutput.StreamDescription.Shards[0].ShardId,
		ShardIteratorType: shardIteratorType,
	})
	if err != nil {
		log.Fatalf("failed to get shard iterator, %v", err)
	}

	shardIterator := shardIteratorOutput.ShardIterator
	for {

		output, err := client.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{
			ShardIterator: shardIterator,
		})
		if err != nil {
			log.Fatalf("failed to get records, %v", err)
		}

		if 0 < len(output.Records) {
			for _, record := range output.Records {
				var event events.DynamoDBEvent

				change, _ := convertStreamRecordUsingJSON(record.Dynamodb)
				event.Records = append(event.Records, events.DynamoDBEventRecord{
					EventID:      *record.EventID,
					EventName:    string(record.EventName),
					EventVersion: *record.EventVersion,
					EventSource:  *record.EventSource,
					AWSRegion:    *record.AwsRegion,
					Change:       change,
				})

				if err := debugHandle(ctx, event); err != nil {
					log.Fatalf("failed to handle request, %v", err)
				}
			}
		}

		if output.NextShardIterator == nil {
			break
		}
		shardIterator = output.NextShardIterator
		time.Sleep(1 * time.Second)
	}
}

func convertStreamRecordUsingJSON(streamRecord *types.StreamRecord) (dynamoStreamRecord events.DynamoDBStreamRecord, error error) {
	// First, ensure that the StreamRecord is correctly populated with all required fields
	if streamRecord == nil {
		return events.DynamoDBStreamRecord{}, fmt.Errorf("streamRecord is nil")
	}

	dynamoStreamRecord.Keys = convertAttributeValueMap(streamRecord.Keys)
	dynamoStreamRecord.NewImage = convertAttributeValueMap(streamRecord.NewImage)
	dynamoStreamRecord.OldImage = convertAttributeValueMap(streamRecord.OldImage)
	dynamoStreamRecord.SequenceNumber = *streamRecord.SequenceNumber
	dynamoStreamRecord.SizeBytes = *streamRecord.SizeBytes
	dynamoStreamRecord.StreamViewType = string(streamRecord.StreamViewType)

	return dynamoStreamRecord, nil
}

func convertAttributeValueMap(attributeValueMap map[string]types.AttributeValue) map[string]events.DynamoDBAttributeValue {
	result := make(map[string]events.DynamoDBAttributeValue)
	for k, v := range attributeValueMap {
		result[k] = convertAttributeValue(v)
	}
	return result
}

func convertAttributeValue(attributeValue types.AttributeValue) events.DynamoDBAttributeValue {
	switch v := attributeValue.(type) {
	case *types.AttributeValueMemberS:
		return events.NewStringAttribute(v.Value)
	case *types.AttributeValueMemberN:
		return events.NewNumberAttribute(v.Value)
	case *types.AttributeValueMemberBOOL:
		return events.NewBooleanAttribute(v.Value)
	case *types.AttributeValueMemberB:
		return events.NewBinaryAttribute(v.Value)
	case *types.AttributeValueMemberSS:
		return events.NewStringSetAttribute(v.Value)
	case *types.AttributeValueMemberNS:
		return events.NewNumberSetAttribute(v.Value)
	case *types.AttributeValueMemberBS:
		return events.NewBinarySetAttribute(v.Value)
	case *types.AttributeValueMemberM:
		return events.NewMapAttribute(convertAttributeValueMap(v.Value))
	case *types.AttributeValueMemberL:
		var list []events.DynamoDBAttributeValue
		for _, item := range v.Value {
			list = append(list, convertAttributeValue(item))
		}
		return events.NewListAttribute(list)
	case *types.AttributeValueMemberNULL:
		return events.NewNullAttribute()
	default:
		return events.NewNullAttribute()
	}
}
