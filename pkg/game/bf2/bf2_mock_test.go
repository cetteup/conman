// Code generated by MockGen. DO NOT EDIT.
// Source: ../common.go

package bf2

import (
	reflect "reflect"

	config "github.com/cetteup/conman/pkg/config"
	handler "github.com/cetteup/conman/pkg/handler"
	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// BuildBasePath mocks base method.
func (m *MockHandler) BuildBasePath(game handler.Game) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildBasePath", game)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildBasePath indicates an expected call of BuildBasePath.
func (mr *MockHandlerMockRecorder) BuildBasePath(game interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildBasePath", reflect.TypeOf((*MockHandler)(nil).BuildBasePath), game)
}

// GetProfileKeys mocks base method.
func (m *MockHandler) GetProfileKeys(game handler.Game) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileKeys", game)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileKeys indicates an expected call of GetProfileKeys.
func (mr *MockHandlerMockRecorder) GetProfileKeys(game interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileKeys", reflect.TypeOf((*MockHandler)(nil).GetProfileKeys), game)
}

// ReadConfigFile mocks base method.
func (m *MockHandler) ReadConfigFile(path string) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadConfigFile", path)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadConfigFile indicates an expected call of ReadConfigFile.
func (mr *MockHandlerMockRecorder) ReadConfigFile(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConfigFile", reflect.TypeOf((*MockHandler)(nil).ReadConfigFile), path)
}

// ReadGlobalConfig mocks base method.
func (m *MockHandler) ReadGlobalConfig(game handler.Game) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadGlobalConfig", game)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadGlobalConfig indicates an expected call of ReadGlobalConfig.
func (mr *MockHandlerMockRecorder) ReadGlobalConfig(game interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadGlobalConfig", reflect.TypeOf((*MockHandler)(nil).ReadGlobalConfig), game)
}

// ReadProfileConfig mocks base method.
func (m *MockHandler) ReadProfileConfig(game handler.Game, profile string) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadProfileConfig", game, profile)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadProfileConfig indicates an expected call of ReadProfileConfig.
func (mr *MockHandlerMockRecorder) ReadProfileConfig(game, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadProfileConfig", reflect.TypeOf((*MockHandler)(nil).ReadProfileConfig), game, profile)
}
