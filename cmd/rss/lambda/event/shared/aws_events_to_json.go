package shared

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func SnsEventToJson(event events.SNSEvent) string {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal SNS event: %v", err)
		return ""
	}
	return string(eventJSON)
}

func EventBridgeEventToJson(event events.EventBridgeEvent) string {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal EventBridge event: %v", err)
		return ""
	}
	return string(eventJSON)
}

func DynamoDBEventToJson(event events.DynamoDBEvent) string {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal DynamoDB event: %v", err)
		return ""
	}
	return string(eventJSON)
}

func APIGatewayProxyRequestToJson(event events.APIGatewayProxyRequest) string {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal EventBridge event: %v", err)
		return ""
	}
	return string(eventJSON)
}
