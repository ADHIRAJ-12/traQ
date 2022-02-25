// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/traPtitech/traQ/model"
	ioext "github.com/traPtitech/traQ/utils/ioext"
)

// MockFileStorage is a mock of FileStorage interface.
type MockFileStorage struct {
	ctrl     *gomock.Controller
	recorder *MockFileStorageMockRecorder
}

// MockFileStorageMockRecorder is the mock recorder for MockFileStorage.
type MockFileStorageMockRecorder struct {
	mock *MockFileStorage
}

// NewMockFileStorage creates a new mock instance.
func NewMockFileStorage(ctrl *gomock.Controller) *MockFileStorage {
	mock := &MockFileStorage{ctrl: ctrl}
	mock.recorder = &MockFileStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileStorage) EXPECT() *MockFileStorageMockRecorder {
	return m.recorder
}

// DeleteByKey mocks base method.
func (m *MockFileStorage) DeleteByKey(key string, fileType model.FileType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByKey", key, fileType)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByKey indicates an expected call of DeleteByKey.
func (mr *MockFileStorageMockRecorder) DeleteByKey(key, fileType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByKey", reflect.TypeOf((*MockFileStorage)(nil).DeleteByKey), key, fileType)
}

// GenerateAccessURL mocks base method.
func (m *MockFileStorage) GenerateAccessURL(key string, fileType model.FileType) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAccessURL", key, fileType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAccessURL indicates an expected call of GenerateAccessURL.
func (mr *MockFileStorageMockRecorder) GenerateAccessURL(key, fileType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAccessURL", reflect.TypeOf((*MockFileStorage)(nil).GenerateAccessURL), key, fileType)
}

// OpenFileByKey mocks base method.
func (m *MockFileStorage) OpenFileByKey(key string, fileType model.FileType) (ioext.ReadSeekCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenFileByKey", key, fileType)
	ret0, _ := ret[0].(ioext.ReadSeekCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenFileByKey indicates an expected call of OpenFileByKey.
func (mr *MockFileStorageMockRecorder) OpenFileByKey(key, fileType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFileByKey", reflect.TypeOf((*MockFileStorage)(nil).OpenFileByKey), key, fileType)
}

// SaveByKey mocks base method.
func (m *MockFileStorage) SaveByKey(src io.Reader, key, name, contentType string, fileType model.FileType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveByKey", src, key, name, contentType, fileType)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveByKey indicates an expected call of SaveByKey.
func (mr *MockFileStorageMockRecorder) SaveByKey(src, key, name, contentType, fileType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveByKey", reflect.TypeOf((*MockFileStorage)(nil).SaveByKey), src, key, name, contentType, fileType)
}
