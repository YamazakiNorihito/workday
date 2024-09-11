package main

import (
	"context"
	"fmt"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/feeds/handler"
	"github.com/aws/aws-lambda-go/events"
)

func main() {
	event := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/feeds",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	handler.Handler(context.Background(), event)
	response, err := handler.Handler(context.Background(), event)

	// 結果をコンソールに表示
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Response: %+v\n", response)
	}
}
