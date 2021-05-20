package util

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// RespondOk builds an API OK response with the given data as the response body.
func RespondOk(data interface{}) events.APIGatewayProxyResponse {
	return Respond(http.StatusOK, data)
}

// RespondBadRequest builds an API BadRequest response with the
// given err as the response body.
func RespondBadRequest(err error) events.APIGatewayProxyResponse {
	return Respond(http.StatusBadRequest, Error{Error: err.Error()})
}

// RespondError builds an API InternalServerError response with the
// given err as the response body.
func RespondError(err error) events.APIGatewayProxyResponse {
	return Respond(http.StatusInternalServerError, Error{Error: err.Error()})
}

// Error is a common error response type. This standardizes the API errors.
type Error struct {
	Error string `json:"error"`
}

// Respond builds an APIGatewayProxyResponse with the given statusCode and data.
// The data will be returned as the response body as JSON.
func Respond(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	var body string
	if data != nil {
		bytes, _ := json.Marshal(data)
		body = string(bytes)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		IsBase64Encoded: false,
		Body:            body,
	}
}
