package util

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHeader_WhereNameIsExactCase_ReturnsValue(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Foo": "Bar",
		},
	}

	v := Header(req, "Foo")
	assert.Equal(t, "Bar", v)
}

func TestHeader_WhereNameIsWrongCase_ReturnsValue(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Foo": "Bar",
		},
	}

	v := Header(req, "foo")
	assert.Equal(t, "Bar", v)
}

func TestHeader_WhereHeaderIsNotPresent_ReturnsEmptyString(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Headers: map[string]string{},
	}

	v := Header(req, "Foo")
	assert.Equal(t, "", v)
}
