package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/feeds/app_service"
	apiGatewayResponse "github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/api_gateway/response"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/aws/aws-lambda-go/events"
)

type executer func(ctx context.Context, logger infrastructure.Logger) ([]app_service.RssResponse, error)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg := awsConfig.LoadConfig(ctx)
	dynamodbClient := cfg.NewDynamodbClient()
	rssRepository := rss.NewDynamoDBRssRepository(dynamodbClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("APIGatewayProxyRequest Event", "event", shared.APIGatewayProxyRequestToJson(request))

	executer := func(ctx context.Context, logger infrastructure.Logger) ([]app_service.RssResponse, error) {
		return app_service.Execute(ctx, logger, rssRepository)
	}
	logger.Info("finish")
	return processRecord(ctx, logger, executer, request), nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, executer executer, _ events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	rssFeeds, err := executer(ctx, logger)

	if err != nil {
		logger.Error("Failed", "error", err)
		return apiGatewayResponse.ErrorResponse(http.StatusInternalServerError, err.Error())
	}
	return apiGatewayResponse.OKResponse(rssFeeds)
}
