package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"task_manager1/Domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u Domain.User) (Domain.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (Domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) PromoteToAdmin(ctx context.Context, username string) (Domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
