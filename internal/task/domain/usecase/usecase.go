package usecase

import (
	"context"
	"ggltask/internal/task"
	"ggltask/internal/task/domain/entities"
)

//go:generate mockgen -source=./usecase.go -destination=../../mock/usecasemock/usecase_mock.go -package=usecasemock
type TaskUseCase interface {
	CreateTask(ctx context.Context, param CreateTaskParams) (*entities.Task, error)
	ListTasks(ctx context.Context, param ListTasksParams) (*ListTasksResult, error)
	UpdateTask(ctx context.Context, param UpdateTaskParams) (*entities.Task, error)
	DeleteTask(ctx context.Context, id uint) error
}

type CreateTaskParams struct {
	Name string
}

type UpdateTaskParams struct {
	ID     uint
	Name   string
	Status task.TaskStatus
}

type ListTasksParams struct {
	PageIndex int
	PageSize  int
}

type ListTasksResult struct {
	Tasks []*entities.Task
	Total int
}
