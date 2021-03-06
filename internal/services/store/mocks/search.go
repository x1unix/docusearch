// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/x1unix/docusearch/internal/services/search (interfaces: Provider)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// AddDocumentRef mocks base method.
func (m *MockProvider) AddDocumentRef(arg0 context.Context, arg1 string, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDocumentRef", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDocumentRef indicates an expected call of AddDocumentRef.
func (mr *MockProviderMockRecorder) AddDocumentRef(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDocumentRef", reflect.TypeOf((*MockProvider)(nil).AddDocumentRef), arg0, arg1, arg2)
}

// RemoveDocumentRef mocks base method.
func (m *MockProvider) RemoveDocumentRef(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDocumentRef", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDocumentRef indicates an expected call of RemoveDocumentRef.
func (mr *MockProviderMockRecorder) RemoveDocumentRef(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDocumentRef", reflect.TypeOf((*MockProvider)(nil).RemoveDocumentRef), arg0, arg1)
}

// SearchDocumentsByWord mocks base method.
func (m *MockProvider) SearchDocumentsByWord(arg0 context.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchDocumentsByWord", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchDocumentsByWord indicates an expected call of SearchDocumentsByWord.
func (mr *MockProviderMockRecorder) SearchDocumentsByWord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchDocumentsByWord", reflect.TypeOf((*MockProvider)(nil).SearchDocumentsByWord), arg0, arg1)
}
