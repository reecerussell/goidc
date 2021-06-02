package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	hashMock "github.com/reecerussell/adaptive-password-hasher/mock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/goidc/dal"
	dalMock "github.com/reecerussell/goidc/dal/mock"
	valMock "github.com/reecerussell/goidc/validator/mock"
)

func TestHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEmail := "myemail"
	testPassword := "myPass"

	mockValidator := valMock.NewMockUserValidator(ctrl)
	mockValidator.EXPECT().ValidateUser(testEmail, testPassword).Return(nil)

	mockProvider := dalMock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(nil, dal.ErrUserNotFound)

	mockHasher := hashMock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash([]byte(testPassword)).Return([]byte("234023u4023"))

	mockService := dalMock.NewMockUserService(ctrl)
	mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	h := &Handler{
		uv:  mockValidator,
		up:  mockProvider,
		hsr: mockHasher,
		us:  mockService,
	}

	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{
			"email": "%s",
			"password": "%s"
		}`, testEmail, testPassword),
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data map[string]string
	json.Unmarshal([]byte(resp.Body), &data)

	_, ok := data["id"]
	assert.True(t, ok)
}

func TestHandle_GivenInvalidHTTPMethod_ReturnsMethodNotAllowed(t *testing.T) {
	h := &Handler{}

	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet, // invalid method
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	var data map[string]string
	json.Unmarshal([]byte(resp.Body), &data)
	assert.Equal(t, "method not allowed", data["error"])
}

func TestHandle_GivenInvalidContentType_ReturnsBadRequests(t *testing.T) {
	h := &Handler{}

	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/xml", // invalid type
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var data map[string]string
	json.Unmarshal([]byte(resp.Body), &data)
	assert.Equal(t, "invalid content type", data["error"])
}

func TestHandle_GivenInvalidUserData_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEmail := "myemail"
	testPassword := "myPass"
	testError := errors.New("test error")

	mockValidator := valMock.NewMockUserValidator(ctrl)
	mockValidator.EXPECT().ValidateUser(testEmail, testPassword).Return(testError)

	h := &Handler{
		uv: mockValidator,
	}

	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{
			"email": "%s",
			"password": "%s"
		}`, testEmail, testPassword),
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var data map[string]string
	json.Unmarshal([]byte(resp.Body), &data)
	assert.Equal(t, testError.Error(), data["error"])
}

func TestHandle_WhereUserAlreadyExists_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEmail := "myemail"
	testPassword := "myPass"

	mockValidator := valMock.NewMockUserValidator(ctrl)
	mockValidator.EXPECT().ValidateUser(testEmail, testPassword).Return(nil)

	mockProvider := dalMock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(nil, nil) // nil indicates user exists

	h := &Handler{
		uv: mockValidator,
		up: mockProvider,
	}

	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{
			"email": "%s",
			"password": "%s"
		}`, testEmail, testPassword),
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var data map[string]string
	json.Unmarshal([]byte(resp.Body), &data)
	assert.Equal(t, "user already exists", data["error"])
}

func TestHandle_WhereUserExistsCheckFails_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEmail := "myemail"
	testPassword := "myPass"
	testError := errors.New("test error")

	mockValidator := valMock.NewMockUserValidator(ctrl)
	mockValidator.EXPECT().ValidateUser(testEmail, testPassword).Return(nil)

	mockProvider := dalMock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(nil, testError)

	h := &Handler{
		uv: mockValidator,
		up: mockProvider,
	}

	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{
			"email": "%s",
			"password": "%s"
		}`, testEmail, testPassword),
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var data map[string]string
	json.Unmarshal([]byte(resp.Body), &data)
	assert.Equal(t, testError.Error(), data["error"])
}

func TestHandle_WhereUserCreationFails_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEmail := "myemail"
	testPassword := "myPass"
	testError := errors.New("test error")

	mockValidator := valMock.NewMockUserValidator(ctrl)
	mockValidator.EXPECT().ValidateUser(testEmail, testPassword).Return(nil)

	mockProvider := dalMock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(nil, dal.ErrUserNotFound)

	mockHasher := hashMock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash([]byte(testPassword)).Return([]byte("234023u4023"))

	mockService := dalMock.NewMockUserService(ctrl)
	mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(testError)

	h := &Handler{
		uv:  mockValidator,
		up:  mockProvider,
		hsr: mockHasher,
		us:  mockService,
	}

	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{
			"email": "%s",
			"password": "%s"
		}`, testEmail, testPassword),
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var data map[string]string
	json.Unmarshal([]byte(resp.Body), &data)
	assert.Equal(t, testError.Error(), data["error"])
}
