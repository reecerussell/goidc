package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/awstesting/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/goidc/dal"
	dalMock "github.com/reecerussell/goidc/dal/mock"
	"github.com/reecerussell/goidc/token"
	tokenMock "github.com/reecerussell/goidc/token/mock"
	"github.com/reecerussell/goidc/validator"
	valMock "github.com/reecerussell/goidc/validator/mock"
)

func TestHandler_GivenIdTokenAndTokenTypes_ReturnsRedirectWithTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testEmail := "my@email.com"
	testPassword := "myPassword1"
	testState := "2374923740234"
	testNonce := "2304820340lskfle"
	testClient := &dal.Client{}
	testUser := &dal.User{ID: "testUserId", PasswordHash: "328y9ewhdk"}
	testIdToken := "239y4o24o234"
	testAccessToken := "1ohweory9843"

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockUserProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(testUser, nil)

	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockClientValidator := valMock.NewMockClientValidator(ctrl)
	mockClientValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(nil)

	mockUserValidator := valMock.NewMockUserValidator(ctrl)
	mockUserValidator.EXPECT().ValidatePassword(testUser, testPassword).Return(nil)

	mockTokenService := tokenMock.NewMockService(ctrl)
	mockTokenService.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&token.Token{AccessToken: testAccessToken, TokenType: "Bearer", Expires: 3600}, nil).Times(1)
	mockTokenService.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&token.Token{AccessToken: testIdToken}, nil)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		userVal:   mockUserValidator,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockClientValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"],
			"email": "%s",
			"password": "%s",
			"responseType": "id_token token",
			"state": "%s",
			"nonce": "%s"
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`), testEmail, testPassword, testState, testNonce),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	redirectUri, err := url.Parse(data["redirectUri"].(string))
	assert.NoError(t, err)

	if !strings.HasPrefix(data["redirectUri"].(string), testRedirectUri) {
		t.Errorf("redirect uri should be '%s'", testRedirectUri)
	}

	queryValues := redirectUri.Query()

	assert.Equal(t, testIdToken, queryValues.Get("id_token"))
	assert.Equal(t, testAccessToken, queryValues.Get("access_token"))
	assert.Equal(t, "Bearer", queryValues.Get("token_type"))
	assert.Equal(t, "3600", queryValues.Get("expires_in"))
}

func TestHandler_GivenInvalidHTTPMethod_ReturnsMethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:    mock.Session,
		users:   mockUserProvider,
		tokens:  mockTokenService,
		clients: mockClientProvider,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, "method not allowed", data["error"])
}

func TestHandler_GivenInvalidContentType_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:    mock.Session,
		users:   mockUserProvider,
		tokens:  mockTokenService,
		clients: mockClientProvider,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, "invalid content type", data["error"])
}

func TestHandler_GivenInvalidContent_ReturnsMethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:    mock.Session,
		users:   mockUserProvider,
		tokens:  mockTokenService,
		clients: mockClientProvider,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: "clientId=2394&scopes=hello world", // expecting JSON
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, "invalid content", data["error"])
}

func TestHandler_GivenInvalidClient_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(nil, dal.ErrClientNotFound)

	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:    mock.Session,
		users:   mockUserProvider,
		tokens:  mockTokenService,
		clients: mockClientProvider,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s"	
		}`, testClientId),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, "invalid client id", data["error"])
}

func TestHandler_WhereClientProviderFails_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testError := errors.New("test error")

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(nil, testError)

	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:    mock.Session,
		users:   mockUserProvider,
		tokens:  mockTokenService,
		clients: mockClientProvider,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s"	
		}`, testClientId),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, testError.Error(), data["error"])
}

func TestHandler_WhereClientInfoIsInvalid_ReturnBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testClient := &dal.Client{}
	testError := errors.New("test error")

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(testError)

	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"]
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`)),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, testError.Error(), data["error"])
}

func TestHandler_GivenInvalidEmail_ReturnBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testEmail := "my@email.com"
	testClient := &dal.Client{}

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockUserProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(nil, dal.ErrUserNotFound)

	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(nil)

	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"],
			"email": "%s"
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`), testEmail),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, errInvalidCredentials.Error(), data["error"])
}

func TestHandler_WhereUserProviderFails_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testEmail := "my@email.com"
	testClient := &dal.Client{}
	testError := errors.New("test error")

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockUserProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(nil, testError)

	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockValidator := valMock.NewMockClientValidator(ctrl)
	mockValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(nil)

	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"],
			"email": "%s"
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`), testEmail),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, testError.Error(), data["error"])
}

func TestHandler_GivenInvalidPassword_ReturnBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testEmail := "my@email.com"
	testPassword := "myPassword1"
	testClient := &dal.Client{}
	testUser := &dal.User{PasswordHash: "328y9ewhdk"}

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockUserProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(testUser, nil)

	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockClientValidator := valMock.NewMockClientValidator(ctrl)
	mockClientValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(nil)

	mockUserValidator := valMock.NewMockUserValidator(ctrl)
	mockUserValidator.EXPECT().ValidatePassword(testUser, testPassword).Return(validator.ErrInvalidPassword)

	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		userVal:   mockUserValidator,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockClientValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"],
			"email": "%s",
			"password": "%s"
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`), testEmail, testPassword),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, errInvalidCredentials.Error(), data["error"])
}

func TestHandler_WherePasswordValidationFails_ReturnInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testEmail := "my@email.com"
	testPassword := "myPassword1"
	testClient := &dal.Client{}
	testUser := &dal.User{PasswordHash: "328y9ewhdk"}
	testError := errors.New("error")

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockUserProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(testUser, nil)

	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockClientValidator := valMock.NewMockClientValidator(ctrl)
	mockClientValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(nil)

	mockUserValidator := valMock.NewMockUserValidator(ctrl)
	mockUserValidator.EXPECT().ValidatePassword(testUser, testPassword).Return(testError)

	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		userVal:   mockUserValidator,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockClientValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"],
			"email": "%s",
			"password": "%s"
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`), testEmail, testPassword),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, testError.Error(), data["error"])
}

func TestHandler_WhereIdTokenFails_ReturnInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testEmail := "my@email.com"
	testPassword := "myPassword1"
	testState := "2374923740234"
	testNonce := "2304820340lskfle"
	testClient := &dal.Client{}
	testUser := &dal.User{ID: "testUserId", PasswordHash: "328y9ewhdk"}
	testAccessToken := "1ohweory9843"
	testError := errors.New("error")

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockUserProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(testUser, nil)

	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockClientValidator := valMock.NewMockClientValidator(ctrl)
	mockClientValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(nil)

	mockUserValidator := valMock.NewMockUserValidator(ctrl)
	mockUserValidator.EXPECT().ValidatePassword(testUser, testPassword).Return(nil)

	mockTokenService := tokenMock.NewMockService(ctrl)
	mockTokenService.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&token.Token{AccessToken: testAccessToken, TokenType: "Bearer", Expires: 3600}, nil).Times(1)
	mockTokenService.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, testError)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		userVal:   mockUserValidator,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockClientValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"],
			"email": "%s",
			"password": "%s",
			"responseType": "id_token token",
			"state": "%s",
			"nonce": "%s"
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`), testEmail, testPassword, testState, testNonce),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, testError.Error(), data["error"])
}

func TestHandler_WhereAccessTokenFails_ReturnInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testEmail := "my@email.com"
	testPassword := "myPassword1"
	testState := "2374923740234"
	testNonce := "2304820340lskfle"
	testClient := &dal.Client{}
	testUser := &dal.User{ID: "testUserId", PasswordHash: "328y9ewhdk"}
	testError := errors.New("error")

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockUserProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(testUser, nil)

	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockClientValidator := valMock.NewMockClientValidator(ctrl)
	mockClientValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(nil)

	mockUserValidator := valMock.NewMockUserValidator(ctrl)
	mockUserValidator.EXPECT().ValidatePassword(testUser, testPassword).Return(nil)

	mockTokenService := tokenMock.NewMockService(ctrl)
	mockTokenService.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, testError).Times(1)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		userVal:   mockUserValidator,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockClientValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"],
			"email": "%s",
			"password": "%s",
			"responseType": "id_token token",
			"state": "%s",
			"nonce": "%s"
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`), testEmail, testPassword, testState, testNonce),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, testError.Error(), data["error"])
}

func TestHandler_GivenUnsupportedResponseType_ReturnBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testClientId := "23493234"
	testRedirectUri := "http://localhost:8080"
	testScopes := []string{"openid", "test"}
	testEmail := "my@email.com"
	testPassword := "myPassword1"
	testState := "2374923740234"
	testNonce := "2304820340lskfle"
	testClient := &dal.Client{}
	testUser := &dal.User{ID: "testUserId", PasswordHash: "328y9ewhdk"}

	mockUserProvider := dalMock.NewMockUserProvider(ctrl)
	mockUserProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(testUser, nil)

	mockClientProvider := dalMock.NewMockClientProvider(ctrl)
	mockClientProvider.EXPECT().Get(gomock.Any(), testClientId).Return(testClient, nil)

	mockClientValidator := valMock.NewMockClientValidator(ctrl)
	mockClientValidator.EXPECT().ValidateLoginRequest(testClient, testRedirectUri, testScopes).Return(nil)

	mockUserValidator := valMock.NewMockUserValidator(ctrl)
	mockUserValidator.EXPECT().ValidatePassword(testUser, testPassword).Return(nil)

	mockTokenService := tokenMock.NewMockService(ctrl)

	handler := &Handler{
		sess:      mock.Session,
		users:     mockUserProvider,
		userVal:   mockUserValidator,
		tokens:    mockTokenService,
		clients:   mockClientProvider,
		clientVal: mockClientValidator,
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StageVariables: map[string]string{
			"JWT_KEY_ID": "key id",
		},
		Body: fmt.Sprintf(`{
			"clientId": "%s",
			"redirectUri": "%s",
			"scopes": ["%s"],
			"email": "%s",
			"password": "%s",
			"responseType": "nnot a supported response type",
			"state": "%s",
			"nonce": "%s"
		}`, testClientId, testRedirectUri, strings.Join(testScopes, `","`), testEmail, testPassword, testState, testNonce),
	}

	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.False(t, resp.IsBase64Encoded)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &data)

	assert.Equal(t, errUnsupportedResponseType.Error(), data["error"])
}
