// Code generated by MockGen. DO NOT EDIT.
// Source: gx/ipfs/QmYJtCabf3prS3HKQUGgqDLVxvbT9iDx5mfeVfhtCcJxxE/go-libp2p-discovery (interfaces: Advertiser,Discoverer,Discovery)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	go_libp2p_peerstore "gx/ipfs/QmPiemjiKBC9VA7vZF82m4x1oygtg2c2YVqag8PX7dN1BD/go-libp2p-peerstore"
	go_libp2p_discovery "gx/ipfs/QmemYsfqwAbyvqwFiApk1GfLKhDkMm8ZQK6fCvzDbaRNyX/go-libp2p-discovery"
	reflect "reflect"
	time "time"
)

// MockAdvertiser is a mock of Advertiser interface
type MockAdvertiser struct {
	ctrl     *gomock.Controller
	recorder *MockAdvertiserMockRecorder
}

// MockAdvertiserMockRecorder is the mock recorder for MockAdvertiser
type MockAdvertiserMockRecorder struct {
	mock *MockAdvertiser
}

// NewMockAdvertiser creates a new mock instance
func NewMockAdvertiser(ctrl *gomock.Controller) *MockAdvertiser {
	mock := &MockAdvertiser{ctrl: ctrl}
	mock.recorder = &MockAdvertiserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAdvertiser) EXPECT() *MockAdvertiserMockRecorder {
	return m.recorder
}

// Advertise mocks base method
func (m *MockAdvertiser) Advertise(arg0 context.Context, arg1 string, arg2 ...go_libp2p_discovery.Option) (time.Duration, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Advertise", varargs...)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Advertise indicates an expected call of Advertise
func (mr *MockAdvertiserMockRecorder) Advertise(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Advertise", reflect.TypeOf((*MockAdvertiser)(nil).Advertise), varargs...)
}

// MockDiscoverer is a mock of Discoverer interface
type MockDiscoverer struct {
	ctrl     *gomock.Controller
	recorder *MockDiscovererMockRecorder
}

// MockDiscovererMockRecorder is the mock recorder for MockDiscoverer
type MockDiscovererMockRecorder struct {
	mock *MockDiscoverer
}

// NewMockDiscoverer creates a new mock instance
func NewMockDiscoverer(ctrl *gomock.Controller) *MockDiscoverer {
	mock := &MockDiscoverer{ctrl: ctrl}
	mock.recorder = &MockDiscovererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDiscoverer) EXPECT() *MockDiscovererMockRecorder {
	return m.recorder
}

// FindPeers mocks base method
func (m *MockDiscoverer) FindPeers(arg0 context.Context, arg1 string, arg2 ...go_libp2p_discovery.Option) (<-chan go_libp2p_peerstore.PeerInfo, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindPeers", varargs...)
	ret0, _ := ret[0].(<-chan go_libp2p_peerstore.PeerInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPeers indicates an expected call of FindPeers
func (mr *MockDiscovererMockRecorder) FindPeers(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPeers", reflect.TypeOf((*MockDiscoverer)(nil).FindPeers), varargs...)
}

// MockDiscovery is a mock of Discovery interface
type MockDiscovery struct {
	ctrl     *gomock.Controller
	recorder *MockDiscoveryMockRecorder
}

// MockDiscoveryMockRecorder is the mock recorder for MockDiscovery
type MockDiscoveryMockRecorder struct {
	mock *MockDiscovery
}

// NewMockDiscovery creates a new mock instance
func NewMockDiscovery(ctrl *gomock.Controller) *MockDiscovery {
	mock := &MockDiscovery{ctrl: ctrl}
	mock.recorder = &MockDiscoveryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDiscovery) EXPECT() *MockDiscoveryMockRecorder {
	return m.recorder
}

// Advertise mocks base method
func (m *MockDiscovery) Advertise(arg0 context.Context, arg1 string, arg2 ...go_libp2p_discovery.Option) (time.Duration, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Advertise", varargs...)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Advertise indicates an expected call of Advertise
func (mr *MockDiscoveryMockRecorder) Advertise(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Advertise", reflect.TypeOf((*MockDiscovery)(nil).Advertise), varargs...)
}

// FindPeers mocks base method
func (m *MockDiscovery) FindPeers(arg0 context.Context, arg1 string, arg2 ...go_libp2p_discovery.Option) (<-chan go_libp2p_peerstore.PeerInfo, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindPeers", varargs...)
	ret0, _ := ret[0].(<-chan go_libp2p_peerstore.PeerInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPeers indicates an expected call of FindPeers
func (mr *MockDiscoveryMockRecorder) FindPeers(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPeers", reflect.TypeOf((*MockDiscovery)(nil).FindPeers), varargs...)
}
