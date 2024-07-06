package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/create/app_service"
	apiGatewayResponse "github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/api_gateway/response"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/aws/aws-lambda-go/events"
)

type requestBody struct {
	FeedURL            string `json:"feed_url"`
	SourceLanguageCode string `json:"source_language_code"`
}

type executer func(ctx context.Context, logger infrastructure.Logger, command app_service.CreateCommand) error

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg := awsConfig.LoadConfig(ctx)
	snsClient := cfg.NewSnsClient()
	snsTopicClient := awsConfig.NewSnsTopicClient(snsClient, os.Getenv("OUTPUT_TOPIC_RSS_ARN"))
	publisher := publisher.NewSubscribeMessagePublisher(snsTopicClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("APIGatewayProxyRequest Event", "event", shared.APIGatewayProxyRequestToJson(request))

	executer := func(ctx context.Context, logger infrastructure.Logger, command app_service.CreateCommand) error {
		return app_service.Execute(ctx, logger, *publisher, command)
	}
	logger.Info("finish")
	return processRecord(ctx, logger, executer, request), nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, executer executer, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	requestBody := requestBody{}
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		logger.Error("Failed", "error", "Invalid JSON body")
		return apiGatewayResponse.ErrorResponse(http.StatusBadRequest, "Invalid JSON body")
	}

	cmd := app_service.CreateCommand{
		FeedURL:            requestBody.FeedURL,
		SourceLanguageCode: requestBody.SourceLanguageCode,
	}

	err := executer(ctx, logger, cmd)

	if err != nil {
		logger.Error("Failed", "error", err)
		if _, ok := err.(*app_service.ValidationError); ok {
			return apiGatewayResponse.ErrorResponse(http.StatusBadRequest, err.Error())
		} else {
			return apiGatewayResponse.ErrorResponse(http.StatusInternalServerError, err.Error())
		}
	}
	return apiGatewayResponse.NoContentResponse()
}
