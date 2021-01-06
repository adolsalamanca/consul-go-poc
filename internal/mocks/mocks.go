// Code generated by MockGen. DO NOT EDIT.
// Source: internal/logic.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	api "github.com/hashicorp/consul/api"
	reflect "reflect"
)

// MockClienter is a mock of Clienter interface
type MockClienter struct {
	ctrl     *gomock.Controller
	recorder *MockClienterMockRecorder
}

// MockClienterMockRecorder is the mock recorder for MockClienter
type MockClienterMockRecorder struct {
	mock *MockClienter
}

// NewMockClienter creates a new mock instance
func NewMockClienter(ctrl *gomock.Controller) *MockClienter {
	mock := &MockClienter{ctrl: ctrl}
	mock.recorder = &MockClienterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClienter) EXPECT() *MockClienterMockRecorder {
	return m.recorder
}

// Agent mocks base method
func (m *MockClienter) Agent() *api.Agent {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Agent")
	ret0, _ := ret[0].(*api.Agent)
	return ret0
}

// Agent indicates an expected call of Agent
func (mr *MockClienterMockRecorder) Agent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Agent", reflect.TypeOf((*MockClienter)(nil).Agent))
}

// Session mocks base method
func (m *MockClienter) Session() *api.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Session")
	ret0, _ := ret[0].(*api.Session)
	return ret0
}

// Session indicates an expected call of Session
func (mr *MockClienterMockRecorder) Session() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Session", reflect.TypeOf((*MockClienter)(nil).Session))
}

// KV mocks base method
func (m *MockClienter) KV() *api.KV {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KV")
	ret0, _ := ret[0].(*api.KV)
	return ret0
}

// KV indicates an expected call of KV
func (mr *MockClienterMockRecorder) KV() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KV", reflect.TypeOf((*MockClienter)(nil).KV))
}

// MockKeyValuer is a mock of KeyValuer interface
type MockKeyValuer struct {
	ctrl     *gomock.Controller
	recorder *MockKeyValuerMockRecorder
}

// MockKeyValuerMockRecorder is the mock recorder for MockKeyValuer
type MockKeyValuerMockRecorder struct {
	mock *MockKeyValuer
}

// NewMockKeyValuer creates a new mock instance
func NewMockKeyValuer(ctrl *gomock.Controller) *MockKeyValuer {
	mock := &MockKeyValuer{ctrl: ctrl}
	mock.recorder = &MockKeyValuerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeyValuer) EXPECT() *MockKeyValuerMockRecorder {
	return m.recorder
}

// Acquire mocks base method
func (m *MockKeyValuer) Acquire(p *api.KVPair, q *api.WriteOptions) (bool, *api.WriteMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Acquire", p, q)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*api.WriteMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Acquire indicates an expected call of Acquire
func (mr *MockKeyValuerMockRecorder) Acquire(p, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Acquire", reflect.TypeOf((*MockKeyValuer)(nil).Acquire), p, q)
}

// Release mocks base method
func (m *MockKeyValuer) Release(p *api.KVPair, q *api.WriteOptions) (bool, *api.WriteMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Release", p, q)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*api.WriteMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Release indicates an expected call of Release
func (mr *MockKeyValuerMockRecorder) Release(p, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*MockKeyValuer)(nil).Release), p, q)
}