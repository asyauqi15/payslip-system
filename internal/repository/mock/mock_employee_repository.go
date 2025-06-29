// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/asyauqi15/payslip-system/internal/repository (interfaces: EmployeeRepository)
//
// Generated by this command:
//
//	mockgen -destination=./mock/mock_employee_repository.go -package=mock github.com/asyauqi15/payslip-system/internal/repository EmployeeRepository
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/asyauqi15/payslip-system/internal/entity"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockEmployeeRepository is a mock of EmployeeRepository interface.
type MockEmployeeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEmployeeRepositoryMockRecorder
	isgomock struct{}
}

// MockEmployeeRepositoryMockRecorder is the mock recorder for MockEmployeeRepository.
type MockEmployeeRepositoryMockRecorder struct {
	mock *MockEmployeeRepository
}

// NewMockEmployeeRepository creates a new mock instance.
func NewMockEmployeeRepository(ctrl *gomock.Controller) *MockEmployeeRepository {
	mock := &MockEmployeeRepository{ctrl: ctrl}
	mock.recorder = &MockEmployeeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmployeeRepository) EXPECT() *MockEmployeeRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockEmployeeRepository) Create(ctx context.Context, o *entity.Employee, tx *gorm.DB) (*entity.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, o, tx)
	ret0, _ := ret[0].(*entity.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockEmployeeRepositoryMockRecorder) Create(ctx, o, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEmployeeRepository)(nil).Create), ctx, o, tx)
}

// FindByID mocks base method.
func (m *MockEmployeeRepository) FindByID(ctx context.Context, i uint, tx *gorm.DB) (*entity.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, i, tx)
	ret0, _ := ret[0].(*entity.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockEmployeeRepositoryMockRecorder) FindByID(ctx, i, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockEmployeeRepository)(nil).FindByID), ctx, i, tx)
}

// FindByTemplate mocks base method.
func (m *MockEmployeeRepository) FindByTemplate(ctx context.Context, t *entity.Employee, tx *gorm.DB) ([]entity.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByTemplate", ctx, t, tx)
	ret0, _ := ret[0].([]entity.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByTemplate indicates an expected call of FindByTemplate.
func (mr *MockEmployeeRepositoryMockRecorder) FindByTemplate(ctx, t, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByTemplate", reflect.TypeOf((*MockEmployeeRepository)(nil).FindByTemplate), ctx, t, tx)
}

// FindOneByTemplate mocks base method.
func (m *MockEmployeeRepository) FindOneByTemplate(ctx context.Context, o *entity.Employee, tx *gorm.DB) (*entity.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByTemplate", ctx, o, tx)
	ret0, _ := ret[0].(*entity.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByTemplate indicates an expected call of FindOneByTemplate.
func (mr *MockEmployeeRepositoryMockRecorder) FindOneByTemplate(ctx, o, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByTemplate", reflect.TypeOf((*MockEmployeeRepository)(nil).FindOneByTemplate), ctx, o, tx)
}

// Save mocks base method.
func (m *MockEmployeeRepository) Save(ctx context.Context, o *entity.Employee, tx *gorm.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, o, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockEmployeeRepositoryMockRecorder) Save(ctx, o, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockEmployeeRepository)(nil).Save), ctx, o, tx)
}

// Updates mocks base method.
func (m *MockEmployeeRepository) Updates(ctx context.Context, o *entity.Employee, u entity.Employee, tx *gorm.DB) (*entity.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Updates", ctx, o, u, tx)
	ret0, _ := ret[0].(*entity.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Updates indicates an expected call of Updates.
func (mr *MockEmployeeRepositoryMockRecorder) Updates(ctx, o, u, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Updates", reflect.TypeOf((*MockEmployeeRepository)(nil).Updates), ctx, o, u, tx)
}
