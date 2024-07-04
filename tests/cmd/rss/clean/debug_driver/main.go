package main

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/clean/handler"
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
			  "id": "0a1fd873-6974-427b-a5dc-cbb809aff1aa",
			  "source": "127.0.0.1:8080",
			  "title": "ダミーニュースのフィード",
			  "link": "http://127.0.0.1:8080",
			  "description": "このフィードはダミーニュースを提供します。",
			  "language": "ja",
			  "last_build_date": "2024-07-03T13:00:00Z",
			  "items": {
				"http://www.example.com/dummy-guid1": {
				  "guid": "http://www.example.com/dummy-guid1",
				  "title": "ダミー記事1",
				  "link": "http://www.example.com/dummy-article1",
				  "description": "これはダミー記事1の概要です。詳細はリンクをクリックしてください。",
				  "author": "item1@dummy.com",
				  "pubDate": "2024-07-03T12:00:00Z",
				  "tags": []
				}
			  },
			  "create_by": {
				"id": "",
				"name": ""
			  },
			  "create_at": "0001-01-01T00:00:00Z",
			  "update_by": {
				"id": "",
				"name": ""
			  },
			  "update_at": "0001-01-01T00:00:00Z"
			},
			"compressed": false
		  }`,
				},
			},
		},
	}
	handler.Handler(context.Background(), event)
}
