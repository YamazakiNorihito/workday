package main

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/delete/handler"
	"github.com/aws/aws-lambda-go/events"
)

func main() {
	event := events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					MessageID: "12345",
					Message: `{
			"rss": {
			  "source": "127.0.0.1:8080"
			}
		  }`,
				},
			},
		},
	}
	handler.Handler(context.Background(), event)
}
