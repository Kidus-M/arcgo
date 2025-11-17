package usecases

import (
	"context"
	"errors"

	"task_manager/domain"
	"task_manager/repositories"
)

// TaskUsecase defines task business rules
type TaskUsecase struct {
	repo repositories.TaskRepository
}

func NewTaskUsecase(r repositories.TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: r}
}

func (t *TaskUsecase) List(ctx context.Context) ([]domain.Task, error) {
	return t.repo.FindAll(ctx)
}

func (t *TaskUsecase) GetByID(ctx context.Context, id string) (domain.Task, error) {
	return t.repo.FindByID(ctx, id)
}

func (t *TaskUsecase) Create(ctx context.Context, input domain.Task) (domain.Task, error) {
	if input.Title == "" {
		return domain.Task{}, errors.New("title required")
	}
	return t.repo.Create(ctx, input)
}

func (t *TaskUsecase) Update(ctx context.Context, id string, input domain.Task) (domain.Task, error) {
	return t.repo.Update(ctx, id, input)
}

func (t *TaskUsecase) Delete(ctx context.Context, id string) (bool, error) {
	return t.repo.Delete(ctx, id)
}
