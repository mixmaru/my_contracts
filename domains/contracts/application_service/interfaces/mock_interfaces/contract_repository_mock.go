// Code generated by MockGen. DO NOT EDIT.
// Source: domains/contracts/application_service/interfaces/contract_repository_interface.go

// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

import (
	gomock "github.com/golang/mock/gomock"
	entities "github.com/mixmaru/my_contracts/domains/contracts/entities"
	gorp "gopkg.in/gorp.v2"
	reflect "reflect"
	time "time"
)

// MockIContractRepository is a mock of IContractRepository interface.
type MockIContractRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIContractRepositoryMockRecorder
}

// MockIContractRepositoryMockRecorder is the mock recorder for MockIContractRepository.
type MockIContractRepositoryMockRecorder struct {
	mock *MockIContractRepository
}

// NewMockIContractRepository creates a new mock instance.
func NewMockIContractRepository(ctrl *gomock.Controller) *MockIContractRepository {
	mock := &MockIContractRepository{ctrl: ctrl}
	mock.recorder = &MockIContractRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIContractRepository) EXPECT() *MockIContractRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIContractRepository) Create(contractEntity *entities.ContractEntity, executor gorp.SqlExecutor) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", contractEntity, executor)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIContractRepositoryMockRecorder) Create(contractEntity, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIContractRepository)(nil).Create), contractEntity, executor)
}

// GetById mocks base method.
func (m *MockIContractRepository) GetById(id int, executor gorp.SqlExecutor) (*entities.ContractEntity, *entities.ProductEntity, interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id, executor)
	ret0, _ := ret[0].(*entities.ContractEntity)
	ret1, _ := ret[1].(*entities.ProductEntity)
	ret2, _ := ret[2].(interface{})
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetById indicates an expected call of GetById.
func (mr *MockIContractRepositoryMockRecorder) GetById(id, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIContractRepository)(nil).GetById), id, executor)
}

// GetBillingTargetByBillingDate mocks base method.
func (m *MockIContractRepository) GetBillingTargetByBillingDate(billingDate time.Time, executor gorp.SqlExecutor) ([]*entities.ContractEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBillingTargetByBillingDate", billingDate, executor)
	ret0, _ := ret[0].([]*entities.ContractEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBillingTargetByBillingDate indicates an expected call of GetBillingTargetByBillingDate.
func (mr *MockIContractRepositoryMockRecorder) GetBillingTargetByBillingDate(billingDate, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBillingTargetByBillingDate", reflect.TypeOf((*MockIContractRepository)(nil).GetBillingTargetByBillingDate), billingDate, executor)
}

// GetRecurTargets mocks base method.
func (m *MockIContractRepository) GetRecurTargets(executeDate time.Time, executor gorp.SqlExecutor) ([]*entities.ContractEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecurTargets", executeDate, executor)
	ret0, _ := ret[0].([]*entities.ContractEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecurTargets indicates an expected call of GetRecurTargets.
func (mr *MockIContractRepositoryMockRecorder) GetRecurTargets(executeDate, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecurTargets", reflect.TypeOf((*MockIContractRepository)(nil).GetRecurTargets), executeDate, executor)
}

// Update mocks base method.
func (m *MockIContractRepository) Update(contractEntity *entities.ContractEntity, executor gorp.SqlExecutor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", contractEntity, executor)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIContractRepositoryMockRecorder) Update(contractEntity, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIContractRepository)(nil).Update), contractEntity, executor)
}

// GetHavingExpiredRightToUseContract mocks base method.
func (m *MockIContractRepository) GetHavingExpiredRightToUseContractIds(baseDate time.Time, executor gorp.SqlExecutor) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHavingExpiredRightToUseContractIds", baseDate, executor)
	ret0, _ := ret[0].([]*entities.ContractEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHavingExpiredRightToUseContract indicates an expected call of GetHavingExpiredRightToUseContract.
func (mr *MockIContractRepositoryMockRecorder) GetHavingExpiredRightToUseContract(baseDate, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHavingExpiredRightToUseContractIds", reflect.TypeOf((*MockIContractRepository)(nil).GetHavingExpiredRightToUseContractIds), baseDate, executor)
}
