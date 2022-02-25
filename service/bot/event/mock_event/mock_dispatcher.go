// Code generated by MockGen. DO NOT EDIT.
// Source: dispatcher.go

// Package mock_event is a generated GoMock package.
package mock_event

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/traPtitech/traQ/model"
)

// MockDispatcher is a mock of Dispatcher interface.
type MockDispatcher struct {
	ctrl     *gomock.Controller
	recorder *MockDispatcherMockRecorder
}

// MockDispatcherMockRecorder is the mock recorder for MockDispatcher.
type MockDispatcherMockRecorder struct {
	mock *MockDispatcher
}

// NewMockDispatcher creates a new mock instance.
func NewMockDispatcher(ctrl *gomock.Controller) *MockDispatcher {
	mock := &MockDispatcher{ctrl: ctrl}
	mock.recorder = &MockDispatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDispatcher) EXPECT() *MockDispatcherMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockDispatcher) Send(b *model.Bot, event model.BotEventType, body []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", b, event, body)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDispatcherMockRecorder) Send(b, event, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDispatcher)(nil).Send), b, event, body)
}
