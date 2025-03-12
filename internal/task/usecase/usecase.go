package usecase

import (
	"context"
	"errors"
	"fmt"
	"ggltask/internal/task"
	"ggltask/internal/task/domain/entities"
	"ggltask/internal/task/domain/repository"
	"ggltask/internal/task/domain/usecase"
)

var _ usecase.TaskUseCase = (*TaskUseCaseImpl)(nil)

type TaskUseCaseImpl struct {
	taskRepo repository.Repository
}

func NewTaskUseCaseImpl(taskRepo repository.Repository) *TaskUseCaseImpl {
	return &TaskUseCaseImpl{taskRepo: taskRepo}
}

// CreateTask is responsible for creating a new task.
func (a *TaskUseCaseImpl) CreateTask(ctx context.Context, param usecase.CreateTaskParams) (*entities.Task, error) {
	entityTask := &entities.Task{
		Name:   param.Name,
		Status: task.TaskStatusIncomplete,
	}

	newTask, err := a.taskRepo.CreateTask(ctx, entityTask)
	if err != nil {
		return nil, fmt.Errorf("repo.CreateTask error: %w", err)
	}

	return newTask, nil
}

// ListTasks is responsible for listing tasks by page.
func (a *TaskUseCaseImpl) ListTasks(ctx context.Context, param usecase.ListTasksParams) (*usecase.ListTasksResult, error) {
	tasks, total, err := a.taskRepo.ListTasksByPage(ctx, param.PageIndex, param.PageSize)
	if err != nil {
		return nil, fmt.Errorf("repo.ListTasksByPage error: %w", err)
	}

	return &usecase.ListTasksResult{
		Tasks: tasks,
		Total: total,
	}, nil
}

// UpdateTask is responsible for updating a task.
func (a *TaskUseCaseImpl) UpdateTask(ctx context.Context, param usecase.UpdateTaskParams) (*entities.Task, error) {
	entityTask := &entities.Task{
		ID:     param.ID,
		Name:   param.Name,
		Status: param.Status,
	}

	updatedTask, err := a.taskRepo.UpdateTask(ctx, entityTask)
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, usecase.NotFoundError{
				Resource: "task",
				ID:       param.ID,
			}
		}

		return nil, fmt.Errorf("repo.UpdateTask error: %w", err)
	}

	return updatedTask, nil
}

// DeleteTask is responsible for deleting a task.
func (a *TaskUseCaseImpl) DeleteTask(ctx context.Context, id uint) error {
	if err := a.taskRepo.DeleteTask(ctx, id); err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return usecase.NotFoundError{
				Resource: "task",
				ID:       id,
			}
		}

		return fmt.Errorf("repo.DeleteTask error: %w", err)
	}

	return nil
}
