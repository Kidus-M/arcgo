package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"task_manager1/Domain"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, t Domain.Task) (Domain.Task, error) {
	args := m.Called(ctx, t)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAll(ctx context.Context) ([]Domain.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id string) (Domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, id string, t Domain.Task) (Domain.Task, error) {
	args := m.Called(ctx, id, t)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) (bool, error) {
    args := m.Called(ctx, id)
    return args.Bool(0), args.Error(1)
}
