// Code generated by MockGen. DO NOT EDIT.
// Source: ./shop.go

// Package repo is a generated GoMock package.
package repo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockShopRepo is a mock of ShopRepo interface.
type MockShopRepo struct {
	ctrl     *gomock.Controller
	recorder *MockShopRepoMockRecorder
}

// MockShopRepoMockRecorder is the mock recorder for MockShopRepo.
type MockShopRepoMockRecorder struct {
	mock *MockShopRepo
}

// NewMockShopRepo creates a new mock instance.
func NewMockShopRepo(ctrl *gomock.Controller) *MockShopRepo {
	mock := &MockShopRepo{ctrl: ctrl}
	mock.recorder = &MockShopRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShopRepo) EXPECT() *MockShopRepoMockRecorder {
	return m.recorder
}

// GetItemPrice mocks base method.
func (m *MockShopRepo) GetItemPrice(name string) (int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemPrice", name)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetItemPrice indicates an expected call of GetItemPrice.
func (mr *MockShopRepoMockRecorder) GetItemPrice(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemPrice", reflect.TypeOf((*MockShopRepo)(nil).GetItemPrice), name)
}
