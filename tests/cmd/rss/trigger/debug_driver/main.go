package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/trigger/handler"
	"github.com/aws/aws-lambda-go/events"
)

func main() {
	event := events.EventBridgeEvent{
		Version:    "0",
		ID:         "cdc73f9d-aea9-11e3-9d5a-835b769c0d9c",
		DetailType: "Scheduled Event",
		Source:     "aws.events",
		AccountID:  "123456789012",
		Time:       time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Region:     "us-east-1",
		Resources:  []string{"arn:aws:events:us-east-1:123456789012:rule/ExampleRule"},
		Detail:     json.RawMessage(`{}`),
	}
	handler.Handler(context.Background(), event)
}
