package main

import (
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/notification/handler"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Handler)
}
