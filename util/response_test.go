package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondOk(t *testing.T) {
	data := map[string]string{
		"message": "Hello World",
	}

	resp := RespondOk(data)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(data)
	assert.Equal(t, string(bytes), resp.Body)
}

func TestRespondBadRequest(t *testing.T) {
	err := errors.New("error")

	resp := RespondBadRequest(err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(Error{Error: err.Error()})
	assert.Equal(t, string(bytes), resp.Body)
}

func TestRespondError(t *testing.T) {
	err := errors.New("error")

	resp := RespondError(err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(Error{Error: err.Error()})
	assert.Equal(t, string(bytes), resp.Body)
}

func TestRespondMethodNotAllowed(t *testing.T) {
	err := errors.New("error")

	resp := RespondMethodNotAllowed(err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(Error{Error: err.Error()})
	assert.Equal(t, string(bytes), resp.Body)
}
