// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ndreyserg/ushort/internal/app/storage (interfaces: Storage)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/ndreyserg/ushort/internal/app/models"
	storage "github.com/ndreyserg/ushort/internal/app/storage"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockStorage) Check(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockStorageMockRecorder) Check(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockStorage)(nil).Check), arg0)
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// Get mocks base method.
func (m *MockStorage) Get(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStorageMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStorage)(nil).Get), arg0, arg1)
}

// GetUserUrls mocks base method.
func (m *MockStorage) GetUserUrls(arg0 context.Context, arg1 string) ([]storage.StorageItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserUrls", arg0, arg1)
	ret0, _ := ret[0].([]storage.StorageItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserUrls indicates an expected call of GetUserUrls.
func (mr *MockStorageMockRecorder) GetUserUrls(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserUrls", reflect.TypeOf((*MockStorage)(nil).GetUserUrls), arg0, arg1)
}

// Set mocks base method.
func (m *MockStorage) Set(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockStorageMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStorage)(nil).Set), arg0, arg1, arg2)
}

// SetBatch mocks base method.
func (m *MockStorage) SetBatch(arg0 context.Context, arg1 models.BatchRequest, arg2 string) (models.BatchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBatch", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.BatchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetBatch indicates an expected call of SetBatch.
func (mr *MockStorageMockRecorder) SetBatch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBatch", reflect.TypeOf((*MockStorage)(nil).SetBatch), arg0, arg1, arg2)
}
