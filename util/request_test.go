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

func TestReadJSON_GivenBase64Request_UnmarshalsBody(t *testing.T) {
	const body = "eyJmb28iOiJiYXIifQ=="

	var data map[string]interface{}
	ReadJSON(events.APIGatewayProxyRequest{
		IsBase64Encoded: true,
		Body:            body,
	}, &data)

	assert.Equal(t, "bar", data["foo"])
}

func TestReadJSON_GivenPlainTextRequest_UnmarshalsBody(t *testing.T) {
	const body = `{"foo":"bar"}`

	var data map[string]interface{}
	ReadJSON(events.APIGatewayProxyRequest{
		IsBase64Encoded: false,
		Body:            body,
	}, &data)

	assert.Equal(t, "bar", data["foo"])
}

func TestReadForm_GivenBase64Request_DecodesBody(t *testing.T) {
	const body = "Zm9vPWJhcg=="

	data := ReadForm(events.APIGatewayProxyRequest{
		IsBase64Encoded: true,
		Body:            body,
	})

	assert.Equal(t, "bar", data.Get("foo"))
}

func TestReadForm_GivenPlainTextRequest_DecodesBody(t *testing.T) {
	const body = "foo=bar"

	data := ReadForm(events.APIGatewayProxyRequest{
		IsBase64Encoded: false,
		Body:            body,
	})

	assert.Equal(t, "bar", data.Get("foo"))
}
