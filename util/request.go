package util

import (
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
