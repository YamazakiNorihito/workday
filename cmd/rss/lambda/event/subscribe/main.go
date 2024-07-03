package main

import (
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/subscribe/handler"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Handler)
}
