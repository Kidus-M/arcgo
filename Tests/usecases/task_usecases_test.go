package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task_manager1/Domain"
	"task_manager1/Tests/mocks"
	"task_manager1/Usecases"
)

func TestCreateTask(t *testing.T) {
	repo := new(mocks.MockTaskRepository)
	uc := Usecases.NewTaskUsecase(repo)

	task := Domain.Task{Title: "Hello World"}

	repo.On("Create", mock.Anything, task).Return(task, nil)

	created, err := uc.Create(context.Background(), task)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", created.Title)
}
