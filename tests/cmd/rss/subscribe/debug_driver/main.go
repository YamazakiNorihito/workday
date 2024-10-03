package main

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/subscribe/handler"
	"github.com/aws/aws-lambda-go/events"
)

func main() {
	event := events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					MessageID: "12345",
					Message:   `{"feed_url": "https://connpass.com/explore/ja.atom","item_filter": {"include_keywords":[],"exclude_keywords":[".*会.*",".*もくもく.*",".*道場.*"]}}`,
				},
			},
		},
	}
	handler.Handler(context.Background(), event)
}
