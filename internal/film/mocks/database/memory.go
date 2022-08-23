// Code generated by MockGen. DO NOT EDIT.
// Source: test/internal/film/repository (interfaces: MemoryDB)

// Package mocks_database is a generated GoMock package.
package mocks_database

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockMemoryDB is a mock of MemoryDB interface.
type MockMemoryDB struct {
	ctrl     *gomock.Controller
	recorder *MockMemoryDBMockRecorder
}

// MockMemoryDBMockRecorder is the mock recorder for MockMemoryDB.
type MockMemoryDBMockRecorder struct {
	mock *MockMemoryDB
}

// NewMockMemoryDB creates a new mock instance.
func NewMockMemoryDB(ctrl *gomock.Controller) *MockMemoryDB {
	mock := &MockMemoryDB{ctrl: ctrl}
	mock.recorder = &MockMemoryDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMemoryDB) EXPECT() *MockMemoryDBMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockMemoryDB) Get(arg0 string) (interface{}, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMemoryDBMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMemoryDB)(nil).Get), arg0)
}

// ItemsTTL mocks base method.
func (m *MockMemoryDB) ItemsTTL() map[string]time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemsTTL")
	ret0, _ := ret[0].(map[string]time.Duration)
	return ret0
}

// ItemsTTL indicates an expected call of ItemsTTL.
func (mr *MockMemoryDBMockRecorder) ItemsTTL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemsTTL", reflect.TypeOf((*MockMemoryDB)(nil).ItemsTTL))
}

// Set mocks base method.
func (m *MockMemoryDB) Set(arg0 string, arg1 interface{}, arg2 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1, arg2)
}

// Set indicates an expected call of Set.
func (mr *MockMemoryDBMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockMemoryDB)(nil).Set), arg0, arg1, arg2)
}