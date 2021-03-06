// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	lib "github.com/Bobochka/thumbnail_service/lib"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockStore) Get(key string) []byte {
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockStoreMockRecorder) Get(key interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), key)
}

// Set mocks base method
func (m *MockStore) Set(key string, data []byte) error {
	ret := m.ctrl.Call(m, "Set", key, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockStoreMockRecorder) Set(key, data interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStore)(nil).Set), key, data)
}

// MockTransformation is a mock of Transformation interface
type MockTransformation struct {
	ctrl     *gomock.Controller
	recorder *MockTransformationMockRecorder
}

// MockTransformationMockRecorder is the mock recorder for MockTransformation
type MockTransformationMockRecorder struct {
	mock *MockTransformation
}

// NewMockTransformation creates a new mock instance
func NewMockTransformation(ctrl *gomock.Controller) *MockTransformation {
	mock := &MockTransformation{ctrl: ctrl}
	mock.recorder = &MockTransformationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTransformation) EXPECT() *MockTransformationMockRecorder {
	return m.recorder
}

// Fingerprint mocks base method
func (m *MockTransformation) Fingerprint(data []byte) string {
	ret := m.ctrl.Call(m, "Fingerprint", data)
	ret0, _ := ret[0].(string)
	return ret0
}

// Fingerprint indicates an expected call of Fingerprint
func (mr *MockTransformationMockRecorder) Fingerprint(data interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fingerprint", reflect.TypeOf((*MockTransformation)(nil).Fingerprint), data)
}

// Perform mocks base method
func (m *MockTransformation) Perform(data []byte) ([]byte, error) {
	ret := m.ctrl.Call(m, "Perform", data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Perform indicates an expected call of Perform
func (mr *MockTransformationMockRecorder) Perform(data interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Perform", reflect.TypeOf((*MockTransformation)(nil).Perform), data)
}

// MockDownloader is a mock of Downloader interface
type MockDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockDownloaderMockRecorder
}

// MockDownloaderMockRecorder is the mock recorder for MockDownloader
type MockDownloaderMockRecorder struct {
	mock *MockDownloader
}

// NewMockDownloader creates a new mock instance
func NewMockDownloader(ctrl *gomock.Controller) *MockDownloader {
	mock := &MockDownloader{ctrl: ctrl}
	mock.recorder = &MockDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDownloader) EXPECT() *MockDownloaderMockRecorder {
	return m.recorder
}

// Download mocks base method
func (m *MockDownloader) Download(url string) ([]byte, error) {
	ret := m.ctrl.Call(m, "Download", url)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download
func (mr *MockDownloaderMockRecorder) Download(url interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockDownloader)(nil).Download), url)
}

// MockLocker is a mock of Locker interface
type MockLocker struct {
	ctrl     *gomock.Controller
	recorder *MockLockerMockRecorder
}

// MockLockerMockRecorder is the mock recorder for MockLocker
type MockLockerMockRecorder struct {
	mock *MockLocker
}

// NewMockLocker creates a new mock instance
func NewMockLocker(ctrl *gomock.Controller) *MockLocker {
	mock := &MockLocker{ctrl: ctrl}
	mock.recorder = &MockLockerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLocker) EXPECT() *MockLockerMockRecorder {
	return m.recorder
}

// NewMutex mocks base method
func (m *MockLocker) NewMutex(name string) lib.Mutex {
	ret := m.ctrl.Call(m, "NewMutex", name)
	ret0, _ := ret[0].(lib.Mutex)
	return ret0
}

// NewMutex indicates an expected call of NewMutex
func (mr *MockLockerMockRecorder) NewMutex(name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMutex", reflect.TypeOf((*MockLocker)(nil).NewMutex), name)
}
