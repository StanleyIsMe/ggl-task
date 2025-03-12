package http

import (
	"encoding/json"
	"errors"
	"flag"
	"ggltask/internal/task"
	"ggltask/internal/task/domain/entities"
	"ggltask/internal/task/domain/usecase"
	"ggltask/internal/task/mock/usecasemock"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func TestTaskHandler_CreateTask(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name           string
		requestBody    string
		wantResponse   interface{}
		getUsecaseMock func(ctrl *gomock.Controller) usecase.TaskUseCase
		wantErr        bool
		wantStatusCode int
	}{
		{
			name:        "success",
			requestBody: `{"name": "test_name"}`,
			wantResponse: CreateTaskResponse{
				Task: &entities.Task{
					ID:        1,
					Name:      "test_name",
					Status:    task.TaskStatusIncomplete,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().CreateTask(gomock.Any(), usecase.CreateTaskParams{
					Name: "test_name",
				}).Return(&entities.Task{
					ID:        1,
					Name:      "test_name",
					Status:    task.TaskStatusIncomplete,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)

				return mockUsecase
			},
			wantErr:        false,
			wantStatusCode: 200,
		},
		{
			name:         "request body is invalid",
			requestBody:  `{"name": ""}`,
			wantResponse: InvalidRequestError(),
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				return usecasemock.NewMockTaskUseCase(ctrl)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:        "create task failed",
			requestBody: `{"name": "test_name"}`,
			wantResponse: ErrorResponse{
				ErrorCode:    "INTERNAL_SERVER_ERROR",
				ErrorMessage: "Internal Server Error",
			},
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().CreateTask(gomock.Any(), usecase.CreateTaskParams{
					Name: "test_name",
				}).Return(nil, errors.New("expected error"))

				return mockUsecase
			},
			wantErr:        true,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:         "task name too long",
			requestBody:  `{"name": "012345678901234567890123456789012345678901234567891"}`,
			wantResponse: InvalidRequestError(),
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				return usecasemock.NewMockTaskUseCase(ctrl)
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockUsecase := tt.getUsecaseMock(gomock.NewController(t))
			handler := NewTaskHandler(mockUsecase)
			
			router := gin.Default()
			router.POST("/tasks", handler.CreateTask)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/tasks", strings.NewReader(tt.requestBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)

			wantResponseJson, err := json.Marshal(tt.wantResponse)
			if err != nil {
				t.Fatalf("Failed to marshal wantResponse: %v", err)
			}

			assert.Equal(t, string(wantResponseJson), w.Body.String())
		})
	}
}

func TestTaskHandler_ListTasks(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name           string
		requestBody    string
		wantResponse   interface{}
		getUsecaseMock func(ctrl *gomock.Controller) usecase.TaskUseCase
		wantStatusCode int
	}{
		{
			name:        "success",
			requestBody: `page_index=1&page_size=5`,
			wantResponse: ListTasksResponse{
				Tasks: []*entities.Task{
					{
						ID:        1,
						Name:      "test_name",
						Status:    task.TaskStatusIncomplete,
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				Total: 1,
			},
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().ListTasks(gomock.Any(), usecase.ListTasksParams{
					PageIndex: 1,
					PageSize:  5,
				}).Return(&usecase.ListTasksResult{
					Tasks: []*entities.Task{
						{
							ID:        1,
							Name:      "test_name",
							Status:    task.TaskStatusIncomplete,
							CreatedAt: now,
							UpdatedAt: now,
						},
					},
					Total: 1,
				}, nil)

				return mockUsecase
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:         "request body is invalid",
			requestBody:  `page_index=3&page_size=-2`,
			wantResponse: InvalidRequestError(),
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				return usecasemock.NewMockTaskUseCase(ctrl)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:        "request body is empty",
			requestBody: ``,
			wantResponse: ListTasksResponse{
				Tasks: []*entities.Task{
					{
						ID:        1,
						Name:      "test_name",
						Status:    task.TaskStatusIncomplete,
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				Total: 1,
			},
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().ListTasks(gomock.Any(), usecase.ListTasksParams{
					PageIndex: 1,
					PageSize:  10,
				}).Return(&usecase.ListTasksResult{
					Tasks: []*entities.Task{
						{
							ID:        1,
							Name:      "test_name",
							Status:    task.TaskStatusIncomplete,
							CreatedAt: now,
							UpdatedAt: now,
						},
					},
					Total: 1,
				}, nil)

				return mockUsecase
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "internal server error",
			requestBody: `page_index=1&page_size=5`,
			wantResponse: ErrorResponse{
				ErrorCode:    "INTERNAL_SERVER_ERROR",
				ErrorMessage: "Internal Server Error",
			},
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().ListTasks(gomock.Any(), usecase.ListTasksParams{
					PageIndex: 1,
					PageSize:  5,
				}).Return(nil, errors.New("expected error"))

				return mockUsecase
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUsecase := tt.getUsecaseMock(gomock.NewController(t))
			handler := NewTaskHandler(mockUsecase)
			
			router := gin.Default()
			router.GET("/tasks", handler.ListTasks)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/tasks", strings.NewReader(tt.requestBody))
			req.URL.RawQuery = tt.requestBody

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)

			wantResponseJson, err := json.Marshal(tt.wantResponse)
			if err != nil {
				t.Fatalf("Failed to marshal wantResponse: %v", err)
			}

			assert.Equal(t, string(wantResponseJson), w.Body.String())
		})
	}
}

func TestTaskHandler_UpdateTask(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name           string
		requestBody    string
		wantResponse   interface{}
		getUsecaseMock func(ctrl *gomock.Controller) usecase.TaskUseCase
		wantStatusCode int
		url string
	}{
		{
			name:        "success",
			url:         "/tasks/1",
			requestBody: `{"name": "test_name", "status": 0}`,
			wantResponse: UpdateTaskResponse{
				Task: &entities.Task{
					ID:        1,
					Name:      "test_name",
					Status:    task.TaskStatusIncomplete,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().UpdateTask(gomock.Any(), usecase.UpdateTaskParams{
					ID:     1,
					Name:    "test_name",
					Status:  task.TaskStatusIncomplete,
				}).Return(&entities.Task{
					ID:        1,
					Name:      "test_name",
					Status:    task.TaskStatusIncomplete,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)

				return mockUsecase
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:         "request body is invalid",
			url:         "/tasks/1",
			requestBody:  `{"name": "test_name", "status": 3}`,
			wantResponse: InvalidRequestError(),
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				return usecasemock.NewMockTaskUseCase(ctrl)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:         "task id is not a number",
			url:         "/tasks/a",
			requestBody:  `{"name": "test_name", "status": 1}`,
			wantResponse: InvalidRequestError(),
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				return usecasemock.NewMockTaskUseCase(ctrl)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:         "internal server error",
			url:         "/tasks/1",
			requestBody:  `{"name": "test_name", "status": 1}`,
			wantResponse: ErrorResponse{
				ErrorCode:    "INTERNAL_SERVER_ERROR",
				ErrorMessage: "Internal Server Error",
			},
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().UpdateTask(gomock.Any(), usecase.UpdateTaskParams{
					ID:     1,
					Name:    "test_name",
					Status:  task.TaskStatusCompleted,
				}).Return(nil, errors.New("expected error"))

				return mockUsecase
			},	
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:         "task not found",
			url:         "/tasks/999",
			requestBody:  `{"name": "test_name", "status": 1}`,
			wantResponse: ErrorResponse{
				ErrorCode:    "NOT_FOUND",
				ErrorMessage: "task 999 not found",
			},
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().UpdateTask(gomock.Any(), usecase.UpdateTaskParams{
					ID:     999,
					Name:    "test_name",
					Status:  task.TaskStatusCompleted,
				}).Return(nil, usecase.NotFoundError{
					Resource: "task",
					ID:       999,
				})

				return mockUsecase
			},	
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUsecase := tt.getUsecaseMock(gomock.NewController(t))
			handler := NewTaskHandler(mockUsecase)
			
			router := gin.Default()
			router.PUT("/tasks/:id", handler.UpdateTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", tt.url, strings.NewReader(tt.requestBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)

			wantResponseJson, err := json.Marshal(tt.wantResponse)
			if err != nil {
				t.Fatalf("Failed to marshal wantResponse: %v", err)
			}

			assert.Equal(t, string(wantResponseJson), w.Body.String())
		})
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		url            string
		wantStatusCode int
		getUsecaseMock func(ctrl *gomock.Controller) usecase.TaskUseCase
	}{
		{
			name:         "success",
			url:         "/tasks/1",
			wantStatusCode: http.StatusOK,
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().DeleteTask(gomock.Any(), uint(1)).Return(nil)

				return mockUsecase
			},
		},
		{
			name:         "task id is not a number",
			url:         "/tasks/a",
			wantStatusCode: http.StatusBadRequest,
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				return usecasemock.NewMockTaskUseCase(ctrl)
			},
		},
		{
			name:         "internal server error",
			url:         "/tasks/1",
			wantStatusCode: http.StatusInternalServerError,
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().DeleteTask(gomock.Any(), uint(1)).Return(errors.New("expected error"))

				return mockUsecase
			},
		},
		{
			name:         "task not found",
			url:         "/tasks/999",
			wantStatusCode: http.StatusNotFound,
			getUsecaseMock: func(ctrl *gomock.Controller) usecase.TaskUseCase {
				mockUsecase := usecasemock.NewMockTaskUseCase(ctrl)
				mockUsecase.EXPECT().DeleteTask(gomock.Any(), uint(999)).Return(usecase.NotFoundError{
					Resource: "task",
					ID:       999,
				})

				return mockUsecase
			},
		},
	}	
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUsecase := tt.getUsecaseMock(gomock.NewController(t))
			handler := NewTaskHandler(mockUsecase)
			
			router := gin.Default()
			router.DELETE("/tasks/:id", handler.DeleteTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", tt.url, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
		})
	}
}
