package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task_manager1/Domain"
	"task_manager1/Infrastructure/security"
	"task_manager1/Tests/mocks"
	"task_manager1/Usecases"
)

func TestRegisterFirstUserBecomesAdmin(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	pw := security.NewPasswordService()
	uc := Usecases.NewUserUsecase(repo, pw)

	repo.On("Count", mock.Anything).Return(int64(0), nil)
	repo.On("Create", mock.Anything, mock.Anything).
		Return(func(_ context.Context, u Domain.User) Domain.User {
			return u
		}, nil)

	user, err := uc.Register(context.Background(), "kidus", "123")
	assert.NoError(t, err)
	assert.Equal(t, "admin", user.Role)
}
