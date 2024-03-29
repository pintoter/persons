// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/pintoter/persons/services/command/internal/entity"
	service "github.com/pintoter/persons/services/command/internal/service"
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

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, person entity.Person) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, person)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, person interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, person)
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, id)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, id int, params *service.UpdateParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, id, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, id, params)
}

// MockGenerator is a mock of Generator interface.
type MockGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockGeneratorMockRecorder
}

// MockGeneratorMockRecorder is the mock recorder for MockGenerator.
type MockGeneratorMockRecorder struct {
	mock *MockGenerator
}

// NewMockGenerator creates a new mock instance.
func NewMockGenerator(ctrl *gomock.Controller) *MockGenerator {
	mock := &MockGenerator{ctrl: ctrl}
	mock.recorder = &MockGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenerator) EXPECT() *MockGeneratorMockRecorder {
	return m.recorder
}

// GenerateAge mocks base method.
func (m *MockGenerator) GenerateAge(ctx context.Context, name string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAge", ctx, name)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAge indicates an expected call of GenerateAge.
func (mr *MockGeneratorMockRecorder) GenerateAge(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAge", reflect.TypeOf((*MockGenerator)(nil).GenerateAge), ctx, name)
}

// GenerateGender mocks base method.
func (m *MockGenerator) GenerateGender(ctx context.Context, name string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateGender", ctx, name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateGender indicates an expected call of GenerateGender.
func (mr *MockGeneratorMockRecorder) GenerateGender(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateGender", reflect.TypeOf((*MockGenerator)(nil).GenerateGender), ctx, name)
}

// GenerateNationalize mocks base method.
func (m *MockGenerator) GenerateNationalize(ctx context.Context, name string) ([]entity.Nationality, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateNationalize", ctx, name)
	ret0, _ := ret[0].([]entity.Nationality)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateNationalize indicates an expected call of GenerateNationalize.
func (mr *MockGeneratorMockRecorder) GenerateNationalize(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateNationalize", reflect.TypeOf((*MockGenerator)(nil).GenerateNationalize), ctx, name)
}
