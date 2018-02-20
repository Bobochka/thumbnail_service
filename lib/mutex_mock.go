// Code generated by MockGen. DO NOT EDIT.
// Source: mutex.go

// Package lib is a generated GoMock package.
package lib

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMutex is a mock of Mutex interface
type MockMutex struct {
	ctrl     *gomock.Controller
	recorder *MockMutexMockRecorder
}

// MockMutexMockRecorder is the mock recorder for MockMutex
type MockMutexMockRecorder struct {
	mock *MockMutex
}

// NewMockMutex creates a new mock instance
func NewMockMutex(ctrl *gomock.Controller) *MockMutex {
	mock := &MockMutex{ctrl: ctrl}
	mock.recorder = &MockMutexMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMutex) EXPECT() *MockMutexMockRecorder {
	return m.recorder
}

// Lock mocks base method
func (m *MockMutex) Lock() error {
	ret := m.ctrl.Call(m, "Lock")
	ret0, _ := ret[0].(error)
	return ret0
}

// Lock indicates an expected call of Lock
func (mr *MockMutexMockRecorder) Lock() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockMutex)(nil).Lock))
}

// Unlock mocks base method
func (m *MockMutex) Unlock() bool {
	ret := m.ctrl.Call(m, "Unlock")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Unlock indicates an expected call of Unlock
func (mr *MockMutexMockRecorder) Unlock() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlock", reflect.TypeOf((*MockMutex)(nil).Unlock))
}