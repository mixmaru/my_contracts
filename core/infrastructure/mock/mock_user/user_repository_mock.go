// Code generated by MockGen. DO NOT EDIT.
// Source: core/application/users/i_user_repository.go

// Package mock_users is a generated GoMock package.
package mock_user

import (
	gomock "github.com/golang/mock/gomock"
	user "github.com/mixmaru/my_contracts/core/domain/models/user"
	gorp "gopkg.in/gorp.v2"
	reflect "reflect"
)

// MockIUserRepository is a mock of IUserRepository interface.
type MockIUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUserRepositoryMockRecorder
}

// MockIUserRepositoryMockRecorder is the mock recorder for MockIUserRepository.
type MockIUserRepositoryMockRecorder struct {
	mock *MockIUserRepository
}

// NewMockIUserRepository creates a new mock instance.
func NewMockIUserRepository(ctrl *gomock.Controller) *MockIUserRepository {
	mock := &MockIUserRepository{ctrl: ctrl}
	mock.recorder = &MockIUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserRepository) EXPECT() *MockIUserRepositoryMockRecorder {
	return m.recorder
}

// GetUserById mocks base method.
func (m *MockIUserRepository) GetUserById(id int, executor gorp.SqlExecutor) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", id, executor)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockIUserRepositoryMockRecorder) GetUserById(id, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockIUserRepository)(nil).GetUserById), id, executor)
}

// SaveUserIndividual mocks base method.
func (m *MockIUserRepository) SaveUserIndividual(userEntity *user.UserIndividualEntity, executor gorp.SqlExecutor) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserIndividual", userEntity, executor)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUserIndividual indicates an expected call of SaveUserIndividual.
func (mr *MockIUserRepositoryMockRecorder) SaveUserIndividual(userEntity, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserIndividual", reflect.TypeOf((*MockIUserRepository)(nil).SaveUserIndividual), userEntity, executor)
}

// GetUserIndividualById mocks base method.
func (m *MockIUserRepository) GetUserIndividualById(id int, executor gorp.SqlExecutor) (*user.UserIndividualEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIndividualById", id, executor)
	ret0, _ := ret[0].(*user.UserIndividualEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIndividualById indicates an expected call of GetUserIndividualById.
func (mr *MockIUserRepositoryMockRecorder) GetUserIndividualById(id, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIndividualById", reflect.TypeOf((*MockIUserRepository)(nil).GetUserIndividualById), id, executor)
}

// SaveUserCorporation mocks base method.
func (m *MockIUserRepository) SaveUserCorporation(userEntity *user.UserCorporationEntity, executor gorp.SqlExecutor) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserCorporation", userEntity, executor)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUserCorporation indicates an expected call of SaveUserCorporation.
func (mr *MockIUserRepositoryMockRecorder) SaveUserCorporation(userEntity, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserCorporation", reflect.TypeOf((*MockIUserRepository)(nil).SaveUserCorporation), userEntity, executor)
}

// GetUserCorporationById mocks base method.
func (m *MockIUserRepository) GetUserCorporationById(id int, executor gorp.SqlExecutor) (*user.UserCorporationEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCorporationById", id, executor)
	ret0, _ := ret[0].(*user.UserCorporationEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCorporationById indicates an expected call of GetUserCorporationById.
func (mr *MockIUserRepositoryMockRecorder) GetUserCorporationById(id, executor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCorporationById", reflect.TypeOf((*MockIUserRepository)(nil).GetUserCorporationById), id, executor)
}