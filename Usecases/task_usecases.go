package Usecases

import (
	"context"
	"errors"

	"task_manager1/Domain"
	"task_manager1/Repositories"
)

// TaskUsecase defines task business rules
type TaskUsecase struct {
	repo Repositories.TaskRepository
}

func NewTaskUsecase(r Repositories.TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: r}
}

func (t *TaskUsecase) List(ctx context.Context) ([]Domain.Task, error) {
	return t.repo.FindAll(ctx)
}

func (t *TaskUsecase) GetByID(ctx context.Context, id string) (Domain.Task, error) {
	return t.repo.FindByID(ctx, id)
}

func (t *TaskUsecase) Create(ctx context.Context, input Domain.Task) (Domain.Task, error) {
	if input.Title == "" {
		return Domain.Task{}, errors.New("title required")
	}
	return t.repo.Create(ctx, input)
}

func (t *TaskUsecase) Update(ctx context.Context, id string, input Domain.Task) (Domain.Task, error) {
	return t.repo.Update(ctx, id, input)
}

func (t *TaskUsecase) Delete(ctx context.Context, id string) (bool, error) {
	return t.repo.Delete(ctx, id)
}
