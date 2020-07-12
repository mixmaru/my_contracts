// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domains/contracts/application_service/interfaces/product_repository_interface.go

// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

import (
	gomock "github.com/golang/mock/gomock"
	entities "github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	gorp "gopkg.in/gorp.v2"
	reflect "reflect"
)

// MockIProductRepository is a mock of IProductRepository interface.
type MockIProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIProductRepositoryMockRecorder
}

// MockIProductRepositoryMockRecorder is the mock recorder for MockIProductRepository.
type MockIProductRepositoryMockRecorder struct {
	mock *MockIProductRepository
}

// NewMockIProductRepository creates a new mock instance.
func NewMockIProductRepository(ctrl *gomock.Controller) *MockIProductRepository {
	mock := &MockIProductRepository{ctrl: ctrl}
	mock.recorder = &MockIProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIProductRepository) EXPECT() *MockIProductRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockIProductRepository) Save(productEntity *entities.ProductEntity, executor gorp.SqlExecutor) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", productEntity, executor)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockIProductRepositoryMockRecorder) Save(productEntity, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIProductRepository)(nil).Save), productEntity, executor)
}

// GetById mocks base method.
func (m *MockIProductRepository) GetById(id int, executor gorp.SqlExecutor) (*entities.ProductEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id, executor)
	ret0, _ := ret[0].(*entities.ProductEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIProductRepositoryMockRecorder) GetById(id, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIProductRepository)(nil).GetById), id, executor)
}

// GetByName mocks base method.
func (m *MockIProductRepository) GetByName(name string, executor gorp.SqlExecutor) (*entities.ProductEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", name, executor)
	ret0, _ := ret[0].(*entities.ProductEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockIProductRepositoryMockRecorder) GetByName(name, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockIProductRepository)(nil).GetByName), name, executor)
}
