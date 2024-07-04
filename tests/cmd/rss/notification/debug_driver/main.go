package main

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/notification/handler"
	"github.com/YamazakiNorihito/workday/pkg/localdebug"
)

func main() {
	streamARN := "arn:aws:dynamodb:ddblocal:000000000000:table/Rss/stream/2024-06-04T23:55:01.203"
	ctx := context.Background()
	localdebug.PollStreamAndInvokeHandler(ctx, streamARN, handler.Handler)
}
