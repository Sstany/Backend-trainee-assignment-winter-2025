// Code generated by MockGen. DO NOT EDIT.
// Source: ./user_transaction.go

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"
	entity "shop/internal/app/entity"
	port "shop/internal/app/port"

	gomock "github.com/golang/mock/gomock"
)

// MockUserTransactionRepo is a mock of UserTransactionRepo interface.
type MockUserTransactionRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserTransactionRepoMockRecorder
}

// MockUserTransactionRepoMockRecorder is the mock recorder for MockUserTransactionRepo.
type MockUserTransactionRepoMockRecorder struct {
	mock *MockUserTransactionRepo
}

// NewMockUserTransactionRepo creates a new mock instance.
func NewMockUserTransactionRepo(ctrl *gomock.Controller) *MockUserTransactionRepo {
	mock := &MockUserTransactionRepo{ctrl: ctrl}
	mock.recorder = &MockUserTransactionRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserTransactionRepo) EXPECT() *MockUserTransactionRepoMockRecorder {
	return m.recorder
}

// GetRecievedOperations mocks base method.
func (m *MockUserTransactionRepo) GetRecievedOperations(ctx context.Context, username string) ([]entity.Received, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecievedOperations", ctx, username)
	ret0, _ := ret[0].([]entity.Received)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecievedOperations indicates an expected call of GetRecievedOperations.
func (mr *MockUserTransactionRepoMockRecorder) GetRecievedOperations(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecievedOperations", reflect.TypeOf((*MockUserTransactionRepo)(nil).GetRecievedOperations), ctx, username)
}

// GetSentOperations mocks base method.
func (m *MockUserTransactionRepo) GetSentOperations(ctx context.Context, username string) ([]entity.Sent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSentOperations", ctx, username)
	ret0, _ := ret[0].([]entity.Sent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSentOperations indicates an expected call of GetSentOperations.
func (mr *MockUserTransactionRepoMockRecorder) GetSentOperations(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSentOperations", reflect.TypeOf((*MockUserTransactionRepo)(nil).GetSentOperations), ctx, username)
}

// SetUserTransaction mocks base method.
func (m *MockUserTransactionRepo) SetUserTransaction(tx port.Transaction, sendCoin entity.SendCoinRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserTransaction", tx, sendCoin)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUserTransaction indicates an expected call of SetUserTransaction.
func (mr *MockUserTransactionRepoMockRecorder) SetUserTransaction(tx, sendCoin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserTransaction", reflect.TypeOf((*MockUserTransactionRepo)(nil).SetUserTransaction), tx, sendCoin)
}
