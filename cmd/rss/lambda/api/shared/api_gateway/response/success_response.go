package response

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func NoContentResponse() events.APIGatewayProxyResponse {

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization, x-user-id, x-user-name, x-hospital-code",
		},
	}
}

func OKResponse(data interface{}) events.APIGatewayProxyResponse {
	var body string
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return ErrorResponse(http.StatusInternalServerError, "Failed to marshal success response")
		}
		body = string(jsonData)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization, x-user-id, x-user-name, x-hospital-code",
		},
	}
}
