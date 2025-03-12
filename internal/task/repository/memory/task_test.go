package memory

import (
	"context"
	"flag"
	"fmt"
	"ggltask/internal/task"
	"ggltask/internal/task/domain/entities"
	"ggltask/internal/task/domain/repository"
	"os"
	"reflect"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestNewTaskRepository(t *testing.T) {
	t.Parallel()

	if repo := NewTaskRepository(); reflect.TypeOf(repo) != reflect.TypeOf(&TaskRepository{}) {
		t.Errorf("NewTaskRepository() = %v, want %v", repo, &TaskRepository{})
	}
}

func TestTaskRepository_TestCreateTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		task    *entities.Task
		wantErr bool
	}{
		{
			name: "success",
			task: &entities.Task{
				Name:   "test task",
				Status: task.TaskStatusIncomplete,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			task: &entities.Task{
				Name:   "",
				Status: task.TaskStatusIncomplete,
			},
			wantErr: true,
		},
		{
			name: "very long name",
			task: &entities.Task{
				Name:   string(make([]byte, 1000)),
				Status: task.TaskStatusIncomplete,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			task: &entities.Task{
				Name:   "test task",
				Status: task.TaskStatus(2),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewTaskRepository()
			got, err := r.CreateTask(context.Background(), tt.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Error("CreateTask() got nil task")
				return
			}
			if got.ID == 0 {
				t.Error("CreateTask() got task with ID = 0")
			}
			if got.CreatedAt.IsZero() {
				t.Error("CreateTask() got task with zero CreatedAt")
			}
			if got.UpdatedAt.IsZero() {
				t.Error("CreateTask() got task with zero UpdatedAt")
			}
			if got.Name != tt.task.Name {
				t.Errorf("CreateTask() got task with Name = %v, want %v", got.Name, tt.task.Name)
			}
			if got.Status != tt.task.Status {
				t.Errorf("CreateTask() got task with Status = %v, want %v", got.Status, tt.task.Status)
			}
		})
	}
}

func TestTaskRepository_GetTaskByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      uint
		setup   func(*TaskRepository)
		want    *entities.Task
		wantErr error
	}{
		{
			name: "success",
			id:   1,
			setup: func(r *TaskRepository) {
				r.tasks[1] = &entities.Task{
					ID:     1,
					Name:   "test task",
					Status: task.TaskStatusIncomplete,
				}
			},
			want: &entities.Task{
				ID:     1,
				Name:   "test task",
				Status: task.TaskStatusIncomplete,
			},
		},
		{
			name:    "not found",
			id:      999,
			setup:   func(r *TaskRepository) {},
			wantErr: repository.ErrDataNotFound,
		},
		{
			name: "zero id",
			id:   0,
			setup: func(r *TaskRepository) {
				r.tasks[0] = &entities.Task{
					ID:     0,
					Name:   "test task",
					Status: task.TaskStatusIncomplete,
				}
			},
			want: &entities.Task{
				ID:     0,
				Name:   "test task",
				Status: task.TaskStatusIncomplete,
			},
		},
		{
			name: "max uint id",
			id:   ^uint(0),
			setup: func(r *TaskRepository) {
				r.tasks[^uint(0)] = &entities.Task{
					ID:     ^uint(0),
					Name:   "test task",
					Status: task.TaskStatusIncomplete,
				}
			},
			want: &entities.Task{
				ID:     ^uint(0),
				Name:   "test task",
				Status: task.TaskStatusIncomplete,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewTaskRepository()
			tt.setup(r)

			got, err := r.GetTaskByID(context.Background(), tt.id)
			if err != tt.wantErr {
				t.Errorf("GetTaskByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("GetTaskByID() got = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Name != tt.want.Name {
				t.Errorf("GetTaskByID() got = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Status != tt.want.Status {
				t.Errorf("GetTaskByID() got = %v, want %v", got.Status, tt.want.Status)
			}
		})
	}
}

func TestTaskRepository_ListTasksByPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pageIndex int
		pageSize  int
		setup     func(*TaskRepository)
		want      []*entities.Task
		wantTotal int
		wantErr   error
	}{
		{
			name:      "invalid page index",
			pageIndex: 0,
			pageSize:  10,
			setup:     func(r *TaskRepository) {},
			wantErr:   repository.ErrInvalidData,
		},
		{
			name:      "invalid page size",
			pageIndex: 1,
			pageSize:  0,
			setup:     func(r *TaskRepository) {},
			wantErr:   repository.ErrInvalidData,
		},
		{
			name:      "empty list",
			pageIndex: 1,
			pageSize:  10,
			setup:     func(r *TaskRepository) {},
			want:      []*entities.Task{},
			wantTotal: 0,
		},
		{
			name:      "single page",
			pageIndex: 1,
			pageSize:  10,
			setup: func(r *TaskRepository) {
				r.tasks[1] = &entities.Task{ID: 1, Name: "task 1", Status: task.TaskStatusIncomplete}
				r.tasks[2] = &entities.Task{ID: 2, Name: "task 2", Status: task.TaskStatusCompleted}
			},
			want: []*entities.Task{
				{ID: 1, Name: "task 1", Status: task.TaskStatusIncomplete},
				{ID: 2, Name: "task 2", Status: task.TaskStatusCompleted},
			},
			wantTotal: 2,
		},
		{
			name:      "multiple pages - first page",
			pageIndex: 1,
			pageSize:  2,
			setup: func(r *TaskRepository) {
				for i := uint(1); i <= 5; i++ {
					r.tasks[i] = &entities.Task{ID: i, Name: fmt.Sprintf("task %d", i), Status: task.TaskStatusIncomplete}
				}
			},
			want: []*entities.Task{
				{ID: 1, Name: "task 1", Status: task.TaskStatusIncomplete},
				{ID: 2, Name: "task 2", Status: task.TaskStatusIncomplete},
			},
			wantTotal: 5,
		},
		{
			name:      "multiple pages - last page",
			pageIndex: 3,
			pageSize:  2,
			setup: func(r *TaskRepository) {
				for i := uint(1); i <= 5; i++ {
					r.tasks[i] = &entities.Task{ID: i, Name: fmt.Sprintf("task %d", i), Status: task.TaskStatusIncomplete}
				}
			},
			want: []*entities.Task{
				{ID: 5, Name: "task 5", Status: task.TaskStatusIncomplete},
			},
			wantTotal: 5,
		},
		{
			name:      "page beyond total",
			pageIndex: 4,
			pageSize:  2,
			setup: func(r *TaskRepository) {
				for i := uint(1); i <= 5; i++ {
					r.tasks[i] = &entities.Task{ID: i, Name: fmt.Sprintf("task %d", i), Status: task.TaskStatusIncomplete}
				}
			},
			want:      nil,
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewTaskRepository()
			tt.setup(r)

			got, total, err := r.ListTasksByPage(context.Background(), tt.pageIndex, tt.pageSize)
			if err != tt.wantErr {
				t.Errorf("ListTasksByPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				return
			}

			if total != tt.wantTotal {
				t.Errorf("ListTasksByPage() total = %v, want %v", total, tt.wantTotal)
			}

			if len(got) != len(tt.want) {
				t.Errorf("ListTasksByPage() got = %v items, want %v items", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("ListTasksByPage() got[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Name != tt.want[i].Name {
					t.Errorf("ListTasksByPage() got[%d].Name = %v, want %v", i, got[i].Name, tt.want[i].Name)
				}
				if got[i].Status != tt.want[i].Status {
					t.Errorf("ListTasksByPage() got[%d].Status = %v, want %v", i, got[i].Status, tt.want[i].Status)
				}
			}
		})
	}
}

func TestTaskRepository_UpdateTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*TaskRepository)
		task    *entities.Task
		want    *entities.Task
		wantErr error
	}{
		{
			name: "success",
			setup: func(r *TaskRepository) {
				r.tasks[1] = &entities.Task{ID: 1, Name: "task 1", Status: task.TaskStatusIncomplete}
			},
			task: &entities.Task{ID: 1, Name: "updated task", Status: task.TaskStatusCompleted},
			want: &entities.Task{ID: 1, Name: "updated task", Status: task.TaskStatusCompleted},
		},
		{
			name:    "empty name",
			setup:   func(r *TaskRepository) {},
			task:    &entities.Task{ID: 1, Name: "", Status: task.TaskStatusCompleted},
			wantErr: repository.ErrInvalidData,
		},
		{
			name:    "name too long",
			setup:   func(r *TaskRepository) {},
			task:    &entities.Task{ID: 1, Name: string(make([]byte, 1000)), Status: task.TaskStatusCompleted},
			wantErr: repository.ErrInvalidData,
		},
		{
			name:    "invalid status",
			setup:   func(r *TaskRepository) {},
			task:    &entities.Task{ID: 1, Name: "task", Status: task.TaskStatus(2)},
			wantErr: repository.ErrInvalidData,
		},
		{
			name:    "task not found",
			setup:   func(r *TaskRepository) {},
			task:    &entities.Task{ID: 999, Name: "task", Status: task.TaskStatusCompleted},
			wantErr: repository.ErrDataNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewTaskRepository()
			tt.setup(r)

			got, err := r.UpdateTask(context.Background(), tt.task)
			if err != tt.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("UpdateTask() got.ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Name != tt.want.Name {
				t.Errorf("UpdateTask() got.Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Status != tt.want.Status {
				t.Errorf("UpdateTask() got.Status = %v, want %v", got.Status, tt.want.Status)
			}
			if got.UpdatedAt.IsZero() {
				t.Error("UpdateTask() got.UpdatedAt is zero")
			}
		})
	}
}

func TestTaskRepository_DeleteTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(r *TaskRepository)
		id      uint
		wantErr error
	}{
		{
			name: "success",
			setup: func(r *TaskRepository) {
				r.tasks[1] = &entities.Task{
					ID:     1,
					Name:   "task",
					Status: task.TaskStatusCompleted,
				}
			},
			id:      1,
			wantErr: nil,
		},
		{
			name:    "task not found",
			setup:   func(r *TaskRepository) {},
			id:      999,
			wantErr: repository.ErrDataNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewTaskRepository()
			tt.setup(r)

			err := r.DeleteTask(context.Background(), tt.id)
			if err != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				// Verify task was actually deleted
				_, err := r.GetTaskByID(context.Background(), tt.id)
				if err != repository.ErrDataNotFound {
					t.Errorf("DeleteTask() task still exists after deletion")
				}
			}
		})
	}
}

