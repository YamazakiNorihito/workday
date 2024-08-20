package main

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/translate/handler"
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
			  "language": "en",
			  "last_build_date": "2024-07-03T13:00:00Z",
			  "items": {
				"http://www.example.com/dummy-guid1": {
				  "guid": "http://www.example.com/dummy-guid1",
				  "title": "ダミー記事1",
				  "link": "http://www.example.com/dummy-article1",
				  "description": "Taking an average of 10 years and $1.3 billion to develop a single new medication, pharmaceutical companies often focus their drug discovery efforts on a high return on investment, developing drugs for diseases prevalent in high-income countries—and leaving lower- and middle-income countries behind.In response, investments in building AI/ML models for drug discovery have soared in the last five years. By using these models, scientists can shorten their research and development timeline by getting better at identifying drug prospects. However, access to these models is limited by data science expertise and computational resources.The nonprofit Ersilia Open Source Initiative is tackling this problem with the Ersilia Model Hub.Through the hub, Ersilia aims to disseminate AI/ML models and computational power to researchers focused on drug discovery for infectious diseases in regions outside of Europe and North America.",
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
