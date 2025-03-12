package memory

import (
	"context"
	"ggltask/internal/task/domain/entities"
	"ggltask/internal/task/domain/repository"
	"sort"
	"sync"
	"time"
)

var _ repository.Repository = (*TaskRepository)(nil)

// TaskRepository is a repository for tasks.
// It is a memory repository that uses a map to store tasks.
type TaskRepository struct {
	mu     sync.RWMutex
	tasks  map[uint]*entities.Task
	lastID uint
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[uint]*entities.Task),
	}
}

// CreateTask is creating a new task.
func (r *TaskRepository) CreateTask(_ context.Context, taskEntity *entities.Task) (*entities.Task, error) {
	if taskEntity.Name == "" || len(taskEntity.Name) > 50 || !taskEntity.Status.Valid() {
		return nil, repository.ErrInvalidData
	}

	r.mu.Lock()
	taskEntity.ID = r.lastID + 1
	taskEntity.CreatedAt = time.Now()
	taskEntity.UpdatedAt = time.Now()

	r.lastID++
	r.tasks[taskEntity.ID] = taskEntity
	r.mu.Unlock()

	return taskEntity, nil
}

// GetTaskByID is getting a task by id.
func (r *TaskRepository) GetTaskByID(_ context.Context, id uint) (*entities.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, repository.ErrDataNotFound
	}

	return task, nil
}

// ListTasksByPage is listing tasks by page.
func (r *TaskRepository) ListTasksByPage(_ context.Context, pageIndex, pageSize int) ([]*entities.Task, int, error) {
	if pageIndex < 1 || pageSize < 1 {
		return nil, 0, repository.ErrInvalidData
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*entities.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	total := len(tasks)

	start := (pageIndex - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = len(tasks)
	}

	if start > total {
		return nil, 0, nil
	}

	return tasks[start:end], total, nil
}

// UpdateTask is updating a task.
func (r *TaskRepository) UpdateTask(_ context.Context, taskEntity *entities.Task) (*entities.Task, error) {
	if taskEntity.Name == "" || len(taskEntity.Name) > 50 || !taskEntity.Status.Valid() {
		return nil, repository.ErrInvalidData
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.tasks[taskEntity.ID]
	if !ok {
		return nil, repository.ErrDataNotFound
	}

	task.Name = taskEntity.Name
	task.Status = taskEntity.Status
	task.UpdatedAt = time.Now()
	r.tasks[taskEntity.ID] = task

	return task, nil
}

// DeleteTask is deleting a task.
func (r *TaskRepository) DeleteTask(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.tasks[id]
	if !ok {
		return repository.ErrDataNotFound
	}

	delete(r.tasks, id)

	return nil
}
