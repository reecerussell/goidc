// Code generated by MockGen. DO NOT EDIT.
// Source: ../client_provider.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	dal "github.com/reecerussell/goidc/dal"
	reflect "reflect"
)

// MockClientProvider is a mock of ClientProvider interface.
type MockClientProvider struct {
	ctrl     *gomock.Controller
	recorder *MockClientProviderMockRecorder
}

// MockClientProviderMockRecorder is the mock recorder for MockClientProvider.
type MockClientProviderMockRecorder struct {
	mock *MockClientProvider
}

// NewMockClientProvider creates a new mock instance.
func NewMockClientProvider(ctrl *gomock.Controller) *MockClientProvider {
	mock := &MockClientProvider{ctrl: ctrl}
	mock.recorder = &MockClientProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientProvider) EXPECT() *MockClientProviderMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockClientProvider) Get(ctx context.Context, id string) (*dal.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*dal.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockClientProviderMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClientProvider)(nil).Get), ctx, id)
}
