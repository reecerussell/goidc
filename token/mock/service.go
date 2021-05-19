// Code generated by MockGen. DO NOT EDIT.
// Source: ../service.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	token "github.com/reecerussell/goidc/token"
	reflect "reflect"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockService) GenerateToken(claims map[string]interface{}, expirySeconds int64, audience string) (*token.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", claims, expirySeconds, audience)
	ret0, _ := ret[0].(*token.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockServiceMockRecorder) GenerateToken(claims, expirySeconds, audience interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockService)(nil).GenerateToken), claims, expirySeconds, audience)
}
