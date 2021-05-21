package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/goidc/dal"
	dalMock "github.com/reecerussell/goidc/dal/mock"
	"github.com/reecerussell/goidc/token"
	tokenMock "github.com/reecerussell/goidc/token/mock"
	valMock "github.com/reecerussell/goidc/validator/mock"
)

func TestHandler(t *testing.T) {
	testClientId := "3247023"
	testClientSecret := "2934uldnf"
	testGrantType := "code"
	testRedirectUri := "http://test.io"
	testScopes := "openid"

	testClient := &dal.Client{
		ID:           testClientId,
		RedirectUris: []string{"http://test.io"},
		Scopes:       []string{"openid"},
		Secrets:      []string{"my secret"},
		GrantTypes:   []string{"code"},
	}
	testToken := &token.Token{
		AccessToken: "my.jwt.token",
		TokenType:   "Bearer",
		Expires:     3600,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := dalMock.NewMockClientProvider(ctrl)
	mockProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockValidator.EXPECT().ValidateRequest(testClient, testClientSecret, testRedirectUri, testGrantType, gomock.Any()).Return(nil)

	mockTokenService := tokenMock.NewMockService(ctrl)
	mockTokenService.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(testToken, nil)

	h := &Handler{
		tokens:    mockTokenService,
		clients:   mockProvider,
		validator: mockValidator,
	}

	testBody := url.Values{
		"client_id":     {testClientId},
		"client_secret": {testClientSecret},
		"grant_type":    {testGrantType},
		"redirect_uri":  {testRedirectUri},
		"scope":         {testScopes},
	}

	resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: testBody.Encode(),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(testToken)
	assert.Equal(t, string(bytes), resp.Body)
}

func TestHandler_GivenInvalidHTTPMethod_ReturnsMethodNotSupported(t *testing.T) {
	testClientId := "3247023"
	testClientSecret := "2934uldnf"
	testGrantType := "code"
	testRedirectUri := "http://test.io"
	testScopes := "openid"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := dalMock.NewMockClientProvider(ctrl)
	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockTokenService := tokenMock.NewMockService(ctrl)

	h := &Handler{
		tokens:    mockTokenService,
		clients:   mockProvider,
		validator: mockValidator,
	}

	testBody := url.Values{
		"client_id":     {testClientId},
		"client_secret": {testClientSecret},
		"grant_type":    {testGrantType},
		"redirect_uri":  {testRedirectUri},
		"scope":         {testScopes},
	}

	resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: testBody.Encode(),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(map[string]string{"error": "method not allowed"})
	assert.Equal(t, string(bytes), resp.Body)
}

func TestHandler_GivenInvalidContentType_ReturnsBadRequest(t *testing.T) {
	testClientId := "3247023"
	testClientSecret := "2934uldnf"
	testGrantType := "code"
	testRedirectUri := "http://test.io"
	testScopes := "openid"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := dalMock.NewMockClientProvider(ctrl)
	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockTokenService := tokenMock.NewMockService(ctrl)

	h := &Handler{
		tokens:    mockTokenService,
		clients:   mockProvider,
		validator: mockValidator,
	}

	testBody := url.Values{
		"client_id":     {testClientId},
		"client_secret": {testClientSecret},
		"grant_type":    {testGrantType},
		"redirect_uri":  {testRedirectUri},
		"scope":         {testScopes},
	}

	resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: testBody.Encode(),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(map[string]string{"error": "invalid content type"})
	assert.Equal(t, string(bytes), resp.Body)
}

func TestHandler_GivenInvalidClientId_ReturnsBadRequest(t *testing.T) {
	testClientId := "3247023"
	testClientSecret := "2934uldnf"
	testGrantType := "code"
	testRedirectUri := "http://test.io"
	testScopes := "openid"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := dalMock.NewMockClientProvider(ctrl)
	mockProvider.EXPECT().Get(gomock.Any(), testClientId).Return(nil, dal.ErrClientNotFound)

	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockTokenService := tokenMock.NewMockService(ctrl)

	h := &Handler{
		tokens:    mockTokenService,
		clients:   mockProvider,
		validator: mockValidator,
	}

	testBody := url.Values{
		"client_id":     {testClientId},
		"client_secret": {testClientSecret},
		"grant_type":    {testGrantType},
		"redirect_uri":  {testRedirectUri},
		"scope":         {testScopes},
	}

	resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: testBody.Encode(),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(map[string]string{"error": "invalid client id"})
	assert.Equal(t, string(bytes), resp.Body)
}

func TestHandler_WhereClientProviderFails_ReturnsInternalServerError(t *testing.T) {
	testClientId := "3247023"
	testClientSecret := "2934uldnf"
	testGrantType := "code"
	testRedirectUri := "http://test.io"
	testScopes := "openid"

	testError := errors.New("an error occured")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := dalMock.NewMockClientProvider(ctrl)
	mockProvider.EXPECT().Get(gomock.Any(), testClientId).Return(nil, testError)

	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockTokenService := tokenMock.NewMockService(ctrl)

	h := &Handler{
		tokens:    mockTokenService,
		clients:   mockProvider,
		validator: mockValidator,
	}

	testBody := url.Values{
		"client_id":     {testClientId},
		"client_secret": {testClientSecret},
		"grant_type":    {testGrantType},
		"redirect_uri":  {testRedirectUri},
		"scope":         {testScopes},
	}

	resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: testBody.Encode(),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(map[string]string{"error": testError.Error()})
	assert.Equal(t, string(bytes), resp.Body)
}

func TestHandler_GivenInvalidClient_ReturnsBadRequest(t *testing.T) {
	testClientId := "3247023"
	testClientSecret := "2934uldnf"
	testGrantType := "code"
	testRedirectUri := "http://test.io"
	testScopes := "openid"

	testClient := &dal.Client{
		ID:           testClientId,
		RedirectUris: []string{"http://test.io"},
		Scopes:       []string{"openid"},
		Secrets:      []string{"my secret"},
		GrantTypes:   []string{"code"},
	}
	testError := errors.New("invalid client")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := dalMock.NewMockClientProvider(ctrl)
	mockProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockValidator.EXPECT().ValidateRequest(testClient, testClientSecret, testRedirectUri, testGrantType, gomock.Any()).Return(testError)

	mockTokenService := tokenMock.NewMockService(ctrl)

	h := &Handler{
		tokens:    mockTokenService,
		clients:   mockProvider,
		validator: mockValidator,
	}

	testBody := url.Values{
		"client_id":     {testClientId},
		"client_secret": {testClientSecret},
		"grant_type":    {testGrantType},
		"redirect_uri":  {testRedirectUri},
		"scope":         {testScopes},
	}

	resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: testBody.Encode(),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(map[string]string{"error": testError.Error()})
	assert.Equal(t, string(bytes), resp.Body)
}

func TestHandler_GivenTokenGenerationFails_ReturnsInternalServerError(t *testing.T) {
	testClientId := "3247023"
	testClientSecret := "2934uldnf"
	testGrantType := "code"
	testRedirectUri := "http://test.io"
	testScopes := "openid"

	testClient := &dal.Client{
		ID:           testClientId,
		RedirectUris: []string{"http://test.io"},
		Scopes:       []string{"openid"},
		Secrets:      []string{"my secret"},
		GrantTypes:   []string{"code"},
	}
	testError := errors.New("invalid client")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := dalMock.NewMockClientProvider(ctrl)
	mockProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockValidator.EXPECT().ValidateRequest(testClient, testClientSecret, testRedirectUri, testGrantType, gomock.Any()).Return(nil)

	mockTokenService := tokenMock.NewMockService(ctrl)
	mockTokenService.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, testError)

	h := &Handler{
		tokens:    mockTokenService,
		clients:   mockProvider,
		validator: mockValidator,
	}

	testBody := url.Values{
		"client_id":     {testClientId},
		"client_secret": {testClientSecret},
		"grant_type":    {testGrantType},
		"redirect_uri":  {testRedirectUri},
		"scope":         {testScopes},
	}

	resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: testBody.Encode(),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])

	bytes, _ := json.Marshal(map[string]string{"error": testError.Error()})
	assert.Equal(t, string(bytes), resp.Body)
}
