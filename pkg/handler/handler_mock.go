// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package handler is a generated GoMock package.
package handler

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockfileRepository is a mock of fileRepository interface.
type MockfileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockfileRepositoryMockRecorder
}

// MockfileRepositoryMockRecorder is the mock recorder for MockfileRepository.
type MockfileRepositoryMockRecorder struct {
	mock *MockfileRepository
}

// NewMockfileRepository creates a new mock instance.
func NewMockfileRepository(ctrl *gomock.Controller) *MockfileRepository {
	mock := &MockfileRepository{ctrl: ctrl}
	mock.recorder = &MockfileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockfileRepository) EXPECT() *MockfileRepositoryMockRecorder {
	return m.recorder
}

// ReadFile mocks base method.
func (m *MockfileRepository) ReadFile(path string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", path)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile.
func (mr *MockfileRepositoryMockRecorder) ReadFile(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockfileRepository)(nil).ReadFile), path)
}
