package usecase

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"context"
	"errors"
	"ggltask/internal/task"
	"ggltask/internal/task/domain/entities"
	"ggltask/internal/task/domain/repository"
	"ggltask/internal/task/domain/usecase"
	"ggltask/internal/task/mock/repositorymock"

	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

func TestNewTaskUseCaseImpl(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymock.NewMockRepository(ctrl)
	uc := NewTaskUseCaseImpl(mockRepo)

	if reflect.TypeOf(uc) != reflect.TypeOf(&TaskUseCaseImpl{}) {
		t.Errorf("NewTaskUseCaseImpl() = %v, want %v", uc, &TaskUseCaseImpl{})
	}
}

func TestTaskUseCaseImpl_CreateTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		param    usecase.CreateTaskParams
		mockRepo func(ctrl *gomock.Controller) repository.Repository
		want     *entities.Task
		wantErr  bool
	}{
		{
			name: "success",
			param: usecase.CreateTaskParams{
				Name: "test task",
			},
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().CreateTask(gomock.Any(), &entities.Task{
					Name:   "test task",
					Status: task.TaskStatusIncomplete,
				}).Return(&entities.Task{
					ID:        1,
					Name:      "test task",
					Status:    task.TaskStatusIncomplete,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)

				return mockRepo
			},
			want: &entities.Task{
				ID:        1,
				Name:      "test task",
				Status:    task.TaskStatusIncomplete,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "repository error",
			param: usecase.CreateTaskParams{
				Name: "test task",
			},
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(nil, errors.New("repository error"))

				return mockRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := tt.mockRepo(gomock.NewController(t))
			uc := NewTaskUseCaseImpl(mockRepo)

			got, err := uc.CreateTask(context.Background(), tt.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Status, got.Status)
		})
	}
}

func TestTaskUseCaseImpl_ListTasks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		param    usecase.ListTasksParams
		mockRepo func(ctrl *gomock.Controller) repository.Repository
		want     *usecase.ListTasksResult
		wantErr  bool
	}{
		{
			name: "success",
			param: usecase.ListTasksParams{
				PageIndex: 1,
				PageSize:  10,
			},
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().ListTasksByPage(gomock.Any(), 1, 10).Return([]*entities.Task{
					{
						ID:        1,
						Name:      "test task",
						Status:    task.TaskStatusIncomplete,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, 1, nil)

				return mockRepo
			},
			want: &usecase.ListTasksResult{
				Tasks: []*entities.Task{
					{
						ID:        1,
						Name:      "test task",
						Status:    task.TaskStatusIncomplete,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
				Total: 1,
			},
			wantErr: false,
		},
		{
			name: "repository error",
			param: usecase.ListTasksParams{
				PageIndex: 1,
				PageSize:  10,
			},
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().ListTasksByPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, 0, errors.New("repository error"))

				return mockRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := tt.mockRepo(gomock.NewController(t))
			uc := NewTaskUseCaseImpl(mockRepo)

			got, err := uc.ListTasks(context.Background(), tt.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.Total, got.Total)
			assert.Equal(t, len(tt.want.Tasks), len(got.Tasks))
			for i := range tt.want.Tasks {
				assert.Equal(t, tt.want.Tasks[i].ID, got.Tasks[i].ID)
				assert.Equal(t, tt.want.Tasks[i].Name, got.Tasks[i].Name)
				assert.Equal(t, tt.want.Tasks[i].Status, got.Tasks[i].Status)
			}
		})
	}
}

func TestTaskUseCaseImpl_UpdateTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		param    usecase.UpdateTaskParams
		mockRepo func(ctrl *gomock.Controller) repository.Repository
		want     *entities.Task
		wantErr  bool
	}{
		{
			name: "success",
			param: usecase.UpdateTaskParams{
				ID:     1,
				Name:   "updated task",
				Status: task.TaskStatusCompleted,
			},
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().UpdateTask(gomock.Any(), &entities.Task{
					ID:     1,
					Name:   "updated task",
					Status: task.TaskStatusCompleted,
				}).Return(&entities.Task{
					ID:        1,
					Name:      "updated task",
					Status:    task.TaskStatusCompleted,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)

				return mockRepo
			},
			want: &entities.Task{
				ID:        1,
				Name:      "updated task",
				Status:    task.TaskStatusCompleted,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "not found",
			param: usecase.UpdateTaskParams{
				ID:     1,
				Name:   "updated task",
				Status: task.TaskStatusCompleted,
			},
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil, repository.ErrDataNotFound)

				return mockRepo
			},
			wantErr: true,
		},
		{
			name: "repository error",
			param: usecase.UpdateTaskParams{
				ID:     1,
				Name:   "updated task",
				Status: task.TaskStatusCompleted,
			},
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil, errors.New("repository error"))

				return mockRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := tt.mockRepo(gomock.NewController(t))
			uc := NewTaskUseCaseImpl(mockRepo)

			got, err := uc.UpdateTask(context.Background(), tt.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Status, got.Status)
		})
	}
}

func TestTaskUseCaseImpl_DeleteTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       uint
		mockRepo func(ctrl *gomock.Controller) repository.Repository
		wantErr  bool
	}{
		{
			name: "success",
			id:   1,
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().DeleteTask(gomock.Any(), uint(1)).Return(nil)

				return mockRepo
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   1,
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().DeleteTask(gomock.Any(), gomock.Any()).Return(repository.ErrDataNotFound)

				return mockRepo
			},
			wantErr: true,
		},
		{
			name: "repository error",
			id:   1,
			mockRepo: func(ctrl *gomock.Controller) repository.Repository {
				mockRepo := repositorymock.NewMockRepository(ctrl)
				mockRepo.EXPECT().DeleteTask(gomock.Any(), gomock.Any()).Return(errors.New("repository error"))

				return mockRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := tt.mockRepo(gomock.NewController(t))
			uc := NewTaskUseCaseImpl(mockRepo)

			err := uc.DeleteTask(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
