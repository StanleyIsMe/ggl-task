// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package repositorymock is a generated GoMock package.
package repositorymock

import (
	context "context"
	entities "ggltask/internal/task/domain/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockRepository) CreateTask(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, task)
	ret0, _ := ret[0].(*entities.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockRepositoryMockRecorder) CreateTask(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockRepository)(nil).CreateTask), ctx, task)
}

// DeleteTask mocks base method.
func (m *MockRepository) DeleteTask(ctx context.Context, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockRepositoryMockRecorder) DeleteTask(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockRepository)(nil).DeleteTask), ctx, id)
}

// GetTaskByID mocks base method.
func (m *MockRepository) GetTaskByID(ctx context.Context, id uint) (*entities.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskByID", ctx, id)
	ret0, _ := ret[0].(*entities.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskByID indicates an expected call of GetTaskByID.
func (mr *MockRepositoryMockRecorder) GetTaskByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskByID", reflect.TypeOf((*MockRepository)(nil).GetTaskByID), ctx, id)
}

// ListTasksByPage mocks base method.
func (m *MockRepository) ListTasksByPage(ctx context.Context, pageIndex, pageSize int) ([]*entities.Task, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTasksByPage", ctx, pageIndex, pageSize)
	ret0, _ := ret[0].([]*entities.Task)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListTasksByPage indicates an expected call of ListTasksByPage.
func (mr *MockRepositoryMockRecorder) ListTasksByPage(ctx, pageIndex, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasksByPage", reflect.TypeOf((*MockRepository)(nil).ListTasksByPage), ctx, pageIndex, pageSize)
}

// UpdateTask mocks base method.
func (m *MockRepository) UpdateTask(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, task)
	ret0, _ := ret[0].(*entities.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockRepositoryMockRecorder) UpdateTask(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockRepository)(nil).UpdateTask), ctx, task)
}
