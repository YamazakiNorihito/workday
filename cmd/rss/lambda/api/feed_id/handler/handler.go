package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/feed_id/app_service"
	apiGatewayResponse "github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/api_gateway/response"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/aws/aws-lambda-go/events"
)

type executer func(ctx context.Context, logger infrastructure.Logger, command app_service.GetCommand) (app_service.RssResponse, error)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg := awsConfig.LoadConfig(ctx)
	dynamodbClient := cfg.NewDynamodbClient()
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("APIGatewayProxyRequest Event", "event", shared.APIGatewayProxyRequestToJson(request))

	executer := func(ctx context.Context, logger infrastructure.Logger, command app_service.GetCommand) (app_service.RssResponse, error) {
		return app_service.Execute(ctx, logger, rssRepository, command)
	}
	logger.Info("finish")
	return processRecord(ctx, logger, executer, request), nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, executer executer, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	source := request.PathParameters["source"]

	command := app_service.GetCommand{
		Source: source,
	}

	rssFeeds, err := executer(ctx, logger, command)

	if err != nil {
		logger.Error("Failed", "error", err)
		if _, ok := err.(*app_service.ValidationError); ok {
			return apiGatewayResponse.ErrorResponse(http.StatusBadRequest, err.Error())
		} else {
			return apiGatewayResponse.ErrorResponse(http.StatusInternalServerError, err.Error())
		}
	}
	return apiGatewayResponse.OKResponse(rssFeeds)
}
