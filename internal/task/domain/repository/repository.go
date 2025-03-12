package repository

import (
	"context"
	"ggltask/internal/task/domain/entities"
)

//go:generate mockgen -source=./repository.go -destination=../../mock/repositorymock/repository_mock.go -package=repositorymock
type Repository interface {
	CreateTask(ctx context.Context, task *entities.Task) (*entities.Task, error)
	GetTaskByID(ctx context.Context, id uint) (*entities.Task, error)
	ListTasksByPage(ctx context.Context, pageIndex, pageSize int) ([]*entities.Task, int, error)
	UpdateTask(ctx context.Context, task *entities.Task) (*entities.Task, error)
	DeleteTask(ctx context.Context, id uint) error
}
