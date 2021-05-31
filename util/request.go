package util

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// Header returns the header value from req, with the given name.
// The name of the header's case is ignored. If there is no header
// an empty string is returned.
func Header(req events.APIGatewayProxyRequest, name string) string {
	for k, v := range req.Headers {
		if strings.ToLower(k) == strings.ToLower(name) {
			return v
		}
	}

	return ""
}

// Body is used to read a JSON request body. If the request
// body is base64 encoded, it will be decoded and the unmarshalled.
func Body(req events.APIGatewayProxyRequest, v interface{}) {
	body := []byte(req.Body)
	if req.IsBase64Encoded {
		body, _ = base64.StdEncoding.DecodeString(req.Body)
	}

	json.Unmarshal(body, v)
}
