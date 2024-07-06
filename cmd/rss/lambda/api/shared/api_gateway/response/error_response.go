package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type errorResponse struct {
	Error errorDetails `json:"error"`
}

type errorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	errorResponse := errorResponse{
		Error: errorDetails{
			Code:    statusCode,
			Message: message,
		},
	}
	body, err := json.Marshal(errorResponse)
	if err != nil {
		body = []byte(`{"error": {"code": 500, "message": "Failed to marshal error response"}}`)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization, x-user-id, x-user-name, x-hospital-code",
		},
	}
}
